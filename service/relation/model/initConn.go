package model

import (
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func NewBloomFilter(c cache.CacheConf) *bloom.Filter {
	store := redis.MustNewRedis(c[0].RedisConf)
	filter := bloom.New(store, "isFollowFilter", 1<<20)
	return filter
}

func NewRedisConn(c cache.CacheConf) *redis.Redis {
	redisConn := redis.MustNewRedis(c[0].RedisConf)
	return redisConn
}
