package message

import (
	"fmt"
	"strings"
)

// NormalizeName 规范化域名（小写、去尾点）。
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.TrimSuffix(name, ".")
	return name
}

// CacheKey 构造缓存键。
func CacheKey(name string, qtype, class uint16) string {
	return fmt.Sprintf("%s|%d|%d", NormalizeName(name), qtype, class)
}

// SplitLabels 拆标签。
func SplitLabels(name string) []string {
	name = NormalizeName(name)
	if name == "" {
		return nil
	}
	return strings.Split(name, ".")
}

// JoinLabels 合标签。
func JoinLabels(labels []string) string {
	return strings.Join(labels, ".")
}

// IsSubdomain 判断 name 是否为 parent 的子域或相等。
func IsSubdomain(name, parent string) bool {
	name = NormalizeName(name)
	parent = NormalizeName(parent)
	if name == parent {
		return true
	}
	return strings.HasSuffix(name, "."+parent)
}
