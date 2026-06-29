package server

import "net/http"

// PathHealthz 是健康检查端点的 URL 路径
const PathHealthz = "/healthz"

func healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// 写状态码
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	// 字符串 "ok\n" 转成 byte 切片，写入 HTTP response body。
	w.Write([]byte("OK\n"))
}

// 创建请求分发器
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(PathHealthz, healthz)

	return mux
}
