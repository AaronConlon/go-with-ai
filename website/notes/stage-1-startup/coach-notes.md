# 阶段一教练笔记

这篇文档记录阶段一真正开始写代码前，值得保存下来的判断、概念和操作顺序。代码由学习者自己创建；教练负责解释、提供样例、review 结果和帮助 debug。

## 本阶段不要急着做什么

阶段一先不要调用 Hacker News API，也不要引入数据库、AI API、Web server 或复杂 CLI 框架。

当前只做一件事：建立 Go 项目最小骨架，并让它能运行、能测试。

这样做的原因是：如果一开始同时学习 HTTP、JSON、并发、测试和服务化，问题会混在一起。阶段一应该先把 Go 的基本工程循环跑通。

## 前三阶段代码注释原则

前三阶段默认学习者对 Go 完全陌生，所以文档里的示例代码要写成“教学版代码”：

- 关键语法旁边要有注释，例如 `package`、`import`、`struct`、`func`、`return`。
- 标准库第一次出现时要解释用途，例如 `fmt`、`os`、`testing`。
- Go 特有写法要解释意图，例如 table-driven test、`t.Run`、`os.Exit`。
- 注释服务理解，不追求生产代码的简洁程度。

等进入第四阶段以后，再逐渐把注释密度降到真实项目常见水平。

## 工作区和项目目录

仓库根目录是学习工作区，负责保存文档、背景资料和实战记录。真正的 Hacker News Go 项目放在独立目录：

```text
hn-agent/
```

这样做有几个好处：

- 根目录不会被某一个 Go module 绑定。
- 文档站、学习记录和实战代码职责清晰。
- 后续如果再增加其他练习项目，可以继续在根目录下并列创建。
- `hn-agent/` 可以像一个真实 Go 项目一样独立构建、测试和部署。

## 本阶段交付物

完成阶段一后，项目应该至少包含：

```text
hn-agent/go.mod
hn-agent/cmd/hnctl/main.go
hn-agent/internal/hn/story.go
hn-agent/internal/hn/story_test.go
```

对应能力：

- `hn-agent/go.mod`：Hacker News 项目有明确 module 边界。
- `hn-agent/cmd/hnctl/main.go`：有一个可运行的 CLI 入口。
- `hn-agent/internal/hn/story.go`：有一个最小领域模型。
- `hn-agent/internal/hn/story_test.go`：有第一条 table-driven test。

## 第一步：确认 Go 版本

执行：

```bash
go version
```

当前机器已确认是：

```text
go version go1.25.1 darwin/arm64
```

这个版本足够完成当前项目。阶段一不会用到依赖版本敏感的高级能力。

## 第二步：创建独立项目目录

推荐执行：

```bash
mkdir -p hn-agent
cd hn-agent
```

这里的关键点是：后续所有 Go 命令默认都在 `hn-agent/` 内执行，而不是在仓库根目录执行。

## 第三步：初始化 Go module

在 `hn-agent/` 目录内执行：

```bash
go mod init github.com/aaron/go-with-ai/hn-agent
```

`module path` 不是随便的项目名字，它会影响代码中的 import 路径。例如后续 CLI 引用内部包时会写：

```go
import "github.com/aaron/go-with-ai/hn-agent/internal/hn"
```

如果暂时不想绑定 GitHub 路径，也可以用：

```bash
go mod init hn-agent
```

但学习阶段建议用 `github.com/aaron/go-with-ai/hn-agent`，因为它表达了“这个 Go module 位于学习仓库下的 `hn-agent/` 子项目”。

## 第四步：创建目录

在 `hn-agent/` 内执行：

```bash
mkdir -p cmd/hnctl internal/hn
```

目录含义：

| 目录 | 作用 |
| --- | --- |
| `cmd/hnctl` | 命令行工具入口，用于学习、调试和手动触发功能 |
| `internal/hn` | Hacker News 领域相关代码，先放模型，下一阶段再写 HTTP client |

## 为什么先有 `cmd/hnctl`

`hnctl` 是学习期间的控制台工具。它让我们不用一开始就搭 HTTP server，也能验证代码是否能运行。

后续可以逐步扩展：

```bash
go run ./cmd/hnctl
go run ./cmd/hnctl version
go run ./cmd/hnctl top --limit=10
go run ./cmd/hnctl summarize --id=...
```

这条路径很适合学习：每一步都能通过命令马上看到结果。

## 为什么使用 `internal/`

`internal/` 是 Go 的特殊目录。放在里面的包只能被当前 module 的父级目录树内部导入，外部项目不能直接 import。

这对应用项目很合适：

- 避免过早承诺公共 API。
- 降低随意复用导致的耦合。
- 让内部实现可以随着学习和项目演进调整。

阶段一先把 `hn` 包放在 `internal/hn`，代表它是本项目内部的 Hacker News 领域能力。

## 第五步：写最小领域模型

创建：

```text
internal/hn/story.go
```

建议内容：

```go
// package 声明这个文件属于哪个包。
// 同一个目录下的 Go 文件通常使用同一个 package 名。
package hn

// Story 表示一条 Hacker News story。
// 现在字段还很少，只保留阶段一需要理解的基础字段。
type Story struct {
	// ID 是 HN story 的唯一编号。
	// 使用 int64 是因为外部 API 的 id 本质是整数，且未来可能比较大。
	ID int64

	// Title 是 story 标题。
	// 阶段一我们把“有标题”当成一条 story 是否有效的最低要求之一。
	Title string

	// URL 是 story 指向的外部链接。
	// 阶段一暂时不校验 URL，只先把字段放进模型里。
	URL string

	// Score 是 HN 上的分数。
	// 阶段一暂时不把分数作为有效性条件。
	Score int
}

// IsValidStory 判断一条 story 是否满足最小有效条件。
// Go 函数的写法是：func 函数名(参数名 参数类型) 返回值类型。
func IsValidStory(story Story) bool {
	// return 后面是这个函数的返回值。
	// && 表示“并且”：ID 必须大于 0，并且 Title 不能为空字符串。
	return story.ID > 0 && story.Title != ""
}
```

这里暂时不追求完整，只练几个关键点：

- `package hn`：声明当前文件属于 `hn` 包。
- `type Story struct`：定义领域数据结构。
- `func IsValidStory`：写一个可测试的业务规则。
- `int64`：HN 的 id 后续来自外部 API，用整数承载。
- `Title != ""`：最小有效 story 至少应该有标题。

## 第六步：写 table-driven test

创建：

```text
internal/hn/story_test.go
```

建议内容：

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

### 为什么 Go 常用 table-driven test

因为很多函数的测试本质是“一组输入，对应一组期望输出”。用表格组织测试，可以让新增 case 很轻：

- 增加一行输入。
- 增加一个期望值。
- 不需要复制整段测试逻辑。

这个模式在 Go 标准库和大量生产项目里都很常见。

### 什么是领域函数

“领域函数”不是 Go 语言的语法词，而是工程建模里的说法。

领域可以理解成“这个项目正在解决的问题范围”。在这个项目里，领域就是 Hacker News 摘要推送，所以 story、item、digest、summary、push 这些概念都属于项目领域。

领域函数就是直接表达业务规则或项目概念的函数。例如：

```go
func IsValidStory(story Story) bool {
	return story.ID > 0 && story.Title != ""
}
```

它表达的不是“怎么打印日志”“怎么读取命令行参数”“怎么发 HTTP 请求”，而是：

> 一条 story 在我们的项目里，什么情况下算有效？

所以 `IsValidStory` 可以叫领域函数。

和它对比：

```go
fmt.Printf("story valid: %v\n", hn.IsValidStory(story))
```

这行是输出逻辑，不是领域函数。

```go
os.Exit(run(os.Args[1:]))
```

这行是程序入口逻辑，也不是领域函数。

阶段一先写一个很小的领域函数，是为了先练习“把业务判断写成可测试代码”。

### `t.Run` 的意义

`t.Run(tt.name, ...)` 会把每个 case 作为子测试运行。好处是失败时能看到具体是哪一个 case 失败，而不是只知道整个测试失败。

## 第六步之后：立刻运行单元测试

你说得对：写完 `internal/hn/story_test.go` 之后，应该先运行测试，再进入第七步写 CLI。

这一刻要验证的是 `internal/hn` 这个包本身是否正确，不需要等 CLI 写完。

在 `hn-agent/` 目录内执行：

```bash
go test ./internal/hn
```

如果想看到每个子测试的名字，加 `-v`：

```bash
go test -v ./internal/hn
```

如果只想运行 `TestIsValidStory`：

```bash
go test -v ./internal/hn -run TestIsValidStory
```

### 这三个命令有什么区别

| 命令 | 用途 |
| --- | --- |
| `go test ./internal/hn` | 跑 `internal/hn` 包里的所有测试，输出简洁 |
| `go test -v ./internal/hn` | 跑所有测试，并显示每个测试和子测试 |
| `go test -v ./internal/hn -run TestIsValidStory` | 只跑名字匹配 `TestIsValidStory` 的测试 |

### 预期结果

如果测试通过，可能看到类似：

```text
ok  	github.com/aaron/go-with-ai/hn-agent/internal/hn	0.123s
```

使用 `-v` 时，可能看到类似：

```text
=== RUN   TestIsValidStory
=== RUN   TestIsValidStory/valid_story
=== RUN   TestIsValidStory/missing_id
=== RUN   TestIsValidStory/missing_title
--- PASS: TestIsValidStory (0.00s)
    --- PASS: TestIsValidStory/valid_story (0.00s)
    --- PASS: TestIsValidStory/missing_id (0.00s)
    --- PASS: TestIsValidStory/missing_title (0.00s)
PASS
ok  	github.com/aaron/go-with-ai/hn-agent/internal/hn	0.123s
```

### 为什么第六步后要先测试

因为这是最小反馈循环：

```text
写领域函数 -> 写单元测试 -> 运行 go test -> 确认包正确
```

如果现在就有错误，问题范围很小，只可能在：

- `internal/hn/story.go`
- `internal/hn/story_test.go`
- package 名是否一致
- 函数名或字段名是否拼错

如果等 CLI 写完再测试，错误来源会变多：可能是 import path、`main.go`、命令行参数、`os.Exit`，排查成本更高。

### 为什么文档原来直接进入第七步

原来的文档把阶段一最终验收放在第八步统一做：

```bash
go test ./...
go run ./cmd/hnctl
```

这个顺序对有经验的开发者可以接受，但对第一次学 Go 不够友好。更好的学习节奏是每完成一个最小单元就运行一次：

```text
第六步：写测试 -> 立即 go test ./internal/hn
第七步：写 CLI -> 再 go run ./cmd/hnctl
第八步：整体 go test ./...
```

## 第七步：写最小 CLI

创建：

```text
cmd/hnctl/main.go
```

如果 module path 使用 `github.com/aaron/go-with-ai/hn-agent`，建议内容：

```go
// main package 表示这是一个可执行程序入口。
// 只有 package main 里的 main 函数才能被 go run / go build 当成程序启动点。
package main

import (
	// fmt 用来格式化输出，例如打印字符串、数字、布尔值。
	"fmt"

	// os 用来访问操作系统能力，这里用它读取命令行参数和设置退出码。
	"os"

	// 这是我们自己项目里的内部包。
	// import 路径来自 hn-agent/go.mod 里的 module path。
	"github.com/aaron/go-with-ai/hn-agent/internal/hn"
)

// const 定义常量。
// version 暂时写死，后续可以在构建时注入。
const version = "0.1.0"

// main 是程序入口。
// 它应该尽量保持很薄，把真正逻辑交给 run。
func main() {
	// os.Args 是命令行参数。
	// os.Args[0] 是程序名，所以传给 run 的是 os.Args[1:]。
	// os.Exit 用 run 返回的数字作为进程退出码。
	os.Exit(run(os.Args[1:]))
}

// run 承担主要 CLI 逻辑。
// 参数 args 是命令行参数切片。
// 返回 int 是进程退出码：0 通常表示成功，非 0 表示失败。
func run(args []string) int {
	// len(args) 表示参数数量。
	// 先判断 len(args) > 0，是为了避免访问 args[0] 时越界。
	if len(args) > 0 && (args[0] == "version" || args[0] == "--version") {
		fmt.Println(version)
		return 0
	}

	// 创建一个 Story struct 的值。
	// 这种写法叫 composite literal，常用于初始化结构体。
	story := hn.Story{
		ID:    1,
		Title: "Hello Go",
		URL:   "https://news.ycombinator.com/",
		Score: 100,
	}

	// Printf 可以用占位符格式化输出。
	// %s 表示字符串，%v 表示使用默认格式输出任意值。
	fmt.Printf("hnctl %s\n", version)
	fmt.Printf("story valid: %v\n", hn.IsValidStory(story))
	return 0
}
```

### 为什么把逻辑放进 `run`

`main` 函数不方便直接测试，因为它没有返回值，通常还会调用 `os.Exit`。

把主要逻辑放到：

```go
func run(args []string) int
```

有几个好处：

- 后续可以直接测试 `run`。
- 可以把命令行参数作为普通 slice 传入。
- 可以用返回码表达成功或失败。
- `main` 保持很薄，只负责连接系统入口。

阶段一不一定马上测试 `run`，但这个结构会为后续扩展留下空间。

## 第八步：验证

在 `hn-agent/` 目录内执行：

```bash
go test ./...
go run ./cmd/hnctl
go run ./cmd/hnctl version
```

期望看到：

```text
ok  	github.com/aaron/go-with-ai/hn-agent/internal/hn	...
```

以及类似：

```text
hnctl 0.1.0
story valid: true
```

`version` 命令期望输出：

```text
0.1.0
```

## 遇到问题时先贴这些输出

完成或卡住时，把下面三段输出贴出来：

```bash
pwd
cat go.mod
go test ./...
go run ./cmd/hnctl
```

如果是 import path 报错，也贴：

```bash
find cmd internal -type f -maxdepth 4 -print
```

## 本阶段判断

### 为什么先写模型和测试

先写一个 `Story` 和 `IsValidStory`，不是因为它们复杂，而是为了把 Go 的基础循环跑通：

```text
写代码 -> 写测试 -> go test ./... -> go run ./cmd/hnctl
```

这个循环比单纯阅读语法更重要。后续所有复杂能力都会套在这个循环上。

### 为什么阶段一不引入 CLI 框架

阶段一只需要 `os.Args` 和标准库。过早引入 Cobra、urfave/cli 之类框架，会把学习焦点从 Go 的基本工程模型转移到框架配置上。

等命令变复杂，例如有多个子命令、flag、配置文件和帮助文档时，再考虑 CLI 框架。

### 为什么不马上做 HN API

HN API 属于阶段二。它会引入：

- HTTP client。
- JSON 解码。
- timeout。
- status code。
- 外部错误。
- `httptest`。

这些内容应该在项目骨架稳定之后再加，否则调试时很难判断问题来自 Go 基础、目录结构、网络调用还是测试方式。
