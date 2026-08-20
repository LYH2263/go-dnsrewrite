package dnsrewrite_test

import (
	"errors"
	"os"
	"path/filepath"
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

// TestBug06_ReloadDiskFailureNoApply 覆盖真实磁盘写入失败：persistPath 指向一个
// 父目录是普通文件的路径，使 persist.Save 的 MkdirAll 失败。Reload 不得生效新表。
func TestBug06_ReloadDiskFailureNoApply(t *testing.T) {
	dir := t.TempDir()
	// blocker 是普通文件，blocker/rules.json 的父目录无法创建 → Save 失败。
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	unwritable := filepath.Join(blocker, "rules.json")

	e := dnsrewrite.New(dnsrewrite.WithPersistPath(unwritable))
	defer e.Close()
	if err := e.AddRule(dnsrewrite.RuleSpec{
		ID: "old", Pattern: "old.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"1.1.1.1"}, Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	// 旧规则仅在内存（该路径不可写，无法 Sync）；Reload 新规则时磁盘写入失败。
	err := e.ReloadRules([]dnsrewrite.RuleSpec{{
		ID: "new", Pattern: "new.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionRewrite, Targets: []string{"2.2.2.2"}, Enabled: true,
	}}, nil)
	if err == nil {
		t.Fatal("expected persist error from disk write failure")
	}
	if _, ok := e.GetRule("new"); ok {
		t.Fatal("failed disk write still applied new rules")
	}
	if _, ok := e.GetRule("old"); !ok {
		t.Fatal("old rules should remain after failed reload")
	}
}
