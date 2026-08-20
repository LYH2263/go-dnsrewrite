package dnsrewrite_test

import (
	"testing"

	"example.com/dnsrewrite"
)

func TestBug04_NilMatcherNoPanic(t *testing.T) {
	e := dnsrewrite.New(dnsrewrite.WithMatcher(nil))
	defer e.Close()
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("nil matcher panicked: %v", rec)
		}
	}()
	ans, err := e.Resolve(dnsrewrite.Question{Name: "none.local", Type: dnsrewrite.TypeA})
	if err != nil {
		// 允许返回错误，但不得 panic
		return
	}
	if ans == nil {
		t.Fatal("nil answer")
	}
}
