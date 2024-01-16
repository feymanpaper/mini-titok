package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"mini-titok/common/consul"
	"mini-titok/common/xcode"
	"mini-titok/service/api/internal/config"
	"mini-titok/service/api/internal/handler"
	"mini-titok/service/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/titok_api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}

	// 自定义错误处理方法
	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
