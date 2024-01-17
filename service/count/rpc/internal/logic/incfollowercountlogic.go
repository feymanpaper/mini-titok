package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"
	"mini-titok/service/count/rpc/internal/types"

	"mini-titok/service/count/rpc/countclient"
	"mini-titok/service/count/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IncFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIncFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncFollowerCountLogic {
	return &IncFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IncFollowerCountLogic) IncFollowerCount(in *countclient.IncFollowerCountRequest) (*countclient.IncFollowerCountResponse, error) {
	err := l.svcCtx.CountModel.IncCacheFollowerCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	msg := types.AddCountMsg{
		UserId: in.UserId,
	}
	// 发送kafka消息, 异步
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			return
		}
		err = l.svcCtx.KqPusherClient.Push(string(data))
		if err != nil {
			return
		}
	})
	return &countclient.IncFollowerCountResponse{}, nil
}
