# 阶段 2：网络期

目标：能稳定调用 Hacker News API，并把外部 JSON 转成 Go struct。

## 学习目标

- 使用 `net/http` 发起请求。
- 使用 `http.NewRequestWithContext` 绑定取消和超时。
- 使用 `encoding/json` 解析响应。
- 对 status code、空响应和 JSON 错误做明确处理。
- 写出可测试的 HN client。

## 完成任务

- [ ] 实现 `TopStories(ctx)`。
- [ ] 实现 `Item(ctx, id)`。
- [ ] 增加 HTTP timeout。
- [ ] 为 HN client 写 `httptest`。
- [ ] 记录 Hacker News item 字段含义。

## 知识记录

- [阶段二教练笔记](/notes/stage-2-network/coach-notes)
- [Go Client 构造函数语法拆解](/notes/stage-2-network/go-client-constructor-syntax)
- [Go Method Receiver 语法拆解](/notes/stage-2-network/go-method-receiver-syntax)
- [Go httptest 使用说明](/notes/stage-2-network/go-httptest-syntax)

### 外部 I/O 的基本规则

所有外部请求都应该带 timeout，并且错误要返回给调用方，而不是只打印日志。

### HN API 草案

```text
GET https://hacker-news.firebaseio.com/v0/topstories.json
GET https://hacker-news.firebaseio.com/v0/item/{id}.json
```

## 问题清单

- 单个 story 拉取失败时，是整批失败还是跳过？
- HN API 是否需要本地缓存？

## 验收标准

```bash
cd hn-agent
go test ./internal/hn
go run ./cmd/hnctl top --limit=10
```

说明：`top --limit=10` 是本阶段把 HN client 接入 CLI 后的目标命令。接入前，先以 `go test ./internal/hn` 验收 `Client.TopStories`。

## 本阶段记录

- [阶段二教练笔记](/notes/stage-2-network/coach-notes)：`context`、`net/http`、JSON 解码、错误返回和 `httptest`。
- [Go Client 构造函数语法拆解](/notes/stage-2-network/go-client-constructor-syntax)：解释 `NewClient`、`*Client`、`&Client{}`、嵌套初始化和 `time.Second`。
- [Go Method Receiver 语法拆解](/notes/stage-2-network/go-method-receiver-syntax)：解释 `(c *Client)`、method、receiver、返回值 `[]int64, error`。
- [Go httptest 使用说明](/notes/stage-2-network/go-httptest-syntax)：解释 `httptest.NewServer` 应该写在哪个文件，以及如何配合 `Client` 测 HTTP 请求。
