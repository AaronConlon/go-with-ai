# 学习推进记录

这份文档用于记录六个阶段学习任务的实际完成过程。它不是知识总结，而是执行台账：每完成一个任务，就记录“怎么做”和“结果是什么”。

## 记录规则

每条记录都尽量包含：

- 所属阶段。
- 对应任务。
- 操作命令或操作方式。
- 实际结果。
- 判断或备注。
- 下一步。

推荐格式：

```md
## YYYY-MM-DD 任务标题

- 阶段：
- 对应任务：
- 操作：
- 结果：
- 判断：
- 下一步：
```

## 阶段任务总览

### 阶段 1：启动期

- [x] 确认本机 Go 版本。
- [x] 创建 `hn-agent/` 项目目录。
- [x] 在 `hn-agent/` 内创建项目 module。
- [x] 建立 `hn-agent/cmd/hnctl` 和 `hn-agent/internal/hn`。
- [x] 写出第一个 CLI 输出。
- [x] 写出第一个单元测试。

### 阶段 2：网络期

- [ ] 实现 `TopStories(ctx)`。
- [ ] 实现 `Item(ctx, id)`。
- [ ] 增加 HTTP timeout。
- [ ] 为 HN client 写 `httptest`。
- [ ] 记录 Hacker News item 字段含义。

### 阶段 3：并发期

- [ ] 实现 `FetchItems(ctx, ids, concurrency)`。
- [ ] 控制最大并发数。
- [ ] 为超时和取消写测试。
- [ ] 记录错误传播策略。
- [ ] 记录与 `Promise.all()` 的差异。

### 阶段 4：服务期

- [ ] 建立 `cmd/hn-agent`。
- [ ] 实现 `/healthz`。
- [ ] 实现 `/api/v1/digests/latest`。
- [ ] 增加配置读取。
- [ ] 增加 graceful shutdown。

### 阶段 5：AI 整合期

- [ ] 建立 `internal/summary`。
- [ ] 定义 `StorySummary`。
- [ ] 实现单条 story 摘要。
- [ ] 增加结构化输出校验。
- [ ] 记录失败重试策略。

### 阶段 6：验证与部署期

- [ ] `go test ./...` 稳定通过。
- [ ] HTTP client 覆盖异常响应。
- [ ] handler 覆盖主要 API。
- [ ] 写 Dockerfile。
- [ ] 写 GitHub Actions workflow。

## 记录

## 2026-06-27 修正项目目录策略

- 阶段：阶段 1：启动期。
- 对应任务：创建项目 module 前的项目结构判断。
- 操作：重新确认仓库职责划分。
- 结果：决定不在仓库根目录执行 `go mod init`。根目录作为学习工作区，保存 `website/`、`background/`、`practice/` 等资料；Hacker News 实战代码放在独立目录 `hn-agent/`。
- 判断：这是更合理的结构。因为当前仓库不是单一 Go 项目，而是“学习资料 + 文档站 + 实战项目”的组合工作区。独立的 `hn-agent/` 目录可以让 Go module 边界更清晰，也方便后续构建、测试和部署。
- 下一步：创建 `hn-agent/`，进入该目录执行 `go mod init github.com/aaron/go-with-ai/hn-agent`。

## 2026-06-27 确认本机 Go 版本

- 阶段：阶段 1：启动期。
- 对应任务：确认本机 Go 版本。
- 操作：在项目根目录执行 `go version`。
- 结果：

```text
go version go1.25.1 darwin/arm64
```

- 判断：当前 Go 版本足够完成本项目学习和阶段一代码。阶段一只需要 Go module、基础语法、测试和 CLI 入口，不依赖特殊版本能力。
- 下一步：创建独立项目目录 `hn-agent/`，并在该目录内初始化 Go module。

## 2026-06-27 下一步：创建项目 module

- 阶段：阶段 1：启动期。
- 对应任务：创建项目 module。
- 操作：先在学习工作区根目录创建独立项目目录，然后进入该目录初始化 module：

```bash
mkdir -p hn-agent
cd hn-agent
go mod init github.com/aaron/go-with-ai/hn-agent
```

- 预期结果：`hn-agent/` 目录中生成 `go.mod`，内容类似：

```go
module github.com/aaron/go-with-ai/hn-agent

go 1.25.1
```

- 判断：`go.mod` 出现在 `hn-agent/` 后，表示 Hacker News 实战项目有了独立 Go module 边界。根目录仍然是学习工作区，用来保存 `website/`、`background/`、`practice/` 等资料。后续内部代码可以使用这个 module path 进行 import，例如：

```go
import "github.com/aaron/go-with-ai/hn-agent/internal/hn"
```

- 验收方式：

```bash
cd hn-agent
cat go.mod
go test ./...
```

- 注意：如果 `go test ./...` 暂时提示没有包或没有可测试文件，不一定是错误。因为此时还没创建 `cmd/hnctl` 和 `internal/hn`。真正的阶段一完整验收要等 CLI 和测试文件都创建后再做。
- 下一步：在 `hn-agent/` 目录内创建 `cmd/hnctl` 和 `internal/hn` 目录。

## 2026-06-27 待验证：运行第一个单元测试

- 阶段：阶段 1：启动期。
- 对应任务：写出第一个单元测试。
- 操作：`internal/hn/story_test.go` 写好后，先在 `hn-agent/` 目录内运行包级测试。

```bash
cd hn-agent
go test ./internal/hn
```

- 可选操作：如果想看到子测试详情，运行：

```bash
go test -v ./internal/hn
```

- 预期结果：`TestIsValidStory` 和三个子测试通过。
- 判断：单元测试通过后，说明 `internal/hn` 包的最小领域逻辑已经跑通，可以继续写 `cmd/hnctl/main.go`。
- 下一步：把 `go test ./internal/hn` 的输出贴给教练 review。

## 2026-06-27 阶段一完成验收

- 阶段：阶段 1：启动期。
- 对应任务：完成 Go module、最小领域模型、单元测试和 CLI。
- 操作：在 `hn-agent/` 内运行 Go 测试和 CLI；在 `website/` 内构建文档站。

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
go run ./cmd/hnctl version
```

```bash
cd website
npm run build
```

- 结果：

```text
?   	github.com/aaron/go-with-ai/hn-agent/cmd/hnctl	[no test files]
ok  	github.com/aaron/go-with-ai/hn-agent/internal/hn	(cached)
```

```text
hnctl 0.1.0
story valid: true
hnctl version is 0.1.0
```

```text
vitepress v1.6.4
build complete
```

- 判断：阶段一目标达成。`hn-agent/` 已经是独立 Go module；`internal/hn` 有最小领域模型和单元测试；`cmd/hnctl` 可以运行；文档站构建通过。
- 下一步：提交并推送阶段一代码与文档，然后进入阶段二：实现 Hacker News API client。
