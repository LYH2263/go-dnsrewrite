package rewrite

import "example.com/dnsrewrite/internal/message"

// Refuse 构造拒绝/否定应答。
func Refuse(q Q, rcode uint16) *A {
	if q.Class == 0 {
		q.Class = message.ClassIN
	}
	return &A{RCode: rcode, Question: q, Answers: nil}
}
