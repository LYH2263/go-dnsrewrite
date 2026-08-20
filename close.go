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
	if e.dirty || e.persistPath != "" {
		if err := e.flushLocked(); err != nil {
			first = err
		}
	}
	// Sync 完成后再丢表与上游，避免丢配置
	if e.table != nil {
		e.table.Replace(nil)
	}
	if e.up != nil {
		e.up.Close()
		e.up = nil // 置 nil：关停后 Resolve 不再触碰转发客户端
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
