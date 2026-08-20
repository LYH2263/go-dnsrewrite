package dnsrewrite

import (
	"time"

	"example.com/dnsrewrite/internal/clock"
	"example.com/dnsrewrite/internal/rule"
)

// Option 配置 Engine。
type Option func(*Engine)

// WithClock 注入时钟（测试用）。
func WithClock(c clock.Clock) Option {
	return func(e *Engine) {
		if c != nil {
			e.clk = c
		}
	}
}

// WithPersistPath 设置规则快照路径。
func WithPersistPath(path string) Option {
	return func(e *Engine) {
		e.persistPath = path
	}
}

// WithUpstreamTimeout 上游超时。
func WithUpstreamTimeout(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.upstreamTimeout = d
		}
	}
}

// WithCacheTTL 短缓存默认 TTL。
func WithCacheTTL(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.cacheTTL = d
		}
	}
}

// WithDefaultUpstream 默认上游地址（无规则匹配时）。
func WithDefaultUpstream(addr string) Option {
	return func(e *Engine) {
		e.defaultUpstream = addr
	}
}

// WithMatcher 注入自定义匹配器；传 nil 表示故意不安装（测缺省行为时慎用）。
func WithMatcher(m rule.Matcher) Option {
	return func(e *Engine) {
		e.matcher = m
		e.matcherOverride = true
	}
}

// WithMaxRules 规则上限。
func WithMaxRules(n int) Option {
	return func(e *Engine) {
		if n > 0 {
			e.maxRules = n
		}
	}
}
