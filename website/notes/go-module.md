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

## `go mod tidy` 做什么

`go mod tidy` 会根据当前源码重新整理依赖。

可以把它理解成：

```text
扫描当前 module 的 Go 代码
-> 看所有 import 到底需要哪些依赖
-> 补上缺少的依赖
-> 删除已经不用的依赖
-> 更新 go.sum 校验信息
-> 重新判断 direct / indirect
```

它不是单纯格式化 `go.mod`，而是让 `go.mod` 和 `go.sum` 回到“和当前代码一致”的状态。

### 为什么 `// indirect` 会消失

假设刚执行：

```bash
go get golang.org/x/sync/errgroup
```

`go.mod` 可能先出现：

```go
require golang.org/x/sync v0.21.0 // indirect
```

`// indirect` 表示 Go 当前认为这个依赖不是源码直接 import 的依赖。

当代码里真正写了：

```go
import "golang.org/x/sync/errgroup"
```

并实际使用：

```go
g, ctx := errgroup.WithContext(ctx)
```

再执行：

```bash
go mod tidy
```

Go 会重新扫描源码，发现当前 module 直接 import 了 `golang.org/x/sync/errgroup`。于是 `golang.org/x/sync` 会变成直接依赖，`// indirect` 通常会被移除：

```go
require golang.org/x/sync v0.21.0
```

### 和 JavaScript 的类比

可以粗略类比：

```text
go.mod       package.json
go.sum       lockfile + checksum
go mod tidy  根据当前源码整理依赖
```

但 Go 更强调“源码 import 决定依赖”。也就是说，代码里 import 了什么，`go mod tidy` 就会尽量把 `go.mod` 整理成对应状态。

### 什么时候执行

建议在这些时候执行：

- 新增第三方 import 后。
- 删除某个第三方 import 后。
- 执行过 `go get` 后。
- 准备提交 Go 代码前。

在本项目里，Go 命令默认在 `hn-agent/` 内执行：

```bash
cd hn-agent
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
