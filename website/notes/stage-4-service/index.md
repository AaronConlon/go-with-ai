# 阶段 4：服务期

目标：把 CLI 能力组织成一个可长期运行的 Go 服务。

## 学习目标

- 使用 `net/http` 启动 HTTP server。
- 实现 `/healthz`。
- 实现读取最新 digest 的 API。
- 使用 `log/slog` 做结构化日志。
- 实现 graceful shutdown。
- 设计基础配置加载方式。

## 完成任务

- [ ] 建立 `cmd/hn-agent`。
- [ ] 实现 `/healthz`。
- [ ] 实现 `/api/v1/digests/latest`。
- [ ] 增加配置读取。
- [ ] 增加 graceful shutdown。

## 知识记录

### 服务入口

服务入口应该尽量薄：读取配置、初始化依赖、注册路由、启动 server。

### API 草案

```text
GET /healthz
GET /api/v1/digests/latest
POST /api/v1/jobs/digest
```

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

- 后续记录追加在本页，或按主题拆出子文档。
