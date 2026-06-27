# Go 测试代码语法拆解

这篇文档拆解阶段一第一段测试代码。目标不是一次性记住所有语法，而是第一次看到 Go 代码时，知道每个符号大概在做什么。

## 原始代码

```go
// 测试文件可以和被测试代码使用同一个 package。
// 这样测试可以直接访问当前包里的导出函数和类型。
package hn

// testing 是 Go 标准库里的测试包。
// 只要函数名以 Test 开头，并接收 *testing.T，go test 就会自动运行它。
import "testing"

func TestIsValidStory(t *testing.T) {
	// tests 是测试用例表。
	// []struct 表示“一个匿名 struct 的切片”，可以理解成一张测试表。
	tests := []struct {
		// name 用来描述当前测试 case，失败时更容易定位。
		name  string

		// story 是输入。
		story Story

		// want 是期望输出。
		want  bool
	}{
		{
			name: "valid story",
			story: Story{
				ID:    1,
				Title: "Learning Go",
			},
			want: true,
		},
		{
			name: "missing id",
			story: Story{
				Title: "Learning Go",
			},
			want: false,
		},
		{
			name: "missing title",
			story: Story{
				ID: 1,
			},
			want: false,
		},
	}

	// range 遍历 tests 里的每一个测试用例。
	// tt 是当前这条测试用例。
	for _, tt := range tests {
		// t.Run 会创建一个子测试。
		// 子测试名就是 tt.name，失败时输出会更清楚。
		t.Run(tt.name, func(t *testing.T) {
			// got 是实际得到的结果。
			got := IsValidStory(tt.story)

			// 如果实际结果和期望结果不一致，测试失败。
			if got != tt.want {
				// t.Fatalf 会标记测试失败，并停止当前子测试。
				t.Fatalf("IsValidStory() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

## `package hn`

```go
package hn
```

每个 Go 文件开头都必须声明自己属于哪个 package。

这里的意思是：这个测试文件属于 `hn` 包。

如果被测试代码在：

```text
internal/hn/story.go
```

并且里面也是：

```go
package hn
```

那么测试文件可以直接使用同一个包里的类型和函数，例如：

```go
Story
IsValidStory
```

不需要再 import 自己。

## `import "testing"`

```go
import "testing"
```

`import` 用来引入其他 package。

`testing` 是 Go 标准库自带的测试包。写测试时，只要用到 `*testing.T`、`t.Run`、`t.Fatalf`，就需要 import 它。

如果 import 了但没使用，Go 会报错。这是 Go 的风格：不允许没用到的 import 留在代码里。

## 测试函数签名

```go
func TestIsValidStory(t *testing.T) {
```

拆开看：

```text
func              定义函数
TestIsValidStory  函数名
(t *testing.T)    参数列表
{                 函数体开始
```

Go 的测试函数有固定约定：

- 函数名以 `Test` 开头。
- 参数是 `t *testing.T`。
- 没有返回值。

`go test` 会自动找到这些函数并执行。

### `*testing.T` 是什么

可以先粗略理解为：`testing.T` 是测试运行时给你的“测试控制器”。

它能做这些事：

- 标记测试失败。
- 输出失败信息。
- 创建子测试。
- 控制当前测试流程。

前面的 `*` 表示这是一个指针。现阶段可以先记住：测试函数固定写 `t *testing.T`。

## `:=` 短变量声明

```go
tests := []struct {
```

`:=` 是 Go 的短变量声明。

它等价于“创建一个新变量，并让 Go 根据右边的值自动推断类型”。

例如：

```go
name := "valid story"
```

Go 会推断 `name` 是 `string`。

这里：

```go
tests := ...
```

表示创建一个变量 `tests`。

## `[]struct { ... }` 是什么

```go
tests := []struct {
	name  string
	story Story
	want  bool
}{
	...
}
```

这段第一次看会有点密。分成两层：

```go
struct {
	name  string
	story Story
	want  bool
}
```

表示一个匿名结构体类型。它有三个字段：

- `name string`
- `story Story`
- `want bool`

前面的 `[]` 表示“切片”：

```go
[]struct { ... }
```

意思是：很多个这种结构体组成的列表。

所以整体意思是：

> tests 是一个测试用例列表，每个测试用例都有 name、story、want 三个字段。

## 什么是切片

切片可以先理解成 Go 里最常用的动态列表。

类似 JavaScript 里的 array：

```js
const tests = [...]
```

Go 里写成：

```go
tests := []SomeType{...}
```

这个例子里，`SomeType` 是匿名 struct。

## 字段声明语法

```go
name  string
story Story
want  bool
```

Go 的字段声明顺序是：

```text
字段名 类型
```

这和 TypeScript 不一样。

TypeScript 常见写法：

```ts
name: string
```

Go 写法：

```go
name string
```

## 复合字面量

```go
{
	name: "valid story",
	story: Story{
		ID:    1,
		Title: "Learning Go",
	},
	want: true,
}
```

这是在创建一个测试用例。

这种用 `{ 字段名: 值 }` 初始化结构体的写法，叫 composite literal，可以翻译成“复合字面量”。

里面的：

```go
Story{
	ID:    1,
	Title: "Learning Go",
}
```

是在创建一个 `Story` 值。

## 为什么有逗号

Go 的多行复合字面量里，最后一项后面也要有逗号：

```go
want: true,
```

这是 Go 语法要求。好处是后续新增字段时，diff 更干净。

## `for _, tt := range tests`

```go
for _, tt := range tests {
```

这是 Go 的遍历写法。

`range tests` 会产生两个值：

```text
索引, 当前元素
```

例如：

```go
for index, value := range tests {
```

但这里我们不需要索引，所以用 `_` 忽略：

```go
for _, tt := range tests {
```

`_` 叫 blank identifier，空白标识符。它表示“这个值我知道存在，但我不用”。

`tt` 是当前测试用例。名字通常来自 `test case` 的缩写。

## `t.Run`

```go
t.Run(tt.name, func(t *testing.T) {
```

`t.Run` 会创建一个子测试。

第一个参数：

```go
tt.name
```

是子测试名字。

第二个参数：

```go
func(t *testing.T) {
	...
}
```

是一个匿名函数，表示这个子测试要执行的内容。

好处是：如果某个 case 失败，输出会显示具体 case 名，例如：

```text
--- FAIL: TestIsValidStory/missing_title
```

## 函数调用

```go
got := IsValidStory(tt.story)
```

意思是：

- 调用 `IsValidStory`。
- 把当前测试用例里的 `story` 传进去。
- 把返回值保存到 `got`。

`got` 是测试里常见命名，表示“实际得到的结果”。

`want` 表示“期望结果”。

## `if got != tt.want`

```go
if got != tt.want {
```

Go 的 `if` 不需要小括号。

JavaScript：

```js
if (got !== want) {
}
```

Go：

```go
if got != want {
}
```

`!=` 表示不相等。

## `t.Fatalf`

```go
t.Fatalf("IsValidStory() = %v, want %v", got, tt.want)
```

`t.Fatalf` 会让当前测试失败，并打印错误信息。

`Fatalf` 里的 `f` 表示 format，格式化。

这里的 `%v` 是占位符，表示“用默认格式打印这个值”。

所以如果 `got` 是 `false`，`tt.want` 是 `true`，输出可能类似：

```text
IsValidStory() = false, want true
```

## 这段代码整体在做什么

它做了三件事：

1. 定义一张测试表 `tests`。
2. 遍历每个测试用例。
3. 调用 `IsValidStory`，比较实际结果和期望结果。

用一句话说：

> 用多组输入验证 `IsValidStory` 是否按预期判断 story 有效性。

## 对 JavaScript 的类比

大概可以类比成：

```js
const tests = [
  {
    name: "valid story",
    story: { id: 1, title: "Learning Go" },
    want: true,
  },
  {
    name: "missing id",
    story: { title: "Learning Go" },
    want: false,
  },
]

for (const tt of tests) {
  test(tt.name, () => {
    const got = isValidStory(tt.story)
    expect(got).toBe(tt.want)
  })
}
```

区别是 Go 更显式：

- 类型写在结构里。
- 测试函数由 `go test` 自动发现。
- 不需要额外测试框架。
- 错误输出通过 `testing.T` 控制。

## 现阶段先记住这些

不用一次记住所有语法。现在先记住：

- `package`：文件属于哪个包。
- `import`：引入别的包。
- `func TestXxx(t *testing.T)`：Go 测试函数。
- `:=`：声明并赋值。
- `[]struct{...}`：一组匿名结构体，常用于测试表。
- `for _, tt := range tests`：遍历列表。
- `t.Run`：子测试。
- `got / want`：实际结果和期望结果。
- `t.Fatalf`：测试失败并输出信息。

