package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RuleJSON 持久化规则。
type RuleJSON struct {
	ID       string   `json:"id"`
	Pattern  string   `json:"pattern"`
	Kind     string   `json:"kind"`
	Action   string   `json:"action"`
	Targets  []string `json:"targets,omitempty"`
	TTL      uint32   `json:"ttl,omitempty"`
	Priority int      `json:"priority,omitempty"`
	Enabled  bool     `json:"enabled"`
	Upstream string   `json:"upstream,omitempty"`
}

// Snapshot 快照。
type Snapshot struct {
	Version int        `json:"version"`
	Rules   []RuleJSON `json:"rules"`
}

// Save 原子写盘（临时文件 + rename）。
func Save(path string, snap Snapshot) error {
	if path == "" {
		return errPath
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && !os.IsExist(err) {
		// Dir may be "." 
		if filepath.Dir(path) != "." {
			return err
		}
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Load 读盘。
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

type persistError string

func (e persistError) Error() string { return string(e) }

const errPath persistError = "persist: empty path"
