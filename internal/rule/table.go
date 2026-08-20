package rule

import "sync"

// Table 规则表，实现 Matcher。
type Table struct {
	mu    sync.RWMutex
	byID  map[string]Rule
	order []string
}

// NewTable 空表。
func NewTable() *Table {
	return &Table{byID: make(map[string]Rule)}
}

// Upsert 插入或更新。
func (t *Table) Upsert(r Rule) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.byID == nil {
		t.byID = make(map[string]Rule)
	}
	r = CloneRule(r)
	if _, ok := t.byID[r.ID]; !ok {
		t.order = append(t.order, r.ID)
	}
	t.byID[r.ID] = r
}

// Remove 删除。
func (t *Table) Remove(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.byID[id]; !ok {
		return false
	}
	delete(t.byID, id)
	out := t.order[:0]
	for _, x := range t.order {
		if x != id {
			out = append(out, x)
		}
	}
	t.order = out
	return true
}

// Has 是否存在。
func (t *Table) Has(id string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.byID[id]
	return ok
}

// Get 取副本。
func (t *Table) Get(id string) (Rule, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.byID[id]
	if !ok {
		return Rule{}, false
	}
	return CloneRule(r), true
}

// List 全部副本（按插入序）。
func (t *Table) List() []Rule {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Rule, 0, len(t.order))
	for _, id := range t.order {
		if r, ok := t.byID[id]; ok {
			out = append(out, CloneRule(r))
		}
	}
	return out
}

// Len 规则数。
func (t *Table) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.byID)
}

// Replace 整体替换。
func (t *Table) Replace(rules []Rule) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.byID = make(map[string]Rule, len(rules))
	t.order = make([]string, 0, len(rules))
	for _, r := range rules {
		r = CloneRule(r)
		t.byID[r.ID] = r
		t.order = append(t.order, r.ID)
	}
}

// Match 选取最高优先级启用规则。
func (t *Table) Match(name string, qtype uint16) (Rule, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var best Rule
	found := false
	for _, id := range t.order {
		r, ok := t.byID[id]
		if !ok || !r.Enabled {
			continue
		}
		if !MatchName(r.Kind, r.Pattern, name) {
			continue
		}
		if !TypeAllowed(r.Types, qtype) {
			continue
		}
		if !found || HigherPriority(r, best) {
			best = CloneRule(r)
			found = true
		}
	}
	return best, found
}
