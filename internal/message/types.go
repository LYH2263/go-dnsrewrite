package message

// 简化 DNS 常量，与门面 RRType/RCode 对齐。
const (
	TypeA     uint16 = 1
	TypeNS    uint16 = 2
	TypeCNAME uint16 = 5
	TypeTXT   uint16 = 16
	TypeAAAA  uint16 = 28

	RCodeNoError  uint16 = 0
	RCodeFormErr  uint16 = 1
	RCodeServFail uint16 = 2
	RCodeNXDomain uint16 = 3
	RCodeRefused  uint16 = 5

	ClassIN uint16 = 1
)

// Question 简化问题。
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

// RR 简化记录。
type RR struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []byte
}

// Message 简化报文。
type Message struct {
	ID       uint16
	RCode    uint16
	Question Question
	Answers  []RR
}
