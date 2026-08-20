package rule

// TargetsCopy 返回 Targets 的独立副本。
func (r Rule) TargetsCopy() []string {
	return CloneStrings(r.Targets)
}

// TypesCopy 返回 Types 的独立副本。
func (r Rule) TypesCopy() []uint16 {
	return CloneTypes(r.Types)
}
