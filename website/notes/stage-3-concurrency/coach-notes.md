# 阶段三教练笔记

阶段三开始学习 Go 的并发。对 Go 新手来说，不要先背 channel 的全部语法，而要先建立一个工程直觉：

> 并发不是“同时开很多任务”这么简单，而是要控制数量、知道何时取消、知道错误怎么返回。

## 本阶段注释策略

阶段三文档代码继续使用高密度注释。重点解释：

- goroutine 是什么。
- 为什么不能无限制启动 goroutine。
- channel 在这里为什么能当 semaphore。
- `errgroup.WithContext` 如何收集错误和传播取消。
- 为什么循环里要写 `i, id := i, id`。

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
		// 这行非常重要。
		// Go 的循环变量会被复用；如果不重新绑定，goroutine 可能读到错误的 i/id。
		i, id := i, id

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

测试不一定一开始写完整，但至少要覆盖：

- `concurrency <= 0` 会返回错误。
- 正常情况下能抓取多条 item。
- 返回结果顺序和输入 ids 一致。
- 其中一个 item 请求失败时，`FetchItems` 返回 error。

## 阶段三验收

在 `hn-agent/` 内执行：

```bash
go test ./internal/hn
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

如果 CLI 还没接好，先只验收：

```bash
go test ./internal/hn
```

