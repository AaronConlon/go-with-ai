# 实战目录

这个目录用于从零开始，一步步学习并完成 Hacker News 摘要推送 Agent 项目。

它和 `website/` 的关系：

- `practice/`：实战推进台账，记录每一步要做什么、做到了哪里、如何验收。
- `website/`：对外可浏览的知识库和学习记录。
- `background/`：背景研究、路线判断和方案沉淀。

## 推进规则

1. 每一步都要有明确目标。
2. 每一步都要能用命令或结果验收。
3. 先完成最小可运行版本，再扩展能力。
4. 学到的知识回填到 `website/notes/` 对应阶段。
5. 设计判断回填到 `background/`。

## 实战路线

| 步骤 | 文档 | 目标 |
| --- | --- | --- |
| 00 | [从零开始](./00-from-zero.md) | 明确项目目标、边界和学习节奏 |
| 01 | [项目骨架](./01-project-skeleton.md) | 创建 Go module、CLI 和测试入口 |
| 02 | [HN Client](./02-hn-client.md) | 拉取 `topstories` 和 item 详情 |
| 03 | [并发抓取](./03-concurrent-fetch.md) | 用有限并发批量抓取 story |
| 04 | [服务化](./04-service-api.md) | 提供 health check 和 digest API |
| 05 | [AI 摘要](./05-ai-summary.md) | 生成结构化中文摘要 |
| 06 | [验证与部署](./06-verify-deliver.md) | 补测试、CI 和部署流程 |

## 当前状态

当前处于步骤 00：从零开始。

