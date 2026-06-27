# 阶段 5：AI 整合期

目标：稳定调用 LLM API，把 story 转成结构化中文摘要。

## 学习目标

- 后端读取并保护 API key。
- 为 LLM 请求设置 timeout。
- 对 429 和 5xx 做有限重试。
- 对模型输出做 JSON 校验。
- 设计 summary 数据结构。

## 完成任务

- [ ] 建立 `internal/summary`。
- [ ] 定义 `StorySummary`。
- [ ] 实现单条 story 摘要。
- [ ] 增加结构化输出校验。
- [ ] 记录失败重试策略。

## 知识记录

### 结构化输出

摘要不应该只保存一段自然语言文本。第一版至少应该包含：

```json
{
  "summary": "两到四句中文摘要",
  "why_it_matters": "一行价值说明",
  "tags": ["Go", "Infra"]
}
```

## 问题清单

- 摘要语言固定简体中文，还是允许配置？
- 是否需要把 prompt 版本写入数据库？
- 失败后是否进入重试队列？

## 验收标准

```bash
go test ./internal/summary
go run ./cmd/hnctl summarize --id=...
```

## 本阶段记录

- 后续记录追加在本页，或按主题拆出子文档。

