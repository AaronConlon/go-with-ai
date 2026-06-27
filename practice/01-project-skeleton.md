# 01 项目骨架

## 本步目标

创建 Go module、最小 CLI 和第一条测试，为后续 HN client 做准备。

## 项目目录约定

仓库根目录是学习工作区，不直接作为 Go module。真正的 Hacker News 实战项目放在独立目录：

```text
hn-agent/
```

因此阶段一所有 Go 命令默认在 `hn-agent/` 目录内执行。

## 教练笔记

本步骤的完整解释、推荐代码、验收方式和常见判断已沉淀到：

```text
website/notes/stage-1-startup/coach-notes.md
```

站点中可阅读：

[阶段一教练笔记](/notes/stage-1-startup/coach-notes)

## 计划目录

```text
hn-agent/
  cmd/hnctl/
  cmd/hn-agent/
  internal/hn/
```

## 操作步骤

1. 确认 Go 版本。
2. 创建 `hn-agent/` 项目目录。
3. 在 `hn-agent/` 内初始化 module。
4. 创建 `cmd/hnctl`。
5. 创建 `internal/hn`。
6. 写第一条 table-driven test。

## 学习者执行原则

- 学习者自己创建文件和写代码。
- 教练不直接改阶段代码，只提供步骤、解释、样例和 review。
- 每次完成后，把关键命令输出贴出来做检查。

## 预计命令

```bash
go version
mkdir -p hn-agent
cd hn-agent
go mod init github.com/aaron/go-with-ai/hn-agent
go test ./...
go run ./cmd/hnctl
```

## 验收标准

- `go test ./...` 通过。
- `go run ./cmd/hnctl` 能输出版本或帮助信息。
- 相关学习记录回填到 `website/notes/stage-1-startup/`。

## 完成后贴给教练

```bash
cat go.mod
go test ./...
go run ./cmd/hnctl
```
