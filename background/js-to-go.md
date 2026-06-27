# 从 JavaScript / Node.js 迁移到 Go

## 核心心智差异

Go 不是“带类型的 JavaScript”。更合适的理解是：

> Go 是一种偏显式、偏工程默认值、偏服务端长期运行任务的语言。

JavaScript 写服务时，很多复杂度来自运行时、依赖和框架组合。Go 则把更多工程默认值放进语言、标准库和工具链里。

## 对照表

| JavaScript / Node.js | Go | 实践建议 |
| --- | --- | --- |
| `async/await` | goroutine + `context` | 不要机械翻译 Promise，先想任务生命周期和取消边界。 |
| `Promise.all()` | `errgroup` | 一组会失败的并发任务，用 `errgroup.WithContext`。 |
| `try/catch` | `if err != nil` | 在每个 I/O 边界显式处理错误。 |
| `package.json` | `go.mod` | 依赖由 Go modules 管理，脚本逻辑尽量放到 CI 或 Makefile。 |
| Express/Fastify | `net/http` / chi / Gin | MVP 先用标准库，路由复杂后再考虑 chi。 |
| `console.log` | `log/slog` | 从一开始就用结构化日志。 |
| npm scripts | `go test` / `go build` | 优先掌握 Go 官方命令。 |

## 最容易踩的坑

1. 把 interface 设计得太早、太大。
2. 用全局变量传配置和依赖。
3. 忽略 `context`，导致请求无法取消。
4. 捕获错误后只打印，不向上返回。
5. 为了“像 Node.js”而过早引入大型框架。

## 写 Go 时的自检问题

- 这个函数的错误是不是被调用方看见了？
- 这个外部 I/O 有没有 timeout？
- 这个 goroutine 有没有退出条件？
- 这个 interface 是调用方需要的，还是实现方臆想出来的？
- 这个 package 的职责能不能用一句话说清楚？

