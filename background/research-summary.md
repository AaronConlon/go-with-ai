# 基于第一性原理的 Go 学习方案

## 核心结论

学习 Go 的最快路径不是先把语法从头背完，而是先抓住系统本质：Go 很适合构建可靠的 I/O 型服务、后台任务、网络程序和小型可部署工具。

当前最适合的练习项目是一个 Hacker News 摘要推送 Agent。它天然覆盖 Go 后端开发中最重要的能力：

- 定时任务和后台 worker。
- HTTP client 和 JSON 编解码。
- `context` 超时、取消和请求链路传递。
- goroutine、channel、`errgroup` 和有限并发。
- 数据去重、存储和幂等推送。
- 调用外部 LLM API 生成结构化摘要。
- 健康检查、日志、指标和部署。

## 为什么适合从这个项目开始

如果已经熟悉 JavaScript / Node.js、异步编程、HTTP / REST 和基础 DevOps，那么 Go 的学习重点应该放在心智迁移上：

- 从 `async/await` 和 `Promise.all()` 迁移到 goroutine、`context` 和 `errgroup`。
- 从 `try/catch` 迁移到显式的 `if err != nil`。
- 从 `package.json` 工具链迁移到 `go mod`、`go test`、`go build` 和 `gofmt`。
- 从框架优先迁移到标准库优先。
- 从运行时拼装迁移到编译期类型约束和清晰的数据结构。

Go 的优势不是“语法更多”，而是把工程默认值做得更收敛：格式化、测试、模块、HTTP、JSON、并发和部署路径都足够直接。

## 建议主线

学习线和产品线并行推进：

- 学习线：Go module、语法基础、错误处理、测试、并发、HTTP、数据库和部署。
- 产品线：HN 拉取、摘要生成、存储、推送、API、监控和 CI。

每个学习阶段都应该有一个可运行成果，而不是只留下阅读笔记。

## 技术判断

MVP 阶段建议保持克制：

- HTTP：先用标准库 `net/http`。
- JSON：先用标准库 `encoding/json`。
- 日志：用标准库 `log/slog`。
- 并发：goroutine + `context` + `errgroup`。
- 数据库：本地优先 SQLite，部署后可扩展到 PostgreSQL。
- 文档：VitePress + Markdown。
- AI 调用：后端服务持有 API key，避免把密钥暴露到前端或客户端。

