package upstream

// Question 上游查询。
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

// RR 上游记录。
type RR struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []byte
}

// Answer 上游应答。
type Answer struct {
	ID      uint16
	RCode   uint16
	Answers []RR
	TTLHint uint32
}
