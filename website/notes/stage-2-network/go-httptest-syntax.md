# Go `httptest` 使用说明

阶段二会写 HN API client。它会访问外部 HTTP API，但测试时不应该每次真的请求 Hacker News。

这时就用 Go 标准库里的 `net/http/httptest`。

## 一句话

`httptest` 可以在测试里临时启动一个假的 HTTP server，让你的 HTTP client 以为自己正在请求真实 API。

它的用途是：

```text
用假的 HTTP 响应，测试真的 HTTP client 代码。
```

## 应该写在哪个文件

生产代码写在：

```text
hn-agent/internal/hn/client.go
```

测试代码写在：

```text
hn-agent/internal/hn/client_test.go
```

`httptest` 只应该出现在测试文件里：

```go
import "net/http/httptest"
```

不要在 `client.go` 里 import `httptest`。因为 `client.go` 是真实业务代码，未来线上运行时不需要假的测试 server。

## 文件分工

### `client.go`

负责真实能力：

- 定义 `Client`。
- 定义 `Item`。
- 发 HTTP 请求。
- 解析 JSON。
- 返回结果和错误。

例如：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
	// 真实 HTTP 请求逻辑
}
```

### `client_test.go`

负责验证：

- 当 API 返回正常 JSON 时，能解析出结果。
- 当 API 返回 500 时，能返回错误。
- 当 API 返回错误 JSON 时，能返回错误。
- 请求路径是否符合预期。

例如：

```go
func TestTopStories(t *testing.T) {
	// 用 httptest 创建假 server
	// 用 Client 请求这个假 server
	// 检查结果是否正确
}
```

## 为什么需要 `BaseURL`

为了让测试能替换真实 API 地址，`Client` 里应该有：

```go
type Client struct {
	BaseURL string
	HTTP    *http.Client
}
```

真实运行时：

```go
client := NewClient("https://hacker-news.firebaseio.com/v0")
```

测试时：

```go
client := NewClient(server.URL)
```

也就是说，同一份 `TopStories` 代码：

- 真实运行时请求 HN。
- 测试时请求 `httptest` 临时 server。

这就是为什么 `BaseURL` 不要写死在 `TopStories` 里。

## 最小测试示例

假设生产代码里有：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
	// 请求 c.BaseURL + "/topstories.json"
}
```

那么测试文件 `internal/hn/client_test.go` 可以这样写：

```go
package hn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTopStories(t *testing.T) {
	// httptest.NewServer 会启动一个临时 HTTP server。
	// 这个 server 只在测试期间存在。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r 是请求。
		// 这里检查 client 是否请求了正确路径。
		if r.URL.Path != "/topstories.json" {
			t.Fatalf("path = %s, want /topstories.json", r.URL.Path)
		}

		// w 是响应写入器。
		// 这里告诉 client：响应内容是 JSON。
		w.Header().Set("Content-Type", "application/json")

		// 写 HTTP 状态码 200。
		w.WriteHeader(http.StatusOK)

		// 写响应 body。
		// 这模拟了 HN topstories API 返回 story id 列表。
		_, _ = w.Write([]byte(`[1,2,3]`))
	}))
	defer server.Close()

	// 用测试 server 的 URL 创建 Client。
	// 这样 TopStories 请求的不是 HN，而是上面的假 server。
	client := NewClient(server.URL)

	ids, err := client.TopStories(context.Background())
	if err != nil {
		t.Fatalf("TopStories() error = %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("len(ids) = %d, want 3", len(ids))
	}
}
```

## 这段测试如何运行

在 `hn-agent/` 目录内执行：

```bash
go test ./internal/hn
```

如果想看到测试名：

```bash
go test -v ./internal/hn
```

如果只跑这个测试：

```bash
go test -v ./internal/hn -run TestTopStories
```

## `httptest.NewServer` 做了什么

```go
server := httptest.NewServer(...)
```

它会启动一个本地 HTTP server，并给你一个 URL：

```go
server.URL
```

这个 URL 可能长这样：

```text
http://127.0.0.1:54321
```

测试期间，你的 client 请求：

```text
server.URL + "/topstories.json"
```

就会进入你写的 handler：

```go
http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	...
})
```

## `defer server.Close()` 为什么需要

```go
defer server.Close()
```

意思是：

> 当前测试函数结束时，关闭这个临时 server。

如果不关闭，测试 server 可能占着端口和资源。

`defer` 是 Go 的延迟执行语法。它会在当前函数返回前执行。

## `http.HandlerFunc` 是什么

```go
http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	...
})
```

可以先理解成：

> 把一个函数变成 HTTP handler。

这个函数会在 server 收到请求时运行。

参数：

- `w http.ResponseWriter`：用来写响应。
- `r *http.Request`：表示收到的请求。

## 为什么测试里要检查路径

```go
if r.URL.Path != "/topstories.json" {
	t.Fatalf("path = %s, want /topstories.json", r.URL.Path)
}
```

这能确认生产代码真的请求了正确 endpoint。

如果你的 `TopStories` 写错成：

```text
/topstory.json
```

测试会立刻失败。

这比手动看代码靠谱。

## `_, _ = w.Write([]byte(`[1,2,3]`))` 是什么

这一行：

```go
_, _ = w.Write([]byte(`[1,2,3]`))
```

是在给假 server 写响应 body。

它可以读成：

> 把字符串 `[1,2,3]` 转成 bytes，写到 HTTP response body 里；`Write` 返回的两个结果我都暂时忽略。

拆开看有四层。

### 第一层：`[1,2,3]` 是模拟 JSON

```go
`[1,2,3]`
```

这是一个字符串，内容是 JSON 数组。

它模拟 HN 的 `topstories.json` 返回：

```json
[1,2,3]
```

这里用反引号是 Go 的原始字符串写法。也可以写成双引号：

```go
"[1,2,3]"
```

这个例子里两种都可以。

### 第二层：`[]byte(...)` 是类型转换

```go
[]byte(`[1,2,3]`)
```

`w.Write` 需要的参数类型是：

```go
[]byte
```

也就是 byte 切片。

但我们手里的是字符串：

```go
`[1,2,3]`
```

所以要把 string 转成 `[]byte`：

```go
[]byte(`[1,2,3]`)
```

可以先理解成：

```text
把这段文本变成可以写入 HTTP 响应体的字节数据。
```

### 第三层：`w.Write(...)` 写响应体

```go
w.Write([]byte(`[1,2,3]`))
```

`w` 是：

```go
http.ResponseWriter
```

它负责写 HTTP response。

前面已经写了响应头和状态码：

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
```

这一行就是写 body：

```go
w.Write([]byte(`[1,2,3]`))
```

最终 client 收到的响应大概是：

```http
HTTP/1.1 200 OK
Content-Type: application/json

[1,2,3]
```

### 第四层：`_, _ =` 忽略两个返回值

`w.Write` 的函数签名大概是：

```go
Write([]byte) (int, error)
```

它返回两个值：

- 第一个 `int`：实际写入了多少个 bytes。
- 第二个 `error`：写入过程有没有出错。

所以严格写法可以是：

```go
n, err := w.Write([]byte(`[1,2,3]`))
```

但在这个简单测试里，我们暂时不关心写入了几个 byte，也不处理写入错误，所以用 `_` 忽略：

```go
_, _ = w.Write([]byte(`[1,2,3]`))
```

`_` 是 Go 的空白标识符，意思是：

```text
这个返回值我知道存在，但我不使用。
```

因为 Go 不允许声明了变量却不用，所以不能写：

```go
n, err := w.Write([]byte(`[1,2,3]`))
```

然后完全不使用 `n` 和 `err`。这样会编译失败。

### 更严谨的写法

如果想更严谨，可以处理 error：

```go
if _, err := w.Write([]byte(`[1,2,3]`)); err != nil {
	t.Fatalf("write response: %v", err)
}
```

这里：

- `_` 仍然忽略写入 byte 数。
- `err` 保留下来。
- 如果写入失败，就让测试失败。

阶段二刚开始可以先用：

```go
_, _ = w.Write([]byte(`[1,2,3]`))
```

理解成本最低。等熟悉后，再逐步改成处理 error 的写法。

## 测试错误状态码

除了成功情况，还应该测试失败情况。

例如 HN 返回 500：

```go
func TestTopStoriesStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	_, err := client.TopStories(context.Background())
	if err == nil {
		t.Fatal("TopStories() error = nil, want error")
	}
}
```

这个测试验证：

> 如果服务端返回 500，`TopStories` 不应该假装成功，而应该返回 error。

## 测试错误 JSON

也可以测试 JSON 格式错误：

```go
func TestTopStoriesInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	_, err := client.TopStories(context.Background())
	if err == nil {
		t.Fatal("TopStories() error = nil, want error")
	}
}
```

这个测试验证：

> 如果响应不是合法 JSON，`TopStories` 应该返回 error。

## 和真实请求的区别

真实请求：

```text
Client -> Hacker News API
```

测试请求：

```text
Client -> httptest.NewServer
```

核心是：测试仍然经过你的真实 `Client` 代码，只是把外部服务替换成假的。

## 现阶段先记住

- `httptest` 只写在 `_test.go` 文件里。
- 它用来启动假的 HTTP server。
- 测试时把 `Client.BaseURL` 指向 `server.URL`。
- 这样可以稳定测试 HTTP client，不依赖真实网络。
- 至少测试三类情况：成功响应、错误状态码、错误 JSON。
