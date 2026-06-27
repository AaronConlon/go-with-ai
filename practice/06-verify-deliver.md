# 06 验证与部署

## 本步目标

建立测试、构建、CI 和部署闭环。

## 关键要求

- `go test ./...` 稳定通过。
- HN client 使用 `httptest`。
- handler 有基础测试。
- 增加 Dockerfile。
- 增加 GitHub Actions。
- 文档记录环境变量和 secret。

## 验收标准

```bash
cd hn-agent
go test ./...
docker build -t hn-agent .
```

## 回填记录

完成后把测试策略、CI、Docker 和部署记录到：

```text
website/notes/stage-6-delivery/
```
