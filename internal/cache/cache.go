package cache

import "sync"

// Stats 缓存统计。
type Stats struct {
	Entries int
	Hits    uint64
	Misses  uint64
	Puts    uint64
}

// Cache 短缓存。
type Cache struct {
	mu      sync.Mutex
	data    map[string]*Entry
	limit   int
	hits    uint64
	misses  uint64
	puts    uint64
}

// New 构造。
func New(limit int) *Cache {
	if limit < 1 {
		limit = 128
	}
	return &Cache{data: make(map[string]*Entry), limit: limit}
}

// Put 写入深拷贝。
func (c *Cache) Put(key string, e *Entry) *Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil {
		c.data = make(map[string]*Entry)
	}
	cp := CloneEntry(e)
	if len(c.data) >= c.limit {
		// 简单淘汰：删一个任意键
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}
	c.data[key] = cp
	c.puts++
	return CloneEntry(cp)
}

// Get 读取深拷贝。
func (c *Cache) Get(key string) (*Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.data[key]
	if !ok {
		c.misses++
		return nil, false
	}
	c.hits++
	return CloneEntry(e), true
}

// Clear 清空。
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*Entry)
}

// Stats 快照。
func (c *Cache) Stats() Stats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Stats{Entries: len(c.data), Hits: c.hits, Misses: c.misses, Puts: c.puts}
}
