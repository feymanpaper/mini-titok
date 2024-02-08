package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"sync"
	"testing"
	"time"
)

var TestFollowModel FollowModel
var TestFanModel FanModel

const GoroutineNum = 100
const OneGoroutineTaskNum = 100

func TestMain(m *testing.M) {
	fmt.Println("Test Init...")
	sqlconn := sqlx.NewMysql("root:351681578wdp@tcp(127.0.0.1:3306)/titok_relation?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	// 创建 RedisConf 实例
	redisConf := redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
		Pass: "",
		Tls:  false,
	}
	// 创建 NodeConf 实例
	nodeConf := cache.NodeConf{
		RedisConf: redisConf,
		Weight:    100,
	}

	// 创建 ClusterConf 实例
	clusterConf := cache.ClusterConf{nodeConf}
	TestFollowModel = NewFollowModel(sqlconn, clusterConf)
	TestFanModel = NewFanModel(sqlconn, clusterConf)
	m.Run()
}

func TestFollowModelBloom(t *testing.T) {
	fmt.Println(TestFollowModel.FindBloomIsFollow(context.Background(), 1, 2))
	fmt.Println(TestFollowModel.AddBloomIsFollow(context.Background(), 1, 2))
	fmt.Println(TestFollowModel.FindBloomIsFollow(context.Background(), 1, 2))
	fmt.Println(TestFollowModel.FindBloomIsFollow(context.Background(), 3, 4))
	fmt.Println(TestFollowModel.FindBloomIsFollow(context.Background(), 5, 6))
	fmt.Println(TestFollowModel.FindBloomIsFollow(context.Background(), 6, 7))
}

func TestCustomFollowModel_FindFollowIdListByFollowTime(t *testing.T) {
	followList, err := TestFollowModel.FindDBFollowPairListByFollowTime(context.Background(), 3, "2025-01-23 10:52:22", 20)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println(followList[0].CreateTime)
}

func TestCustomFollowModel_FindCacheFollowPairListByFollowTime(t *testing.T) {
	cursor := time.Now().Unix()
	followList, err := TestFollowModel.FindCacheFollowPairListByFollowTime(context.Background(), 3, cursor, 20)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	for _, followInfo := range followList {
		fmt.Println(*followInfo)
	}
}

func TestCustomFollowModel_AddCacheFollowPairList(t *testing.T) {
	var followPairList []*FollowPair
	followPairList = append(followPairList, &FollowPair{
		ToId:       11,
		CreateTime: time.Now(),
	})
	followPairList = append(followPairList, &FollowPair{
		ToId:       12,
		CreateTime: time.Now(),
	})
	followPairList = append(followPairList, &FollowPair{
		ToId:       13,
		CreateTime: time.Now(),
	})
	err := TestFollowModel.AddCacheFollowPairList(context.Background(), 3, followPairList)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func TestCustomFollowModel_AddCacheFollowingCountHash(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(GoroutineNum)
	for i := 0; i < GoroutineNum; i++ {
		go func(i int) {
			defer wg.Done()
			for j := i*OneGoroutineTaskNum + 1; j <= (i+1)*OneGoroutineTaskNum; j++ {
				err := TestFollowModel.AddCacheFollowingCountHash(context.Background(), int64(j))
				if err != nil {
					logx.Error("error")
				}
			}
			logx.Info("goroutine finish")
		}(i)
	}
	wg.Wait()
	logx.Info("finish")
}

// qps 50000
func TestCustomFollowModel_GetCacheFollowingCountHash(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(GoroutineNum)
	for i := 0; i < GoroutineNum; i++ {
		go func(i int) {
			defer wg.Done()
			for j := i*OneGoroutineTaskNum + 1; j <= (i+1)*OneGoroutineTaskNum; j++ {
				_, err := TestFollowModel.GetCacheFollowingCountHash(context.Background(), 200009)
				if err != nil {
					logx.Error("error")
				}
			}
			logx.Info("goroutine finish")
		}(i)
	}
	wg.Wait()
	logx.Info("finish")
}

//func TestCustomFollowModel_GetCacheFollowingCount(t *testing.T) {
//	wg := sync.WaitGroup{}
//	wg.Add(GoroutineNum)
//	for i := 0; i < GoroutineNum; i++ {
//		go func(i int) {
//			defer wg.Done()
//			for j := i*OneGoroutineTaskNum + 1; j < (i+1)*OneGoroutineTaskNum; j++ {
//				_, err := TestFollowModel.GetCacheFollowingCount(context.Background(), 1)
//				if err != nil {
//					logx.Error("error")
//				}
//			}
//			logx.Info("goroutine finish")
//		}(i)
//	}
//	wg.Wait()
//	logx.Info("finish")
//}

// 并发协程写redis
//func TestCustomFollowModel_AddCacheFollowingCount(t *testing.T) {
//	wg := sync.WaitGroup{}
//	wg.Add(GoroutineNum)
//	for i := 0; i < GoroutineNum; i++ {
//		go func(i int) {
//			defer wg.Done()
//			for j := i*OneGoroutineTaskNum + 1; j < (i+1)*OneGoroutineTaskNum; j++ {
//				err := TestFollowModel.AddCacheFollowingCount(context.Background(), 1)
//				if err != nil {
//					logx.Error("error")
//				}
//			}
//			logx.Info("goroutine finish")
//		}(i)
//	}
//	wg.Wait()
//	logx.Info("finish")
//}
