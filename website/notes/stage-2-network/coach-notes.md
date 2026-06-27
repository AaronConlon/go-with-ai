# 阶段二教练笔记

阶段二开始真正接触外部世界：HTTP、JSON、timeout、status code 和错误处理。对 Go 新手来说，这一阶段最重要的不是“把 API 调通”，而是理解 Go 如何把外部 I/O 写成可取消、可测试、可返回错误的代码。

## 本阶段注释策略

阶段二的文档代码继续使用高密度注释。重点解释：

- `context.Context` 为什么几乎每个外部 I/O 都要带。
- `http.Client` 和 `http.Request` 分别负责什么。
- `json.Decoder` 如何把响应体转成 Go struct。
- 为什么错误要返回给调用者，而不是只 `fmt.Println`。
- 为什么测试要用 `httptest`，而不是每次真的请求 HN。

## 目标文件

阶段二建议在 `hn-agent/` 内逐步创建或扩展：

```text
internal/hn/client.go
internal/hn/client_test.go
cmd/hnctl/main.go
```

## 第一步：定义 HN item 模型

建议先在 `internal/hn/client.go` 中放一个接近 HN API 的结构体：

```go
package hn

// Item 表示 Hacker News API 返回的一条 item。
// HN 里的 item 可以是 story、comment、job、poll 等。
// 阶段二我们主要关心 type=story 的 item。
type Item struct {
	// json tag 告诉 encoding/json：
	// JSON 里的 "id" 字段应该映射到 Go 的 ID 字段。
	ID int64 `json:"id"`

	// Type 表示 item 类型，例如 "story"。
	Type string `json:"type"`

	// By 是作者用户名。
	By string `json:"by"`

	// Time 是 Unix 时间戳。
	Time int64 `json:"time"`

	// Title 是 story 标题。
	Title string `json:"title"`

	// URL 是 story 指向的外部链接。
	URL string `json:"url"`

	// Score 是 story 分数。
	Score int `json:"score"`

	// Descendants 是评论数量。
	Descendants int `json:"descendants"`
}
```

### 为什么字段名是大写

Go 里大写开头的标识符是 exported，可以被其他 package 访问。`encoding/json` 也只能给可导出的字段赋值，所以 `ID`、`Title`、`URL` 都要大写。

## 第二步：定义 Client

```go
package hn

import (
	"net/http"
	"time"
)

// Client 封装访问 Hacker News API 所需的状态。
// 这里先只放 BaseURL 和 HTTP client。
type Client struct {
	// BaseURL 是 HN API 根地址。
	// 测试时可以把它替换成 httptest server 的地址。
	BaseURL string

	// HTTP 是真正发请求的标准库客户端。
	// 把它放进结构体，是为了后续可以设置 timeout，也方便测试替换。
	HTTP *http.Client
}

// NewClient 创建一个默认 HN client。
// Go 里常用 NewXxx 作为构造函数命名。
func NewClient() *Client {
	return &Client{
		BaseURL: "https://hacker-news.firebaseio.com/v0",
		HTTP: &http.Client{
			// Timeout 是整个请求的最大耗时。
			// 外部 I/O 不设置 timeout，程序可能一直卡住。
			Timeout: 10 * time.Second,
		},
	}
}
```

## 第三步：实现 `TopStories`

```go
package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TopStories 拉取 HN top story 的 id 列表。
// ctx 用来传递超时、取消和请求范围信息。
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
	// 创建一个带 context 的 HTTP request。
	// 如果 ctx 被取消，请求也会跟着取消。
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/topstories.json", nil)
	if err != nil {
		// 创建 request 失败也要返回 error。
		return nil, err
	}

	// 真正发起 HTTP 请求。
	resp, err := c.HTTP.Do(req)
	if err != nil {
		// 网络错误、超时、DNS 问题等会走这里。
		return nil, err
	}
	defer resp.Body.Close()

	// HTTP 请求成功不代表业务成功。
	// 例如 500、404 也会返回 resp，所以要检查 status code。
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("hn topstories status: %s", resp.Status)
	}

	var ids []int64
	// json.NewDecoder(resp.Body).Decode(&ids)
	// 会从响应体读取 JSON，并写入 ids。
	// &ids 是传指针，因为 Decode 需要修改 ids 的值。
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}
```

## 第四步：实现 `Item`

```go
package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Item 根据 id 拉取单条 HN item。
func (c *Client) Item(ctx context.Context, id int64) (Item, error) {
	url := fmt.Sprintf("%s/item/%d.json", c.BaseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Item{}, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Item{}, fmt.Errorf("hn item %d status: %s", id, resp.Status)
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return Item{}, err
	}

	return item, nil
}
```

### 为什么错误时返回 `Item{}`

Go 函数如果返回 `(Item, error)`，出错时也必须返回一个 `Item` 值。`Item{}` 是 `Item` 的零值，表示“没有有效结果”。调用方应该先检查 `err`，再使用 `item`。

## 第五步：用 `httptest` 测试

```go
package hn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTopStories(t *testing.T) {
	// httptest.NewServer 会启动一个只在测试期间存在的本地 HTTP server。
	// 这样测试不依赖真实 HN API，速度更快，也更稳定。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求路径是否符合预期。
		if r.URL.Path != "/topstories.json" {
			t.Fatalf("path = %s, want /topstories.json", r.URL.Path)
		}

		// 返回一段假的 HN JSON。
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[1,2,3]`))
	}))
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:   server.Client(),
	}

	ids, err := client.TopStories(context.Background())
	if err != nil {
		t.Fatalf("TopStories() error = %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("len(ids) = %d, want 3", len(ids))
	}
}
```

### `context.Background()` 是什么

`context.Background()` 是一个空的根 context。测试或 main 入口里经常用它作为起点。后续如果需要 timeout，可以从它派生：

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

## 阶段二验收

在 `hn-agent/` 内执行：

```bash
go test ./internal/hn
go run ./cmd/hnctl top --limit=10
```

如果还没接 CLI，先只验收：

```bash
go test ./internal/hn
```

