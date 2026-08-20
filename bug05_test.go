package dnsrewrite_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/dnsrewrite"
)

func TestBug05_UpstreamErrorWrapsSentinel(t *testing.T) {
	e := dnsrewrite.New()
	defer e.Close()
	if err := e.AddRule(dnsrewrite.RuleSpec{
		ID: "fwd", Pattern: "up.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionForward, Upstream: "127.0.0.1:1", Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := e.TryResolve(ctx, "up.local", dnsrewrite.TypeA)
	if err == nil {
		t.Fatal("expected upstream error")
	}
	if !errors.Is(err, dnsrewrite.ErrUpstream) {
		t.Fatalf("TryResolve must errors.Is ErrUpstream, got %v", err)
	}
}
