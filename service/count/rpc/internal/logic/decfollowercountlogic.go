package logic

import (
	"context"

	"mini-titok/service/count/rpc/countclient"
	"mini-titok/service/count/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecFollowerCountLogic {
	return &DecFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecFollowerCountLogic) DecFollowerCount(in *countclient.DecFollowerCountRequest) (*countclient.DecFollowerCountResponse, error) {
	err := l.svcCtx.CountModel.DecCacheFollowerCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &countclient.DecFollowerCountResponse{}, nil
}
