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

这里有两个概念叠在一起：Go 的 package 可见性，以及 `encoding/json` 如何给 struct 字段赋值。

## Go 的大小写可见性规则

Go 没有 `public`、`private` 这种关键字。它用首字母大小写控制一个名字能不能被别的 package 使用。

| 写法 | 名称 | 其他 package 能访问吗 |
| --- | --- | --- |
| `Story` | exported | 可以 |
| `Item` | exported | 可以 |
| `ID` | exported | 可以 |
| `Title` | exported | 可以 |
| `story` | unexported | 不可以 |
| `item` | unexported | 不可以 |
| `id` | unexported | 不可以 |
| `title` | unexported | 不可以 |

例如在 `internal/hn` 包里：

```go
package hn

type Story struct {
	ID    int64
	Title string
}
```

其他 package 可以这样用：

```go
story := hn.Story{
	ID:    1,
	Title: "Learning Go",
}
```

但如果写成小写：

```go
package hn

type Story struct {
	id    int64
	title string
}
```

其他 package 不能直接访问：

```go
story := hn.Story{
	id:    1,              // 不允许：id 没有导出
	title: "Learning Go",  // 不允许：title 没有导出
}
```

## `encoding/json` 为什么也需要大写字段

`encoding/json` 是标准库里的另一个 package。它不属于你的 `hn` package。

当你写：

```go
var item Item
err := json.NewDecoder(resp.Body).Decode(&item)
```

`encoding/json` 需要把 JSON 里的字段写进 `item`：

```json
{
  "id": 1,
  "title": "Learning Go",
  "url": "https://example.com"
}
```

如果 struct 是这样：

```go
type Item struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
```

`encoding/json` 可以写入，因为 `ID`、`Title`、`URL` 都是大写开头，是 exported 字段。

如果 struct 写成这样：

```go
type Item struct {
	id    int64  `json:"id"`
	title string `json:"title"`
	url   string `json:"url"`
}
```

`encoding/json` 看得到 tag，但不能给这些字段赋值，因为 `id`、`title`、`url` 是 unexported 字段。

结果通常会是：JSON 解码不报错，但字段保持零值。

```go
// id    仍然是 0
// title 仍然是 ""
// url   仍然是 ""
```

这会很隐蔽，所以写 API 响应 struct 时，字段一般都用大写开头。

## 那 `json:"id"` 是干什么的

Go 字段要大写：

```go
ID
Title
URL
```

但 JSON 字段通常是小写或 snake_case：

```json
{
  "id": 1,
  "title": "Learning Go",
  "url": "https://example.com"
}
```

所以用 json tag 做映射：

```go
ID    int64  `json:"id"`
Title string `json:"title"`
URL   string `json:"url"`
```

可以理解成：

```text
Go 代码里叫 ID，JSON 里叫 id。
Go 代码里叫 Title，JSON 里叫 title。
Go 代码里叫 URL，JSON 里叫 url。
```

## `json:"url"` 这种写法到底是什么

```go
URL string `json:"url"`
```

这整行可以分成三部分：

```text
URL           字段名
string        字段类型
`json:"url"`  struct tag
```

`json:"url"` 不是字段值，也不是注释。它叫 struct tag，可以理解成“贴在字段上的说明标签”。

Go 语言本身不会主动使用这个 tag，但标准库里的 `encoding/json` 会读取它。

### 为什么用反引号

Go 里有两种常见字符串写法：

```go
"普通字符串"
`原始字符串`
```

struct tag 通常用反引号，因为里面经常包含双引号：

```go
`json:"url"`
```

如果不用反引号，就需要写很多转义字符，不好读。

### `json:"url"` 的意思

```go
URL string `json:"url"`
```

意思是：

> 当使用 `encoding/json` 编码或解码时，这个 Go 字段对应 JSON 里的 `url` 字段。

解码时：

```json
{
  "url": "https://example.com"
}
```

会填进：

```go
item.URL
```

编码时：

```go
Item{URL: "https://example.com"}
```

会变成：

```json
{
  "url": "https://example.com"
}
```

### 为什么不是 `URL string json:"url"`

因为 tag 是一个整体的字符串元信息，Go 语法要求它写在字段类型后面，并且用反引号包起来：

```go
字段名 字段类型 `tag`
```

所以完整格式是：

```go
URL string `json:"url"`
```

### 常见 tag 写法

忽略空值：

```go
URL string `json:"url,omitempty"`
```

意思是：如果 `URL` 是空字符串，编码成 JSON 时可以省略这个字段。

忽略字段：

```go
Secret string `json:"-"`
```

意思是：这个字段不参与 JSON 编码和解码。

JSON 字段名和 Go 字段名完全不同：

```go
Descendants int `json:"descendants"`
```

意思是：Go 里叫 `Descendants`，JSON 里叫 `descendants`。

### 先记住一句话

`json:"url"` 是给 `encoding/json` 看的字段映射说明：Go 里字段叫 `URL`，JSON 里字段叫 `url`。

## 为什么是 `ID` 和 `URL`，不是 `Id` 和 `Url`

这是 Go 的命名惯例。常见首字母缩写一般保持全大写：

```go
ID
URL
HTTP
JSON
API
```

所以更推荐：

```go
ID  int64
URL string
```

而不是：

```go
Id  int64
Url string
```

## 先记住一句话

只要一个 struct 要给别的 package 用，或者要被 `encoding/json` 自动填充字段，字段名就应该大写开头。

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

### 怎么判断 `Item` 是否真的有效

不要靠 `Item{}` 判断成功失败，要靠 `err`。

::: danger 先看 `err`
Go 里 `(结果, error)` 的判断顺序是：**先看 `err`，再使用结果**。

`err != nil` 时，前面的结果值不可信；`err == nil` 时，前面的结果才表示成功结果。
:::

调用方应该这样写：

```go
item, err := client.Item(ctx, id)
if err != nil {
	return err
}

// 到这里才说明 Item 调用成功，可以使用 item。
fmt.Println(item.Title)
```

这条规则很重要：

```text
err != nil：结果值不可信，不要使用 item。
err == nil：函数成功，item 才是有效结果。
```

为什么不靠 `item == Item{}` 判断？

因为 `Item{}` 只是 `Item` 的零值，不是状态标志。理论上，一个正常返回的值也可能有一些字段是零值，例如 `Score` 是 `0`、`URL` 是空字符串。业务是否有效应该由函数的错误返回或领域校验来表达，而不是靠猜测 struct 是否等于零值。

如果未来确实需要表达“没有找到，但这不是系统错误”，可以选择更明确的 API 设计：

```go
func (c *Client) Item(ctx context.Context, id int64) (*Item, error)
```

这种设计里，`nil, nil` 可以表示“没找到且不是错误”。也可以设计成：

```go
func (c *Client) Item(ctx context.Context, id int64) (Item, bool, error)
```

其中 `bool` 专门表示有没有找到。但阶段二先用 `(Item, error)`，并约定：有错误就返回 `Item{}, err`；没错误就返回 `item, nil`。

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

这里的：

```go
ids, err := client.TopStories(context.Background())
```

可以读成：

```text
调用 TopStories，成功时拿到 ids，失败时拿到 err。
```

这和 JavaScript 里常见的 `try/catch` 不一样。Go 不把这里的错误藏到外层 `catch`，而是把错误作为第二个返回值交给当前调用点。

所以测试里紧接着写：

```go
if err != nil {
	t.Fatalf("TopStories() error = %v", err)
}
```

意思是：如果拉取 top stories 失败，这个测试立刻失败；只有 `err == nil` 时，后面才继续检查 `ids`。

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

这里的 `top --limit=10` 表示阶段二希望达到的 CLI 目标：

- `top`：调用 HN client 拉取 top stories。
- `--limit=10`：最多展示 10 个 story id，后续再扩展成展示 story 详情。

如果当前 `cmd/hnctl` 还没有接入 `Client.TopStories`，这条命令不会真的使用 `limit`。此时先用 `go test ./internal/hn` 验收 HTTP client，等 client 稳定后再把它接到 CLI。

如果还没接 CLI，先只验收：

```bash
go test ./internal/hn
```
