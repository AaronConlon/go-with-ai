package hn

import (
	"context"
	// fmt 在实现里用来创建 error，在测试里用 fmt.Fprint 写 response body。
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试并发为 0 时应该报错

func TestZeroConcurrency(t *testing.T) {
	client := NewClient()

	items, err := client.FetchItems(context.Background(), []int64{101, 102}, 0)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if items != nil {
		t.Fatalf("expected nil items, got %#v", items)
	}
}

func TestFetchBatchItemsAllSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch req.URL.Path {
		case "/item/101.json":
			_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
		case "/item/102.json":
			_, _ = fmt.Fprint(w, `{"id":102,"type":"story","title":"second"}`)
		case "/item/103.json":
			_, _ = fmt.Fprint(w, `{"id":103,"type":"story","title":"third"}`)
		default:
			http.NotFound(w, req)
		}
	}))
	defer server.Close()

	// 开启服务器之后，测试 client 指向本地
	client := &Client{
		BaseURL: server.URL,
		HTTP:    server.Client(),
	}

	// 开启并发量
	const concurrency = 2
	items, err := client.FetchItems(context.Background(), []int64{101, 102, 103}, concurrency)
	if err != nil {
		// 说明有错误
		t.Fatalf("FetchItems returned error: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	// 保证并发完成返回的顺序
	wantIDs := []int64{101, 102, 103}

	for idx, wantID := range wantIDs {
		if items[idx].ID != wantID {
			t.Fatalf("items[%d].ID = %d, want %d,", idx, items[idx].ID, wantID)
		}
	}
}

func TestFetchItemsReturnsErrorWhenOneItemFails(t *testing.T) {
	// 这个测试模拟“其中一个 item 请求失败”。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/item/102.json" {
			http.Error(w, "Ops", http.StatusInternalServerError)
			return
		}

		// 其他请求直接返回
		w.Header().Set("Content-Type", "application/json")
		// 使用 fmt 写入 res
		fmt.Fprint(w, `{"id":101,"type":"story","title":"ok"}`)
	}))

	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:    server.Client(),
	}

	// 任意失败，整体失败
	items, err := client.FetchItems(context.Background(), []int64{101, 102}, 2)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if items != nil {
		t.Fatalf("expected nil items on error, got %#v", items)
	}
}
