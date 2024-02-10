package main

import (
	"context"
	"flag"
	"fmt"
	"mini-titok/service/relation/mq/internal/config"
	"mini-titok/service/relation/mq/internal/logic"
	"mini-titok/service/relation/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/relation_mq.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	for _, mq := range logic.Consumers(ctx, svcCtx) {
		serviceGroup.Add(mq)
	}
	fmt.Printf("relation mq server start...\n")
	serviceGroup.Start()
}
