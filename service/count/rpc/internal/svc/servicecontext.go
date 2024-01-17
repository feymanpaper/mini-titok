package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"mini-titok/service/count/model"
	"mini-titok/service/count/rpc/internal/config"
)

type ServiceContext struct {
	Config         config.Config
	CountModel     model.CountModel
	KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		CountModel:     model.NewCountModel(conn, c.CacheRedis),
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
