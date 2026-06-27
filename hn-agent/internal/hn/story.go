// 声明当前文件属于 hn 包
package hn

// 数据模型 Story 的结构体
type Story struct {
	ID    int
	Title string
	URL   string
	// 评分
	Score int
}

// 可测试的函数，定义参数类型和返回值类型
func IsValidStory(story Story) bool {
	// 存在有效 id 和标题即可
	return story.ID > 0 && story.Title != ""
}
