package model

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func NewSqlConn() sqlx.SqlConn {
	dsn := "root:351681578wdp@tcp(127.0.0.1:3306)/titok_count?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	conn := sqlx.NewSqlConn("mysql", dsn)
	return conn
}

func NewRedisConn() *redis.Redis {
	conf := redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
		Pass: "",
		Tls:  false,
	}
	return redis.MustNewRedis(conf)
}
