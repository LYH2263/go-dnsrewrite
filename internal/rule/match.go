package rule

import (
	"path"
	"strings"
)

// MatchName 按 Kind 匹配域名。
func MatchName(kind Kind, pattern, name string) bool {
	pattern = strings.ToLower(strings.TrimSuffix(pattern, "."))
	name = strings.ToLower(strings.TrimSuffix(name, "."))
	switch kind {
	case KindExact, "":
		return pattern == name
	case KindSuffix:
		if pattern == name {
			return true
		}
		return strings.HasSuffix(name, "."+pattern)
	case KindGlob:
		ok, _ := path.Match(pattern, name)
		return ok
	default:
		return false
	}
}

// TypeAllowed 检查查询类型是否在规则允许列表；空列表表示全部。
func TypeAllowed(types []uint16, qtype uint16) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if t == qtype {
			return true
		}
	}
	return false
}
