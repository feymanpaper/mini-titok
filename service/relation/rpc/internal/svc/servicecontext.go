package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/sync/singleflight"
	"mini-titok/common/interceptors"
	"mini-titok/service/relation/model"
	"mini-titok/service/relation/rpc/internal/config"
	"mini-titok/service/user/rpc/user"
)

type ServiceContext struct {
	Config            config.Config
	FollowModel       model.FollowModel
	UserRpc           user.User
	SingleFlightGroup singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		FollowModel: model.NewFollowModel(conn, c.CacheRedis),
		UserRpc:     user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))),
	}
}
