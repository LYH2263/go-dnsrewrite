package dnsrewrite

import "time"

// MatchKind 域名匹配方式。
type MatchKind string

const (
	MatchExact  MatchKind = "exact"
	MatchSuffix MatchKind = "suffix"
	MatchGlob   MatchKind = "glob"
)

// Action 改写动作。
type Action string

const (
	ActionRewrite Action = "rewrite"
	ActionForward Action = "forward"
	ActionRefuse  Action = "refuse"
)

// RRType 简化资源记录类型。
type RRType uint16

const (
	TypeA     RRType = 1
	TypeNS    RRType = 2
	TypeCNAME RRType = 5
	TypeAAAA  RRType = 28
	TypeTXT   RRType = 16
)

// RCode 简化应答码。
type RCode uint16

const (
	RCodeNoError  RCode = 0
	RCodeFormErr  RCode = 1
	RCodeServFail RCode = 2
	RCodeNXDomain RCode = 3
	RCodeRefused  RCode = 5
)

// Question 查询问题。
type Question struct {
	Name  string
	Type  RRType
	Class uint16
}

// RR 简化资源记录。
type RR struct {
	Name  string
	Type  RRType
	Class uint16
	TTL   uint32
	Data  []byte
}

// Answer 应答。
type Answer struct {
	ID       uint16
	RCode    RCode
	Question Question
	Answers  []RR
	TTLHint  uint32
	Source   string
}

// RuleSpec 对外规则视图。
type RuleSpec struct {
	ID        string
	Pattern   string
	Kind      MatchKind
	Action    Action
	Types     []RRType
	Targets   []string
	TTL       uint32
	Priority  int
	Enabled   bool
	Upstream  string
	UpdatedAt time.Time
}

// ResolveResult 解析摘要（管理页试查询）。
type ResolveResult struct {
	MatchedRule string
	Action      Action
	RCode       RCode
	Answers     []RR
	DurationMs  int64
	FromCache   bool
	Upstream    string
}

// Stats 引擎统计。
type Stats struct {
	Rules       int
	Resolves    uint64
	Rewrites    uint64
	Forwards    uint64
	Refuses     uint64
	CacheHits   uint64
	UpstreamErr uint64
	Closed      bool
}
