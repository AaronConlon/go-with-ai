package config

import "os"

// 配置对象结构体
type Config struct {
	Addr string
}

func Load() Config {
	addr := os.Getenv("HN_AGENT_ADDR")

	if addr == "" {
		addr = ":8080"
	}

	return Config{
		Addr: addr,
	}
}
