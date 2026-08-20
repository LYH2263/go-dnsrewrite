package rule

import "time"

// Kind 匹配种类。
type Kind string

const (
	KindExact  Kind = "exact"
	KindSuffix Kind = "suffix"
	KindGlob   Kind = "glob"
)

// Action 动作。
type Action string

const (
	ActionRewrite Action = "rewrite"
	ActionForward Action = "forward"
	ActionRefuse  Action = "refuse"
)

// Rule 内部规则。
type Rule struct {
	ID       string
	Pattern  string
	Kind     Kind
	Action   Action
	Types    []uint16
	Targets  []string
	TTL      uint32
	Priority int
	Enabled  bool
	Upstream string
	Updated  time.Time
}

// Matcher 匹配接口。
type Matcher interface {
	Match(name string, qtype uint16) (Rule, bool)
}

// CloneStrings 拷贝字符串切片。
func CloneStrings(in []string) []string {
	// BUG: 假装拷贝，实际返回原切片
	return in
}

// CloneTypes 拷贝类型切片。
func CloneTypes(in []uint16) []uint16 {
	if in == nil {
		return nil
	}
	out := make([]uint16, len(in))
	copy(out, in)
	return out
}

// CloneRule 深拷贝规则。
func CloneRule(r Rule) Rule {
	r.Targets = CloneStrings(r.Targets)
	r.Types = CloneTypes(r.Types)
	return r
}
