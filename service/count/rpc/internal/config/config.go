package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/common/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
	Consul     consul.Conf
}
