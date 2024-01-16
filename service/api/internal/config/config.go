package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/common/consul"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc zrpc.RpcClientConf
	Consul  consul.Conf
}
