package cache

// RR 缓存中的资源记录。
type RR struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []byte
}

// Entry 缓存条目。
type Entry struct {
	ID      uint16
	RCode   uint16
	Name    string
	QType   uint16
	QClass  uint16
	Answers []RR
	TTLHint uint32
	Source  string
}

// CloneEntry 深拷贝（含 RDATA）。
func CloneEntry(e *Entry) *Entry {
	if e == nil {
		return nil
	}
	out := *e
	if e.Answers != nil {
		out.Answers = make([]RR, len(e.Answers))
		for i := range e.Answers {
			// BUG: RDATA 与库存共享
			out.Answers[i] = e.Answers[i]
		}
	}
	return &out
}
