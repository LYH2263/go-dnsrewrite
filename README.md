# go-dnsrewrite

DNS 改写引擎（库 + `cmd/dnsd` 管理服务）。简化 Question/Answer 模型，支持精确/后缀/通配规则、改写与上游转发、短缓存与规则持久化。

```bash
go test ./... -count=1
go run ./cmd/dnsd -addr :8097 -web web
```
