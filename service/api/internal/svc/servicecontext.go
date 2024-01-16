package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/common/interceptors"
	"mini-titok/service/user/api/internal/config"
	"mini-titok/service/user/api/internal/middleware"
	"mini-titok/service/user/rpc/user"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.User
	JwtAuth rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))),
		JwtAuth: middleware.NewJwtAuthMiddleware(c).Handle,
	}
}
