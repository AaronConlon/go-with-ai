# Go Module

## 它解决什么问题

Go module 是 Go 官方的依赖和版本管理机制。一个 module 通常对应一个仓库或一个可独立发布的代码单元。

常用文件：

- `go.mod`：声明 module path、Go 版本和依赖。
- `go.sum`：记录依赖校验信息。

## 常用命令

初始化：

```bash
go mod init example.com/hn-agent
```

整理依赖：

```bash
go mod tidy
```

运行测试：

```bash
go test ./...
```

构建：

```bash
go build ./...
```

## 学习重点

- module path 是代码导入路径，不只是本地目录名。
- `internal/` 下的包只能被父级目录树内部导入。
- 不要手动编辑 `go.sum`。
- 提交代码时通常需要一起提交 `go.mod` 和 `go.sum`。

