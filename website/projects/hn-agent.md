# HN 摘要推送 Agent

## 项目定位

这是学习 Go 的贯穿项目：用 Go 构建一个后台服务，定时抓取 Hacker News 内容，调用 LLM 生成中文摘要，并推送到指定渠道。

它能覆盖后端开发的核心能力：

- HTTP client。
- JSON 编解码。
- 并发抓取。
- 数据去重。
- 存储。
- 后台任务。
- API server。
- AI API 调用。
- 测试、日志、指标和部署。

## MVP 能力

1. 拉取 HN `topstories`。
2. 拉取 story item 详情。
3. 去重并保存 story。
4. 调用 LLM 生成中文摘要。
5. 生成 digest。
6. 暴露 `/healthz` 和 `/api/v1/digests/latest`。

## 推荐包结构

```text
hn-agent/
  cmd/hn-agent/
  internal/config/
  internal/hn/
  internal/summary/
  internal/digest/
  internal/store/
  internal/notifier/
  internal/httpapi/
  internal/obs/
```

## 第一版验收

- `go test ./...` 通过。
- HN client 有 `httptest` 覆盖。
- 所有外部请求都有 timeout。
- 摘要生成失败不会影响已抓取数据保存。
- digest 具备幂等键，避免重复推送。

## 后续扩展

- 支持 `newstories` 和 `beststories`。
- 增加评论摘要。
- 增加主题过滤。
- 增加 Telegram、Slack、Email 推送。
- 增加 ChatGPT Action 查询入口。
