# 03 并发抓取

## 本步目标

在 HN client 基础上实现有限并发批量抓取 story 详情。

## 目标函数

```go
FetchItems(ctx context.Context, ids []int64, concurrency int) ([]Item, error)
```

## 关键要求

- 使用 goroutine。
- 使用 `errgroup.WithContext`。
- 控制最大并发数。
- 明确失败策略。
- 对 timeout 和取消写测试。

## 验收标准

先验收 `FetchItems` 本身：

```bash
cd hn-agent
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
```

测试至少覆盖：

- `concurrency <= 0` 返回 error。
- 成功批量抓取多个 item。
- 返回结果顺序和输入 ids 一致。
- 任意一个 item 请求失败时，整体返回 error。
- 进阶：并发上限不超过 `concurrency`。
- 进阶：timeout 或取消能通过 `context` 传播。

CLI 接好后再执行：

```bash
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

## 回填记录

完成后把 goroutine、`errgroup`、bounded concurrency 记录到：

```text
website/notes/stage-3-concurrency/
```
