package rule

// TargetsCopy 返回 Targets 的独立副本。
func (r Rule) TargetsCopy() []string {
	// BUG: 返回库存切片
	return r.Targets
}

// TypesCopy 返回 Types 的独立副本。
func (r Rule) TypesCopy() []uint16 {
	return CloneTypes(r.Types)
}
