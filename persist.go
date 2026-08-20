package dnsrewrite

import (
	"example.com/dnsrewrite/internal/persist"
	"example.com/dnsrewrite/internal/rule"
)

// PersistFunc 自定义持久化钩子；失败时 Reload 不得生效新表。
type PersistFunc func([]RuleSpec) error

// ReloadRules 热加载规则；持久化失败不替换内存表。
func (e *Engine) ReloadRules(specs []RuleSpec, persistFn PersistFunc) error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	candidates := make([]rule.Rule, 0, len(specs))
	views := make([]RuleSpec, 0, len(specs))
	for _, s := range specs {
		if err := ValidateRule(s); err != nil {
			return err
		}
		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = e.clk.Now()
		}
		candidates = append(candidates, toInternal(s))
		views = append(views, fromInternal(toInternal(s)))
	}
	// 先持久化候选规则，成功后才替换内存表；持久化失败时旧规则保持不变。
	if persistFn != nil {
		if err := persistFn(views); err != nil {
			return wrapPersist(err)
		}
	} else if e.persistPath != "" {
		if err := e.persistRulesLocked(candidates); err != nil {
			return err
		}
	}
	e.table.Replace(candidates)
	e.dirty = persistFn == nil && e.persistPath == ""
	return nil
}

// LoadPersist 从磁盘加载。
func (e *Engine) LoadPersist() error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.persistPath == "" {
		return nil
	}
	snap, err := persist.Load(e.persistPath)
	if err != nil {
		return wrapPersist(err)
	}
	rules := make([]rule.Rule, 0, len(snap.Rules))
	for _, s := range snap.Rules {
		rules = append(rules, toInternal(RuleSpec{
			ID: s.ID, Pattern: s.Pattern, Kind: MatchKind(s.Kind),
			Action: Action(s.Action), Targets: s.Targets, TTL: s.TTL,
			Priority: s.Priority, Enabled: s.Enabled, Upstream: s.Upstream,
		}))
	}
	e.table.Replace(rules)
	e.dirty = false
	return nil
}

// Sync 将脏规则刷盘。
func (e *Engine) Sync() error {
	if e == nil {
		return ErrNilEngine
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	return e.flushLocked()
}

func (e *Engine) flushLocked() error {
	if e.persistPath == "" {
		return nil
	}
	if err := e.persistRulesLocked(e.table.List()); err != nil {
		return err
	}
	e.dirty = false
	return nil
}

// persistRulesLocked 将给定规则快照刷盘；不改动内存表与脏标记。调用方持锁。
func (e *Engine) persistRulesLocked(rules []rule.Rule) error {
	snap := persist.Snapshot{Version: 1}
	for _, r := range rules {
		snap.Rules = append(snap.Rules, persist.RuleJSON{
			ID: r.ID, Pattern: r.Pattern, Kind: string(r.Kind),
			Action: string(r.Action), Targets: append([]string(nil), r.Targets...),
			TTL: r.TTL, Priority: r.Priority, Enabled: r.Enabled, Upstream: r.Upstream,
		})
	}
	if err := persist.Save(e.persistPath, snap); err != nil {
		return wrapPersist(err)
	}
	return nil
}
