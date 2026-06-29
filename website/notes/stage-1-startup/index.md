# 阶段 1：启动期

目标：能从零创建、运行和测试一个 Go module，并建立最小项目结构。

## 项目目录约定

仓库根目录是学习工作区，不直接作为 Go module。Hacker News 实战项目放在独立目录：

```text
hn-agent/
```

阶段一涉及 Go 代码的命令，默认都在 `hn-agent/` 目录内执行。

## 学习目标

- 理解 `go mod init`、`go mod tidy`、`go test ./...`。
- 理解 package、module path、`internal/`。
- 掌握基础类型、函数、slice、map、struct。
- 能写一个最小 CLI。
- 能写一个最小 table-driven test。

## 完成任务

- [x] 确认本机 Go 版本。
- [x] 创建 `hn-agent/` 项目目录。
- [x] 在 `hn-agent/` 内创建项目 module。
- [x] 建立 `hn-agent/cmd/hnctl` 和 `hn-agent/internal/hn`。
- [x] 写出第一个 CLI 输出。
- [x] 写出第一个单元测试。

## 知识记录

- [阶段一教练笔记](/notes/stage-1-startup/coach-notes)
- [Go 测试代码语法拆解](/notes/stage-1-startup/go-test-syntax)
- [Go 变量定义语法](/notes/stage-1-startup/go-variable-syntax)
- [Go 切片与数组](/notes/stage-1-startup/go-slice-vs-array)
- [Go if 语法](/notes/stage-1-startup/go-if-syntax)

### Go module

Go module 是依赖管理和版本边界。它由 `go.mod` 描述，通常对应一个仓库或一个独立发布单元。

### 目录结构

```text
hn-agent/
  cmd/hnctl/
  internal/hn/
```

`cmd/` 放可执行入口，`internal/` 放只允许当前 module 内部使用的包。

## 问题清单

- module path 暂时用什么命名？
- CLI 是否先只服务学习，还是后续保留为正式工具？

## 验收标准

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
```

## 本阶段记录

- [阶段一教练笔记](/notes/stage-1-startup/coach-notes)：项目骨架、module path、`internal/`、table-driven test、最小 CLI 和验收命令。
- [Go 测试代码语法拆解](/notes/stage-1-startup/go-test-syntax)：逐段解释第一段 Go 测试代码的语法。
- [Go 变量定义语法](/notes/stage-1-startup/go-variable-syntax)：解释 `var`、`const`、`:=` 和 `=` 的区别。

## 完成状态

阶段一已完成。当前已具备：

- 独立 Go module：`hn-agent/go.mod`。
- 最小领域模型：`hn-agent/internal/hn/story.go`。
- 第一个单元测试：`hn-agent/internal/hn/story_test.go`。
- 最小 CLI：`hn-agent/cmd/hnctl/main.go`。

下一阶段进入 [阶段 2：网络期](/notes/stage-2-network/)。
