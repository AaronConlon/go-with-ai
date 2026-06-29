# Go Client 构造函数语法拆解

这篇文档拆解阶段二里第一次看到的 `NewClient` 代码。它包含几个 Go 新概念：函数返回指针、`&` 取地址、结构体初始化、嵌套结构体初始化、标准库时间常量。

## 原始代码

```go
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
```

## 这段代码整体在做什么

它创建并返回一个 `Client`。

用一句话读：

> NewClient 接收一个 baseURL 字符串，创建一个带默认 HTTP timeout 的 Client，并返回它的指针。

也就是调用方可以这样用：

```go
client := NewClient("https://hacker-news.firebaseio.com/v0")
```

之后：

```go
client.TopStories(ctx)
```

## `func NewClient(baseURL string) *Client`

拆开看：

```text
func       定义函数
NewClient  函数名
baseURL    参数名
string     参数类型
*Client    返回值类型
```

完整意思：

```text
定义一个函数 NewClient。
它接收一个 string 类型的参数 baseURL。
它返回一个 *Client。
```

## 为什么函数名叫 `NewClient`

Go 里没有 class constructor 语法。

常见习惯是用：

```go
func NewXxx(...) *Xxx
```

作为构造函数。

例如：

```go
func NewClient(baseURL string) *Client
```

意思是：

```text
创建一个新的 Client。
```

这不是语言强制规则，而是 Go 社区常见命名习惯。

## `*Client` 是什么

```go
*Client
```

表示 `Client` 的指针。

可以先这样理解：

```text
Client   是一个值
*Client  是指向 Client 的地址
```

为什么返回指针？

因为 `Client` 通常代表一个可复用的对象，里面有 HTTP client、配置等。返回指针有几个好处：

- 避免复制整个结构体。
- 后续方法可以共享同一个 Client 状态。
- 这是 Go 里构造复杂对象的常见写法。

阶段二先不用深入内存模型，先记住：

> `*Client` 表示“返回一个 Client 的引用/地址”。

## `&Client{...}` 是什么

```go
return &Client{
	BaseURL: baseURL,
	HTTP: &http.Client{
		Timeout: 10 * time.Second,
	},
}
```

这里的 `&` 是取地址符。

```go
Client{...}
```

会创建一个 `Client` 值。

```go
&Client{...}
```

会创建一个 `Client` 值，并返回它的地址，也就是 `*Client`。

所以：

```go
return &Client{...}
```

刚好匹配函数返回值：

```go
func NewClient(...) *Client
```

## 为什么不是直接返回 `Client{...}`

如果函数签名是：

```go
func NewClient(baseURL string) Client
```

那就可以返回：

```go
return Client{...}
```

但现在函数签名是：

```go
func NewClient(baseURL string) *Client
```

返回值要求是 `*Client`，也就是指针，所以要返回：

```go
return &Client{...}
```

## `Client{ BaseURL: baseURL }` 是什么

这是结构体初始化，也叫 composite literal。

假设有：

```go
type Client struct {
	BaseURL string
	HTTP    *http.Client
}
```

那就可以这样创建：

```go
Client{
	BaseURL: baseURL,
	HTTP:    something,
}
```

左边的 `BaseURL` 是 struct 字段名。

右边的 `baseURL` 是函数参数。

```go
BaseURL: baseURL
```

可以读成：

```text
把参数 baseURL 的值，放进 Client 的 BaseURL 字段。
```

大小写不同不是偶然：

- `BaseURL`：结构体字段名，大写开头，外部可访问。
- `baseURL`：函数参数名，小写开头，只在当前函数里使用。

## `HTTP: &http.Client{...}` 是什么

```go
HTTP: &http.Client{
	Timeout: 10 * time.Second,
},
```

这里是嵌套初始化。

`Client` 结构体里有一个字段：

```go
HTTP *http.Client
```

它的类型是 `*http.Client`，也就是标准库 HTTP client 的指针。

所以给它赋值时，也要给一个 `*http.Client`：

```go
&http.Client{...}
```

这和外层的：

```go
&Client{...}
```

是同一种思路：

```text
http.Client{...}  创建 http.Client 值
&http.Client{...} 返回 http.Client 的指针
```

## `Timeout: 10 * time.Second`

```go
Timeout: 10 * time.Second
```

`Timeout` 是 `http.Client` 的字段。

`time.Second` 是 Go 标准库 `time` 包里的常量，表示“一秒”。

所以：

```go
10 * time.Second
```

表示 10 秒。

这比直接写数字更清楚：

```go
Timeout: 10 * time.Second
```

读起来就是：

```text
超时时间是 10 秒。
```

## 为什么最后有逗号

```go
Timeout: 10 * time.Second,
```

Go 的多行结构体初始化里，最后一项也要写逗号。

这是语法要求，也是 Go 格式化工具 `gofmt` 喜欢的风格。

## 一层一层读这段代码

从内到外读：

```go
&http.Client{
	Timeout: 10 * time.Second,
}
```

创建一个 HTTP client，并设置 10 秒超时。

再看外层：

```go
&Client{
	BaseURL: baseURL,
	HTTP:   上面那个 HTTP client,
}
```

创建一个我们自己的 HN client。

最后：

```go
return ...
```

把这个 client 返回给调用方。

## 对 JavaScript 的类比

粗略类比：

```js
function newClient(baseURL) {
  return {
    baseURL,
    http: {
      timeout: 10_000,
    },
  }
}
```

但 Go 里多了两件事：

- 类型明确：`baseURL string`、返回 `*Client`。
- 指针明确：`&Client{...}` 表示返回地址。

## 现阶段先记住

- `NewClient` 是 Go 里常见的构造函数命名。
- `*Client` 表示返回 Client 指针。
- `&Client{...}` 表示创建 Client 并返回它的地址。
- `BaseURL: baseURL` 表示把参数值放进结构体字段。
- `&http.Client{...}` 是嵌套创建一个标准库 HTTP client。
- `10 * time.Second` 表示 10 秒。

