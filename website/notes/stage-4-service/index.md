# 阶段 4：服务期

目标：把 CLI 能力组织成一个可长期运行的 Go 服务。

::: tip 当前状态
阶段四已完成。当前服务已经具备服务入口、健康检查、配置读取、结构化日志、graceful shutdown 和开发期 live reload 方案。
:::

## 学习目标

- 使用 `net/http` 启动 HTTP server。
- 实现 `/healthz`。
- 使用 `log/slog` 做结构化日志。
- 实现 graceful shutdown。
- 设计基础配置加载方式。
- 理解 `context` 在服务关闭中的作用。

## 完成任务

- [x] 建立 `cmd/hn-agent`。
- [x] 实现 `/healthz`。
- [x] 增加配置读取。
- [x] 增加结构化日志。
- [x] 增加 graceful shutdown。
- [x] 补充 `context` 专题和开发期 live reload。

## 知识记录

- [阶段四教练笔记](/notes/stage-4-service/coach-notes)
- [Go context 入门：服务如何知道“该停了”](/notes/stage-4-service/go-context)

### 服务入口

服务入口应该尽量薄：读取配置、初始化依赖、注册路由、启动 server。

### API 草案

```text
GET /healthz
GET /api/v1/digests/latest
POST /api/v1/jobs/digest
```

`/api/v1/digests/latest` 和 `POST /api/v1/jobs/digest` 暂列为后续扩展，会随阶段五的 AI 摘要和 digest 数据结构一起推进。

## 问题清单

- 第一版是否需要 router 库，还是先使用标准库？
- 配置先只用环境变量，还是同时支持配置文件？

## 验收标准

```bash
cd hn-agent
go test ./...
go run ./cmd/hn-agent
curl http://localhost:8080/healthz
```

## 本阶段记录

- 阶段四完成：服务骨架、健康检查、配置、日志、优雅关闭、`context` 专题和开发期 live reload 已整理完成。
- 后续进入阶段五：AI 整合期。
