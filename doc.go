// Package dnsrewrite 提供域名规则匹配与 DNS 应答改写引擎。
//
// 能力概览：
//   - 精确 / 后缀 / 通配规则匹配，按优先级选取
//   - Resolve / ResolveContext：本地改写或转发上游
//   - 应答装配（A/AAAA/CNAME/拒绝）与 TTL 覆盖
//   - 规则热加载；持久化失败不应用
//   - Close 后拒绝 Resolve
//
// DNS 模型为简化 Question/Answer 结构体，不追求完整 RFC 报文编解码，
// 但匹配、改写、转发与缓存语义真实可测。
package dnsrewrite
