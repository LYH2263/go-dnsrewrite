package message

// CloneRR 深拷贝单条 RR。
func CloneRR(rr RR) RR {
	out := rr
	if rr.Data != nil {
		out.Data = append([]byte(nil), rr.Data...)
	}
	return out
}

// CloneRRs 深拷贝 RR 切片。
func CloneRRs(in []RR) []RR {
	if in == nil {
		return nil
	}
	out := make([]RR, len(in))
	for i := range in {
		out[i] = CloneRR(in[i])
	}
	return out
}

// CloneMessage 深拷贝报文。
func CloneMessage(m *Message) *Message {
	if m == nil {
		return nil
	}
	out := *m
	out.Answers = CloneRRs(m.Answers)
	return &out
}
