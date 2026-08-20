package dnsrewrite

// Stats 返回运行统计快照。
func (e *Engine) Stats() Stats {
	if e == nil {
		return Stats{Closed: true}
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	n := 0
	if e.table != nil {
		n = e.table.Len()
	}
	return Stats{
		Rules:       n,
		Resolves:    e.resolves,
		Rewrites:    e.rewrites,
		Forwards:    e.forwards,
		Refuses:     e.refuses,
		CacheHits:   e.cacheHits,
		UpstreamErr: e.upstreamErr,
		Closed:      e.closed,
	}
}
