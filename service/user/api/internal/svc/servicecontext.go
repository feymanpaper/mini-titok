package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/service/user/api/internal/config"
	"mini-titok/service/user/rpc/user"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
