# 阶段 3：并发期完整代码

阶段 3 新增有限并发批量抓取：`FetchItems`、`errgroup`、semaphore 和批量抓取测试。

阶段 1、阶段 2 的文件继续保留。本页保存阶段 3 新增或重点修改的完整文件。代码块使用学习版注释，重点解释 goroutine、channel、semaphore、错误传播和测试策略。

## 文件清单

```text
hn-agent/internal/hn/fetch_batch.go
hn-agent/internal/hn/fetch_batch_test.go
```

## 批量抓取实现

### `context` 在并发里做什么

阶段三里，`context` 不只是给单个 HTTP 请求设置取消，还会把“一组并发任务”串起来。

这一行：

```go
g, ctx := errgroup.WithContext(ctx)
```

会基于旧的 `ctx` 派生出一个新的 `ctx`。只要某个 goroutine 返回 error，`errgroup` 就会取消这个新的 `ctx`。后面的 HTTP 请求都使用这个新 `ctx`：

```go
item, err := c.Item(ctx, id)
```

效果是：

```text
任意一个 item 请求失败
-> goroutine 返回 error
-> errgroup 取消 ctx
-> 其他还没完成的 HTTP 请求感知取消
-> FetchItems 返回 error
```

所以阶段三可以把 `context` 记成：并发任务之间传递“该停了”的信号。

```go
// hn-agent/internal/hn/fetch_batch.go
package hn

import (
	// context 负责传递取消和超时信号。
	// 在阶段三里，它让一组 goroutine 可以一起被取消。
	"context"
	"fmt"

	// errgroup 不在标准库里，来自 golang.org/x/sync。
	// 它用来管理一组会返回 error 的 goroutine。
	"golang.org/x/sync/errgroup"
)

// FetchItems 根据一组 id 批量拉取 HN item。
// concurrency 控制最多同时发起多少个请求。
func (c *Client) FetchItems(ctx context.Context, ids []int64, concurrency int) ([]Item, error) {
	// 并发数必须是正数。
	// 如果是 0 或负数，后面的 semaphore 就没有意义，甚至可能导致任务一直阻塞。
	if concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be positive")
	}

	// 预先创建固定长度的结果切片。
	// 每个 goroutine 按自己的下标写入 items[i]，这样返回顺序和输入 ids 顺序一致。
	items := make([]Item, len(ids))

	// errgroup.WithContext 返回两个值：
	// g：用来启动 goroutine、等待结束、收集错误。
	// ctx：派生 context；任意 goroutine 返回 error 时，这个 ctx 会被取消。
	g, ctx := errgroup.WithContext(ctx)

	// sem 是 semaphore，中文可以理解成“并发许可证”。
	// 这里用带缓冲 channel 实现，容量就是最大并发数。
	sem := make(chan struct{}, concurrency)

	// 遍历所有待抓取 id。
	// 当前学习契约以 Go 1.27 为基准，for range 循环变量每轮独立，
	// 所以这里不需要旧教程里的 i, id := i, id。
	for i, id := range ids {
		// g.Go 会启动一个 goroutine。
		// 这个函数返回 error，最后由 g.Wait 统一收集。
		g.Go(func() error {
			// 获取一个并发许可证。
			// 如果 sem 已满，这里会阻塞，直到别的 goroutine 归还许可证。
			sem <- struct{}{}
			defer func() {
				// 当前 goroutine 结束时归还许可证。
				// defer 能保证成功或失败路径都会执行这一步。
				<-sem
			}()

			// 真正拉取单条 item。
			// 这里传入 errgroup 返回的新 ctx，让错误取消能传播到 HTTP 请求。
			item, err := c.Item(ctx, id)
			if err != nil {
				// 返回非 nil error 会让 g.Wait 最终返回错误，
				// 同时 errgroup 派生的 ctx 会被取消。
				return err
			}

			// 每个 goroutine 写不同下标，因此不需要 append，也不需要 mutex。
			items[i] = item
			return nil
		})
	}

	// Wait 会等待所有 g.Go 启动的 goroutine 结束。
	// 如果其中任意一个返回 error，这里会拿到一个 error。
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// 所有 item 都成功时，返回完整列表和 nil error。
	return items, nil
}
```

## 批量抓取测试

### `fmt.Fprint` 在测试里做什么

`fmt.Fprint(w, "...")` 可以读成：把字符串写入 `w`。

这里的 `w` 是 `http.ResponseWriter`，所以它不是打印到终端，而是在构造测试服务器的 HTTP response body。函数会返回“写入了多少字节”和“写入是否出错”，所以示例里写成：

```go
_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
```

两个 `_` 表示暂时忽略这两个返回值。

你也可以直接写：

```go
fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
```

Go 允许把函数调用作为单独一行语句，返回值会被直接丢弃。所以 `_, _ =` 不是语法必需。

这里写 `_, _ =` 是为了让意图更清楚：这行函数确实有返回值，但在这个测试场景里我们暂时不关心“写入了多少字节”和“写 response 是否出错”。很多 linter 也更喜欢这种显式忽略返回值的写法。

如果这个写入错误会影响测试判断，更严谨的写法是检查 error：

```go
if _, err := fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`); err != nil {
	t.Fatalf("write response: %v", err)
}
```

### 为什么 `t.Fatal('expected error, got nil')` 不行

Go 里的字符串必须用双引号或反引号：

```go
t.Fatal("expected error, got nil")
```

单引号表示 `rune`，也就是一个 Unicode 字符：

```go
letter := 'a'
```

所以这行会报错：

```go
t.Fatal('expected error, got nil') // illegal rune literal
```

因为单引号里有很多字符，已经不是一个 `rune`。

`gofmt` 不能自动修复这个问题。它只负责格式化已经能解析的 Go 代码；`illegal rune literal` 是语法解析阶段就失败了，`gofmt` 没有合法 AST 可以格式化。并且单引号和双引号代表不同类型，自动替换可能改变语义。

### `%d` 和 `%#v` 有什么区别

`t.Fatalf` 里的第一个参数是格式化字符串。里面的 `%d`、`%#v` 叫格式化动词。

```go
t.Fatalf("expected 3 items, got %d", len(items))
```

这里用 `%d`，因为 `len(items)` 返回的是整数，`%d` 表示按十进制整数输出。

```go
t.Fatalf("expected nil items on error, got %#v", items)
```

这里用 `%#v`，因为 `items` 是 `[]Item`，属于 slice。我们想在测试失败时看清它到底是什么：

```text
[]hn.Item(nil)
[]hn.Item{}
[]hn.Item{hn.Item{ID:101, ...}}
```

`%v` 是默认格式，够日常；`%#v` 更偏调试，会尽量按 Go 语法形式把值打印出来，适合 slice、map、struct 这类复合值。

```go
// hn-agent/internal/hn/fetch_batch_test.go
package hn

import (
	"context"
	// fmt 在实现里用来创建 error，在测试里用 fmt.Fprint 写 response body。
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchItemsInvalidConcurrency(t *testing.T) {
	// 这个测试只验证输入校验，不需要真的访问网络。
	client := NewClient()

	// concurrency 为 0，应该直接返回错误。
	items, err := client.FetchItems(context.Background(), []int64{101}, 0)
	if err == nil {
		// 测试失败信息是字符串，所以必须用双引号。
		// 单引号在 Go 里表示 rune，例如 'a'，不能包多字符字符串。
		t.Fatal("expected error, got nil")
	}
	if items != nil {
		// %#v 会尽量用 Go 语法形式打印 items。
		// 这里 items 是 []Item，比起 %d 这种整数格式，%#v 更适合调试复合值。
		t.Fatalf("expected nil items, got %#v", items)
	}
}

func TestFetchItemsSuccessPreservesOrder(t *testing.T) {
	// 用本地测试服务器模拟多个 item 接口。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 根据请求路径返回不同 item。
		// 这样可以确认 FetchItems 确实按 id 发起了多个请求。
		switch r.URL.Path {
		case "/item/101.json":
			// fmt.Fprint 会把内容写入第一个参数指定的 writer。
			// 这里的 w 是 http.ResponseWriter，所以不是打印到终端，
			// 而是把这段 JSON 写进测试 HTTP 响应体。
			// 它返回两个值：写入的字节数和 error。
			// Go 允许直接调用 fmt.Fprint(...) 并丢弃返回值；
			// 这里写 _, _ = 是显式告诉读者和 linter：这两个返回值是有意忽略的。
			_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
		case "/item/102.json":
			_, _ = fmt.Fprint(w, `{"id":102,"type":"story","title":"second"}`)
		case "/item/103.json":
			_, _ = fmt.Fprint(w, `{"id":103,"type":"story","title":"third"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// 测试 client 指向本地 server。
	// server.Client() 会返回适合访问这个测试服务器的 http.Client。
	client := &Client{
		BaseURL: server.URL,
		HTTP:    server.Client(),
	}

	// concurrency 为 2，表示最多同时抓两个 item。
	items, err := client.FetchItems(context.Background(), []int64{101, 102, 103}, 2)
	if err != nil {
		t.Fatalf("FetchItems returned error: %v", err)
	}

	if len(items) != 3 {
		// len(items) 是 int，所以用 %d 按十进制整数输出。
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	// 阶段三非常重要的一点：
	// 并发完成顺序可能不同，但返回切片顺序必须和输入 ids 一致。
	wantIDs := []int64{101, 102, 103}
	for i, wantID := range wantIDs {
		if items[i].ID != wantID {
			t.Fatalf("items[%d].ID = %d, want %d", i, items[i].ID, wantID)
		}
	}
}

func TestFetchItemsReturnsErrorWhenOneItemFails(t *testing.T) {
	// 这个测试模拟“其中一个 item 请求失败”。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/item/102.json" {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// 这里同样用 fmt.Fprint 把 JSON 写入 HTTP response body。
		_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"ok"}`)
	}))
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:    server.Client(),
	}

	// 阶段三先采用“任意一个失败，整体失败”的策略。
	items, err := client.FetchItems(context.Background(), []int64{101, 102}, 2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if items != nil {
		// 如果失败时 items 不是 nil，%#v 能帮助我们看清它到底装了什么。
		t.Fatalf("expected nil items on error, got %#v", items)
	}
}
```

## 验收命令

```bash
cd hn-agent
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
go test -race ./internal/hn
```
