# 04 服务化

## 本步目标

把 CLI 能力组织成长期运行的 Go HTTP 服务。

## 目标接口

```text
GET /healthz
GET /api/v1/digests/latest
POST /api/v1/jobs/digest
```

## 关键要求

- 使用 `net/http`。
- 服务入口保持薄。
- 使用 `log/slog`。
- 支持 graceful shutdown。
- 配置通过环境变量读取。

## 验收标准

```bash
cd hn-agent
go test ./...
go run ./cmd/hn-agent
curl http://localhost:8080/healthz
```

## 回填记录

完成后把服务入口、路由、配置和 shutdown 记录到：

```text
website/notes/stage-4-service/
```
