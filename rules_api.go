package dnsrewrite

import (
	"example.com/dnsrewrite/internal/rule"
)

// AddRule 添加或替换规则。
func (e *Engine) AddRule(spec RuleSpec) error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	if err := ValidateRule(spec); err != nil {
		return err
	}
	if e.table.Len() >= e.maxRules && !e.table.Has(spec.ID) {
		return ErrInvalidRule
	}
	if spec.UpdatedAt.IsZero() {
		spec.UpdatedAt = e.clk.Now()
	}
	r := toInternal(spec)
	if r.Kind == "" {
		r.Kind = rule.KindExact
	}
	e.table.Upsert(r)
	e.dirty = true
	return nil
}

// RemoveRule 删除规则。
func (e *Engine) RemoveRule(id string) error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	if !e.table.Remove(id) {
		return ErrNotFound
	}
	e.dirty = true
	return nil
}

// GetRule 返回规则副本。
func (e *Engine) GetRule(id string) (RuleSpec, bool) {
	if e == nil {
		return RuleSpec{}, false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	r, ok := e.table.Get(id)
	if !ok {
		return RuleSpec{}, false
	}
	return fromInternal(r), true
}

// ListRules 返回全部规则副本（含 Targets 深拷贝）。
func (e *Engine) ListRules() []RuleSpec {
	if e == nil {
		return nil
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	list := e.table.List()
	out := make([]RuleSpec, 0, len(list))
	for _, r := range list {
		out = append(out, fromInternal(r))
	}
	return out
}

// EnableRule 启用/禁用规则。
func (e *Engine) EnableRule(id string, enabled bool) error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	r, ok := e.table.Get(id)
	if !ok {
		return ErrNotFound
	}
	r.Enabled = enabled
	r.Updated = e.clk.Now()
	e.table.Upsert(r)
	e.dirty = true
	return nil
}
