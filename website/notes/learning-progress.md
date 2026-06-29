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

- [x] 实现 `TopStories(ctx)`。
- [x] 实现 `Item(ctx, id)`。
- [x] 增加 HTTP timeout。
- [x] 为 HN client 写 `httptest`。
- [x] 记录 Hacker News item 字段含义。

### 阶段 3：并发期

- [x] 实现 `FetchItems(ctx, ids, concurrency)`。
- [x] 控制最大并发数。
- [x] 为阶段三核心行为写测试。
- [x] 记录错误传播策略。
- [x] 记录与 `Promise.all()` 的差异。

### 阶段 4：服务期

- [x] 建立 `cmd/hn-agent`。
- [x] 实现 `/healthz`。
- [x] 增加配置读取。
- [x] 增加结构化日志。
- [x] 增加 graceful shutdown。
- [x] 记录 `context`、`slog` 和开发期 live reload。

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

## 2026-06-29 阶段三完成验收并同步远程

- 阶段：阶段 3：并发期。
- 对应任务：实现有限并发批量抓取、错误传播、顺序稳定和阶段三测试。
- 操作：学习者完成 `FetchItems(ctx, ids, concurrency)` 及相关测试，修正阶段三测试断言后重新验收，并将代码同步到远程仓库。
- 结果：阶段三实战代码已完成并同步远程；文档站同步进入阶段四学习上下文。
- 判断：阶段三目标达成。当前项目已经具备基于 `errgroup`、goroutine 和 channel semaphore 的有限并发批量抓取能力，并明确了“任意一个 item 失败，整体返回 error”的错误策略。
- 下一步：进入阶段 4：服务期，开始设计长期运行的 HTTP service，优先实现 `/healthz` 和服务启动/关闭骨架。

## 2026-06-29 阶段四 `/healthz` path 常量化

- 阶段：阶段 4：服务期。
- 对应任务：实现 `/healthz` 并编写 handler 测试。
- 操作：将阶段四教练笔记里的路由字符串推荐写法调整为 `const healthzPath = "/healthz"`。
- 结果：`NewMux()` 注册路由和 `server_test.go` 构造请求都使用同一个常量。
- 判断：这是比重复写字符串更清晰的学习写法，也能顺便理解 Go 的 package-level `const`、小写 unexported 名字和大写 exported 名字。
- 下一步：学习者在 `hn-agent/internal/server/server.go` 中按推荐写法实现，并运行 `go test ./internal/server` 验收。

## 2026-06-29 阶段四函数调用参数规则

- 阶段：阶段 4：服务期。
- 对应任务：理解 `httptest.NewRequest(http.MethodGet, healthzPath, nil)`。
- 操作：补充 Go 函数调用规则：普通函数调用必须和函数签名参数个数一致。
- 结果：阶段四教练笔记解释了第三个参数 `nil` 的含义，并说明 Go 不会自动补 `undefined` 或 `nil`。
- 判断：这是从 JavaScript 切换到 Go 时很重要的基础心智，尤其容易在标准库函数调用中反复遇到。
- 下一步：实现测试时先查函数签名，再逐个确认参数含义和是否允许传 `nil`。

## 2026-06-29 阶段四配置章节前置

- 阶段：阶段 4：服务期。
- 对应任务：组织 `internal/config` 和 `cmd/hn-agent/main.go` 的学习顺序。
- 操作：把“配置先怎么做”章节移动到“服务入口怎么写”之前。
- 结果：读者先看到 `internal/config/config.go` 的 `Config` 和 `Load()`，再进入 `main.go` 中的 `cfg := config.Load()`。
- 判断：这个顺序更符合第一次实现的理解路径，减少入口代码里突然出现未解释 package 的割裂感。
- 下一步：按文档顺序先实现 `internal/config`，再实现服务入口。

## 2026-06-29 阶段四配置返回值设计

- 阶段：阶段 4：服务期。
- 对应任务：理解 `func Load() Config` 的返回值设计。
- 操作：在配置章节补充为什么第一版返回 `Config` 值，而不是 `*Config` 指针。
- 结果：文档说明了小配置、默认值、只读语义和避免 nil 判断的理由。
- 判断：这个点能帮助学习者建立 Go 的值语义心智，不把“struct”误解成必须像 JS object 一样通过引用传递。
- 下一步：阶段四先保持 `func Load() Config`；如果后续配置读取会失败，再升级为 `func Load() (Config, error)`。

## 2026-06-29 阶段四配置函数可见性

- 阶段：阶段 4：服务期。
- 对应任务：理解 `Load` 和 `load` 的 package 可见性差异。
- 操作：在配置章节补充为什么入口必须调用大写 exported 的 `config.Load()`。
- 结果：文档说明了小写 `load` 会在 `config` package 内部未使用时触发 unused warning，且不能被 `package main` 调用。
- 判断：这是 Go package 边界的基础规则，也是阶段四第一次把 `internal/config` 接入 `cmd/hn-agent` 时容易踩到的点。
- 下一步：学习者保持 `func Load() Config`，并在服务入口中 import `internal/config` 后使用 `cfg.Addr`。

## 2026-06-29 阶段四新增新手错误范例

- 阶段：阶段 4：服务期。
- 对应任务：把配置接入服务入口。
- 操作：把 `load` unused、`undefined: config`、`cfg` unused 整理成“新手容易错”范例。
- 结果：阶段四教练笔记新增错误代码、原因解释、修正代码和三步排查顺序。
- 判断：这个范例能帮助学习者把 Go 的 exported 名字、package import 和 unused 规则串起来理解。
- 下一步：学习者按范例修正 `cmd/hn-agent/main.go`，再运行 `go test ./...` 验证。

## 2026-06-29 阶段四解释 `<-ctx.Done()`

- 阶段：阶段 4：服务期。
- 对应任务：理解服务入口等待退出信号的写法。
- 操作：在阶段四教练笔记中补充 `<-ctx.Done()` 的拆解。
- 结果：文档说明了 `ctx.Done()` 返回 channel、`<-` 表示接收并阻塞等待、`signal.NotifyContext` 收到系统信号后会取消 context。
- 判断：这是服务期第一次把 `context`、channel receive 和系统信号连在一起，必须单独解释。
- 下一步：学习者实现入口时把 `<-ctx.Done()` 当成 main goroutine 的暂停点，收到 Ctrl+C 后再进入 graceful shutdown。

## 2026-06-29 阶段四服务入口逐段导读

- 阶段：阶段 4：服务期。
- 对应任务：理解 `cmd/hn-agent/main.go` 的完整 `main()`。
- 操作：在阶段四教练笔记中增加 `main()` 逐段导读。
- 结果：文档按执行顺序解释了配置、logger、mux、`http.Server`、`signal.NotifyContext`、goroutine、`ListenAndServe` 错误判断、`<-ctx.Done()`、shutdown timeout 和 `srv.Shutdown`。
- 判断：这是阶段四新概念最密集的一段代码，需要比普通示例更丰富的解释。
- 下一步：学习者实现入口时按导读逐段写，不需要一次记住所有语法。

## 2026-06-29 阶段四新增 `context` 独立专题

- 阶段：阶段 4：服务期。
- 对应任务：理解服务入口里的 `context`。
- 操作：新增 `go-context` 专题页，生成三张小黑正文配图，并同步中英文入口和侧边栏。
- 结果：文档独立解释了 `context.Background()`、`signal.NotifyContext`、`ctx.Done()`、`<-ctx.Done()`、`context.WithTimeout` 和 shutdown context 的职责区别。
- 判断：`context` 是阶段四服务生命周期的核心概念，单独成篇比塞在 `main()` 导读里更适合复习。
- 下一步：学习者先阅读 `go-context`，再回到服务入口实现 graceful shutdown。

## 2026-06-29 阶段四 graceful shutdown 错误判断

- 阶段：阶段 4：服务期。
- 对应任务：验证 `Ctrl+C` 后的 graceful shutdown 输出。
- 操作：根据学习者输出定位 `ListenAndServe()` 错误判断写反的问题，并补充到教练笔记。
- 结果：文档新增 `http.ErrServerClosed` 常见误判范例，说明少写 `!` 会把正常关闭记录成失败。
- 判断：`http.ErrServerClosed` 是 graceful shutdown 的预期返回值，应该被忽略；真正失败是“返回了 error 且不是 `http.ErrServerClosed`”。
- 下一步：学习者把条件改成 `err != nil && !errors.Is(err, http.ErrServerClosed)` 后重新运行服务，用 `Ctrl+C` 验证输出。

## 2026-06-29 阶段四自定义日志时间

- 阶段：阶段 4：服务期。
- 对应任务：调整 `log/slog` 输出格式。
- 操作：在阶段四教练笔记补充 `slog.HandlerOptions.ReplaceAttr` 示例。
- 结果：文档说明了如何把默认 `time` 字段转换为上海时区，并格式化为 `2006-01-02 15:04:05`。
- 判断：这是阶段四日志能力的自然延伸，也顺便引入 Go 特有的 time layout 心智。
- 下一步：学习者可以把 logger 初始化抽成 `newLogger()`，再运行服务观察日志时间格式。

## 2026-06-29 阶段四开发期 live reload

- 阶段：阶段 4：服务期。
- 对应任务：提高 `go run ./cmd/hn-agent` 的开发迭代效率。
- 操作：在阶段四教练笔记补充 Go 开发期 live reload 方案。
- 结果：文档推荐使用 `air`，并给出 `hn-agent/.air.toml` 示例，指定 `go build -o ./tmp/hn-agent ./cmd/hn-agent`。
- 判断：Go 本身没有内置 JS 风格 dev server 热加载；阶段四用 watcher 自动重编译和重启服务即可。
- 下一步：学习者在 `hn-agent/` 内配置 `air`，保存代码后自动重启服务。

## 2026-06-30 阶段四澄清 `Shutdown(ctx)` 因果关系

- 阶段：阶段 4：服务期。
- 对应任务：理解 `srv.Shutdown(shutdownCtx)` 中 context 的作用。
- 操作：在 `go-context` 专题补充说明：调用 `Shutdown` 才触发关闭，`shutdownCtx` 只限制关闭等待时间。
- 结果：文档明确了 `Shutdown` 的流程：停止接新连接、关闭空闲连接、等待活跃请求，直到完成或 context 超时。
- 判断：这是避免误解 `context` 的关键点；context 传控制信号和期限，但不主动执行业务动作。
- 下一步：学习者把 `shutdownCtx` 理解成 graceful shutdown 的计时器，而不是关闭触发器。

## 2026-06-30 阶段四补充 context 超时状态

- 阶段：阶段 4：服务期。
- 对应任务：理解 `context.WithTimeout` 超时后的状态变化。
- 操作：在 `go-context` 专题补充 `Done()`、`Err()`、`context deadline exceeded` 和 `context canceled` 的区别。
- 结果：文档说明了 5 秒 deadline 到期后 context 会关闭 `Done()` channel，并让 `Err()` 返回 `context deadline exceeded`。
- 判断：这能帮助学习者理解 context 是通知机制，不是强制杀掉工作。
- 下一步：学习者可以用小片段打印 `shutdownCtx.Err()`，观察超时取消和手动取消的区别。

## 2026-06-30 阶段四完成验收并准备同步远程

- 阶段：阶段 4：服务期。
- 对应任务：完成 HTTP service 骨架、健康检查、配置、日志、优雅关闭和开发体验补充。
- 操作：运行 Go 测试、文档站构建，并手动请求 `/healthz` 验收。
- 结果：

```text
go test ./...
?   	github.com/aaron/go-with-ai/hn-agent/cmd/hn-agent	[no test files]
?   	github.com/aaron/go-with-ai/hn-agent/cmd/hnctl	[no test files]
?   	github.com/aaron/go-with-ai/hn-agent/internal/config	[no test files]
ok  	github.com/aaron/go-with-ai/hn-agent/internal/hn
ok  	github.com/aaron/go-with-ai/hn-agent/internal/server
```

```text
npm run build
vitepress v1.6.4
build complete
```

```text
GET /healthz -> OK
```

- 判断：阶段四学习目标达成。当前服务已经具备可运行入口、健康检查、配置读取、结构化日志、graceful shutdown、`context` 专题说明和 Air 开发期 live reload 方案。`/api/v1/digests/latest` 暂不作为本阶段完成项，后续随 AI 摘要能力和 digest 数据结构一起推进。
- 下一步：提交并推送阶段四代码与文档，然后进入阶段 5：AI 整合期。
