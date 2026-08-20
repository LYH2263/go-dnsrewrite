package upstream

import "sync"

// Pool 简单地址池。
type Pool struct {
	mu    sync.Mutex
	addrs []string
	idx   int
}

// NewPool 构造。
func NewPool(addrs ...string) *Pool {
	p := &Pool{}
	p.addrs = append(p.addrs, addrs...)
	return p
}

// Add 添加地址。
func (p *Pool) Add(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.addrs = append(p.addrs, addr)
}

// Next 轮询取地址。
func (p *Pool) Next() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.addrs) == 0 {
		return "", false
	}
	a := p.addrs[p.idx%len(p.addrs)]
	p.idx++
	return a, true
}

// Len 数量。
func (p *Pool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.addrs)
}
