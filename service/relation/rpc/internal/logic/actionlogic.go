package logic

import (
	"context"
	"mini-titok/common/xcode"

	"mini-titok/service/relation/rpc/internal/svc"
	"mini-titok/service/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActionLogic) Action(in *relationclient.ActionRequest) (*relationclient.ActionResponse, error) {
	if in.ActionType == 1 {
		err := l.followAction(in.FromUserId, in.ToUserId)
		if err != nil {
			return &relationclient.ActionResponse{}, err
		}
		return &relationclient.ActionResponse{}, nil
	} else if in.ActionType == 2 {
		//
	}
	return &relationclient.ActionResponse{}, xcode.InvalidFollowActionType
}

func (l *ActionLogic) followAction(fromId, toId int64) error {
	//TODO 多协程优化
	go func() {
		// 增加关注人数
		err := l.svcCtx.FollowModel.AddCacheFollowingCountHash(l.ctx, fromId)
		if err != nil {
			logx.Error(err)
		}
		// 增加粉丝数
		err = l.svcCtx.FanModel.AddCacheFanCountHash(l.ctx, toId)
		if err != nil {
			logx.Error(err)
		}
	}()

	go func() {
		// 增加fromId-->toId 关注关系
		err := l.svcCtx.FollowModel.AddBloomIsFollow(l.ctx, fromId, toId)
		if err != nil {
			logx.Error(err)
		}
	}()
	return nil
}

func (l *ActionLogic) unfollowAction() {

}
