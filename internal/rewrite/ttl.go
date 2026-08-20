package rewrite

import "example.com/dnsrewrite/internal/message"

// OverrideTTL 覆盖应答中所有 RR 的 TTL。
func OverrideTTL(ans *A, ttl uint32) *A {
	if ans == nil || ttl == 0 {
		return ans
	}
	cp := *ans
	cp.Answers = message.CloneRRs(ans.Answers)
	for i := range cp.Answers {
		cp.Answers[i].TTL = ttl
	}
	cp.TTLHint = ttl
	return &cp
}
