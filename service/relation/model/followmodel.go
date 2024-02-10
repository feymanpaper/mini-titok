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
		// isfollow
		AddBloomIsFollow(ctx context.Context, fromId int64, toId int64) error
		AddDBIsFollow(ctx context.Context, fromId, toId int64, time time.Time) error
		FindBloomOrDBIsFollow(ctx context.Context, fromId, toId int64) (bool, error)
		// follow_list
		AddCacheFollowPair(ctx context.Context, userId int64, followPair *FollowPair) error
		FindCacheFollowPairListByFollowTime(ctx context.Context, uid, cursor, ps int64) ([]*FollowPair, error)
		FindDBFollowPairListByFollowTime(ctx context.Context, fromId int64, followTime string, limit int) ([]*FollowPair, error)
		AddCacheFollowPairList(ctx context.Context, userId int64, followPairList []*FollowPair) error
		// follow_count
		AddCacheFollowingCountHash(ctx context.Context, userId int64) error
		GetCacheFollowingCountHash(ctx context.Context, userId int64) (int64, error)
		//fan_list
		AddCacheFanPair(ctx context.Context, userId int64, followPair *FollowPair) error
		FindDBFanPairListByFollowTime(ctx context.Context, fromId int64, followTime string, limit int) ([]*FollowPair, error)
		FindCacheFanPairListByFollowTime(ctx context.Context, uid, cursor, ps int64) ([]*FollowPair, error)
		AddCacheFanPairList(ctx context.Context, userId int64, followPairList []*FollowPair) error
		// fan_count
		AddCacheFanCountHash(ctx context.Context, userId int64) error
		GetCacheFanCountHash(ctx context.Context, userId int64) (int64, error)
	}

	customFollowModel struct {
		*defaultFollowModel
		redisConn   *redis.Redis
		bloomFilter *bloom.Filter
		followTable string
		fanTable    string
	}
)

// NewFollowModel returns a model for the database table.
func NewFollowModel(conn sqlx.SqlConn, c cache.CacheConf) FollowModel {
	return &customFollowModel{
		defaultFollowModel: newFollowModel(conn),
		redisConn:          NewRedisConn(c),
		bloomFilter:        NewBloomFilter(c),
		followTable:        "`follow`",
		fanTable:           "`fan`",
	}
}

func (m *customFollowModel) findBloomIsFollow(ctx context.Context, fromId int64, toId int64) (bool, error) {
	key := formatFollowKey(fromId, toId)
	exists, err := m.bloomFilter.Exists([]byte(key))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *customFollowModel) AddBloomIsFollow(ctx context.Context, fromId int64, toId int64) error {
	key := formatFollowKey(fromId, toId)
	err := m.bloomFilter.Add([]byte(key))
	if err != nil {
		return err
	}
	return nil
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
	selectFieldsStr := []string{"create_time", "to_id"}
	selectFields := strings.Join(selectFieldsStr, ",")
	sql = fmt.Sprintf("select %s from %s where from_id=? and create_time < ? order by %s desc limit ?", followRows, m.followTable, selectFields)
	err = m.conn.QueryRowsCtx(ctx, &followList, sql, fromId, followTime, limit)
	if err != nil {
		return nil, err
	}
	return followList, nil
}

func (m *customFollowModel) FindDBFanPairListByFollowTime(ctx context.Context, fromId int64, followTime string, limit int) ([]*FollowPair, error) {
	var (
		err        error
		sql        string
		followList []*FollowPair
	)
	selectFieldsStr := []string{"create_time", "to_id"}
	selectFields := strings.Join(selectFieldsStr, ",")
	sql = fmt.Sprintf("select %s from %s where from_id=? and create_time < ? order by %s desc limit ?", followRows, m.fanTable, selectFields)
	err = m.conn.QueryRowsCtx(ctx, &followList, sql, fromId, followTime, limit)
	if err != nil {
		return nil, err
	}
	return followList, nil
}

func (m *customFollowModel) FindCacheFollowPairListByFollowTime(ctx context.Context, uid, cursor, ps int64) ([]*FollowPair, error) {
	key := formatFollowListKey(uid)
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

func (m *customFollowModel) FindCacheFanPairListByFollowTime(ctx context.Context, uid, cursor, ps int64) ([]*FollowPair, error) {
	key := formatFanListKey(uid)
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
	key := formatFollowListKey(userId)
	for _, followPair := range followPairList {
		var score int64
		if followPair.ToId != -1 {
			score = followPair.CreateTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := m.redisConn.Zadd(key, score, strconv.FormatInt(followPair.ToId, 10))
		if err != nil {
			return err
		}
	}
	return m.redisConn.Expire(key, FollowListExpire)
}

func (m *customFollowModel) AddCacheFollowPair(ctx context.Context, userId int64, followPair *FollowPair) error {
	key := formatFollowListKey(userId)
	_, err := m.redisConn.Zadd(key, followPair.CreateTime.Local().Unix(), strconv.FormatInt(followPair.ToId, 10))
	logx.Error(err)
	return err
}

func (m *customFollowModel) AddCacheFanPairList(ctx context.Context, userId int64, followPairList []*FollowPair) error {
	if len(followPairList) == 0 {
		return nil
	}
	key := formatFanListKey(userId)
	for _, followPair := range followPairList {
		var score int64
		if followPair.ToId != -1 {
			score = followPair.CreateTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := m.redisConn.Zadd(key, score, strconv.Itoa(int(followPair.ToId)))
		if err != nil {
			return err
		}
	}
	return m.redisConn.Expire(key, FollowListExpire)
}

func (m *customFollowModel) AddCacheFanPair(ctx context.Context, userId int64, followPair *FollowPair) error {
	key := formatFanListKey(userId)
	_, err := m.redisConn.Zadd(key, followPair.CreateTime.Local().Unix(), strconv.FormatInt(followPair.ToId, 10))
	return err
}

func (m *customFollowModel) AddCacheFollowingCountHash(ctx context.Context, userId int64) error {
	idstr := strconv.FormatInt(userId, 10)
	num := hashBucketNum(idstr)
	key := formatFollowCountKeyHash(num)
	_, err := m.redisConn.Hincrby(key, idstr, 1)
	if err != nil {
		return err
	}
	return nil
}

func (m *customFollowModel) GetCacheFollowingCountHash(ctx context.Context, userId int64) (int64, error) {
	idstr := strconv.FormatInt(userId, 10)
	num := hashBucketNum(idstr)
	key := formatFollowCountKeyHash(num)
	followCountStr, err := m.redisConn.Hget(key, idstr)
	if err != nil {
		return 0, err
	}
	followCount, err := strconv.ParseInt(followCountStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return followCount, err
}

func (m *customFollowModel) AddDBIsFollow(ctx context.Context, fromId, toId int64, curTime time.Time) error {
	fieldSet := strings.Join([]string{"from_id", "to_id", "create_time"}, ",")
	followQuery := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.followTable, fieldSet)
	fanQuery := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.fanTable, fieldSet)
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := session.ExecCtx(ctx, followQuery, fromId, toId, curTime)
		if err != nil {
			return err
		}
		_, err = session.ExecCtx(ctx, fanQuery, toId, fromId, curTime)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (m *customFollowModel) findDBIsFollow(ctx context.Context, fromId, toId int64) (bool, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE from_id=? AND to_id=? AND delete_time IS NULL", m.table)
	var ans int64
	err := m.conn.QueryRowCtx(ctx, &ans, query, fromId, toId)
	switch err {
	case nil:
		return true, nil
	case sqlx.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (m *customFollowModel) FindBloomOrDBIsFollow(ctx context.Context, fromId, toId int64) (bool, error) {
	isFollow, err := m.findBloomIsFollow(ctx, fromId, toId)
	if err != nil {
		return false, err
	}
	// bloom filter查不存在是准确的
	if !isFollow {
		return false, nil
	}
	// 如果存在, 则查DB
	isFollow, err = m.findDBIsFollow(ctx, fromId, toId)
	if err != nil {
		return false, err
	}
	return isFollow, nil
}

func (m *customFollowModel) AddCacheFanCountHash(ctx context.Context, userId int64) error {
	idstr := strconv.FormatInt(userId, 10)
	num := hashBucketNum(idstr)
	key := formatFanCountKeyHash(num)
	_, err := m.redisConn.Hincrby(key, idstr, 1)
	if err != nil {
		return err
	}
	return nil
}

func (m *customFollowModel) GetCacheFanCountHash(ctx context.Context, userId int64) (int64, error) {
	idstr := strconv.FormatInt(userId, 10)
	num := hashBucketNum(idstr)
	key := formatFanCountKeyHash(num)
	fanCountStr, err := m.redisConn.Hget(key, idstr)
	if err != nil {
		return 0, err
	}
	fanCount, err := strconv.ParseInt(fanCountStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return fanCount, err
}

//func (m *customFollowModel) AddCacheFollowingCount(ctx context.Context, userId int64) error {
//	key := formatFollowCountKey(userId)
//	_, err := m.redisConn.Incr(key)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func (m *customFollowModel) GetCacheFollowingCount(ctx context.Context, userId int64) (int64, error) {
//	key := formatFollowCountKey(userId)
//	followCountStr, err := m.redisConn.Get(key)
//	if err != nil {
//		return -1, err
//	}
//	followCount, err := strconv.ParseInt(followCountStr, 10, 64)
//	if err != nil {
//		return -1, err
//	}
//	return followCount, nil
//}
