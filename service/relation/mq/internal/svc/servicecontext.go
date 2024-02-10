package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"mini-titok/service/relation/model"
	"mini-titok/service/relation/mq/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	FollowModel model.FollowModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		FollowModel: model.NewFollowModel(conn, c.CacheRedis),
	}
}
