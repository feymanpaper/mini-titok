package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/service/relation/rpc/relation"
	"sync"
	"testing"
)

// qps 5000
func Test_Action(t *testing.T) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "dns:///127.0.0.1:8082",
	})
	client := relation.NewRelation(conn)
	wg := sync.WaitGroup{}
	goroutineNum := 100
	oneGoroutineTaskNum := 1000
	wg.Add(goroutineNum)
	ctx := context.Background()
	for i := 0; i < goroutineNum; i++ {
		go func(i int) {
			defer wg.Done()
			for j := i*oneGoroutineTaskNum + 1; j <= (i+1)*oneGoroutineTaskNum; j++ {
				_, err := client.Action(ctx, &relation.ActionRequest{
					FromUserId: int64(j),
					ToUserId:   1,
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
