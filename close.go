package dnsrewrite

// Close 关闭引擎：先 Sync 再释放规则表与上游客户端。
func (e *Engine) Close() error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil
	}
	var first error
	// 先把当前规则 Sync 落盘，再释放内存表，避免重启等于裸机。
	if e.dirty || e.persistPath != "" {
		if err := e.flushLocked(); err != nil {
			first = err
		}
	}
	if e.table != nil {
		e.table.Replace(nil)
	}
	if e.up != nil {
		e.up.Close()
		e.up = nil
	}
	if e.cache != nil {
		e.cache.Clear()
	}
	e.matcher = nil
	e.closed = true
	return first
}

// Closed 是否已关闭。
func (e *Engine) Closed() bool {
	if e == nil {
		return true
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.closed
}
