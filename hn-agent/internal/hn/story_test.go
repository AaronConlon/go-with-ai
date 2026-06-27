package hn

// 标准库
import "testing"

// 只要函数名以 Test 开头，并接受 *testing.T 类型的参数, go test 会自动调用该函数
// 测试是否是有效 Story
func TestIsValidStory(t *testing.T) {
	// 保存测试用例数组
	// []struct 表示“一个匿名 struct 的切片”，可以理解成一张测试表。
	tests := []struct {
		name     string
		story    Story
		expected bool
	}{
		{
			name: "valid story",
			story: Story{
				ID:    1,
				Title: "a valid story",
			},
			expected: true,
		},
		{
			name: "invalid story",
			story: Story{
				ID:    0,
				Title: "",
			},
			expected: false,
		},
	}

	// 遍历数组进行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行测试，传入名字作为子测试名称，失败时看得清楚
			got := IsValidStory(tt.story)

			// 比较结果
			if got != tt.expected {
				// 打印实际结果和期望结果，便于调试
				//  t.Fatalf 会标记测试失败，并停止当前子测试。
				t.Fatalf("IsValidStory() = %v, want %v", got, tt.expected)
			}
		})
	}
}
