# Go 变量定义语法

Go 里常见的“定义名字”的方式有三类：

- `var`：定义变量。
- `const`：定义常量。
- `:=`：在函数内部快速定义变量。

它们看起来都像“创建一个名字”，但含义不一样。

## 一句话区分

| 写法 | 用途 | 能不能改 | 常见位置 |
| --- | --- | --- | --- |
| `var name string` | 定义变量 | 可以改 | package 级别或函数内部 |
| `const version = "0.1.0"` | 定义常量 | 不可以改 | package 级别或函数内部 |
| `name := "Aaron"` | 短变量声明 | 可以改 | 只能在函数内部 |

## `var`：定义变量

最完整写法：

```go
var title string
```

意思是：

```text
定义一个变量 title，它的类型是 string。
```

如果没有手动赋值，Go 会给它一个零值。

```go
var title string // 零值是 ""
var score int    // 零值是 0
var ok bool      // 零值是 false
```

也可以定义时直接赋值：

```go
var title string = "Learning Go"
```

Go 可以根据右边的值推断类型，所以也可以省略类型：

```go
var title = "Learning Go"
```

这里 Go 会推断 `title` 是 `string`。

## `const`：定义常量

```go
const version = "0.1.0"
```

意思是：

```text
定义一个常量 version，它的值永远是 "0.1.0"。
```

常量不能被重新赋值：

```go
const version = "0.1.0"
version = "0.2.0" // 不允许
```

在 `cmd/hnctl/main.go` 里：

```go
const version = "0.1.0"
```

很合理，因为程序版本在运行期间不应该变化。

## `:=`：短变量声明

```go
got := IsValidStory(tt.story)
```

意思是：

```text
创建一个新变量 got，并把 IsValidStory(tt.story) 的返回值放进去。
```

Go 会自动推断 `got` 的类型。

如果 `IsValidStory` 返回 `bool`，那么：

```go
got := IsValidStory(tt.story)
```

就等价于：

```go
var got bool = IsValidStory(tt.story)
```

只是更短。

## `:=` 只能在函数内部用

可以这样：

```go
func run() {
	title := "Learning Go"
	fmt.Println(title)
}
```

不能这样：

```go
package hn

title := "Learning Go" // 不允许
```

package 级别只能用 `var` 或 `const`：

```go
package hn

var defaultTitle = "Learning Go"
const version = "0.1.0"
```

## `=`：重新赋值

`:=` 是定义新变量。

`=` 是给已有变量重新赋值。

```go
score := 1 // 定义新变量
score = 2  // 修改已有变量
```

如果变量已经存在，再写：

```go
score := 2
```

通常会报错，因为这不是修改，而是又想定义一个同名变量。

## 为什么有时候“不加 var”

你看到的“不加 var”，通常是 `:=`。

例如：

```go
tests := []struct {
	name string
	want bool
}{}
```

这里不是没定义变量，而是用了短变量声明。

Go 的意思是：

```text
我创建一个新变量 tests，类型从右边推断。
```

这种写法在函数内部非常常见。

## 多变量声明

Go 可以一次定义多个变量：

```go
var id int64
var title string
```

也可以写成一组：

```go
var (
	id    int64
	title string
)
```

函数内部也常见多个返回值：

```go
got, err := FetchStory()
```

这里一次定义了两个变量：

- `got`
- `err`

这类写法后面会经常出现。

## `const` 和 `var` 的选择

如果值不会变，用 `const`：

```go
const version = "0.1.0"
```

如果值会变，用 `var` 或 `:=`：

## 字符串为什么用双引号

Go 里字符串要用双引号或反引号：

```go
t.Fatal("expected error, got nil")
```

或者：

```go
json := `{"id":101,"title":"ok"}`
```

单引号不是字符串。单引号表示 `rune`，也就是一个 Unicode 字符：

```go
letter := 'a'
```

所以这行是错的：

```go
t.Fatal('expected error, got nil') // illegal rune literal
```

因为单引号里放了很多字符，而 Go 只允许它表示一个 `rune`。

`gofmt` 也不会自动把它修成双引号。原因是：`gofmt` 只能格式化已经能被 Go parser 解析的代码。`illegal rune literal` 属于语法解析失败，代码还没有形成 AST，`gofmt` 没有东西可以安全地格式化。

另外，即使它能猜，也不应该随便把单引号改成双引号，因为这会改变字面量类型：

```go
'a' // rune
"a" // string
```

这两者在 Go 里不是同一个类型。

```go
score := 0
score = score + 1
```

## `var` 和 `:=` 的选择

函数内部，优先用 `:=`：

```go
story := Story{ID: 1, Title: "Learning Go"}
```

需要明确零值时，用 `var`：

```go
var item Item
```

这在 JSON 解码里很常见：

```go
var item Item
err := json.NewDecoder(resp.Body).Decode(&item)
```

这里先创建一个空的 `Item`，再让 JSON 解码器把数据填进去。

package 级别，不能用 `:=`，只能用 `var` 或 `const`。

## 对 JavaScript 的类比

粗略类比：

| JavaScript | Go |
| --- | --- |
| `let score = 1` | `score := 1` 或 `var score = 1` |
| `const version = "0.1.0"` | `const version = "0.1.0"` |
| `let item` | `var item Item` |

但注意：Go 是静态类型语言。变量类型一旦确定，不能随便换。

```go
score := 1
score = "one" // 不允许，score 已经是 int
```

## 现阶段先记住

- `const`：不会变的值。
- `var`：变量，package 级别和函数内部都能用。
- `:=`：函数内部快速创建变量。
- `=`：给已有变量重新赋值。
- Go 变量一旦有类型，后面不能改成别的类型。
