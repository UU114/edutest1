package svc

import (
	"ai-education-platform/services/course/course/internal/config"
	"ai-education-platform/services/course/course/internal/models"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config        config.Config
	DB            sqlx.SqlConn
	Redis         *redis.Redis
	CourseModel   *models.CourseModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	r := redis.MustNewRedis(c.Redis.Host, c.Redis.Pass, c.Redis.Key)
	
	return &ServiceContext{
		Config:      c,
		DB:          conn,
		Redis:       r,
		CourseModel: models.NewCourseModel(conn),
	}
}