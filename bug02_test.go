package dnsrewrite_test

import (
	"testing"

	"example.com/dnsrewrite"
)

func TestBug02_ListRulesTargetsShared(t *testing.T) {
	e := dnsrewrite.New()
	defer e.Close()
	if err := e.AddRule(dnsrewrite.RuleSpec{
		ID: "r1", Pattern: "x.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"1.1.1.1"},
		Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	list := e.ListRules()
	if len(list) != 1 {
		t.Fatal(list)
	}
	list[0].Targets[0] = "9.9.9.9"
	again, ok := e.GetRule("r1")
	if !ok {
		t.Fatal("missing")
	}
	if again.Targets[0] != "1.1.1.1" {
		t.Fatalf("Targets slice shared: %v", again.Targets)
	}
}
