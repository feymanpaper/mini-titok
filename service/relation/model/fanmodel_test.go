package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"testing"
)

func TestCustomFanModel_AddCacheFanCountHash(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(GoroutineNum)
	for i := 0; i < GoroutineNum; i++ {
		go func(i int) {
			defer wg.Done()
			for j := i*OneGoroutineTaskNum + 1; j <= (i+1)*OneGoroutineTaskNum; j++ {
				err := TestFanModel.AddCacheFanCountHash(context.Background(), int64(j))
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

func TestCustomFanModel_GetCacheFanCountHash(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(GoroutineNum)
	for i := 0; i < GoroutineNum; i++ {
		go func(i int) {
			defer wg.Done()
			for j := i*OneGoroutineTaskNum + 1; j <= (i+1)*OneGoroutineTaskNum; j++ {
				cnt, err := TestFanModel.GetCacheFanCountHash(context.Background(), 2)
				if err != nil {
					logx.Error("error")
				}
				logx.Info(cnt)
			}
			logx.Info("goroutine finish")
		}(i)
	}
	wg.Wait()
	logx.Info("finish")
}
