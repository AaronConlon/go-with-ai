package hn

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

/**
 * 批量获取 hn 的数据
 * 定义为 Clint 结构体的方法
 */
func (c *Client) FetchItems(
	ctx context.Context,
	ids []int64,
	concurrency int,
) ([]Item, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be positive!")
	}

	// 预先创建一个长度和 ids 一样的 []Item 结果列表。
	items := make([]Item, len(ids))

	// errgroup 来管理一组 goroutine
	// WithContext 返回一个新的 ctx
	//  只要其中一个 goroutine 返回 error，这个 ctx 就会被取消
	//
	// 通过这个方式创建上下文，并且可以统一管理内部的 goroutine（协程，轻量线程）
	// errgroup = 用来管理一组会返回 error 的 goroutine 的工具。
	// g：一个 group，用来启动 goroutine、等待 goroutine、收集 error。
	// ctx：一个新的 context，基于传入的旧 ctx 派生出来。
	g, ctx := errgroup.WithContext(ctx)

	// sem = semaphore 信号量，相当于并发“许可证”
	sem := make(chan struct{}, concurrency)

	// 遍历所有 id 切片
	for i, id := range ids {
		g.Go(func() error {
			// 获取一个并发许可证，记住这种写法
			// 如果信号量已满，则会阻塞
			sem <- struct{}{}
			defer func() {
				// 函数结束，一定要释放信号量，相当于归还许可证
				<-sem
			}()

			// 真正去请求单挑 item
			item, err := c.Item(ctx, id)
			if err != nil {
				// 出错了
				return err
			}

			items[i] = item
			return nil
		})
	}

	// 如果等待整个 group 都完成，出现错误，则直接返回 0 值和错误信息
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return items, nil

}
