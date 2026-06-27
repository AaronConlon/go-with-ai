# 阶段 6：验证与部署期

目标：让项目具备可测试、可构建、可部署、可回归的工程闭环。

## 学习目标

- 使用 table-driven test 组织测试用例。
- 使用 `httptest` 测 HTTP client 和 handler。
- 构建 Docker multi-stage image。
- 使用 GitHub Actions 跑 CI。
- 记录部署环境变量和 secret。

## 完成任务

- [ ] `go test ./...` 稳定通过。
- [ ] HTTP client 覆盖异常响应。
- [ ] handler 覆盖主要 API。
- [ ] 写 Dockerfile。
- [ ] 写 GitHub Actions workflow。

## 知识记录

### 验证优先级

先覆盖最容易回归的边界：

- 外部 API 返回错误。
- JSON 字段缺失。
- context timeout。
- 摘要输出不是合法 JSON。
- digest 重复推送。

## 问题清单

- 第一版部署到哪里？
- SQLite 文件如何备份？
- secret 如何在本地和 CI 中分别管理？

## 验收标准

```bash
go test ./...
docker build -t hn-agent .
```

## 本阶段记录

- 后续记录追加在本页，或按主题拆出子文档。

