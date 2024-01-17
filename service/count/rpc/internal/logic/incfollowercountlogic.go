package logic

import (
	"context"

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
	return &countclient.IncFollowerCountResponse{}, nil
}
