package rewrite

import (
	"example.com/dnsrewrite/internal/message"
	"example.com/dnsrewrite/internal/rule"
)

// Q 查询视图。
type Q struct {
	Name  string
	Type  uint16
	Class uint16
}

// A 应答视图。
type A struct {
	ID       uint16
	RCode    uint16
	Question Q
	Answers  []message.RR
	TTLHint  uint32
	Source   string
}

// Apply 按规则装配改写应答。
func Apply(q Q, r rule.Rule) (*A, error) {
	if q.Class == 0 {
		q.Class = message.ClassIN
	}
	rrs, rcode, err := BuildFromRule(q.Name, q.Type, q.Class, r)
	if err != nil {
		return nil, err
	}
	ttl := r.TTL
	if ttl == 0 {
		ttl = 300
	}
	return &A{
		RCode:    rcode,
		Question: q,
		Answers:  rrs,
		TTLHint:  ttl,
	}, nil
}

// BuildFromRule 根据 Targets 生成 RR。
func BuildFromRule(name string, qtype, class uint16, r rule.Rule) ([]message.RR, uint16, error) {
	ttl := r.TTL
	if ttl == 0 {
		ttl = 300
	}
	var rrs []message.RR
	for _, t := range r.Targets {
		switch {
		case qtype == message.TypeAAAA:
			if data, err := message.BuildAAAA(t); err == nil {
				rrs = append(rrs, message.RR{Name: name, Type: message.TypeAAAA, Class: class, TTL: ttl, Data: data})
			} else {
				rrs = append(rrs, message.RR{Name: name, Type: message.TypeCNAME, Class: class, TTL: ttl, Data: message.BuildCNAME(t)})
			}
		case qtype == message.TypeCNAME || looksHost(t):
			rrs = append(rrs, message.RR{Name: name, Type: message.TypeCNAME, Class: class, TTL: ttl, Data: message.BuildCNAME(t)})
		case qtype == message.TypeTXT:
			rrs = append(rrs, message.RR{Name: name, Type: message.TypeTXT, Class: class, TTL: ttl, Data: message.BuildTXT(t)})
		default:
			if data, err := message.BuildA(t); err == nil {
				rrs = append(rrs, message.RR{Name: name, Type: message.TypeA, Class: class, TTL: ttl, Data: data})
			} else {
				rrs = append(rrs, message.RR{Name: name, Type: message.TypeCNAME, Class: class, TTL: ttl, Data: message.BuildCNAME(t)})
			}
		}
	}
	return rrs, message.RCodeNoError, nil
}

func looksHost(s string) bool {
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			if !isIPLike(s) {
				return true
			}
		}
	}
	return false
}

func isIPLike(s string) bool {
	dots, colons := 0, 0
	for _, c := range s {
		switch {
		case c == '.':
			dots++
		case c == ':':
			colons++
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return dots == 3 || colons > 0
}
