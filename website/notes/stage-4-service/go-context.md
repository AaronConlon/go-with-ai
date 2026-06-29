# Go context 入门：服务如何知道“该停了”

阶段四的 `main()` 里第一次密集出现 `context`：

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

<-ctx.Done()
stop()

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
	...
}
```

这些代码看起来很 Go，因为它把三个概念连在了一起：

- `context`：传递取消、超时这类控制信号。
- channel：用 `<-ctx.Done()` 等待信号发生。
- server lifecycle：收到退出信号后优雅关闭。

这篇先不追求覆盖 `context` 的所有 API，只建立阶段四需要的心智。

## 一句话理解

`context` 不是业务数据容器。它更像一张随工作流传递的“控制票据”，上面写着：

```text
这件事还要不要继续？
什么时候必须停？
有没有人已经按下取消按钮？
```

![小黑拉住 context 取消信号](/images/go-context/01-context-cancel-signal.png)

在阶段四里，`context` 的核心作用是：让服务知道什么时候该从“运行中”切换到“准备关闭”。

## `context.Background()` 是起点

你会看到：

```go
context.Background()
```

它可以先理解成一个空白起点：

```text
还没有取消
还没有超时
没有额外信息
```

很多 context 都从它派生出来：

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

`Background()` 本身不会自动取消，也不会自动超时。它只是告诉 Go：我要从这里开始创建一个新的控制信号链。

## `signal.NotifyContext`：把系统信号变成取消信号

阶段四入口里这一行很重要：

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
```

它的意思是：

```text
创建一个 ctx。
当进程收到 Ctrl+C 或 SIGTERM 时，取消这个 ctx。
```

两个信号先这样理解：

- `os.Interrupt`：通常是你在终端按 `Ctrl+C`。
- `syscall.SIGTERM`：通常是系统、容器、部署平台要求进程停止。

它返回两个值：

```go
ctx, stop
```

- `ctx`：后面用来等待取消。
- `stop`：停止接收这些信号，并释放相关资源。

所以写：

```go
defer stop()
```

表示 `main()` 结束前一定清理 signal 监听。

## `ctx.Done()` 返回的是 channel

这行是很多 JS 背景读者第一次会卡住的地方：

```go
<-ctx.Done()
```

先拆开看：

```go
ctx.Done()
```

它返回一个 channel。你可以先把类型读成：

```go
<-chan struct{}
```

意思是：这是一个只能接收的 channel。它不负责传业务数据，只负责发信号。

`ctx` 没被取消时，这个 channel 没有动静。

`ctx` 被取消时，这个 channel 会关闭。

## `<-ctx.Done()`：停在这里等取消

前面的 `<-` 是 channel receive：

```go
<-ctx.Done()
```

可以读成：

```text
从 ctx.Done() 这个 channel 等一个信号。
```

如果信号还没来，当前 goroutine 就会阻塞。

在阶段四入口里，这个 goroutine 是 main goroutine。也就是说：

```text
server 已经在另一个 goroutine 里运行
main goroutine 停在 <-ctx.Done()
用户按 Ctrl+C
ctx 被取消
ctx.Done() 关闭
main goroutine 继续往下走
进入 Shutdown
```

![小黑等待 Done 门打开](/images/go-context/03-context-done-channel.png)

注意区别：

```go
ctx.Done()
```

只是拿到 channel，不会等待。

```go
<-ctx.Done()
```

才是在这里等待。

## 为什么 `Shutdown` 要用新的 context

收到退出信号后，代码会继续往下执行：

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
	...
}
```

这里有个容易忽略的点：为什么不直接把前面的 `ctx` 传给 `Shutdown`？

因为前面的 `ctx` 已经被取消了。它的职责是告诉 main goroutine：

```text
有人要求退出了。
```

而 `Shutdown` 需要的是另一个控制信号：

```text
给 server 最多 5 秒收尾。
```

所以这里创建了一个新的 `shutdownCtx`。

注意这里的因果关系：

```text
不是 shutdownCtx 触发关闭。
是你调用 srv.Shutdown(shutdownCtx) 触发关闭。
shutdownCtx 只是限制这次关闭最多等多久。
```

![小黑给 shutdown 套上 5 秒沙漏](/images/go-context/02-context-timeout-shutdown.png)

## `context.WithTimeout`：给一段工作设截止时间

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

这行可以读成：

```text
创建一个新的 context。
它最多存活 5 秒。
5 秒后自动取消。
```

这里的“自动取消”具体会发生三件事：

```text
5 秒倒计时结束
-> shutdownCtx 被标记为已取消
-> shutdownCtx.Done() 这个 channel 被关闭
-> shutdownCtx.Err() 返回 context deadline exceeded
```

你可以用一小段代码观察：

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

<-shutdownCtx.Done()

fmt.Println(shutdownCtx.Err())
// context deadline exceeded
```

如果不是等超时，而是你主动调用：

```go
cancel()
```

那么 `shutdownCtx.Err()` 通常会是：

```text
context canceled
```

所以 `WithTimeout` 创建的 context 有两种常见结束方式：

| 结束方式 | `Done()` | `Err()` |
| --- | --- | --- |
| 5 秒到了 | 关闭 | `context deadline exceeded` |
| 手动调用 `cancel()` | 关闭 | `context canceled` |

`srv.Shutdown(shutdownCtx)` 会用它来决定最多等多久：

```text
停止接收新请求
等待正在处理的请求完成
如果 5 秒内完成，正常返回
如果超过 5 秒，返回 error
```

更准确地说，`Shutdown` 做事的顺序是：

```text
你调用 srv.Shutdown(shutdownCtx)
-> server 立刻开始进入关闭流程
-> 不再接收新连接
-> 关闭空闲连接
-> 等正在处理的请求结束
-> 如果 shutdownCtx 先超时或取消，就停止等待并返回 error
```

所以 `shutdownCtx` 不是“发消息让 server 关闭”的人。它像一个计时器和取消开关，告诉 `Shutdown`：

```text
你可以等，但不能无限等。
```

`cancel` 是清理函数：

```go
defer cancel()
```

即使没有等满 5 秒，也要释放 timer 相关资源。

## `context` 不应该拿来放业务数据

Go 的 `context` 也有 `WithValue`，但阶段四先不要使用它。

初学阶段可以记住这条规则：

```text
context 主要传取消、超时、deadline。
不要把它当成普通 map 或全局变量。
```

例如不要这样想：

```text
把用户资料放进 context
把配置放进 context
把业务结果放进 context
```

配置应该像阶段四这样显式传：

```go
cfg := config.Load()

srv := &http.Server{
	Addr: cfg.Addr,
}
```

## 阶段四入口里的 context 心智图

可以把阶段四入口想成两条线：

```text
server goroutine:
ListenAndServe() 一直运行，直到 server 关闭

main goroutine:
等待 ctx.Done()
收到退出信号
创建 shutdownCtx
调用 srv.Shutdown(shutdownCtx)
```

对应到代码：

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

go func() {
	_ = srv.ListenAndServe()
}()

<-ctx.Done()
stop()

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

_ = srv.Shutdown(shutdownCtx)
```

这里有两个 context：

| 名字 | 谁创建 | 作用 | 什么时候取消 |
| --- | --- | --- | --- |
| `ctx` | `signal.NotifyContext` | 等退出信号 | 收到 `Ctrl+C` 或 `SIGTERM` |
| `shutdownCtx` | `context.WithTimeout` | 限制关闭时间 | 5 秒超时或手动 `cancel()` |

## 新手检查清单

看到 `context` 时，先问四个问题：

- 这个 `ctx` 是从哪里来的？
- 它什么时候会被取消？
- 谁在等它的 `Done()`？
- 取消后代码下一步做什么？

阶段四的答案是：

```text
ctx 来自 signal.NotifyContext。
收到退出信号时取消。
main goroutine 用 <-ctx.Done() 等它。
取消后执行 srv.Shutdown。
```

这就是服务期最重要的 context 心智。
