package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
	"mini-titok/service/relation/model"
	"mini-titok/service/relation/rpc/internal/config"
)

type ServiceContext struct {
	Config            config.Config
	FollowModel       model.FollowModel
	SingleFlightGroup singleflight.Group
	KqPusherClient    *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		FollowModel:    model.NewFollowModel(conn, c.CacheRedis),
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
