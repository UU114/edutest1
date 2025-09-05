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
	Exam struct {
		DefaultTimeLimit int     `json:"default_time_limit"` // 默认考试时间限制(分钟)
		MaxQuestionCount int     `json:"max_question_count"` // 最大题目数量
		AllowedQuestionTypes []string `json:"allowed_question_types"` // 允许的题目类型
		ScorePrecision   float64 `json:"score_precision"`   // 分数精度
	}
}