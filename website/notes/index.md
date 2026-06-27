# 知识库

知识库按学习路径拆成六个模块。每个模块都用于保存这个阶段的完整记录：目标、任务、概念、代码片段、问题、复盘和验收结果。

## 六个模块

| 模块 | 阶段 | 记录重点 |
| --- | --- | --- |
| [阶段 1：启动期](/notes/stage-1-startup/) | Go module 和基础语法 | 环境、命令、项目结构、第一批测试 |
| [阶段 2：网络期](/notes/stage-2-network/) | HTTP 和 JSON | HN API、`context`、错误边界 |
| [阶段 3：并发期](/notes/stage-3-concurrency/) | goroutine 和有限并发 | `errgroup`、取消、并发数控制 |
| [阶段 4：服务期](/notes/stage-4-service/) | 长期运行服务 | API server、配置、日志、健康检查 |
| [阶段 5：AI 整合期](/notes/stage-5-ai-integration/) | LLM API 调用 | 结构化输出、重试、限流、密钥管理 |
| [阶段 6：验证与部署期](/notes/stage-6-delivery/) | 测试、CI 和部署 | `httptest`、Docker、GitHub Actions |

## 记录方式

每个阶段都保持同一套结构：

- `学习目标`：这个阶段要真正掌握什么。
- `完成任务`：以可运行结果为准。
- `知识记录`：概念、命令、代码片段和心智模型。
- `问题清单`：不确定、踩坑、需要验证的点。
- `验收结果`：用命令输出、测试结果或截图确认完成。

通用模板：

- [学习推进记录](/notes/learning-progress)
- [学习问题对话记录](/notes/learning-dialogues)
- [阶段记录模板](/notes/stage-record-template)
- [学习日志模板](/notes/learning-log-template)
