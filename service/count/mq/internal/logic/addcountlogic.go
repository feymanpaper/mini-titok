package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"mini-titok/service/count/mq/internal/svc"
	"mini-titok/service/count/mq/internal/types"
)

type AddCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCountLogic {
	return &AddCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCountLogic) Consume(key, val string) error {
	fmt.Printf("get key: %s val: %s\n", key, val)
	msg := types.AddCountMsg{}
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		return err
	}
	err = l.svcCtx.CountModel.IncDBFollowerCount(l.ctx, msg.UserId)
	if err != nil {
		return err
	}
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewAddCountLogic(ctx, svcCtx)),
	}
}
