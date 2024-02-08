package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
)

var _ FanModel = (*customFanModel)(nil)

type (
	// FanModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFanModel.
	FanModel interface {
		fanModel
		AddCacheFanCountHash(ctx context.Context, userId int64) error
		GetCacheFanCountHash(ctx context.Context, userId int64) (int64, error)
	}

	customFanModel struct {
		redisConn *redis.Redis
		*defaultFanModel
	}
)

// NewFanModel returns a model for the database table.
func NewFanModel(conn sqlx.SqlConn, c cache.CacheConf) FanModel {
	return &customFanModel{
		defaultFanModel: newFanModel(conn),
		redisConn:       NewRedisConn(c),
	}
}

func (m *customFanModel) AddCacheFanCountHash(ctx context.Context, userId int64) error {
	idstr := strconv.FormatInt(userId, 10)
	num := hashBucketNum(idstr)
	key := formatFanCountKeyHash(num)
	_, err := m.redisConn.Hincrby(key, idstr, 1)
	if err != nil {
		return err
	}
	return nil
}

func (m *customFanModel) GetCacheFanCountHash(ctx context.Context, userId int64) (int64, error) {
	idstr := strconv.FormatInt(userId, 10)
	num := hashBucketNum(idstr)
	key := formatFanCountKeyHash(num)
	fanCountStr, err := m.redisConn.Hget(key, idstr)
	if err != nil {
		return -1, err
	}
	fanCount, err := strconv.ParseInt(fanCountStr, 10, 64)
	if err != nil {
		return -1, err
	}
	return fanCount, err
}
