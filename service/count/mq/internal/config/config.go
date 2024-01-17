package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	KqConsumerConf kq.KqConf
	Mysql          struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
}
