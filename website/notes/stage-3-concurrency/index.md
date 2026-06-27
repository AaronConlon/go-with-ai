# 阶段 3：并发期

目标：能用有限并发批量拉取 story 详情，并正确处理取消、错误和资源边界。

## 学习目标

- 理解 goroutine 的生命周期。
- 理解 channel 的同步和通信语义。
- 使用 `errgroup.WithContext` 管理一组会失败的任务。
- 使用 semaphore 或 worker pool 控制并发数。
- 明确局部失败和整体失败策略。

## 完成任务

- [ ] 实现 `FetchItems(ctx, ids, concurrency)`。
- [ ] 控制最大并发数。
- [ ] 为超时和取消写测试。
- [ ] 记录错误传播策略。
- [ ] 记录与 `Promise.all()` 的差异。

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

```bash
cd hn-agent
go test ./internal/hn
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

## 本阶段记录

- [阶段三教练笔记](/notes/stage-3-concurrency/coach-notes)：goroutine、`errgroup`、channel semaphore、取消和错误传播。
