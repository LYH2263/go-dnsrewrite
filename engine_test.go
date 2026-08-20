package dnsrewrite_test

import (
	"testing"

	"example.com/dnsrewrite"
)

func TestRewriteARecord(t *testing.T) {
	e := dnsrewrite.New()
	defer e.Close()
	if err := e.AddRule(dnsrewrite.RuleSpec{
		ID: "a1", Pattern: "app.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"10.0.0.1"},
		TTL: 60, Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	ans, err := e.Resolve(dnsrewrite.Question{Name: "app.local", Type: dnsrewrite.TypeA})
	if err != nil {
		t.Fatal(err)
	}
	if ans.RCode != dnsrewrite.RCodeNoError || len(ans.Answers) != 1 {
		t.Fatalf("bad answer: %+v", ans)
	}
}

func TestRefuseRule(t *testing.T) {
	e := dnsrewrite.New()
	defer e.Close()
	_ = e.AddRule(dnsrewrite.RuleSpec{
		ID: "deny", Pattern: "blocked.test", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRefuse, Enabled: true,
	})
	ans, err := e.Resolve(dnsrewrite.Question{Name: "blocked.test", Type: dnsrewrite.TypeA})
	if err != nil {
		t.Fatal(err)
	}
	if ans.RCode != dnsrewrite.RCodeRefused {
		t.Fatalf("want refused got %v", ans.RCode)
	}
}
