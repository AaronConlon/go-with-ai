# 阶段 1：启动期完整代码

阶段 1 的目标是建立最小 Go 项目：有 module、有 CLI、有领域函数、有测试。

这里保存的是学习版完整代码。注释会比真实生产代码更密，目的是帮助第一次接触 Go 的读者理解每一块在做什么。

## 文件清单

```text
hn-agent/go.mod
hn-agent/cmd/hnctl/main.go
hn-agent/internal/hn/story.go
hn-agent/internal/hn/story_test.go
```

## `go.mod`

```go
// hn-agent/go.mod
// module 声明当前 Go module 的身份。
// 后续项目内部 import 会以这个 module path 作为前缀。
module github.com/aaron/go-with-ai/hn-agent

// go 表示当前 module 面向的 Go 语言版本。
// 本学习项目的文档解释以 Go 1.27 为基准。
go 1.27
```

## CLI 入口

```go
// hn-agent/cmd/hnctl/main.go
// package main 表示这个 package 会编译成可执行程序。
// 如果不是 main package，Go 通常会把它当成可被 import 的库。
package main

import (
	// fmt 用来格式化输出，例如 Println、Printf。
	"fmt"
	// os 用来访问操作系统能力，这里用 os.Args 和 os.Exit。
	"os"

	// 这是项目内部 package。
	// import 路径来自 go.mod 里的 module path 加上目录路径。
	"github.com/aaron/go-with-ai/hn-agent/internal/hn"
)

// const 定义常量。
// version 暂时写死，后面阶段可以再改成构建时注入。
const version = "0.1.0"

// run 把真正的 CLI 逻辑从 main 里拆出来。
// 这样以后可以对 run 写测试，而不是直接测试 os.Exit。
//
// args 是命令行参数，但不包含程序名本身。
// 返回 int 是进程退出码：0 通常表示成功，非 0 表示失败。
func run(args []string) int {
	// len(args) > 0 先确认至少有一个参数，避免访问 args[0] 越界。
	// 这里暂时只支持 version 和 -v。
	if len(args) > 0 && (args[0] == "version" || args[0] == "-v") {
		fmt.Println("hnctl version is", version)
		return 0
	}

	// 创建一个 Story 结构体值。
	// 字段名写在左边，值写在右边，这是 Go 的复合字面量写法。
	story := hn.Story{
		ID:    1,
		Title: "Hello",
		URL:   "https://www.example.com",
		Score: 0,
	}

	// %s 用来格式化 string。
	fmt.Printf("hnctl %s\n", version)
	// %v 用来按默认形式打印任意值，这里会打印 true 或 false。
	fmt.Printf("story valid: %v\n", hn.IsValidStory(story))
	return 0
}

// main 是 Go 可执行程序的入口函数。
func main() {
	// os.Args[0] 是程序名，os.Args[1:] 才是真正传给 CLI 的参数。
	// os.Exit 会用 run 返回的退出码结束进程。
	os.Exit(run(os.Args[1:]))
}
```

## Story 领域模型

```go
// hn-agent/internal/hn/story.go
// package hn 表示这个文件属于 hn package。
// 同一个目录下的 .go 文件通常属于同一个 package。
package hn

// Story 是 Hacker News story 在项目里的最小领域模型。
// struct 是一组字段的集合，类似 TypeScript 里的 interface 加运行时数据形状。
type Story struct {
	// ID 是 story 的唯一标识。
	ID int
	// Title 是标题。阶段一先用它判断 story 是否有效。
	Title string
	// URL 是 story 指向的外部链接。
	URL string
	// Score 是 HN 上的分数。
	Score int
}

// IsValidStory 是一个可测试的领域函数。
// 它不依赖网络、文件、命令行，所以阶段一适合先用它练习 go test。
func IsValidStory(story Story) bool {
	// 阶段一先定义最小规则：
	// 有有效 ID，并且标题不为空，就认为是有效 story。
	return story.ID > 0 && story.Title != ""
}
```

## Story 测试

```go
// hn-agent/internal/hn/story_test.go
// 测试文件通常和被测试代码使用同一个 package。
package hn

// testing 是 Go 标准库测试包。
// go test 会自动寻找 TestXxx 形式的测试函数。
import "testing"

// TestIsValidStory 测试 IsValidStory。
// 测试函数必须以 Test 开头，并接收 *testing.T。
func TestIsValidStory(t *testing.T) {
	// tests 是一张测试表。
	// []struct 表示“匿名 struct 的切片”，适合保存多组输入和期望输出。
	tests := []struct {
		name string
		story Story
		want bool
	}{
		{
			name: "valid story",
			story: Story{
				ID:    1,
				Title: "a valid story",
			},
			want: true,
		},
		{
			name: "missing id and title",
			story: Story{
				ID:    0,
				Title: "",
			},
			want: false,
		},
	}

	// range 遍历测试表。
	// tt 是当前测试用例。
	for _, tt := range tests {
		// t.Run 会创建子测试。
		// 测试失败时，输出里能看到具体是哪一个 case 失败。
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidStory(tt.story)
			if got != tt.want {
				// Fatalf 会标记当前子测试失败，并停止这个子测试。
				t.Fatalf("IsValidStory() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

## 验收命令

```bash
cd hn-agent
go test ./internal/hn
go run ./cmd/hnctl
go run ./cmd/hnctl version
```

