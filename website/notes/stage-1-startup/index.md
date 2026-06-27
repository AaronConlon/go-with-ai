# 阶段 1：启动期

目标：能从零创建、运行和测试一个 Go module，并建立最小项目结构。

## 学习目标

- 理解 `go mod init`、`go mod tidy`、`go test ./...`。
- 理解 package、module path、`internal/`。
- 掌握基础类型、函数、slice、map、struct。
- 能写一个最小 CLI。
- 能写一个最小 table-driven test。

## 完成任务

- [ ] 确认本机 Go 版本。
- [ ] 创建项目 module。
- [ ] 建立 `cmd/hnctl` 和 `internal/hn`。
- [ ] 写出第一个 CLI 输出。
- [ ] 写出第一个单元测试。

## 知识记录

### Go module

Go module 是依赖管理和版本边界。它由 `go.mod` 描述，通常对应一个仓库或一个独立发布单元。

### 目录结构

```text
cmd/hnctl/
internal/hn/
```

`cmd/` 放可执行入口，`internal/` 放只允许当前 module 内部使用的包。

## 问题清单

- module path 暂时用什么命名？
- CLI 是否先只服务学习，还是后续保留为正式工具？

## 验收标准

```bash
go test ./...
go run ./cmd/hnctl
```

## 本阶段记录

- 后续记录追加在本页，或按主题拆出子文档。

