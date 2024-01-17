package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

var _ CountModel = (*customCountModel)(nil)

type (
	// CountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCountModel.
	CountModel interface {
		countModel
		IncCacheFollowerCount(ctx context.Context, id int64) error
		DecCacheFollowerCount(ctx context.Context, id int64) error
		FindCahceorDBFollowerCount(ctx context.Context, id int64) (int64, error)
		InsertDBFollowerCount(ctx context.Context, id int64) error
		IncDBFollowerCount(ctx context.Context, id int64) error
	}

	customCountModel struct {
		*defaultCountModel
		sqlConn   sqlx.SqlConn
		redisConn *redis.Redis
		barrier   syncx.SingleFlight
	}
)

// NewCountModel returns a model for the database table.
func NewCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CountModel {
	return &customCountModel{
		defaultCountModel: newCountModel(conn, c, opts...),
		sqlConn:           conn,
		redisConn:         NewRedisConn(),
		barrier:           syncx.NewSingleFlight(),
	}
}

var cacheFollowerCountIdPrefix = "follower:"

func (m *customCountModel) FindCahceorDBFollowerCount(ctx context.Context, id int64) (int64, error) {
	countkey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, id)
	var resp int64
	countvalStr := "count_val"
	err := m.QueryRowCtx(ctx, &resp, countkey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `count_key` = ? limit 1", countvalStr, m.table)
		return conn.QueryRowCtx(ctx, v, query, countkey)
	})
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return NotFoundCount, ErrNotFound
	default:
		return NotFoundCount, err
	}
}

func (m *customCountModel) IncCacheFollowerCount(ctx context.Context, id int64) error {
	// 先把数据从db读到缓存中来
	_, err := m.FindCahceorDBFollowerCount(ctx, id)
	if err != nil {
		return err
	}
	countkey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, id)
	_, err = m.redisConn.Incr(countkey)
	if err != nil {
		return err
	}
	return nil
}

func (m *customCountModel) DecCacheFollowerCount(ctx context.Context, id int64) error {
	// 先把数据从db读到缓存中来
	_, err := m.FindCahceorDBFollowerCount(ctx, id)
	if err != nil {
		return err
	}
	countkey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, id)
	_, err = m.redisConn.Decr(countkey)
	if err != nil {
		return err
	}
	return nil
}

func (m *customCountModel) InsertDBFollowerCount(ctx context.Context, id int64) error {
	countkey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, id)
	query := fmt.Sprintf("insert into %s (`count_key`, `count_val`) values (?, ?)", m.table)
	_, err := m.sqlConn.ExecCtx(ctx, query, countkey, 0)
	if err != nil {
		return err
	}
	return nil
}

func (m *customCountModel) IncDBFollowerCount(ctx context.Context, id int64) error {
	countkey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, id)
	query := fmt.Sprintf("update %s set `count_val` = `count_val` + 1 where `count_key` = ?", m.table)
	_, err := m.sqlConn.ExecCtx(ctx, query, countkey)
	if err != nil {
		return err
	}
	return nil
}
