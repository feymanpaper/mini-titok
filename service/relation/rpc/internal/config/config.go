package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/common/consul"
)

type Config struct {
	zrpc.RpcServerConf
	CacheRedis cache.CacheConf
	Mysql      struct {
		DataSource string
	}
	Consul       consul.Conf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
