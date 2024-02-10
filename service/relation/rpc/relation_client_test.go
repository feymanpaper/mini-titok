package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/service/relation/rpc/relation"
	"mini-titok/service/relation/rpc/relationclient"
	"sync"
	"testing"
)

// qps 5000
func Test_FollowAction(t *testing.T) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "dns:///127.0.0.1:8082",
	})
	client := relation.NewRelation(conn)
	wg := sync.WaitGroup{}
	goroutineNum := 1
	oneGoroutineTaskNum := 10
	wg.Add(goroutineNum)
	ctx := context.Background()
	for i := 0; i < goroutineNum; i++ {
		go func(i int) {
			defer wg.Done()
			for j := i*oneGoroutineTaskNum + 1; j <= (i+1)*oneGoroutineTaskNum; j++ {
				_, err := client.FollowAction(ctx, &relation.FollowActionRequest{
					FromUserId: int64(j),
					ToUserId:   100,
					ActionType: 1,
				})
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

func Test_SingeFollowAction(t *testing.T) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "dns:///127.0.0.1:8082",
	})
	client := relation.NewRelation(conn)
	_, err := client.FollowAction(context.Background(), &relation.FollowActionRequest{
		FromUserId: 1,
		ToUserId:   2,
		ActionType: 1,
	})
	_, err = client.FollowAction(context.Background(), &relation.FollowActionRequest{
		FromUserId: 2,
		ToUserId:   1,
		ActionType: 1,
	})
	if err != nil {
		logx.Error(err)
	}
}

func Test_FollowList(t *testing.T) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "dns:///127.0.0.1:8082",
	})
	client := relation.NewRelation(conn)
	list, err := client.FollowList(context.Background(), &relationclient.FollowListRequest{
		UserId:   1,
		Cursor:   0,
		PageSize: 20,
		EndId:    0,
	})
	if err != nil {
		logx.Error(err)
	}
	logx.Info(list)
}

func Test_GetRelationCount(t *testing.T) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "dns:///127.0.0.1:8082",
	})
	client := relation.NewRelation(conn)
	count, err := client.GetRelationCount(context.Background(), &relation.GetRelationCountRequest{
		UserId: 50000,
	})
	if err != nil {
		logx.Error(err)
	}
	logx.Info(count.FollowCount, count.FanCount)
}
