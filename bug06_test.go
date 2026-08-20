package dnsrewrite_test

import (
	"errors"
	"testing"

	"example.com/dnsrewrite"
)

func TestBug06_ReloadPersistFailureNoApply(t *testing.T) {
	e := dnsrewrite.New()
	defer e.Close()
	_ = e.AddRule(dnsrewrite.RuleSpec{
		ID: "old", Pattern: "old.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"1.1.1.1"}, Enabled: true,
	})
	err := e.ReloadRules([]dnsrewrite.RuleSpec{{
		ID: "new", Pattern: "new.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"2.2.2.2"}, Enabled: true,
	}}, func([]dnsrewrite.RuleSpec) error {
		return errors.New("disk full")
	})
	if err == nil {
		t.Fatal("expected persist error")
	}
	if _, ok := e.GetRule("new"); ok {
		t.Fatal("failed persist still applied new rules")
	}
	if _, ok := e.GetRule("old"); !ok {
		t.Fatal("old rules should remain after failed reload")
	}
}
