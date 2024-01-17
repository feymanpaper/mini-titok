package logic

import (
	"context"

	"mini-titok/service/count/rpc/countclient"
	"mini-titok/service/count/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerCountLogic {
	return &GetFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerCountLogic) GetFollowerCount(in *countclient.GetFollowerCountRequest) (*countclient.GetFollowerCountResponse, error) {
	followerCount, err := l.svcCtx.CountModel.FindCahceorDBFollowerCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &countclient.GetFollowerCountResponse{
		FollowerCount: followerCount,
	}, nil
}
