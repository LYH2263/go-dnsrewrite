package rule

// HigherPriority 返回 true 表示 a 优先于 b。
// 优先级数值更大者优先；相同时 ID 字典序更小者优先。
func HigherPriority(a, b Rule) bool {
	if a.Priority != b.Priority {
		return a.Priority > b.Priority
	}
	return a.ID < b.ID
}

// SortByPriority 原地按优先级排序（高到低）。
func SortByPriority(rules []Rule) {
	for i := 0; i < len(rules); i++ {
		for j := i + 1; j < len(rules); j++ {
			if HigherPriority(rules[j], rules[i]) {
				rules[i], rules[j] = rules[j], rules[i]
			}
		}
	}
}
