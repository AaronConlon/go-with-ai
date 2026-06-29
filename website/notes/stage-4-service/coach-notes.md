# 阶段四教练笔记

阶段四开始把前面写好的能力组织成一个可以长期运行的 HTTP service。

阶段一到三更像是在练“函数、package、测试和并发”。阶段四开始练的是另一个工程能力：

> 程序启动后不立刻结束，而是持续监听请求、处理请求、记录日志，并能被安全关闭。

## 本阶段心智模型

服务期先不要急着引入 Web 框架。第一版建议只用 Go 标准库：

```text
配置 -> 依赖初始化 -> 路由注册 -> HTTP server -> 等待请求 -> 优雅关闭
```

阶段四最重要的不是 API 多，而是边界清楚：

- `cmd/hn-agent` 只负责启动服务。
- `internal/server` 负责 HTTP handler 和路由。
- `internal/config` 负责读取配置。
- 核心业务仍然放在 `internal/hn` 等领域 package 里。

## 目标目录

建议阶段四先新增：

```text
hn-agent/cmd/hn-agent/main.go
hn-agent/internal/config/config.go
hn-agent/internal/server/server.go
hn-agent/internal/server/server_test.go
```

第一版先不要做太多 API。先把服务骨架跑通：

```text
GET /healthz
```

等 `/healthz` 能测试、能运行、能被 curl 访问后，再扩展：

```text
GET /api/v1/digests/latest
POST /api/v1/jobs/digest
```

## 阶段四术语表

| 英文 | 中文理解 | 阶段四先怎么记 |
| --- | --- | --- |
| HTTP server | HTTP 服务 | 长期运行，监听端口，接收请求并返回响应 |
| handler | 处理器 | 一个接收 request、写 response 的函数或对象 |
| route | 路由 | URL 路径和 handler 的绑定关系 |
| mux / router | 路由分发器 | 根据请求路径把请求交给对应 handler |
| middleware | 中间件 | 包在 handler 外层的通用逻辑，例如日志、鉴权、恢复 panic |
| health check | 健康检查 | 给人或平台判断服务是否活着的轻量接口 |
| graceful shutdown | 优雅关闭 | 收到停止信号后，不粗暴中断正在处理的请求 |
| config | 配置 | 运行时可调整的信息，例如端口、API 地址、超时时间 |
| structured logging | 结构化日志 | 日志用 key/value 记录，方便检索和机器处理 |

## `net/http` 是什么

`net/http` 是 Go 标准库里的 HTTP 包。阶段四会同时用到两部分：

- server 侧：`http.Server`、`http.ServeMux`、`http.HandlerFunc`。
- test 侧：`httptest.NewRecorder`、`httptest.NewRequest`。

先建立最小心智：

```text
浏览器 / curl / 其他 client
-> HTTP request
-> Go HTTP server
-> handler
-> HTTP response
```

## Handler 怎么读

最常见的 handler 函数长这样：

```go
// hn-agent/internal/server/server.go
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
```

这里两个参数很关键：

- `w http.ResponseWriter`：用来写 response，包括 status code、header、body。
- `r *http.Request`：当前请求，里面有 method、path、header、body、context 等信息。

可以读成：

```text
收到一个 request，用 w 写回 response。
```

注意：handler 不是 return response，而是通过 `w` 写 response。这和很多 JavaScript web 框架的心智不太一样。

## 第一版 server 推荐形状

阶段四推荐先把路由注册放进一个函数，方便测试复用。

这里还推荐把 `/healthz` 提成常量。它和 JavaScript 里的：

```js
const HEALTHZ_PATH = "/healthz"
```

是同一种思路：让“稳定、会被重复使用的字符串”有一个名字。Go 里如果只在当前 package 内部使用，通常写成小写开头的 `healthzPath`；如果将来要让别的 package 也使用，才改成大写开头的 `HealthzPath`。

```go
// hn-agent/internal/server/server.go
package server

import (
	"net/http"
)

// healthzPath 是健康检查接口的固定路径。
// 小写开头表示它只在 server package 内部可见。
const healthzPath = "/healthz"

// NewMux 创建并返回服务使用的路由分发器。
// *http.ServeMux 是 Go 标准库自带的 router。
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// HandleFunc 把路径和处理函数绑定起来。
	// 当请求路径是 /healthz 时，交给 healthz 处理。
	mux.HandleFunc(healthzPath, healthz)

	return mux
}

// healthz 是健康检查接口。
// 它应该非常轻量，不依赖外部 API，不访问慢资源。
func healthz(w http.ResponseWriter, r *http.Request) {
	// 健康检查只接受 GET。
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Header 要在 WriteHeader 或 Write 之前设置。
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// w.Write 返回写入字节数和 error。
	// 这里是极小的测试响应，先显式忽略返回值。
	_, _ = w.Write([]byte("ok\n"))
}
```

### `http.ServeMux` 是什么

`ServeMux` 可以理解成标准库自带的路由表：

```text
/healthz -> healthz handler
/api/v1/digests/latest -> latestDigest handler
```

第一版先用 `http.NewServeMux()` 足够，不需要马上引入 `chi`、`gin` 或 `echo`。等你真的遇到路由参数、中间件组合、复杂 API 分组，再考虑框架。

## Handler 测试怎么写

handler 可以不启动真实端口就测试。标准库提供了 `httptest`：

```go
// hn-agent/internal/server/server_test.go
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := NewMux()

	// NewRequest 创建一个测试请求。
	// 因为这个测试文件也是 package server，所以可以直接使用小写的 healthzPath。
	req := httptest.NewRequest(http.MethodGet, healthzPath, nil)

	// NewRecorder 创建一个 response 记录器。
	// handler 写给 w 的 status、header、body 都会被它记录下来。
	rec := httptest.NewRecorder()

	// ServeHTTP 手动把请求交给 mux。
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
```

### Go 调函数时不能省略参数

这行测试代码里有三个参数：

```go
req := httptest.NewRequest(http.MethodGet, healthzPath, nil)
```

`httptest.NewRequest` 的函数签名可以先简化读成：

```go
func NewRequest(method string, target string, body io.Reader) *http.Request
```

所以调用时必须传三个参数：

- `method`：请求方法，这里是 `http.MethodGet`。
- `target`：请求路径，这里是 `healthzPath`。
- `body`：请求 body，这里没有 body，所以显式传 `nil`。

这和 JavaScript 很不一样。JS 里少传参数时，函数内部拿到的通常是 `undefined`：

```js
function request(method, path, body) {}

request("GET", "/healthz") // body 是 undefined
```

Go 不会这样做。少传或多传都会编译失败：

```go
httptest.NewRequest(http.MethodGet, healthzPath)
// not enough arguments in call to httptest.NewRequest

httptest.NewRequest(http.MethodGet, healthzPath, nil, nil)
// too many arguments in call to httptest.NewRequest
```

也就是说，`nil` 不是 Go 自动帮你补的默认值，而是你主动传入的值。并且只有参数类型允许为 `nil` 时才能传 `nil`。这里的 `body` 类型是 `io.Reader` interface，可以用 `nil` 表示“没有请求体”；如果参数类型是 `string`、`int`、`bool` 这类值类型，就不能传 `nil`。

Go 也没有普通函数参数默认值。需要表达“可选配置”时，常见做法是：

- 传一个明确的零值，例如没有 body 就传 `nil`。
- 传配置 struct，例如 `ServerConfig{Addr: ":8080"}`。
- 使用 variadic 参数，例如 `fmt.Println(a ...any)`，这种函数签名里会显式出现 `...`。

### 为什么测试里能直接用 `NewMux`

`server_test.go` 第一行是：

```go
package server
```

这表示测试文件和 `server.go` 属于同一个 package。运行：

```bash
go test ./internal/server
```

时，Go 会把这个 package 里的普通 `.go` 文件和 `_test.go` 文件一起编译成一个测试二进制。可以先理解成：

```text
server.go      -> 生产代码
server_test.go -> 测试代码
go test        -> 临时编译到一起，然后运行 TestXxx
```

所以测试里可以直接调用：

```go
mux := NewMux()
```

这不是“把非测试代码错误地引入测试”，而是测试的目的：调用真实生产代码，然后验证它的行为是否正确。

如果测试不调用 `NewMux`，而是在测试里重新写一份 mux，就会变成“测试自己的假实现”，反而测不到真实服务路由。

### `package server` 和 `package server_test` 的区别

Go 测试有两种常见写法：

```go
package server
```

同包测试。优点是可以直接访问同 package 里的 exported 和 unexported 名字。阶段四刚开始建议用这个，心智更简单。

另一种是：

```go
package server_test
```

外部包测试。它像真正的外部使用者一样 import `server` package，只能访问 exported 名字，例如 `server.NewMux()`。这种写法更适合测试公开 API，但对新手来说会多一层 import 心智。

阶段四先用：

```go
package server
```

等你熟悉 package 边界后，再学习什么时候改成外部包测试。

这段测试的心智是：

```text
造一个假 request
造一个假 response writer
把 request 丢给 handler
检查 response
```

### `httptest.NewRecorder` 是什么

`httptest.NewRecorder()` 会创建一个 `ResponseRecorder`。它假装自己是 `http.ResponseWriter`，但不会真的通过网络返回 response，而是把 handler 写入的内容记录下来，方便测试断言。

这和阶段二的 `httptest.NewServer` 不一样：

| 工具 | 是否启动真实 HTTP server | 适合测试什么 |
| --- | --- | --- |
| `httptest.NewServer` | 会 | 测 HTTP client 真的发请求 |
| `httptest.NewRecorder` | 不会 | 测 handler 如何写 response |

## 配置先怎么做

服务入口里第一步会调用 `config.Load()`，所以在看 `cmd/hn-agent/main.go` 之前，先把配置 package 的第一版形状定下来。

第一版配置可以非常简单：只从环境变量读取监听地址。

```go
// hn-agent/internal/config/config.go
package config

import "os"

// Config 保存服务启动时需要的运行配置。
// 现阶段先只有 Addr，后面可以继续加入 timeout、HN API 地址等字段。
type Config struct {
	Addr string
}

// Load 从环境变量读取配置。
// 如果没有设置 HN_AGENT_ADDR，就使用默认地址 :8080。
func Load() Config {
	addr := os.Getenv("HN_AGENT_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{
		Addr: addr,
	}
}
```

### 为什么 `Load()` 返回 `Config`，不是 `*Config`

这里先写：

```go
func Load() Config
```

而不是：

```go
func Load() *Config
```

原因是第一版 `Config` 很小，而且读取完成后基本只读：

- `Config` 现在只有一个 `Addr string` 字段，复制成本很低。
- `Load()` 总能返回一个可用配置；没有环境变量时会使用默认值 `:8080`，所以不需要用 `nil` 表达“没有配置”。
- 返回值可以减少共享可变状态。调用方拿到的是一份配置值，后面不容易被别的地方偷偷改掉。
- 调用方不需要写 `if cfg == nil { ... }` 这类防御判断，学习阶段心智更简单。

可以把它理解成：`Config` 是启动时读出来的一张小纸条，内容确定后直接交给 `main.go` 使用即可。

什么时候更适合返回指针？

- struct 很大，频繁复制成本高。
- 函数需要返回“可能没有值”的结果，例如 `(*Config, error)`。
- 多个地方必须共享并修改同一个对象。
- 类型的方法需要修改自身状态。

如果以后配置读取可能失败，例如要求必须设置 `OPENAI_API_KEY`，第一选择通常也不是马上返回指针，而是把签名改成：

```go
func Load() (Config, error)
```

成功时返回 `cfg, nil`，失败时返回 `Config{}, err`。这和阶段二、阶段三学过的 Go 显式错误处理保持一致。

这里先不引入 `.env`、配置文件或第三方库。阶段四目标是理解服务生命周期，不是做复杂配置系统。

后面 `main.go` 里的：

```go
cfg := config.Load()
```

就可以读成：从 `internal/config` package 读取服务启动配置，然后用 `cfg.Addr` 作为 HTTP server 的监听地址。

### 为什么必须叫 `Load`，不是 `load`

Go 用首字母大小写控制 package 外可见性：

- `Load`：大写开头，是 exported，其他 package 可以调用。
- `load`：小写开头，是 unexported，只能在 `config` package 内部调用。

`hn-agent/internal/config/config.go` 属于：

```go
package config
```

而服务入口 `hn-agent/cmd/hn-agent/main.go` 属于：

```go
package main
```

它们是两个不同 package。所以 `main.go` 里只能通过 import 后调用大写的 `config.Load()`：

```go
import "github.com/aaron/go-with-ai/hn-agent/internal/config"

func main() {
	cfg := config.Load()
	_ = cfg
}
```

如果你写成：

```go
func load() Config
```

编辑器提示 `function "load" is unused` 是合理的：因为它是小写 unexported 函数，当前 `config` package 内部没有任何代码调用它，外部 package 也不能调用它。

注意：`cfg := config.Load()` 只是读取配置。如果后续没有使用 `cfg`，Go 还会继续提示 `cfg` unused。完整入口里会用 `cfg.Addr`：

```go
srv := &http.Server{
	Addr: cfg.Addr,
}
```

### 新手容易错：`load`、`config`、`cfg` 三个 warning 连在一起

阶段四第一次拆 package 时，常见错误不是一个点，而是一串：

```text
function "load" is unused
undefined: config
declared and not used: cfg
```

它通常来自这样的半成品代码。

`hn-agent/internal/config/config.go`：

```go
package config

type Config struct {
	Addr string
}

func load() Config {
	return Config{Addr: ":8080"}
}
```

`hn-agent/cmd/hn-agent/main.go`：

```go
package main

func main() {
	cfg := config
}
```

这里有三个问题：

- `load` 小写开头，只能在 `config` package 内部使用；`main` package 不能调用它。
- `main.go` 没有 import `internal/config`，所以 `config` 这个名字不存在。
- `cfg` 定义后没有被真正使用，所以 Go 会继续报 unused。

修正方向是：

```go
// hn-agent/internal/config/config.go
package config

import "os"

type Config struct {
	Addr string
}

func Load() Config {
	addr := os.Getenv("HN_AGENT_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{Addr: addr}
}
```

```go
// hn-agent/cmd/hn-agent/main.go
package main

import (
	"net/http"

	"github.com/aaron/go-with-ai/hn-agent/internal/config"
)

func main() {
	cfg := config.Load()

	srv := &http.Server{
		Addr: cfg.Addr,
	}

	_ = srv
}
```

排查顺序可以固定成三步：

1. 先看函数是否需要跨 package 调用。需要就用大写 `Load`。
2. 再看调用方有没有 import 对应 package。
3. 最后看变量是否真的被后续代码使用。

## 服务入口怎么写

`cmd/hn-agent/main.go` 不要写太多业务逻辑。它应该像启动脚本：

1. 读取配置。
2. 创建 mux。
3. 创建 `http.Server`。
4. 启动 server。
5. 等待停止信号。
6. graceful shutdown。

```go
// hn-agent/cmd/hn-agent/main.go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaron/go-with-ai/hn-agent/internal/config"
	"github.com/aaron/go-with-ai/hn-agent/internal/server"
)

func main() {
	cfg := config.Load()

	// slog 是 Go 标准库里的结构化日志包。
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := server.NewMux()

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// NotifyContext 会在收到中断信号时取消 ctx。
	// Ctrl+C 通常会触发 os.Interrupt。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server starting", "addr", cfg.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 阻塞等待停止信号。
	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
```

### `main()` 逐段导读

这段入口代码第一次看会很密。先不要追求一次记住所有语法，可以按“启动一个长期运行服务”的顺序理解。

整体流程是：

```text
读取配置
创建 logger
创建 mux
创建 http.Server
监听系统退出信号
在 goroutine 里启动 server
main goroutine 等待退出信号
收到信号后创建 shutdown timeout
调用 srv.Shutdown
记录服务已停止
```

#### 读取配置

```go
cfg := config.Load()
```

这行调用 `internal/config` package 里的 `Load()`，拿到启动配置。

`cfg` 的类型是 `config.Config`，里面现在只有：

```go
Addr string
```

后面 `http.Server` 会用：

```go
Addr: cfg.Addr
```

决定监听哪个地址，例如默认的 `:8080`。

#### 创建 logger

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
```

这行可以先分成两层看：

```go
slog.NewTextHandler(os.Stdout, nil)
```

创建一个 text 格式的日志 handler，把日志写到 `os.Stdout`，也就是标准输出。第二个参数是配置选项，这里先传 `nil`，表示使用默认选项。

外层：

```go
slog.New(...)
```

用这个 handler 创建一个 logger。后面就可以写：

```go
logger.Info("server starting", "addr", cfg.Addr)
```

如果希望日志时间改成固定格式和上海时区，可以先看后面的 [`log/slog` 自定义时间格式](#自定义日志时间格式和上海时区) 小节，再把这里替换成：

```go
logger := newLogger()
```

#### 创建 mux

```go
mux := server.NewMux()
```

这行调用 `internal/server` package 里的 `NewMux()`，拿到路由表。

可以把它理解成：

```text
/healthz -> healthz handler
```

`main.go` 不关心每个 handler 的内部细节，只负责把 mux 放进 HTTP server。

#### 创建 `http.Server`

```go
srv := &http.Server{
	Addr:              cfg.Addr,
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
}
```

这里有几个新语法。

`&http.Server{...}` 表示创建一个 `http.Server` struct，并取它的地址。`srv` 的类型是：

```go
*http.Server
```

也就是 server 指针。后面调用：

```go
srv.ListenAndServe()
srv.Shutdown(shutdownCtx)
```

都是在使用这个 server 对象。

`Addr: cfg.Addr` 是 struct literal 的字段赋值，表示监听地址。

`Handler: mux` 表示所有请求先交给 mux，再由 mux 分发给具体 handler。

`ReadHeaderTimeout: 5 * time.Second` 表示读取 request header 最多等 5 秒，避免连接一直占着不发完整请求。`time.Second` 是一个时间单位，`5 * time.Second` 的类型是 `time.Duration`。

#### 监听退出信号

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

`signal.NotifyContext(...)` 会创建一个 context。这个 context 平时不做事，等进程收到指定信号时会被取消。

这里监听两个信号：

- `os.Interrupt`：通常来自键盘 `Ctrl+C`。
- `syscall.SIGTERM`：常见于系统、容器或部署平台要求进程停止。

它返回两个值：

```go
ctx, stop
```

- `ctx`：后面用 `<-ctx.Done()` 等待取消。
- `stop`：停止接收 signal，并释放相关资源。

`defer stop()` 表示 `main()` 结束前一定调用 `stop()`。

#### 在 goroutine 里启动 server

```go
go func() {
	logger.Info("server starting", "addr", cfg.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}()
```

`go func() { ... }()` 表示启动一个新的 goroutine 执行这个匿名函数。

为什么要新开 goroutine？因为：

```go
srv.ListenAndServe()
```

会阻塞。它会一直监听端口、处理请求，直到 server 关闭或出错。如果放在 main goroutine 里，后面的 `<-ctx.Done()` 和 `srv.Shutdown(...)` 就执行不到。

这里的错误判断也值得拆开：

```go
if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
```

`ListenAndServe()` 正常运行时会一直阻塞；当 server 被 `Shutdown` 关闭时，它通常会返回 `http.ErrServerClosed`。这个错误在 graceful shutdown 里是预期结果，不应该当成真正失败。

所以条件的意思是：

```text
如果 ListenAndServe 返回了错误，
并且这个错误不是“server 已正常关闭”，
才记录为启动失败并退出。
```

`os.Exit(1)` 表示用非 0 状态码退出程序，通常代表异常退出。

新手容易写反这个条件：

```go
if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
	logger.Error("server failed", "error", err)
	os.Exit(1)
}
```

这段少了 `!`，意思变成：

```text
如果错误正好是“server 已正常关闭”，就当成失败。
```

这会导致按 `Ctrl+C` 后出现类似输出：

```text
level=INFO msg="server stopped"
level=ERROR msg="server failed" error="http: Server closed"
exit status 1
```

这不是我们想要的结果。正确写法是：

```go
if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	logger.Error("server failed", "error", err)
	os.Exit(1)
}
```

这里的 `!errors.Is(...)` 可以读成：“不是正常关闭错误，才算真正失败”。

#### 等待退出信号

```go
<-ctx.Done()
stop()
```

`<-ctx.Done()` 会让 main goroutine 停在这里，直到 `ctx` 被取消。

收到 Ctrl+C 或 SIGTERM 后，`ctx.Done()` 会关闭，这行才会继续往下走。

后面的 `stop()` 是主动停止 signal 通知。虽然前面已经有 `defer stop()`，这里再调用一次也可以；它让收到第一个退出信号后的清理更明确。可以先理解成：信号已经收到了，接下来进入关闭流程，不再继续监听这一轮 signal。

#### 创建 shutdown timeout

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

这里创建一个新的 context，专门给 shutdown 用。

为什么不用前面的 `ctx`？因为前面的 `ctx` 已经被取消了。`Shutdown` 需要一个还没取消、但有超时时间的 context，表示：

```text
最多给 server 5 秒时间优雅关闭。
```

`cancel` 用来释放 timer 等资源，所以写：

```go
defer cancel()
```

#### 执行 graceful shutdown

```go
if err := srv.Shutdown(shutdownCtx); err != nil {
	logger.Error("server shutdown failed", "error", err)
	os.Exit(1)
}
```

`srv.Shutdown(shutdownCtx)` 会：

```text
停止接收新请求
等待正在处理的请求结束
如果超过 shutdownCtx 的 5 秒限制，就返回 error
```

如果 shutdown 失败，就记录错误并异常退出。

#### 记录停止完成

```go
logger.Info("server stopped")
```

能执行到这里，说明 server 已经完成关闭流程。

阶段四先把这个入口写通、跑通、能 `Ctrl+C` 退出，比一开始做复杂 API 更重要。

如果 `context` 这一块还比较陌生，先读独立篇章：[Go context 入门：服务如何知道“该停了”](/notes/stage-4-service/go-context)。那里会把 `context.Background()`、`signal.NotifyContext`、`context.WithTimeout` 和 `<-ctx.Done()` 用图拆开讲。

### 为什么 `ListenAndServe` 要放 goroutine

`ListenAndServe` 会阻塞当前 goroutine。也就是说，如果直接写：

```go
srv.ListenAndServe()
```

后面的“等待信号”和“优雅关闭”代码就没机会执行。

所以阶段四常见结构是：

```text
一个 goroutine 负责运行 server
main goroutine 负责等待退出信号
收到信号后调用 Shutdown
```

### `<-ctx.Done()` 是什么

这行非常有 Go 特色：

```go
<-ctx.Done()
```

它可以先读成：

```text
一直等，直到 ctx 发出“已经取消”的信号。
```

拆开看：

```go
ctx.Done()
```

返回的是一个 channel，类型可以先理解成：

```go
<-chan struct{}
```

也就是“只能接收信号的 channel”。它不是用来传业务数据的，而是用来通知：这个 context 已经结束了。

前面的 `<-` 是 receive 操作：

```go
<-someChannel
```

表示从 channel 里接收一个值。如果现在还没有值、channel 也还没关闭，当前 goroutine 就会阻塞等待。

在这里，`signal.NotifyContext(...)` 创建的 `ctx` 会在收到 `os.Interrupt` 或 `syscall.SIGTERM` 时被取消。取消发生后，`ctx.Done()` 这个 channel 会被关闭，于是：

```go
<-ctx.Done()
```

不再阻塞，代码继续往下执行，进入 `srv.Shutdown(...)`。

所以这行不是“读取一个有用的值”，而是“用 channel 等一个停止信号”。可以把它理解成阶段四服务入口里的暂停点：

```text
server 已经启动
main goroutine 在这里等待
用户按 Ctrl+C
ctx 被取消
ctx.Done() 关闭
<-ctx.Done() 放行
开始 graceful shutdown
```

如果写成：

```go
ctx.Done()
```

只是拿到 channel，但不会等待；程序会继续往下跑。必须加上 `<-`，才表示“在这里等到它完成”。

## `log/slog` 是什么

`log/slog` 是 Go 标准库的结构化日志包。结构化日志不是只拼一个字符串，而是写 key/value：

```go
logger.Info("server starting", "addr", cfg.Addr)
logger.Error("server failed", "error", err)
```

这样后续机器或日志平台更容易检索：

```text
msg="server starting" addr=:8080
```

阶段四先学会两个方法：

- `logger.Info(...)`：记录正常状态。
- `logger.Error(...)`：记录错误状态。

### 自定义日志时间格式和上海时区

默认 `slog.NewTextHandler(os.Stdout, nil)` 会输出类似：

```text
time=2026-06-29T16:18:29.075+08:00 level=INFO msg="server starting" addr=:8080
```

如果你希望把 `time` 改成指定格式，并固定使用上海时区，可以给 `TextHandler` 传 `slog.HandlerOptions`，用 `ReplaceAttr` 改写 `time` 字段。

推荐先把 logger 创建逻辑抽成一个 helper：

```go
// hn-agent/cmd/hn-agent/main.go
func newLogger() *slog.Logger {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 正常情况下 Go 能加载 Asia/Shanghai。
		// 这里保留兜底，避免极端环境缺少时区数据时 logger 创建失败。
		shanghai = time.FixedZone("Asia/Shanghai", 8*60*60)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				t := a.Value.Time().In(shanghai)
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}

			return a
		},
	})

	return slog.New(handler)
}
```

然后 `main()` 里写：

```go
logger := newLogger()
```

输出会变成类似：

```text
time="2026-06-29 16:18:29" level=INFO msg="server starting" addr=:8080
```

注意这里时间会带引号，是因为 `TextHandler` 看到字符串里有空格，会自动加引号。这是正常的。

#### 为什么格式是 `2006-01-02 15:04:05`

Go 的时间格式不是写：

```text
YYYY-MM-DD HH:mm:ss
```

而是用一个固定参考时间：

```text
Mon Jan 2 15:04:05 MST 2006
```

所以常见格式要写成：

```go
t.Format("2006-01-02 15:04:05")
```

可以先死记这个对应关系：

| 想表达 | Go layout |
| --- | --- |
| 年 | `2006` |
| 月 | `01` |
| 日 | `02` |
| 时，24 小时 | `15` |
| 分 | `04` |
| 秒 | `05` |

如果希望保留时区偏移，可以用：

```go
t.Format("2006-01-02 15:04:05 -07:00")
```

输出类似：

```text
time="2026-06-29 16:18:29 +08:00"
```

阶段四先用 `2006-01-02 15:04:05` 足够。等后面真的需要接日志平台，再决定是否改成 RFC3339 或 JSON 日志。

## 开发期热加载怎么做

Go 标准工具链没有像 Vite、Next.js 那样内置开发服务器热加载。常见做法是用一个 watcher 工具监听 `.go` 文件变化，然后自动：

```text
停止旧进程
重新编译
启动新进程
```

这更准确地叫 live reload，或者“自动重编译 + 自动重启”。它不是生产环境的 hot deploy。

阶段四推荐先用 `air`。安装：

```bash
go install github.com/air-verse/air@latest
```

如果安装后提示 `air: command not found`，通常是 Go 的 bin 目录没有加入 `PATH`。可以先检查：

```bash
go env GOPATH
```

常见路径是：

```text
$HOME/go/bin
```

macOS 上也可以用 Homebrew：

```bash
brew install go-air
```

### 给 `hn-agent` 配置 Air

因为我们的服务入口在：

```text
hn-agent/cmd/hn-agent
```

不是 module 根目录，所以建议在 `hn-agent/.air.toml` 写清楚 build 命令：

```toml
# hn-agent/.air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/hn-agent ./cmd/hn-agent"
bin = "./tmp/hn-agent"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
```

然后在 `hn-agent/` 目录运行：

```bash
air
```

之后你修改 `.go` 文件并保存，`air` 会自动重新 build 并重启服务。你不需要每次手动：

```bash
go run ./cmd/hn-agent
```

### 和 `go run` 的关系

`go run ./cmd/hn-agent` 适合确认一次启动是否正常。

`air` 适合开发期持续修改：

```text
写代码 -> 保存 -> air 自动重启 -> curl/browser 验证
```

如果遇到奇怪问题，先回到最朴素的命令排查：

```bash
go test ./...
go run ./cmd/hn-agent
```

这样可以判断问题来自你的 Go 代码，还是来自 watcher 配置。

## Graceful shutdown 是什么

普通关闭像“拔电源”。graceful shutdown 更像“通知打烊”：

```text
收到停止信号
-> 不再接新请求
-> 给正在处理的请求一点时间完成
-> 超时还没完成就退出
```

Go 标准库里对应的是：

```go
srv.Shutdown(ctx)
```

这里传入的 `ctx` 通常带 timeout：

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

意思是：最多等 5 秒。

## 阶段四推荐步骤

1. 新建 `internal/server`，实现 `NewMux()` 和 `/healthz`。
2. 写 `server_test.go`，用 `httptest.NewRecorder` 测 `/healthz`。
3. 新建 `internal/config`，先支持 `HN_AGENT_ADDR`，默认 `:8080`。
4. 新建 `cmd/hn-agent/main.go`，启动 `http.Server`。
5. 用 `curl http://localhost:8080/healthz` 做手动验收。
6. 加 graceful shutdown。
7. 再考虑 `/api/v1/digests/latest`。

## 阶段四测试方向

第一批测试先覆盖：

- `GET /healthz` 返回 200。
- `POST /healthz` 返回 405。
- response body 是 `ok\n`。
- 默认配置地址是 `:8080`。
- 设置 `HN_AGENT_ADDR` 后配置能被覆盖。

示例命令：

```bash
cd hn-agent
go test ./internal/server
go test ./internal/config
go test ./...
```

## 阶段四验收

包级验收：

```bash
cd hn-agent
go test ./...
```

手动运行：

```bash
go run ./cmd/hn-agent
```

另开一个终端访问：

```bash
curl -i http://localhost:8080/healthz
```

预期看到：

```text
HTTP/1.1 200 OK
...
ok
```

如果设置了自定义地址：

```bash
HN_AGENT_ADDR=:9090 go run ./cmd/hn-agent
curl -i http://localhost:9090/healthz
```

## 暂时不要做什么

阶段四第一步先不要急着做：

- 不要马上引入 Web 框架。
- 不要马上做数据库。
- 不要马上做复杂定时任务。
- 不要把 AI 摘要逻辑塞进 handler。

先把服务骨架打稳。服务期最怕一开始把 HTTP、配置、日志、任务调度、AI、存储都混在一起，后面每个测试都会很难写。
