package hn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTopStories(t *testing.T) {
	// 启动一个本地 http 服务器，模拟 hacker news 的 API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求路径是否符合预期
		if r.URL.Path != "/topstories.json" {
			t.Fatalf("expected path /topstories.json, got %s", r.URL.Path)
		}
		// 返回模拟的响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[1, 2, 3]`))
	}))
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}

	ids, err := client.TopStories(context.Background())
	if err != nil {
		t.Fatalf("TopStories returned error: %v", err)
	}
	if len(ids) != 3 {
		t.Fatalf("expected 3 ids, got %d", len(ids))
	}
}
