package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/common/interceptors"
	"mini-titok/service/relation/rpc/relation"
	"mini-titok/service/user/model"
	"mini-titok/service/user/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	UserModel   model.UserModel
	RelationRpc relation.Relation
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		UserModel:   model.NewUserModel(conn, c.CacheRedis),
		RelationRpc: relation.NewRelation(zrpc.MustNewClient(c.RelationRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))),
	}
}
