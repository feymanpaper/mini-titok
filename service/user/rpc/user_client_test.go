package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"mini-titok/service/user/rpc/user"
	"testing"
)

// qps 5000
func Test_GetUserInfo(t *testing.T) {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "dns:///127.0.0.1:8080",
	})
	client := user.NewUser(conn)
	info, err := client.GetUserInfo(context.Background(), &user.UserInfoRequest{
		UserId: 1,
	})
	if err != nil {
		logx.Error(err)
	}
	logx.Info(info.UserInfo)
}
