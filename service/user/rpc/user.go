package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"mini-titok/common/consul"
	"mini-titok/common/interceptors"
	"mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/conf"
	_ "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"mini-titok/service/user/rpc/internal/config"
	"mini-titok/service/user/rpc/internal/server"
	"mini-titok/service/user/rpc/internal/svc"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		userclient.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())
	defer s.Stop()

	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
