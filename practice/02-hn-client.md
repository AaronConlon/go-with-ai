# 02 HN Client

## 本步目标

实现 Hacker News API client，能拉取 story id 列表和单条 item 详情。

## 目标函数

```go
TopStories(ctx context.Context) ([]int64, error)
Item(ctx context.Context, id int64) (Item, error)
```

## 关键要求

- 使用 `net/http`。
- 使用 `http.NewRequestWithContext`。
- 所有请求有 timeout。
- 非 2xx 响应返回错误。
- 使用 `httptest` 写测试。

## 验收标准

```bash
cd hn-agent
go test ./internal/hn
go run ./cmd/hnctl top --limit=10
```

## 回填记录

完成后把 HTTP、JSON、错误处理和 HN API 字段记录到：

```text
website/notes/stage-2-network/
```
