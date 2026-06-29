# 阶段三教练笔记

阶段三开始学习 Go 的并发。对 Go 新手来说，不要先背 channel 的全部语法，而要先建立一个工程直觉：

> 并发不是“同时开很多任务”这么简单，而是要控制数量、知道何时取消、知道错误怎么返回。

## 配图导读：从 JS 心智切到 Go 并发心智

先用三张小黑图建立直觉。它们不是 API 图解，而是帮助你把 JS 的异步心智迁移到 Go 的并发心智。

![JS 事件循环与 Go goroutine 的心智差异](/images/goroutine-js-mindset/01-js-event-loop-vs-go-goroutine.png)

JS 更像“排队回调”：任务在一个调度模型里等待机会回来执行。Go 的 `go` 关键字则可以启动很多轻量任务，但重点不是“能开很多”，而是后面必须管理它们。

![errgroup.WithContext 管理一组 goroutine](/images/goroutine-js-mindset/02-errgroup-context-cancel.png)

`errgroup.WithContext(ctx)` 可以把一组 goroutine 绑在一起：用 `g.Go` 开任务，用 `g.Wait` 收口；任意任务返回 error，派生出来的 `ctx` 会取消，其他任务就有机会停下来。

![make 预先创建结果槽位](/images/goroutine-js-mindset/03-make-items-slots.png)

`make([]Item, len(ids))` 像先摆好一排结果槽位。每个 goroutine 抓到自己的 item 后写回 `items[i]`，这样顺序稳定，也避免多个 goroutine 同时 `append` 同一个 slice。

## 本阶段注释策略

阶段三文档代码继续使用高密度注释。重点解释：

- goroutine 是什么。
- 为什么不能无限制启动 goroutine。
- channel 在这里为什么能当 semaphore。
- `errgroup.WithContext` 如何收集错误和传播取消。
- 为什么旧教程会写 `i, id := i, id`，以及为什么 Go 1.25 中不需要。

## 阶段三术语表

这些词建议保留英文原文，同时建立中文理解。不要为了翻译而把术语翻得很绕。

| 英文 | 常见中文 | 更适合本项目的理解 |
| --- | --- | --- |
| goroutine | 不强行翻译；可说 Go 协程 / 轻量线程 | Go 里的轻量并发任务，用 `go` 启动，重点是启动后要能等待、取消和收错 |
| concurrency | 并发 | 同一段时间内处理多个任务，重点是任务组织和资源边界，不等于真的同时执行 |
| parallelism | 并行 | 多个任务在多个 CPU 核上真的同时执行 |
| channel | 通道 | goroutine 之间传值、同步或做并发控制的管道 |
| semaphore | 信号量 / 并发许可证 | 控制最多同时进入某段代码的任务数量 |
| worker pool | 工作池 | 固定数量的 worker 从任务队列里取活干 |
| errgroup | 错误组 / goroutine 任务组 | 管理一组会返回 error 的 goroutine：启动、等待、收错、取消 |
| context | 上下文 | 在调用链里传递取消、超时和请求范围信息 |
| cancellation | 取消 | 通知任务“别继续做了，尽快收尾” |
| timeout | 超时 | 到时间后自动取消任务 |
| fan-out | 扇出 | 一个任务拆成多个并发子任务，例如多个 story id 同时抓 |
| fan-in | 扇入 | 多个并发子任务的结果汇总回来 |

### `concurrency` 和 `parallelism` 的区别

这两个词经常被混用，但阶段三最好分开理解：

```text
concurrency：有多个任务在同一段时间内被管理。
parallelism：有多个任务在同一时刻真的同时执行。
```

举个接近项目的例子：

```text
FetchItems 同时管理 20 个 id 的抓取，这是 concurrency。
其中某几个 HTTP 请求可能真的同时在不同 CPU 核上处理，这是 parallelism。
```

Go 的 `goroutine` 让你很容易写出 concurrency，但这不代表可以无限制启动任务。阶段三的重点是：

```text
有限并发 + 错误传播 + 取消传播 + 结果顺序稳定
```

### `goroutine` 怎么理解

`goroutine` 可以先理解成：

```text
Go runtime 管理的轻量并发任务。
```

它不是操作系统线程本身。Go runtime 会把很多 goroutine 调度到少量操作系统线程上执行。

所以 `go func() { ... }()` 的重点不是“开了一个线程”，而是：

```text
启动了一个可能和当前函数并发执行的任务。
```

一旦启动，就要考虑：

- 谁等它结束？
- 它失败了怎么办？
- 它什么时候应该取消？
- 它会不会太多，压垮外部 API？

这就是为什么阶段三不用裸写一堆 goroutine，而是引入 `errgroup` 和 semaphore。

## 目标文件

建议继续在 `hn-agent/internal/hn/` 中扩展：

```text
internal/hn/fetch_batch.go
internal/hn/fetch_batch_test.go
```

## 第一版目标函数

```go
FetchItems(ctx context.Context, ids []int64, concurrency int) ([]Item, error)
```

含义：

- `ctx`：控制整批任务的取消和超时。
- `ids`：要抓取的 HN item id。
- `concurrency`：最大并发数。
- `[]Item`：抓到的 item 列表。
- `error`：整批任务是否失败。

## 推荐实现

```go
// hn-agent/internal/hn/fetch_batch.go
package hn

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// FetchItems 根据一组 id 批量拉取 HN item。
// concurrency 控制最多同时发起多少个请求。
func (c *Client) FetchItems(ctx context.Context, ids []int64, concurrency int) ([]Item, error) {
	// 并发数必须是正数。
	// 如果调用方传 0 或负数，直接返回错误，避免后面 channel 卡住。
	if concurrency <= 0 {
		return nil, fmt.Errorf("concurrency must be positive")
	}

	// items 预先创建固定长度。
	// 每个 goroutine 按下标写入自己的位置，这样最后顺序和 ids 一致。
	items := make([]Item, len(ids))

	// errgroup 用来管理一组 goroutine。
	// WithContext 会返回一个新的 ctx：
	// 只要其中一个 goroutine 返回 error，这个 ctx 就会被取消。
	g, ctx := errgroup.WithContext(ctx)

	// sem 是 semaphore，也就是并发“许可证”。
	// channel 容量等于 concurrency，表示最多允许 concurrency 个任务同时执行。
	sem := make(chan struct{}, concurrency)

	for i, id := range ids {
		// g.Go 启动一个 goroutine，并让 errgroup 记录它的返回 error。
		g.Go(func() error {
			// 获取一个并发许可证。
			// 如果 sem 已满，这里会阻塞，直到别的任务释放许可证。
			sem <- struct{}{}
			defer func() {
				// 函数结束时释放许可证。
				<-sem
			}()

			// 真正拉取单条 item。
			// 这里传入 errgroup 返回的 ctx，这样任意任务失败后，其他请求也能被取消。
			item, err := c.Item(ctx, id)
			if err != nil {
				return err
			}

			// 每个 goroutine 写入不同下标，所以这里不需要 mutex。
			items[i] = item
			return nil
		})
	}

	// Wait 会等待所有 goroutine 结束。
	// 如果任意 goroutine 返回 error，Wait 会返回其中一个 error。
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return items, nil
}
```

## `i, id := i, id` 为什么现在不需要

你可能会在旧教程里看到这种写法：

```go
for i, id := range ids {
	i, id := i, id

	g.Go(func() error {
		item, err := c.Item(ctx, id)
		if err != nil {
			return err
		}
		items[i] = item
		return nil
	})
}
```

它的历史作用是：给当前循环轮次重新创建一份 `i` 和 `id`，避免 goroutine 闭包捕获到被后续循环复用的变量。

但当前项目使用：

```go
go 1.25.1
```

Go 1.22 之后，`for range` 循环变量已经是每轮独立的。也就是说，在本项目里可以直接写：

```go
for i, id := range ids {
	g.Go(func() error {
		item, err := c.Item(ctx, id)
		if err != nil {
			return err
		}
		items[i] = item
		return nil
	})
}
```

如果你继续写：

```go
i, id := i, id
```

linter 可能会提示：

```text
copying variable is unneeded
```

这个提示的意思是：当前 Go 版本已经不需要你手动复制循环变量。

::: tip 阶段三先记住
当前学习契约以 Go 1.27 为基准，`for range` 里的 `i` 和 `id` 每一轮都是独立变量，不需要再写 `i, id := i, id`。

如果以后维护 Go 1.21 或更早的老项目，再留意旧的循环变量捕获问题。
:::

## `sem <- struct{}{}` 和 `<-sem` 是什么

这段代码：

```go
sem <- struct{}{}
defer func() {
	// 函数结束时释放许可证。
	<-sem
}()
```

要和前面的 channel 创建放在一起看：

```go
sem := make(chan struct{}, concurrency)
```

`sem` 是一个带缓冲的 channel。缓冲区大小是 `concurrency`。

可以把它想象成一个装许可证的小盒子：

```text
concurrency = 3

sem 最多能放 3 张许可证。
```

### `struct{}{}` 是什么

`struct{}{}` 是一个空 struct 值。

它不携带任何业务数据，可以先读成：

```text
一个空标记。
```

这里我们不关心传递什么内容，只关心 channel 里“占了几个位置”。所以用 `struct{}{}` 很合适，因为它表达：

```text
我要占一个许可证位置，但不需要传具体数据。
```

### `sem <- struct{}{}`：获取许可证

这一行：

```go
sem <- struct{}{}
```

表示往 `sem` 这个 channel 里发送一个空 struct。

在 semaphore 心智里，可以读成：

```text
获取一个并发许可证。
```

如果 `sem` 还没满，发送成功，当前 goroutine 继续执行。

如果 `sem` 已经满了，说明当前已经有 `concurrency` 个任务正在执行。此时这一行会阻塞，当前 goroutine 会停在这里，直到别的任务归还许可证。

例如：

```text
concurrency = 3

任务 A 拿到许可证
任务 B 拿到许可证
任务 C 拿到许可证
任务 D 想拿许可证，但 sem 已满，所以等待
```

### `<-sem`：归还许可证

这一行：

```go
<-sem
```

表示从 `sem` 里接收一个值。

在 semaphore 心智里，可以读成：

```text
归还一个并发许可证。
```

归还后，前面被卡住的某个 goroutine 就有机会继续执行。

### 为什么用 `defer`

这段：

```go
defer func() {
	<-sem
}()
```

表示：

```text
等当前 goroutine 的函数快结束时，再归还许可证。
```

这样无论后面是成功：

```go
items[i] = item
return nil
```

还是失败：

```go
return err
```

都会先执行 defer，把许可证还回去。

如果忘记归还许可证，会出现很隐蔽的问题：

```text
某个任务拿了许可证但没还
-> sem 的空位越来越少
-> 后面的 goroutine 一直卡在 sem <- struct{}{}
-> 整批任务可能永远等不完
```

所以这个模式通常成对出现：

```go
sem <- struct{}{} // 获取许可证
defer func() {
	<-sem // 归还许可证
}()
```

::: tip 阶段三先记住
`sem <- struct{}{}` 是拿许可证；`<-sem` 是还许可证；`defer` 保证函数结束时一定还。
:::

### 更严谨的取消写法

当前推荐实现先用简单写法，方便理解：

```go
sem <- struct{}{}
```

它的特点是：如果 `sem` 满了，goroutine 会等许可证。

更严谨的版本会让“等待许可证”这件事也响应 `ctx` 取消：

```go
select {
case sem <- struct{}{}:
	defer func() {
		<-sem
	}()
case <-ctx.Done():
	return ctx.Err()
}
```

这段意思是：

```text
要么拿到许可证继续执行；
要么 ctx 已经取消，直接返回错误。
```

阶段三可以先掌握简单版本；等你理解 `select` 后，再升级到这个版本。

## `golang.org/x/sync/errgroup` 是什么

这一行 import：

```go
"golang.org/x/sync/errgroup"
```

表示引入 `errgroup` 包。

它不是 Go 标准库的一部分。它来自 Go 官方维护的扩展库：

```text
golang.org/x/sync
```

这意味着：只写 import 还不够，当前 Go module 还需要先记录这个依赖。

在 `hn-agent/` 内执行：

```bash
go get golang.org/x/sync/errgroup
go mod tidy
```

这有点像 JavaScript 里的：

```bash
npm install some-package
```

区别是 Go 会把依赖写进：

```text
go.mod
go.sum
```

其中：

- `go.mod`：记录当前 module 依赖了哪些模块。
- `go.sum`：记录依赖版本的校验信息。

执行后你可能会看到：

```go
module github.com/aaron/go-with-ai/hn-agent

go 1.25.1

require golang.org/x/sync v0.21.0 // indirect
```

这里可以这样读：

- `require`：当前 module 需要这个依赖。
- `golang.org/x/sync`：依赖的 module path。
- `v0.21.0`：依赖版本。
- `// indirect`：Go 当前把它标记为间接依赖。

`// indirect` 不是错误。它常见于两种情况：

1. 这个 module 是被别的依赖间接需要的。
2. 你刚执行了 `go get`，但当前源码还没有真正 import 并使用这个包。

等你在代码里实际写了：

```go
import "golang.org/x/sync/errgroup"
```

并且使用了：

```go
g, ctx := errgroup.WithContext(ctx)
```

再执行：

```bash
go mod tidy
```

Go 会重新扫描源码。如果它发现当前项目直接 import 了 `errgroup`，通常会把依赖整理成直接依赖，`// indirect` 会消失：

```go
require golang.org/x/sync v0.21.0
```

所以阶段三先不用纠结 `// indirect`。关键是：

```text
代码真实 import 了什么，以 go mod tidy 整理后的 go.mod 为准。
```

### `go mod tidy` 为什么会清理 `// indirect`

`go mod tidy` 会重新扫描当前 module 的源码 import。

在这个阶段，它会看到两件事：

1. `fetch_batch.go` 里有没有 `import "golang.org/x/sync/errgroup"`。
2. 代码里有没有实际使用 `errgroup.WithContext`、`g.Go`、`g.Wait`。

如果都存在，Go 就知道：

```text
golang.org/x/sync 是当前项目直接使用的 module。
```

于是它会把：

```go
require golang.org/x/sync v0.21.0 // indirect
```

整理成：

```go
require golang.org/x/sync v0.21.0
```

如果源码里没有实际 import，`go mod tidy` 反而可能把这条依赖删掉，因为它认为当前代码不需要它。

所以 `go mod tidy` 的核心作用是：

```text
让 go.mod / go.sum 和当前源码真正需要的依赖保持一致。
```

如果还没有执行 `go get`，编辑器的 Go language server 可能无法解析这个包，也就不会正确补全。

可以先理解成：

```text
errgroup = 用来管理一组会返回 error 的 goroutine 的工具。
```

如果不用 `errgroup`，你需要自己处理这些事情：

- 启动多个 goroutine。
- 等所有 goroutine 结束。
- 收集其中任意一个 goroutine 的错误。
- 某个 goroutine 失败后，通知其他 goroutine 尽快停止。

这些事情都能手写，但对新手来说容易写错。`errgroup` 把这个模式封装好了。

阶段三用它，是因为 `FetchItems` 要同时抓取多条 HN item，而且任意一条请求都可能失败。

## `errgroup.WithContext(ctx)` 做了什么

这一行：

```go
g, ctx := errgroup.WithContext(ctx)
```

注意大小写：是 `WithContext`，不是 `withContext`。

Go 的包级函数如果要被别的 package 使用，名字必须大写开头。`WithContext` 是 `errgroup` 包导出的函数；小写开头的 `withContext` 如果存在，也只能在 `errgroup` 包内部使用，外部代码访问不到。

返回两个值：

1. `g`：一个 group，用来启动 goroutine、等待 goroutine、收集 error。
2. `ctx`：一个新的 context，基于传入的旧 `ctx` 派生出来。

可以把它读成：

```text
创建一组并发任务，并给这组任务配一个联动取消的 context。
```

后面用：

```go
g.Go(func() error {
	// 做一个并发任务
	return nil
})
```

每调用一次 `g.Go`，就启动一个 goroutine。这个函数必须返回 `error`：

- 返回 `nil`：这个任务成功。
- 返回非 `nil` error：这个任务失败。

最后用：

```go
if err := g.Wait(); err != nil {
	return nil, err
}
```

`g.Wait()` 会做两件事：

1. 等所有通过 `g.Go` 启动的 goroutine 结束。
2. 如果其中任意一个 goroutine 返回了 error，`g.Wait()` 返回其中一个 error。

### 为什么 `WithContext` 还会返回新的 `ctx`

关键在这里：只要某个 goroutine 返回 error，`errgroup` 派生出来的 `ctx` 就会被取消。

所以我们在每个任务里调用：

```go
item, err := c.Item(ctx, id)
```

这里传的是 `errgroup.WithContext` 返回的新 `ctx`。

效果是：

```text
任意一个 item 抓取失败
-> g 记录这个 error
-> ctx 被取消
-> 其他还没完成的 HTTP 请求感知到取消
-> g.Wait() 返回 error
-> FetchItems 返回 nil, err
```

这就是“错误传播”和“取消传播”。

### 和 `sync.WaitGroup` 的区别

Go 标准库里也有 `sync.WaitGroup`，它只能等待一组 goroutine 结束。

但是 `sync.WaitGroup` 不会帮你：

- 收集 error。
- 返回第一个 error。
- 自动取消其他任务。

所以阶段三更适合用 `errgroup`。

| 工具 | 能等待 goroutine | 能收集 error | 能联动取消 |
| --- | --- | --- | --- |
| `sync.WaitGroup` | 可以 | 不可以，需要自己写 | 不可以，需要自己写 |
| `errgroup.WithContext` | 可以 | 可以 | 可以 |

::: tip 阶段三先记住
`errgroup.WithContext(ctx)` 是“批量并发 + 错误返回 + 取消传播”的组合工具。

用 `g.Go` 启动任务，用 `g.Wait` 等结果；任务里返回 error，整组任务就会失败并触发取消。
:::

### 为什么没有自动补全

常见原因有三个：

1. 还没有在 `hn-agent/` 内执行 `go get golang.org/x/sync/errgroup`。
2. 写成了 `errgroup.with...`，但 Go 的导出函数是大写 `WithContext`。
3. 编辑器的 Go language server 还没刷新。可以保存文件、执行 `go mod tidy`，必要时重启编辑器或重启 language server。

## 为什么预先创建固定长度的 `items`

这一行：

```go
items := make([]Item, len(ids))
```

意思是：创建一个长度和 `ids` 一样的 `[]Item`。

如果 `ids` 有 3 个：

```go
ids := []int64{101, 102, 103}
```

那么：

```go
items := make([]Item, len(ids))
```

会创建 3 个结果槽位：

```text
items[0]
items[1]
items[2]
```

它们一开始都是 `Item{}` 零值。后面每个 goroutine 成功抓到 item 后，按自己的下标写回去：

```go
items[i] = item
```

这样设计有三个目的。

第一，保持结果顺序稳定。

输入是：

```text
ids[0] -> 101
ids[1] -> 102
ids[2] -> 103
```

输出也保持：

```text
items[0] -> 101 对应的 item
items[1] -> 102 对应的 item
items[2] -> 103 对应的 item
```

即使 103 比 101 更早请求完成，也不会改变最终顺序。

第二，避免多个 goroutine 同时 `append` 同一个 slice。

不要在这个并发场景里这样写：

```go
var items []Item

g.Go(func() error {
	item, err := c.Item(ctx, id)
	if err != nil {
		return err
	}

	items = append(items, item) // 多个 goroutine 同时 append，会有数据竞争
	return nil
})
```

多个 goroutine 同时修改同一个 slice 的长度和底层数组，是不安全的。要么需要 mutex，要么就像这里一样提前分配好固定位置，让每个 goroutine 写自己的下标。

第三，结果数量一开始就明确。

`FetchItems` 的目标是：输入多少个 id，成功时就返回多少个 item。所以结果 slice 的长度天然就是 `len(ids)`。

::: tip 阶段三先记住
并发批量抓取时，先 `make([]Item, len(ids))`，再让每个 goroutine 写自己的 `items[i]`。

这比并发 `append` 更容易保持顺序，也更容易避免数据竞争。
:::

## `make` 在这里做了什么

`make` 是 Go 的内置函数，专门用来创建并初始化 slice、map 和 channel。

在阶段三这段代码里会看到两次：

```go
items := make([]Item, len(ids))
sem := make(chan struct{}, concurrency)
```

第一行创建结果列表：

```go
items := make([]Item, len(ids))
```

它创建的是 `[]Item`，长度是 `len(ids)`。创建完成后，`items[0]`、`items[1]` 这些位置已经存在，可以直接按下标写入。

第二行创建并发许可证 channel：

```go
sem := make(chan struct{}, concurrency)
```

它创建的是一个带缓冲的 channel，缓冲区大小是 `concurrency`。这个缓冲区就像最多能放多少张“通行证”，从而限制最多同时执行多少个请求。

::: tip 先记住
`make` 不是打印，也不是普通业务函数。它是 Go 创建 slice、map、channel 的内置工具。

看到 `make([]T, n)`，先读成：创建一个长度为 `n` 的 `[]T`。
:::

## `return nil, fmt.Errorf(...)` 是什么

这一段：

```go
if concurrency <= 0 {
	return nil, fmt.Errorf("concurrency must be positive")
}
```

可以拆成三层理解。

第一，Go 的 `if` 条件不需要括号。推荐写：

```go
if concurrency < 1 {
	...
}
```

不需要写成：

```go
if (concurrency < 1) {
	...
}
```

第二，`fmt.Errorf(...)` 不是打印。它会创建一个 `error` 值。只有 `fmt.Println`、`fmt.Printf` 这类函数才会直接输出到终端。

第三，`FetchItems` 的函数签名是：

```go
func (c *Client) FetchItems(ctx context.Context, ids []int64, concurrency int) ([]Item, error)
```

所以它必须返回两个值：

1. `[]Item`：成功时的 item 列表。
2. `error`：失败时的错误。

当 `concurrency` 小于 1 时，函数不能继续执行，因为后面要用这个数字控制并发。如果并发数是 0 或负数，逻辑没有意义，所以直接失败返回：

```go
return nil, fmt.Errorf("concurrency must be positive")
```

可以读成：

```text
没有可用的 item 列表，返回一个错误。
```

调用方通常这样接：

```go
items, err := client.FetchItems(ctx, ids, concurrency)
if err != nil {
	return err
}

// 只有 err == nil 时，才使用 items。
```

如果 `concurrency < 1`，执行顺序是：

1. 进入 `FetchItems`。
2. 判断 `concurrency < 1` 为 true。
3. `fmt.Errorf(...)` 创建一个 error。
4. `return nil, error` 立刻结束函数。
5. 后面的并发抓取逻辑不会执行。

如果 `concurrency >= 1`，这段 `if` 会跳过，函数继续往下执行。注意：Go 函数声明了返回值，就必须保证所有路径最终都有 `return`。所以正常路径最后也要有：

```go
return items, nil
```

`fmt.Errorf` 常用于需要格式化错误信息时：

```go
fmt.Errorf("item %d failed: %w", id, err)
```

如果只是固定错误文案，也可以用：

```go
errors.New("concurrency must be positive")
```

阶段三先用 `fmt.Errorf` 没问题，但错误字符串建议小写开头，不加句号或感叹号。

## `nil` 是什么类型的兜底值

看到这里，很容易产生一个误解：

```go
return nil, nil
```

是不是说明 `nil` 可以当所有类型的兜底值？

答案是：不是。

`nil` 只能用于“可以为空”的类型。现阶段先记住这几类：

| 类型 | 能不能是 `nil` | 例子 |
| --- | --- | --- |
| slice | 可以 | `[]Item`、`[]int64` |
| map | 可以 | `map[string]int` |
| channel | 可以 | `chan Item` |
| function | 可以 | `func() error` |
| pointer | 可以 | `*Client` |
| interface | 可以 | `error` |
| struct | 不可以 | `Item` |
| int / bool / string | 不可以 | `int64`、`bool`、`string` |

所以 `FetchItems` 可以写：

```go
return nil, fmt.Errorf("concurrency must be positive")
```

是因为它的返回值是：

```go
([]Item, error)
```

其中：

- `[]Item` 是 slice，可以是 `nil`。
- `error` 是 interface，可以是 `nil`，也可以装一个具体错误。

但如果函数返回的是：

```go
func (c *Client) Item(ctx context.Context, id int64) (Item, error)
```

失败时就不能写：

```go
return nil, err // 不可以：Item 不是可 nil 类型
```

要写：

```go
return Item{}, err
```

`Item{}` 是 `Item` struct 的零值。

再比如：

```go
func Count() (int, error)
```

失败时不能写：

```go
return nil, err // 不可以：int 不是可 nil 类型
```

要写：

```go
return 0, err
```

### 那 `return nil, nil` 合法吗

对 `FetchItems` 来说，语法上合法：

```go
return nil, nil
```

它表示：

```text
没有 item 列表，也没有错误。
```

但语义上要小心。大多数时候，成功路径应该返回真实结果：

```go
return items, nil
```

只有当“没有结果也是正常情况”时，才考虑返回空结果加 `nil` error。

对 slice 来说，空结果常见有两种写法：

```go
return nil, nil
```

或者：

```go
return []Item{}, nil
```

二者都表示没有元素，`len(items)` 都是 0。区别是：

- `nil` slice 表示没有分配底层数组。
- 空 slice 表示明确返回了一个空列表。

阶段三为了表达“成功抓取结束”，更推荐正常路径返回：

```go
return items, nil
```

失败路径返回：

```go
return nil, err
```

## 需要新增依赖

`errgroup` 来自 Go 官方扩展库，不在标准库里。第一次使用时，在 `hn-agent/` 内执行：

```bash
go get golang.org/x/sync/errgroup
go mod tidy
```

`go.mod` 会记录依赖，`go.sum` 会记录校验信息。这两个文件后续都应该提交。

## 为什么不用 `Promise.all` 心智模型

JavaScript 的 `Promise.all` 主要表达“等待一组异步任务完成”。Go 这里还要额外关心：

- goroutine 数量是否有上限。
- 某个任务失败后，其他任务是否应该取消。
- 结果顺序是否要和输入顺序一致。
- 错误如何返回给调用者。

所以阶段三更重要的心智模型是：

```text
任务列表 -> 有限并发 -> 可取消 -> 可返回错误 -> 结果顺序稳定
```

## 阶段三测试方向

阶段三的测试不要先盯着 CLI。CLI 是外层入口，真正需要先证明的是 `internal/hn` 里的 `FetchItems` 行为正确。

可以分成三层：

| 层级 | 目标 | 命令 |
| --- | --- | --- |
| package 单元测试 | 验证 `FetchItems` 的并发、顺序、错误语义 | `go test ./internal/hn` |
| 精准运行阶段三测试 | 只跑 `FetchItems` 相关测试，便于调试 | `go test -run TestFetchItems -v ./internal/hn` |
| CLI 集成验收 | 验证命令行参数已经接到业务逻辑 | `go run ./cmd/hnctl top --limit=20 --concurrency=5` |

### 必做测试

至少覆盖这四个场景：

1. `concurrency <= 0` 会返回错误。
2. 正常情况下能抓取多条 item。
3. 返回结果顺序和输入 `ids` 一致。
4. 其中一个 item 请求失败时，`FetchItems` 返回 error。

第一条是在验收输入校验：

```go
items, err := client.FetchItems(ctx, []int64{101}, 0)
```

预期是：

- `err != nil`
- `items` 不应该被当成正常结果继续使用

第二条是在验收正常路径：

```go
items, err := client.FetchItems(ctx, []int64{101, 102, 103}, 2)
```

预期是：

- `err == nil`
- `len(items) == 3`
- 每个 `Item` 都是测试服务器返回的数据

第三条非常重要。阶段三不是只要“都抓到了”就行，还要保证输出顺序稳定：

```text
输入 ids:     [101, 102, 103]
返回 items:  [101, 102, 103]
```

即使 `103` 的请求先完成，它也应该写回 `items[2]`，不能因为完成顺序不同就打乱返回结果。

第四条是在验收错误传播：

```text
101 成功
102 返回 500
103 成功或被取消
```

预期是 `FetchItems` 返回 error。阶段三先采用“有一个失败，整体失败”的策略，不做“部分成功结果 + 部分错误”的复杂设计。

### 常见断言写反

测试失败日志不能只看文案，还要看触发 `t.Fatal` 的 `if` 条件。

如果目标是“应该有 error”，断言应该写：

```go
items, err := client.FetchItems(ctx, ids, 0)
if err == nil {
	t.Fatal("expected error, got nil")
}
if items != nil {
	t.Fatalf("expected nil items, got %#v", items)
}
```

这里的意思是：

- `err == nil` 才失败，因为我们期望有 error。
- `items != nil` 才失败，因为失败路径应该返回 `nil` items。

如果写成：

```go
if err != nil {
	t.Fatal("expected error, got nil")
}
```

这就是断言条件和失败文案反了。它会在“真的拿到 error”时失败，但日志却说“got nil”，非常容易误导自己。

同理，测试“任意一个 item 失败时整体返回 error”也应该检查：

```go
if err == nil {
	t.Fatal("expected error, got nil")
}
```

### 进阶测试

这两类测试可以稍后补，但它们更能证明你真的理解了并发：

- 并发上限：测试服务器记录同时处理中的请求数，确认最大值没有超过 `concurrency`。
- 取消和 timeout：使用 `context.WithTimeout` 或手动 cancel，确认 `FetchItems` 能返回 context 相关错误。

并发上限测试的心智大概是：

```text
active++              // 请求开始
记录 maxActive
模拟慢请求
active--              // 请求结束
断言 maxActive <= concurrency
```

这类测试通常需要 `sync.Mutex` 或 `sync/atomic` 保护计数器，因为多个请求会同时修改同一个数字。

## 阶段三验收

阶段三验收也分两步。

第一步是必做的 package 验收。在 `hn-agent/` 内执行：

```bash
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
```

如果想额外检查并发读写风险，可以跑：

```bash
go test -race ./internal/hn
```

`-race` 会启用 Go 的 race detector，用来发现多个 goroutine 同时读写同一块数据导致的数据竞争。它不是每天都必须跑，但阶段三很值得跑一次。

第二步是 CLI 集成验收。只有当 `top` 命令已经真正把 `--limit` 和 `--concurrency` 传给业务逻辑后，再执行：

```bash
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

这条命令验收的是：

- `--limit=20` 控制最多处理多少条 story。
- `--concurrency=5` 控制最多同时抓取多少个 item。
- CLI 最终会调用 `TopStories` 和 `FetchItems`，而不是只根据参数个数打印 version。

如果 CLI 还没接好，不要卡在这条命令上。阶段三可以先用 `go test ./internal/hn` 验收核心能力，等 CLI 接线完成后再补集成验收。
