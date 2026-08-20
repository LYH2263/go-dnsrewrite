package rewrite

import "example.com/dnsrewrite/internal/message"

// ChainCNAME 追加 CNAME 链记录。
func ChainCNAME(ans *A, name, target string, ttl uint32) *A {
	if ans == nil {
		ans = &A{}
	}
	if ttl == 0 {
		ttl = 300
	}
	cp := *ans
	cp.Answers = append(message.CloneRRs(ans.Answers), message.RR{
		Name:  message.NormalizeName(name),
		Type:  message.TypeCNAME,
		Class: message.ClassIN,
		TTL:   ttl,
		Data:  message.BuildCNAME(target),
	})
	return &cp
}
