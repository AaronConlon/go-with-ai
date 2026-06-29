# 阶段 2：网络期完整代码

阶段 2 新增 HN API client，重点是 HTTP 请求、JSON 解析、timeout 和 `httptest`。

阶段 1 的文件继续保留。本页保存阶段 2 新增或重点修改的完整文件。代码块使用学习版注释，帮助理解 HTTP client 和测试服务器的工作方式。

## 文件清单

```text
hn-agent/internal/hn/client.go
hn-agent/internal/hn/client_test.go
```

## `context` 包的作用

`context` 可以先理解成：在一条调用链里传递“控制信号”的工具。

阶段二最重要的信号有两个：

- 取消：调用方已经不想等了，正在执行的请求应该尽快停下来。
- 超时：超过指定时间还没完成，就自动取消。

它不是用来传业务数据的，也不是替代函数参数。更准确地说，`context.Context` 通常放在函数第一个参数，用来告诉函数：

```text
这次工作属于哪个请求？
有没有超时？
有没有被取消？
```

例如：

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
```

这行的意思是：创建一个 HTTP request，并把它和 `ctx` 绑定。如果 `ctx` 被取消，HTTP 请求也会被取消。

测试里常见：

```go
ids, err := client.TopStories(context.Background())
```

`context.Background()` 是一个空的根 context。它没有超时，也不会主动取消，常用于 `main`、测试或最外层入口。后续如果需要 timeout，可以从它派生：

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

`defer cancel()` 是清理动作，用来释放这个 timeout context 关联的资源。

## HN Client

```go
// hn-agent/internal/hn/client.go
package hn

import (
	// context 用来在调用链里传递取消、超时和请求范围信息。
	"context"
	// encoding/json 用来把 JSON response body 解码成 Go 值。
	"encoding/json"
	// fmt 用来创建格式化字符串和错误信息。
	"fmt"
	// net/http 是 Go 标准库 HTTP client/server 包。
	"net/http"
	// time 用来表达超时时间，例如 10 * time.Second。
	"time"
)

// Item 对应 Hacker News API 的 item 结构。
// 字段必须大写开头，encoding/json 才能给它们赋值。
type Item struct {
	// struct tag `json:"id"` 表示 JSON 里的 id 字段映射到 Go 的 ID 字段。
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Descendants int    `json:"descendants"`
}

// Client 封装访问 HN API 需要的信息。
// 这样生产代码和测试代码都可以复用同一套方法。
type Client struct {
	// BaseURL 是 API 基础地址。
	// 测试里会把它替换成 httptest server 的地址。
	BaseURL string
	// HTTP 是真正发请求的标准库客户端。
	HTTP *http.Client
}

// NewClient 创建生产环境使用的默认 HN client。
func NewClient() *Client {
	return &Client{
		BaseURL: "https://hacker-news.firebaseio.com/v0",
		HTTP: &http.Client{
			// Timeout 是整个 HTTP 请求允许的最大耗时。
			// 外部 I/O 不设置 timeout，程序可能一直卡住。
			Timeout: 10 * time.Second,
		},
	}
}

// TopStories 拉取 HN top stories 的 id 列表。
// 返回 ([]int64, error) 是 Go 常见写法：成功时返回数据和 nil error，失败时返回 nil 和 error。
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
	// NewRequestWithContext 创建带 context 的请求。
	// ctx 取消或超时时，请求也会跟着取消，避免外部请求无限等待。
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/topstories.json", nil)
	if err != nil {
		return nil, err
	}

	// Do 真正发起 HTTP 请求。
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	// Body 是网络资源，用完必须关闭。
	// defer 表示当前函数结束前执行 Close。
	defer resp.Body.Close()

	// 只把 2xx 当作成功。
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// ids 先声明成 []int64。
	// Decode 会把 JSON 数组写入 ids。
	var ids []int64
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

// Item 根据 id 拉取单条 HN item 详情。
func (c *Client) Item(ctx context.Context, id int64) (Item, error) {
	// Sprintf 用 id 拼出 /item/{id}.json 的完整地址。
	url := fmt.Sprintf("%s/item/%d.json", c.BaseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// Item 是 struct，不能返回 nil。
		// 失败时返回 Item{} 这个零值，再配合 error 表示失败。
		return Item{}, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Item{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return Item{}, err
	}

	return item, nil
}
```

## HN Client 测试

```go
// hn-agent/internal/hn/client_test.go
package hn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTopStories(t *testing.T) {
	// httptest.NewServer 启动一个本地临时 HTTP server。
	// 测试用它模拟真实 HN API，避免测试依赖外网。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Path 是 client 请求过来的路径。
		// 如果路径错了，说明生产代码拼 URL 的逻辑有问题。
		if r.URL.Path != "/topstories.json" {
			t.Fatalf("expected path /topstories.json, got %s", r.URL.Path)
		}

		// 告诉 client：响应内容是 JSON。
		w.Header().Set("Content-Type", "application/json")
		// 返回 200，表示请求成功。
		w.WriteHeader(http.StatusOK)
		// 写入一个模拟的 JSON 数组。
		_, _ = w.Write([]byte(`[1, 2, 3]`))
	}))
	// 测试结束时关闭 server，释放端口和 goroutine。
	defer server.Close()

	// 测试 client 指向本地 server，而不是真实 HN API。
	client := &Client{
		BaseURL: server.URL,
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}

	ids, err := client.TopStories(context.Background())
	if err != nil {
		t.Fatalf("TopStories returned error: %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("expected 3 ids, got %d", len(ids))
	}
}
```

## 验收命令

```bash
cd hn-agent
go test ./internal/hn
```
