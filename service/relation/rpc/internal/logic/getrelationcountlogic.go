package logic

import (
	"context"

	"mini-titok/service/relation/rpc/internal/svc"
	"mini-titok/service/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationCountLogic {
	return &GetRelationCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationCountLogic) GetRelationCount(in *relationclient.GetRelationCountRequest) (*relationclient.GetRelationCountResponse, error) {
	followCount, err := l.svcCtx.FollowModel.GetCacheFollowingCountHash(l.ctx, in.UserId)
	if err != nil {
		logx.Error(err)
	}
	fanCount, err := l.svcCtx.FollowModel.GetCacheFanCountHash(l.ctx, in.UserId)
	if err != nil {
		logx.Error(err)
	}
	return &relationclient.GetRelationCountResponse{
		FollowCount: followCount,
		FanCount:    fanCount,
	}, nil
}
