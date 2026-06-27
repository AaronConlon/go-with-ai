# 01 项目骨架

## 本步目标

创建 Go module、最小 CLI 和第一条测试，为后续 HN client 做准备。

## 计划目录

```text
cmd/hnctl/
cmd/hn-agent/
internal/hn/
```

## 操作步骤

1. 确认 Go 版本。
2. 初始化 module。
3. 创建 `cmd/hnctl`。
4. 创建 `internal/hn`。
5. 写第一条 table-driven test。

## 预计命令

```bash
go version
go mod init github.com/aaron/go-with-ai
go test ./...
go run ./cmd/hnctl
```

## 验收标准

- `go test ./...` 通过。
- `go run ./cmd/hnctl` 能输出版本或帮助信息。
- 相关学习记录回填到 `website/notes/stage-1-startup/`。

