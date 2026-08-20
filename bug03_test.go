package dnsrewrite_test

import (
	"errors"
	"testing"

	"example.com/dnsrewrite"
)

func TestBug03_ResolveAfterCloseNoPanic(t *testing.T) {
	e := dnsrewrite.New()
	_ = e.AddRule(dnsrewrite.RuleSpec{
		ID: "r1", Pattern: "z.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"1.2.3.4"}, Enabled: true,
	})
	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Resolve after Close panicked: %v", rec)
		}
	}()
	_, err := e.Resolve(dnsrewrite.Question{Name: "z.local", Type: dnsrewrite.TypeA})
	if err == nil {
		t.Fatal("expected closed error")
	}
	if !errors.Is(err, dnsrewrite.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
