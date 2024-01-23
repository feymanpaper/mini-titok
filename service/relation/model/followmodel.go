package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"strings"
	"time"
)

var _ FollowModel = (*customFollowModel)(nil)

type (
	// FollowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowModel.
	FollowModel interface {
		followModel
		FindBloomIsFollow(ctx context.Context, fromId int64, toId int64) (bool, error)
		AddBloomFollow(ctx context.Context, fromId int64, toId int64) error
		FindCacheFollowPairListByFollowTime(ctx context.Context, uid, cursor, ps int64) ([]*FollowPair, error)
		FindDBFollowPairListByFollowTime(ctx context.Context, fromId int64, followTime string, limit int) ([]*FollowPair, error)
		AddCacheFollowPairList(ctx context.Context, userId int64, followPairList []*FollowPair) error
	}

	customFollowModel struct {
		*defaultFollowModel
		redisConn   *redis.Redis
		bloomFilter *bloom.Filter
	}
)

// NewFollowModel returns a model for the database table.
func NewFollowModel(conn sqlx.SqlConn, c cache.CacheConf) FollowModel {
	return &customFollowModel{
		defaultFollowModel: newFollowModel(conn),
		redisConn:          NewRedisConn(c),
		bloomFilter:        NewBloomFilter(c),
	}
}

func (m *customFollowModel) formatFollowKey(fromId int64, toId int64) string {
	fromIdStr := strconv.FormatInt(fromId, 10)
	toIdStr := strconv.FormatInt(toId, 10)
	key := fromIdStr + ":" + toIdStr
	return key
}
func (m *customFollowModel) FindBloomIsFollow(ctx context.Context, fromId int64, toId int64) (bool, error) {
	key := m.formatFollowKey(fromId, toId)
	exists, err := m.bloomFilter.Exists([]byte(key))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *customFollowModel) AddBloomFollow(ctx context.Context, fromId int64, toId int64) error {
	key := m.formatFollowKey(fromId, toId)
	err := m.bloomFilter.Add([]byte(key))
	if err != nil {
		return err
	}
	return nil
}

const (
	prefixFollowList = "followList:%d"
	FollowListExpire = 3600
)

func (m *customFollowModel) formatFollowListKey(uid int64) string {
	return fmt.Sprintf(prefixFollowList, uid)
}

type FollowPair struct {
	ToId       int64     `db:"to_id"`
	CreateTime time.Time `db:"create_time"`
}

func (m *customFollowModel) FindDBFollowPairListByFollowTime(ctx context.Context, fromId int64, followTime string, limit int) ([]*FollowPair, error) {
	var (
		err        error
		sql        string
		followList []*FollowPair
	)
	selectFieldsStr := []string{"to_id", "create_time"}
	selectFields := strings.Join(selectFieldsStr, ",")
	sql = fmt.Sprintf("select %s from %s where from_id=? and create_time < ? order by %s desc limit ?", followRows, m.table, selectFields)
	err = m.conn.QueryRowsCtx(ctx, &followList, sql, fromId, followTime, limit)
	if err != nil {
		return nil, err
	}
	return followList, nil
}

func (m *customFollowModel) FindCacheFollowPairListByFollowTime(ctx context.Context, uid, cursor, ps int64) ([]*FollowPair, error) {
	key := m.formatFollowListKey(uid)
	b, err := m.redisConn.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b {
		err = m.redisConn.ExpireCtx(ctx, key, FollowListExpire)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}
	pairs, err := m.redisConn.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(ps))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimit key: %s error: %v", key, err)
		return nil, err
	}
	var followList []*FollowPair
	for _, pair := range pairs {
		toId, err := strconv.ParseInt(pair.Key, 10, 64)
		followTime := time.Unix(pair.Score, 0)
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		followList = append(followList, &FollowPair{
			ToId:       toId,
			CreateTime: followTime,
		})
	}
	return followList, nil
}

func (m *customFollowModel) AddCacheFollowPairList(ctx context.Context, userId int64, followPairList []*FollowPair) error {
	if len(followPairList) == 0 {
		return nil
	}
	key := m.formatFollowListKey(userId)
	for _, followPair := range followPairList {
		var score int64
		if followPair.ToId != -1 {
			score = followPair.CreateTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := m.redisConn.ZaddCtx(ctx, key, score, strconv.Itoa(int(followPair.ToId)))
		if err != nil {
			return err
		}
	}
	return m.redisConn.ExpireCtx(ctx, key, FollowListExpire)
}