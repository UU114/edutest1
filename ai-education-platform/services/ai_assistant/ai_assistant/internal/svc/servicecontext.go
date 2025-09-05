package svc

import (
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/config"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/models"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config         config.Config
	DB             sqlx.SqlConn
	Redis          *redis.Redis
	AIModel        *models.AIModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	r := redis.MustNewRedis(c.Redis.Host, c.Redis.Pass, c.Redis.Key)
	
	return &ServiceContext{
		Config:  c,
		DB:      conn,
		Redis:   r,
		AIModel: models.NewAIModel(conn),
	}
}