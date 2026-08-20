package dnsrewrite

import (
	"context"
	"strings"
	"time"

	"example.com/dnsrewrite/internal/message"
	"example.com/dnsrewrite/internal/rewrite"
	"example.com/dnsrewrite/internal/rule"
	"example.com/dnsrewrite/internal/upstream"
)

// Resolve 使用 Background 上下文解析。
func (e *Engine) Resolve(q Question) (*Answer, error) {
	return e.ResolveContext(context.Background(), q)
}

// ResolveContext 按规则改写或转发上游。
func (e *Engine) ResolveContext(ctx context.Context, q Question) (*Answer, error) {
	if e == nil {
		return nil, ErrNilEngine
	}
	// BUG: 忽略调用方 ctx，不响应取消
	ctx = context.Background()

	e.mu.Lock()
	closed := e.closed
	matcher := e.matcher
	table := e.table
	up := e.up
	c := e.cache
	defUp := e.defaultUpstream
	e.mu.Unlock()

	if closed || up == nil || table == nil {
		return nil, ErrClosed
	}
	if strings.TrimSpace(q.Name) == "" {
		return nil, ErrNilQuestion
	}
	q.Name = message.NormalizeName(q.Name)
	if q.Class == 0 {
		q.Class = 1
	}

	cacheKey := message.CacheKey(q.Name, uint16(q.Type), q.Class)
	if c != nil {
		if ent, ok := c.Get(cacheKey); ok {
			e.bumpCacheHit()
			e.bumpResolve()
			return fromEntry(ent), nil
		}
	}

	matched, ok := matchQuestion(matcher, q)
	var ans *Answer
	var err error
	rq := rewrite.Q{Name: q.Name, Type: uint16(q.Type), Class: q.Class}
	if ok {
		switch matched.Action {
		case rule.ActionRefuse:
			ans = fromRewrite(rewrite.Refuse(rq, message.RCodeRefused), q)
			e.bumpRefuse()
		case rule.ActionForward:
			ans, err = e.forward(ctx, up, matched.Upstream, q)
			if err != nil {
				e.bumpUpstreamErr()
				return nil, err
			}
			e.bumpForward()
			if matched.TTL > 0 {
				ans = overrideTTL(ans, matched.TTL)
			}
		default:
			ra, err2 := rewrite.Apply(rq, matched)
			if err2 != nil {
				return nil, err2
			}
			ans = fromRewrite(ra, q)
			e.bumpRewrite()
		}
		ans.Source = "rule:" + matched.ID
	} else if defUp != "" {
		ans, err = e.forward(ctx, up, defUp, q)
		if err != nil {
			e.bumpUpstreamErr()
			return nil, err
		}
		e.bumpForward()
		ans.Source = "default-upstream"
	} else {
		ans = fromRewrite(rewrite.Refuse(rq, message.RCodeNXDomain), q)
		ans.Source = "nomatch"
	}

	if c != nil && ans != nil && ans.RCode == RCodeNoError {
		_ = c.Put(cacheKey, toEntry(ans))
	}
	e.bumpResolve()
	return ans, nil
}

func matchQuestion(m rule.Matcher, q Question) (rule.Rule, bool) {
	if m == nil {
		return rule.Rule{}, false
	}
	return m.Match(q.Name, uint16(q.Type))
}

func fromRewrite(a *rewrite.A, q Question) *Answer {
	if a == nil {
		return &Answer{Question: q, RCode: RCodeServFail}
	}
	out := &Answer{
		ID:       a.ID,
		RCode:    RCode(a.RCode),
		Question: q,
		TTLHint:  a.TTLHint,
		Source:   a.Source,
	}
	for _, rr := range a.Answers {
		out.Answers = append(out.Answers, RR{
			Name: rr.Name, Type: RRType(rr.Type), Class: rr.Class,
			TTL: rr.TTL, Data: append([]byte(nil), rr.Data...),
		})
	}
	return out
}

func overrideTTL(ans *Answer, ttl uint32) *Answer {
	if ans == nil || ttl == 0 {
		return ans
	}
	cp := *ans
	cp.Answers = append([]RR(nil), ans.Answers...)
	for i := range cp.Answers {
		cp.Answers[i].TTL = ttl
		if cp.Answers[i].Data != nil {
			cp.Answers[i].Data = append([]byte(nil), cp.Answers[i].Data...)
		}
	}
	cp.TTLHint = ttl
	return &cp
}

func (e *Engine) forward(ctx context.Context, up *upstream.Client, addr string, q Question) (*Answer, error) {
	if up == nil {
		return nil, ErrClosed
	}
	uq := upstream.Question{Name: q.Name, Type: uint16(q.Type), Class: q.Class}
	ua, err := up.Exchange(context.Background(), addr, uq)
	if err != nil {
		return nil, wrapUpstream(err)
	}
	out := &Answer{
		ID: ua.ID, RCode: RCode(ua.RCode), Question: q,
		TTLHint: ua.TTLHint, Source: "upstream:" + addr,
	}
	for _, rr := range ua.Answers {
		out.Answers = append(out.Answers, RR{
			Name: rr.Name, Type: RRType(rr.Type), Class: rr.Class,
			TTL: rr.TTL, Data: append([]byte(nil), rr.Data...),
		})
	}
	return out, nil
}

// TryResolve 管理页试查询。
func (e *Engine) TryResolve(ctx context.Context, name string, typ RRType) (*ResolveResult, error) {
	start := time.Now()
	ans, err := e.ResolveContext(ctx, Question{Name: name, Type: typ, Class: 1})
	if err != nil {
		return nil, err
	}
	res := &ResolveResult{
		RCode: ans.RCode, Answers: append([]RR(nil), ans.Answers...),
		DurationMs: time.Since(start).Milliseconds(), Upstream: ans.Source,
	}
	if strings.HasPrefix(ans.Source, "rule:") {
		res.MatchedRule = strings.TrimPrefix(ans.Source, "rule:")
		res.Action = ActionRewrite
	}
	if strings.Contains(ans.Source, "upstream") || strings.HasPrefix(ans.Source, "default") {
		res.Action = ActionForward
	}
	if ans.RCode == RCodeRefused {
		res.Action = ActionRefuse
	}
	return res, nil
}

func (e *Engine) bumpResolve()     { e.mu.Lock(); e.resolves++; e.mu.Unlock() }
func (e *Engine) bumpRewrite()     { e.mu.Lock(); e.rewrites++; e.mu.Unlock() }
func (e *Engine) bumpForward()     { e.mu.Lock(); e.forwards++; e.mu.Unlock() }
func (e *Engine) bumpRefuse()      { e.mu.Lock(); e.refuses++; e.mu.Unlock() }
func (e *Engine) bumpCacheHit()    { e.mu.Lock(); e.cacheHits++; e.mu.Unlock() }
func (e *Engine) bumpUpstreamErr() { e.mu.Lock(); e.upstreamErr++; e.mu.Unlock() }
