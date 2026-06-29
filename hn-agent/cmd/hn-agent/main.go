package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaron/go-with-ai/hn-agent/internal/config"
	"github.com/aaron/go-with-ai/hn-agent/internal/server"
)

// hn-agent/cmd/hn-agent/main.go
func newLogger() *slog.Logger {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 正常情况下 Go 能加载 Asia/Shanghai。
		// 这里保留兜底，避免极端环境缺少时区数据时 logger 创建失败。
		shanghai = time.FixedZone("Asia/Shanghai", 8*60*60)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				t := a.Value.Time().In(shanghai)
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}

			return a
		},
	})

	return slog.New(handler)
}

func main() {
	cfg := config.Load()

	// 使用标准的 slog 包处理日志
	logger := newLogger()

	mux := server.NewMux()

	// 启动 http 服务器
	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	// NotifyContext 会在收到中断信号时取消 ctx。
	// Ctrl+C 通常会触发 os.Interrupt。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// 确保服务关闭，避免端口、进程占用内存
	defer stop()

	// 自动开启一个轻量线程
	go func() {
		logger.Info("Server starting", "addr", cfg.Addr)

		// 错误监听，如果错误是网络关闭，则退出
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 等待停止信号,当监听到错误或者外部信号之后，这里就会执行，否则就阻塞
	<-ctx.Done()

	// 一个退出超时的上下文，专门做优雅退出的，srv 的 shutdown 会接收这个上下文，5 秒超时会导致返回错误，从而触发强制退出
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown fail", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")
}
