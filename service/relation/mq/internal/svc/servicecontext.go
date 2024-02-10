package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"mini-titok/service/count/model"
	"mini-titok/service/count/mq/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	CountModel model.CountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		CountModel: model.NewCountModel(conn, c.CacheRedis),
	}
}
