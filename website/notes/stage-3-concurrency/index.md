# 阶段 3：并发期

目标：能用有限并发批量拉取 story 详情，并正确处理取消、错误和资源边界。

::: tip 当前状态
阶段三已完成验收，并已由学习者同步到远程仓库。后续进入阶段四：服务期。
:::

## 学习目标

- 理解 goroutine 的生命周期。
- 理解 channel 的同步和通信语义。
- 使用 `errgroup.WithContext` 管理一组会失败的任务。
- 使用 semaphore 或 worker pool 控制并发数。
- 明确局部失败和整体失败策略。
- 理解 Go 1.22 之后 `for range` 循环变量每轮独立，旧写法 `i, id := i, id` 在本项目中不需要。

## 术语速查

| 英文 | 中文理解 | 阶段三先怎么记 |
| --- | --- | --- |
| goroutine | Go 的轻量并发任务 | 用 `go` 启动，但启动后要管理生命周期 |
| concurrency | 并发 | 同一段时间内处理多个任务，重点是任务管理 |
| parallelism | 并行 | 多个任务真的同时在多个 CPU 核上跑 |
| channel | 通道 | goroutine 之间通信或同步的管道 |
| semaphore | 信号量 / 并发许可证 | 控制最多同时运行几个任务 |
| errgroup | 错误组 / goroutine 任务组 | 管理一组会返回 error 的 goroutine |
| cancellation | 取消 | 通知还没完成的任务尽快停下来 |
| context | 上下文 | 传递取消、超时和请求范围信息 |

## 完成任务

- [x] 实现 `FetchItems(ctx, ids, concurrency)`。
- [x] 控制最大并发数。
- [x] 为阶段三核心行为写测试。
- [x] 记录错误传播策略。
- [x] 记录与 `Promise.all()` 的差异。

## 知识记录

- [阶段三教练笔记](/notes/stage-3-concurrency/coach-notes)

### 有限并发

并发不是越多越好。对外部 API 做 fan-out 时，要用明确的并发上限保护自己和对方服务。

```go
sem := make(chan struct{}, concurrency)
```

## 问题清单

- 应该保留成功结果并返回部分错误，还是直接返回失败？
- 并发数默认值应该放在 config 还是调用参数？

## 验收标准

先验收 package 层的核心行为：

```bash
cd hn-agent
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
```

如果想检查并发读写风险，可以额外执行：

```bash
go test -race ./internal/hn
```

CLI 接好之后，再做集成验收：

```bash
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

阶段三真正要验收的是：有限并发、结果顺序稳定、错误能返回、取消能传播。CLI 命令只是最后一层入口验收。

## 本阶段记录

- [阶段三教练笔记](/notes/stage-3-concurrency/coach-notes)：goroutine、`errgroup`、channel semaphore、取消和错误传播。
- 2026-06-29：阶段三完成验收，代码已同步远程；下一步进入阶段四服务期。
