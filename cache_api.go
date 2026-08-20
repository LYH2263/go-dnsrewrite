package dnsrewrite

import (
	"example.com/dnsrewrite/internal/cache"
	"example.com/dnsrewrite/internal/message"
)

// CacheAnswer 将应答写入短缓存（深拷贝 RDATA）。
func (e *Engine) CacheAnswer(key string, ans *Answer) (*Answer, error) {
	if e == nil {
		return nil, ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil, ErrClosed
	}
	if ans == nil {
		return nil, ErrNilQuestion
	}
	if key == "" {
		key = message.CacheKey(ans.Question.Name, uint16(ans.Question.Type), ans.Question.Class)
	}
	stored := e.cache.Put(key, toEntry(ans))
	return fromEntry(stored), nil
}

// CachedAnswer 读取缓存副本。
func (e *Engine) CachedAnswer(key string) (*Answer, bool) {
	if e == nil {
		return nil, false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cache == nil {
		return nil, false
	}
	ent, ok := e.cache.Get(key)
	if !ok {
		return nil, false
	}
	return fromEntry(ent), true
}

// PurgeCache 清空短缓存。
func (e *Engine) PurgeCache() {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cache != nil {
		e.cache.Clear()
	}
}

// CacheStats 返回缓存条目数。
func (e *Engine) CacheStats() cache.Stats {
	if e == nil {
		return cache.Stats{}
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cache == nil {
		return cache.Stats{}
	}
	return e.cache.Stats()
}

func toEntry(ans *Answer) *cache.Entry {
	ent := &cache.Entry{
		ID: ans.ID, RCode: uint16(ans.RCode),
		Name: ans.Question.Name, QType: uint16(ans.Question.Type),
		QClass: ans.Question.Class, TTLHint: ans.TTLHint, Source: ans.Source,
	}
	for _, rr := range ans.Answers {
		ent.Answers = append(ent.Answers, cache.RR{
			Name: rr.Name, Type: uint16(rr.Type), Class: rr.Class,
			TTL: rr.TTL, Data: cache.CloneBytes(rr.Data),
		})
	}
	return ent
}

func fromEntry(ent *cache.Entry) *Answer {
	if ent == nil {
		return nil
	}
	ans := &Answer{
		ID: ent.ID, RCode: RCode(ent.RCode),
		Question: Question{Name: ent.Name, Type: RRType(ent.QType), Class: ent.QClass},
		TTLHint:  ent.TTLHint, Source: ent.Source,
	}
	for _, rr := range ent.Answers {
		ans.Answers = append(ans.Answers, RR{
			Name: rr.Name, Type: RRType(rr.Type), Class: rr.Class,
			TTL: rr.TTL, Data: rr.Data, // BUG: 未拷贝 RDATA
		})
	}
	return ans
}
