# 学习问题对话记录

这份文档把学习过程中遇到的问题简化成对话。目标不是完整还原聊天，而是保留最有价值的困惑、解释、结论和下一步动作。

## 记录规则

每条记录使用：

```md
## YYYY-MM-DD 问题标题

**我：** 简化后的真实问题。

**教练：** 简洁解释。

**结论：** 可以复用的一句话判断。

**后续动作：** 下一步该做什么。
```

如果某个问题适合用图解释，就配一张正文图，并保存到：

```text
assets/learning-dialogues-illustrations/
```

站点展示图放在：

```text
website/static/images/learning-dialogues/
```

## 2026-06-27 `go mod init` 是什么意思

**我：** `go mod init github.com/aaron/go-with-ai` 是什么意思？效果是什么？

**教练：** 它会在当前目录初始化一个 Go module，并生成 `go.mod`。`module path` 是这个项目的身份，也会成为项目内部 package 的 import 前缀。

**结论：** `go mod init` 是给 Go 项目立“身份证”，让 Go 工具链知道项目边界、import 路径和依赖记录位置。

**后续动作：** 在真正的 Go 项目目录里执行 `go mod init`，不要随便在学习工作区根目录执行。

## 2026-06-27 根目录要不要作为 Go module

![根目录和 Go module 边界](/images/learning-dialogues/01-workspace-vs-module.png)

**我：** 可是我不需要在根目录做这个，这个学习性质的 Hacker News 项目，应该在根目录有一个独立的目录和名称。

**教练：** 对。根目录现在承担的是学习工作区职责，里面有文档站、学习记录、背景研究和实战计划。真正的 Hacker News Go 项目应该放在独立目录，例如 `hn-agent/`。

**结论：** 根目录是学习工作区，`hn-agent/` 才是 Go module 项目。

**后续动作：** 创建 `hn-agent/`，进入该目录后执行：

```bash
mkdir -p hn-agent
cd hn-agent
go mod init github.com/aaron/go-with-ai/hn-agent
```

生成后的 `go.mod` 应该位于：

```text
hn-agent/go.mod
```

## 2026-06-27 前三阶段代码为什么要多注释

**我：** 文档的教练笔记里，内部代码应该多一些注释。前三阶段我对 Go 完全不理解，需要更多解释。

**教练：** 对。前三阶段的代码示例不是生产代码，而是教学代码。它应该解释 `package`、`import`、`struct`、测试函数、`context`、`net/http`、goroutine、channel 等新概念，让学习者看代码时知道每一块为什么存在。

**结论：** 前三阶段文档代码要用高密度教学注释，先帮助理解；第四阶段以后再逐步靠近真实项目的注释密度。

**后续动作：** 阶段一教练笔记补充逐行注释；新增阶段二、阶段三教练笔记，并把 `context`、HTTP、JSON、`httptest`、goroutine、`errgroup`、semaphore 等内容写成注释充分的示例。

## 2026-06-27 第一次看到 Go 测试代码怎么读

**我：** 这是我第一次看到大量 Go 代码，需要介绍一下这段测试代码的语法，并同步到文档里。

**教练：** 这段代码最核心的是 table-driven test。先定义一组测试用例，每个用例包含输入 `story` 和期望输出 `want`，然后用 `for range` 遍历，每个 case 用 `t.Run` 跑成子测试。

**结论：** 第一次读 Go 测试代码时，先看四层结构：`package/import`、测试函数、测试表、遍历和断言。

**后续动作：** 阅读 [Go 测试代码语法拆解](/notes/stage-1-startup/go-test-syntax)，重点记住 `:=`、`[]struct`、复合字面量、`range`、`t.Run`、`got/want`、`t.Fatalf`。

## 2026-06-27 写完 `story_test.go` 后怎么测试

**我：** `internal/hn/story_test.go` 已经写好了，要怎么测试？为什么第六步之后不做任何运行的操作，而是直接进入第七步写最小 CLI？

**教练：** 写完测试后应该立刻运行 `go test ./internal/hn`。第七步写 CLI 是下一层验证，不应该替代单元测试。原文档把测试放到最终验收里，对新手不够细。

**结论：** 每完成一个最小学习单元，就立刻运行对应命令。阶段一第六步之后先跑 `go test ./internal/hn`，再进入 CLI。

**后续动作：** 在 `hn-agent/` 目录内执行：

```bash
go test ./internal/hn
go test -v ./internal/hn
```

通过后再继续第七步：创建 `cmd/hnctl/main.go`。

## 2026-06-27 什么是领域函数

**我：** 什么是领域函数？

**教练：** 领域函数不是 Go 语法词，而是工程建模里的说法。领域就是项目正在解决的问题范围；在这个项目里，Hacker News story、digest、summary 都属于领域。领域函数就是表达这些业务概念或业务规则的函数。

**结论：** `IsValidStory` 是领域函数，因为它表达“一条 story 在本项目里怎样才算有效”，而不是处理打印、命令行或程序启动。

**后续动作：** 阶段一先用 `IsValidStory` 练习把业务判断写成可测试函数；后面再逐步增加 `TopStories`、`FetchItems`、`BuildDigest` 这类更贴近项目的函数。
