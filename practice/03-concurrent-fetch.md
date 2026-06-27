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

```bash
go test ./internal/hn
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

## 回填记录

完成后把 goroutine、`errgroup`、bounded concurrency 记录到：

```text
website/notes/stage-3-concurrency/
```

