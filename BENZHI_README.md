# DNS改写引擎

按域名模式匹配规则，改写应答（CNAME/A/AAAA/拒绝）或转发上游；管理页维护规则与试查询。


## 环境

- 镜像：`benzhi.Dockerfile` 基于 `golang:1.22`（官方多架构）
- `go.mod` 语言版本：go 1.22
- 容器内使用镜像自带工具链即可

## 标准命令

```bash
go build ./...
go test ./... -count=1
go vet ./...
```

## 构建评测镜像（须双架构）

验证请用 `bash -c`（勿用 `bash -lc`）。

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-dnsrewrite linux/amd64
docker run --platform linux/amd64 --rm go-dnsrewrite:latest bash -c 'go build ./...'

./build_benzhi_docker.sh go-dnsrewrite linux/arm64
docker run --platform linux/arm64 --rm go-dnsrewrite:latest bash -c 'go build ./...'
```

构建阶段已 `go mod download`；容器内编译不应再出现 `downloading ...`。
