package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	rest.RestConf
	DataSource string
	Redis      struct {
		Host string
		Pass string
		Type string
		Key  string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Log struct {
		ServiceName string
		Mode        string
		Level       string
	}
	RateLimit struct {
		Seconds int
		Quota   int
	}
	Prometheus struct {
		Host string
		Port int
		Path string
	}
	AI struct {
		APIKey      string
		APIEndpoint string
		Model       string
		MaxTokens   int
		Temperature float64
	}
}