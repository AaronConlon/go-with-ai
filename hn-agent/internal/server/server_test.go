package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := NewMux()

	req := httptest.NewRequest(http.MethodGet, PathHealthz, nil)

	// 创建一颗测试 response record
	// handler 写给 w 的状态、头、body 都会被记录
	rec := httptest.NewRecorder()

	// ServeHTTP 手动把请求交给 Mux 处理
	// 当 mux 接收到请求，分给对应函数之后，响应处理的写入会同步到 record 里
	mux.ServeHTTP(rec, req)

	// 拿到结果
	resp := rec.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

}
