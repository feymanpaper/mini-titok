package logic

import (
	"context"

	"mini-titok/service/count/rpc/countclient"
	"mini-titok/service/count/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InsertUserCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInsertUserCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertUserCountLogic {
	return &InsertUserCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// InsertUserCount 注册时初始化用户的count数据
func (l *InsertUserCountLogic) InsertUserCount(in *countclient.InsertUserCountRequest) (*countclient.InsertUserCountResponse, error) {
	err := l.svcCtx.CountModel.InsertDBFollowerCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &countclient.InsertUserCountResponse{}, nil
}
