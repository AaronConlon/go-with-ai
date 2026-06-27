# Go 学习路线

## 阶段 1：启动期

目标：能创建、运行和测试一个 Go module。

重点：

- 安装 Go。
- 理解 `go mod init`、`go mod tidy`。
- 学会 `go run`、`go test`、`go build`。
- 熟悉 package、函数、变量、slice、map、struct。
- 建立 `cmd/`、`internal/` 的基本项目结构直觉。

交付物：

- 一个 `hnctl` CLI 雏形。
- 一个最小的单元测试。

## 阶段 2：网络期

目标：能稳定调用 HTTP API 并处理 JSON。

重点：

- `net/http` client。
- `http.NewRequestWithContext`。
- `encoding/json`。
- status code 和错误处理。
- `context.WithTimeout`。

交付物：

- 拉取 HN `topstories`。
- 根据 story id 拉取 item 详情。

## 阶段 3：并发期

目标：能安全地做有限并发 I/O。

重点：

- goroutine。
- channel。
- `sync.WaitGroup`。
- `errgroup.WithContext`。
- bounded concurrency。

交付物：

- 并发拉取前 20 条 HN story。
- 设置最大并发数。
- 任意请求失败时能正确取消或降级。

## 阶段 4：服务期

目标：能把功能组织成长期运行服务。

重点：

- HTTP server。
- `/healthz`。
- `/metrics`。
- 配置读取。
- 结构化日志。
- graceful shutdown。

交付物：

- 一个能启动的 Go 服务。
- 一个最新 digest 查询接口。

## 阶段 5：AI 整合期

目标：能稳定调用 LLM 生成结构化摘要。

重点：

- 后端保存 API key。
- 请求超时。
- 速率限制。
- 重试和退避。
- JSON schema 或结构化输出校验。

交付物：

- story 到 summary 的稳定转换。
- 摘要结果可保存、可复查。

## 阶段 6：验证与部署期

目标：能把项目跑进 CI 并部署。

重点：

- table-driven test。
- `httptest`。
- Docker multi-stage build。
- GitHub Actions。
- 环境变量和 secret 管理。

交付物：

- CI 跑通。
- 镜像可构建。
- 服务可部署。

