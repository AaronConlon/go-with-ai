# Go Method Receiver 语法拆解

这篇文档拆解阶段二里第一次看到的 method 写法：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
}
```

第一次看会困惑：`func` 和函数名中间的 `(c *Client)` 是什么？

它叫 receiver，可以理解成“这个方法属于谁”。

## 普通函数 vs 方法

普通函数长这样：

```go
func TopStories(ctx context.Context) ([]int64, error) {
}
```

它是一个独立函数。

方法长这样：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
}
```

它属于 `Client`。

调用时也不一样。

普通函数：

```go
TopStories(ctx)
```

方法：

```go
client.TopStories(ctx)
```

## `(c *Client)` 是 receiver

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
}
```

拆开：

```text
func              定义函数或方法
(c *Client)       receiver，表示这个方法属于 *Client
TopStories        方法名
(ctx context.Context) 参数列表
([]int64, error)  返回值列表
```

`(c *Client)` 可以读成：

> 这个方法绑定在 `*Client` 上。方法内部用名字 `c` 代表当前这个 Client。

类似 JavaScript / TypeScript 里的：

```ts
class Client {
  topStories(ctx) {
    this.baseURL
  }
}
```

Go 没有 class，所以用 receiver 表达“这个函数属于某个类型”。

## 为什么 receiver 叫 `c`

```go
func (c *Client) TopStories(...)
```

`c` 是 receiver 的变量名。

它和普通参数一样，可以在函数体里使用：

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/topstories.json", nil)
```

这里的：

```go
c.BaseURL
```

就是访问当前 Client 的 `BaseURL` 字段。

Go 里 receiver 名通常很短：

```go
func (c *Client) ...
func (s *Store) ...
func (h *Handler) ...
```

不太常写成：

```go
func (client *Client) ...
```

但这不是强制规则，只是 Go 社区风格。

## 为什么是 `*Client`，不是 `Client`

```go
func (c *Client) TopStories(...)
```

这里 receiver 类型是 `*Client`，表示指针 receiver。

如果写成：

```go
func (c Client) TopStories(...)
```

那就是值 receiver。

阶段二这里用 `*Client` 更合适，因为：

- `Client` 里包含 `*http.Client`，通常希望复用同一个对象。
- 以后如果 `Client` 里有缓存、计数器、配置更新等状态，指针 receiver 更自然。
- 避免每次调用方法都复制整个 `Client` 值。

先记住：

> 如果这个类型代表一个有状态的服务对象，比如 client、store、handler，通常用指针 receiver。

## `TopStories(ctx context.Context)`

```go
TopStories(ctx context.Context)
```

这是方法名和参数列表。

- `TopStories` 是方法名。
- `ctx` 是参数名。
- `context.Context` 是参数类型。

`context.Context` 用来传递取消、超时和请求范围信息。阶段二所有外部 I/O 都建议带上 `ctx`。

## `([]int64, error)` 是返回值

```go
([]int64, error)
```

表示这个方法返回两个值：

1. `[]int64`：story id 列表。
2. `error`：错误。

在 Go 里，多个返回值很常见。调用时通常这样写：

```go
ids, err := client.TopStories(ctx)
if err != nil {
	return err
}
```

先检查 `err`，再使用 `ids`。

## 和 JS `try/catch` 的区别

如果用 JavaScript 心智类比，很多异步错误会写成：

```js
try {
  const ids = await client.topStories()
  // 使用 ids
} catch (err) {
  // 处理错误
}
```

这里的错误会通过 throw 或 rejected Promise 进入 `catch`。

Go 通常不这样写。Go 更常见的是：

```go
ids, err := client.TopStories(ctx)
if err != nil {
	return err
}

// 到这里才说明 TopStories 成功了，可以使用 ids。
```

也就是说，Go 把错误放在返回值里，让调用方在很小的范围内立刻处理。

这有几个好处：

- 代码读到哪一行，就能看到这一行可能失败。
- 每一步失败后怎么处理，由当前调用点决定。
- 不容易把“网络错误”“JSON 错误”“业务状态码错误”混成一个很远的统一 `catch`。

阶段二先养成一个习惯：

```text
看到 xxx, err := ...
下一行优先写 if err != nil { ... }
确认 err 是 nil 后，再使用 xxx。
```

## 为什么返回 `[]int64`

HN 的 `topstories.json` 返回的是一组 story id：

```json
[123, 456, 789]
```

Go 里对应：

```go
[]int64
```

`[]` 表示切片，也就是列表。

`int64` 表示整数类型。

所以：

```go
[]int64
```

可以先读成：

```text
一组 int64 整数
```

## 为什么返回 `error`

外部 HTTP 请求可能失败：

- 网络断了。
- DNS 解析失败。
- 请求超时。
- HN 返回 500。
- JSON 格式不对。

Go 的习惯是把错误作为返回值显式返回：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error)
```

调用方必须处理：

```go
ids, err := client.TopStories(ctx)
if err != nil {
	// 处理错误
}
```

## 完整读法

这行：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
```

可以读成：

> 定义一个属于 `*Client` 的方法，方法名是 `TopStories`。它接收一个 `context.Context`，返回一组 story id 和一个错误。

## 和普通函数的区别

如果写成普通函数：

```go
func TopStories(c *Client, ctx context.Context) ([]int64, error) {
}
```

调用时是：

```go
TopStories(client, ctx)
```

写成 method：

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
}
```

调用时是：

```go
client.TopStories(ctx)
```

后者更符合“这个能力属于 Client”的直觉。

## 现阶段先记住

- `func TopStories(...)` 是普通函数。
- `func (c *Client) TopStories(...)` 是 `Client` 的方法。
- `(c *Client)` 叫 receiver。
- 方法内部可以用 `c.BaseURL`、`c.HTTP` 访问当前 Client 的字段。
- `[]int64, error` 表示返回 story id 列表和错误。
