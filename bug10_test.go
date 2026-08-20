package dnsrewrite_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"example.com/dnsrewrite"
)

func TestBug10_CloseSyncsBeforeDropRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	e := dnsrewrite.New(dnsrewrite.WithPersistPath(path))
	if err := e.AddRule(dnsrewrite.RuleSpec{
		ID: "keep", Pattern: "keep.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"8.8.8.8"}, Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var snap struct {
		Rules []struct {
			ID string `json:"id"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatal(err)
	}
	if len(snap.Rules) != 1 || snap.Rules[0].ID != "keep" {
		t.Fatalf("persist lost rules on Close: %s", data)
	}
}
