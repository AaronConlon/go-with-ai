# 六阶段路线

## 阶段 1：启动期

目标：能创建、运行和测试一个 Go module。

必做：

- `go mod init`。
- `go test ./...`。
- 一个最小 CLI。
- 一个 table-driven test。

输出：

- `cmd/hnctl`。
- `internal/hn` 包雏形。

## 阶段 2：网络期

目标：能稳定调用 Hacker News API。

必做：

- 用 `http.NewRequestWithContext` 构造请求。
- 用 `json.Decoder` 解析响应。
- 对非 2xx 响应返回明确错误。
- 所有外部请求都设置 timeout。

输出：

- `TopStories(ctx)`。
- `Item(ctx, id)`。

## 阶段 3：并发期

目标：能有限并发拉取 story 详情。

必做：

- 用 `errgroup.WithContext` 管理一组并发任务。
- 用 channel 或 semaphore 限制并发数。
- 明确失败策略：全部失败、局部失败或降级返回。

输出：

- `FetchItems(ctx, ids, concurrency)`。

## 阶段 4：服务期

目标：能把 CLI 能力变成服务。

必做：

- `/healthz`。
- `/api/v1/digests/latest`。
- graceful shutdown。
- `slog` 结构化日志。

输出：

- `cmd/hn-agent`。

## 阶段 5：AI 整合期

目标：能生成稳定中文摘要。

必做：

- 后端读取 API key。
- 每次请求有 timeout。
- 对 429 / 5xx 做有限重试。
- 对模型输出做 JSON 校验。

输出：

- `internal/summary`。

## 阶段 6：验证与部署期

目标：能在 CI 中验证并部署。

必做：

- HTTP client 使用 `httptest`。
- 数据库逻辑有测试。
- Docker multi-stage build。
- GitHub Actions 运行 `go test ./...`。

输出：

- 可部署的服务。
- 可复现的构建流程。

