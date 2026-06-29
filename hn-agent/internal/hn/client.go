package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Item struct {
	// 唯一标识符，标识 ID 对应的是 json 里的 id
	ID int64 `json:"id"`

	// Type 表示项目的类型，如 story、comment、job 等
	Type string `json:"type"`

	// By 表示项目的作者
	By string `json:"by"`

	// Time 表示项目的发布时间
	Time int64 `json:"time"`

	// URL 表示项目的 URL
	URL string `json:"url"`

	// Score 表示项目的评分
	Score int `json:"score"`

	// 标题
	Title string `json:"title"`

	// Descendants 表示评论数量
	Descendants int `json:"descendants"`
}

type Client struct {
	// BaseURL 表示 HN API 的基础 URL
	BaseURL string

	// 标准库客户端
	HTTP *http.Client
}

// 返回一个客户端指针
func NewClient() *Client {
	// 使用 & 取一个实例化结构体的指针
	return &Client{
		BaseURL: "https://hacker-news.firebaseio.com/v0",
		HTTP: &http.Client{
			// Timeout 是整个请求的最大耗时。
			// 外部 I/O 不设置 timeout，程序可能一直卡住。
			Timeout: 10 * time.Second,
		},
	}
}

// 实现获取 topStories 数据（返回的是 ids）
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
	// 创建一个带有 context 的 https request（并没有发起请求，仅仅是构建 request 作准备）
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/topstories.json", nil)
	if err != nil {
		// 创建失败 request
		return nil, err
	}

	// 发起请求，让客户端执行这个准备好的 request
	resp, err := c.HTTP.Do(req)
	if err != nil {
		// nil 就是错误
		return nil, err
	}

	// 确保响应体被关闭
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 声明变量，未赋值，默认是空或0
	var ids []int64
	// 解构并且将值传到 ids，如果出错，则返回错误，否则就成功将返回值给到 ids
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

// 基于 id 获取单条详情
func (c *Client) Item(ctx context.Context, id int64) (Item, error) {
	// 创建完整 url
	url := fmt.Sprintf("%s/item/%d.json", c.BaseURL, id)
	// 创建 request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return Item{}, err
	}

	// 发起请求
	resp, err := c.HTTP.Do(req)

	if err != nil {
		return Item{}, err
	}

	// 确保关闭
	defer resp.Body.Close()

	// 继续解析状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Item{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 解析 json 并且赋值

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return Item{}, err
	}

	return item, nil
}
