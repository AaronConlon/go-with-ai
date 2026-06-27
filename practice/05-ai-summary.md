# 05 AI 摘要

## 本步目标

调用 LLM API，把 story 转换成结构化简体中文摘要。

## 目标结构

```json
{
  "summary": "两到四句中文摘要",
  "why_it_matters": "一行价值说明",
  "tags": ["Go", "Infra"]
}
```

## 关键要求

- API key 只保存在后端环境变量中。
- 请求有 timeout。
- 对 429 / 5xx 做有限重试。
- 对输出做 JSON 校验。
- 保存 prompt 版本，便于后续复盘。

## 验收标准

```bash
cd hn-agent
go test ./internal/summary
go run ./cmd/hnctl summarize --id=...
```

## 回填记录

完成后把 API 调用、结构化输出、重试和限流记录到：

```text
website/notes/stage-5-ai-integration/
```
