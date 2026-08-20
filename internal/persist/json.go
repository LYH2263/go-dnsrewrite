package persist

import "encoding/json"

// MarshalRules 编码规则列表。
func MarshalRules(rules []RuleJSON) ([]byte, error) {
	return json.Marshal(Snapshot{Version: 1, Rules: rules})
}

// UnmarshalRules 解码。
func UnmarshalRules(data []byte) ([]RuleJSON, error) {
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return snap.Rules, nil
}
