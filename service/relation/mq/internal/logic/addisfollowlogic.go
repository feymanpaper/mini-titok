package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"mini-titok/service/relation/mq/internal/svc"
	"mini-titok/service/relation/mq/internal/types"
)

type AddIsfollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddIsfollowLogic {
	return &AddIsfollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddIsfollowLogic) Consume(key, val string) error {
	fmt.Printf("get key: %s val: %s\n", key, val)
	var msg types.AddIsFollowMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		return err
	}
	err = l.svcCtx.FollowModel.AddDBIsFollow(l.ctx, msg.FromId, msg.ToId)
	if err != nil {
		return err
	}
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewAddIsFollowLogic(ctx, svcCtx)),
	}
}
