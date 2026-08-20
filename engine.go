package dnsrewrite

import (
	"sync"
	"time"

	"example.com/dnsrewrite/internal/cache"
	"example.com/dnsrewrite/internal/clock"
	"example.com/dnsrewrite/internal/rule"
	"example.com/dnsrewrite/internal/upstream"
)

const (
	defaultUpstreamTimeout = 5 * time.Second
	defaultCacheTTL        = 30 * time.Second
	defaultMaxRules        = 4096
)

// Engine DNS 改写引擎门面。零值不可用，须经 New 构造。
type Engine struct {
	mu sync.Mutex

	closed          bool
	clk             clock.Clock
	table           *rule.Table
	matcher         rule.Matcher
	matcherOverride bool

	up              *upstream.Client
	cache           *cache.Cache
	persistPath     string
	upstreamTimeout time.Duration
	cacheTTL        time.Duration
	defaultUpstream string
	maxRules        int
	dirty           bool

	resolves    uint64
	rewrites    uint64
	forwards    uint64
	refuses     uint64
	cacheHits   uint64
	upstreamErr uint64
}

// New 构造 Engine。缺省安装规则表匹配器与上游客户端。
func New(opts ...Option) *Engine {
	e := &Engine{
		clk:             clock.Real{},
		table:           rule.NewTable(),
		up:              upstream.NewClient(defaultUpstreamTimeout),
		cache:           cache.New(256),
		upstreamTimeout: defaultUpstreamTimeout,
		cacheTTL:        defaultCacheTTL,
		maxRules:        defaultMaxRules,
	}
	for _, o := range opts {
		if o != nil {
			o(e)
		}
	}
	if e.clk == nil {
		e.clk = clock.Real{}
	}
	if e.table == nil {
		e.table = rule.NewTable()
	}
	if e.up == nil {
		e.up = upstream.NewClient(e.upstreamTimeout)
	}
	if e.cache == nil {
		e.cache = cache.New(256)
	}
	if e.upstreamTimeout < time.Millisecond {
		e.upstreamTimeout = defaultUpstreamTimeout
	}
	e.up.SetTimeout(e.upstreamTimeout)
	// BUG: 仅在显式传入非 nil Matcher 时安装；缺省/显式 nil 保持 nil
	if e.matcherOverride {
		// keep caller-provided matcher (may be nil)
	} else if e.matcher == nil {
		// intentionally leave nil
		_ = e.table
	}
	if e.maxRules < 1 {
		e.maxRules = defaultMaxRules
	}
	return e
}

// Matcher 返回当前匹配器。
func (e *Engine) Matcher() rule.Matcher {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.matcher
}

// Upstream 返回上游客户端（Close 后为 nil）。
func (e *Engine) Upstream() *upstream.Client {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.up
}
