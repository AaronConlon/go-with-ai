package main

import (
	// 格式化输出
	"fmt"
	// 访问系统
	"os"
	// 内部包, import 路径来自 go.mod 里的 module path
	"github.com/aaron/go-with-ai/hn-agent/internal/hn"
)

// 定义不可变常量
const version = "0.1.0"

func run(args []string) int {
	if len(args) > 0 && (args[0] == "version" || args[0] == "-v") {
		fmt.Println("hnctl version is", version)
		return 0
	}

	story := hn.Story{
		ID:    1,
		Title: "Hello",
		URL:   "https://www.example.com",
		Score: 0,
	}

	fmt.Printf("hnctl %s\n", version)
	// 打印是否是有效结构
	fmt.Printf("story valid: %v\n", hn.IsValidStory(story))
	return 0
}

// 经典的程序入口
func main() {
	os.Exit(run(os.Args[1:]))
}
