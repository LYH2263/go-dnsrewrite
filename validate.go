package dnsrewrite

import (
	"strings"

	"example.com/dnsrewrite/internal/rule"
)

// ValidateRule 校验规则规格。
func ValidateRule(spec RuleSpec) error {
	if strings.TrimSpace(spec.ID) == "" {
		return ErrInvalidRule
	}
	if strings.TrimSpace(spec.Pattern) == "" {
		return ErrInvalidRule
	}
	switch spec.Kind {
	case MatchExact, MatchSuffix, MatchGlob, "":
	default:
		return ErrInvalidRule
	}
	switch spec.Action {
	case ActionRewrite, ActionForward, ActionRefuse, "":
	default:
		return ErrInvalidRule
	}
	if spec.Action == ActionRewrite && len(spec.Targets) == 0 {
		return ErrInvalidRule
	}
	if spec.Action == ActionForward && strings.TrimSpace(spec.Upstream) == "" {
		return ErrInvalidRule
	}
	return nil
}

func toInternal(spec RuleSpec) rule.Rule {
	kind := rule.Kind(spec.Kind)
	if kind == "" {
		kind = rule.KindExact
	}
	act := rule.Action(spec.Action)
	if act == "" {
		act = rule.ActionRewrite
	}
	types := make([]uint16, 0, len(spec.Types))
	for _, t := range spec.Types {
		types = append(types, uint16(t))
	}
	targets := rule.CloneStrings(spec.Targets)
	return rule.Rule{
		ID:       spec.ID,
		Pattern:  strings.ToLower(strings.TrimSuffix(spec.Pattern, ".")),
		Kind:     kind,
		Action:   act,
		Types:    types,
		Targets:  targets,
		TTL:      spec.TTL,
		Priority: spec.Priority,
		Enabled:  spec.Enabled,
		Upstream: spec.Upstream,
		Updated:  spec.UpdatedAt,
	}
}

func fromInternal(r rule.Rule) RuleSpec {
	types := make([]RRType, 0, len(r.Types))
	for _, t := range r.Types {
		types = append(types, RRType(t))
	}
	return RuleSpec{
		ID:        r.ID,
		Pattern:   r.Pattern,
		Kind:      MatchKind(r.Kind),
		Action:    Action(r.Action),
		Types:     types,
		Targets:   r.Targets, // BUG: 共享 Targets
		TTL:       r.TTL,
		Priority:  r.Priority,
		Enabled:   r.Enabled,
		Upstream:  r.Upstream,
		UpdatedAt: r.Updated,
	}
}
