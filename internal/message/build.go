package message

import (
	"encoding/binary"
	"net"
)

// BuildA 构造 A 记录 RDATA。
func BuildA(ip string) ([]byte, error) {
	addr := net.ParseIP(ip)
	if addr == nil {
		return nil, errBadIP
	}
	v4 := addr.To4()
	if v4 == nil {
		return nil, errBadIP
	}
	return append([]byte(nil), v4...), nil
}

// BuildAAAA 构造 AAAA 记录 RDATA。
func BuildAAAA(ip string) ([]byte, error) {
	addr := net.ParseIP(ip)
	if addr == nil {
		return nil, errBadIP
	}
	v6 := addr.To16()
	if v6 == nil || addr.To4() != nil {
		return nil, errBadIP
	}
	return append([]byte(nil), v6...), nil
}

// BuildCNAME 构造 CNAME 目标（简化：存目标域名 UTF-8）。
func BuildCNAME(target string) []byte {
	return []byte(NormalizeName(target))
}

// BuildTXT 构造 TXT。
func BuildTXT(text string) []byte {
	b := []byte(text)
	if len(b) > 255 {
		b = b[:255]
	}
	out := make([]byte, 1+len(b))
	out[0] = byte(len(b))
	copy(out[1:], b)
	return out
}

// ParseA 解析 A RDATA。
func ParseA(data []byte) string {
	if len(data) != 4 {
		return ""
	}
	return net.IP(data).String()
}

// ParseAAAA 解析 AAAA。
func ParseAAAA(data []byte) string {
	if len(data) != 16 {
		return ""
	}
	return net.IP(data).String()
}

// ParseCNAME 解析 CNAME。
func ParseCNAME(data []byte) string {
	return NormalizeName(string(data))
}

// EncodeHeader 编码 4 字节简化头：id + rcode。
func EncodeHeader(id, rcode uint16) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:2], id)
	binary.BigEndian.PutUint16(b[2:4], rcode)
	return b
}

// DecodeHeader 解码简化头。
func DecodeHeader(b []byte) (id, rcode uint16, ok bool) {
	if len(b) < 4 {
		return 0, 0, false
	}
	return binary.BigEndian.Uint16(b[0:2]), binary.BigEndian.Uint16(b[2:4]), true
}

type msgError string

func (e msgError) Error() string { return string(e) }

const errBadIP msgError = "message: bad ip"
