package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"
	"mini-titok/common/xcode"
	"mini-titok/service/relation/rpc/internal/types"

	"mini-titok/service/relation/rpc/internal/svc"
	"mini-titok/service/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowActionLogic {
	return &FollowActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowActionLogic) FollowAction(in *relationclient.FollowActionRequest) (*relationclient.FollowActionResponse, error) {
	if in.ActionType == 1 {
		err := l.followAction(in.FromUserId, in.ToUserId)
		if err != nil {
			return &relationclient.FollowActionResponse{}, err
		}
		return &relationclient.FollowActionResponse{}, nil
	} else if in.ActionType == 2 {
		//
	}
	return &relationclient.FollowActionResponse{}, xcode.InvalidFollowActionType
}

func (l *FollowActionLogic) followAction(fromId, toId int64) error {
	//TODO 多协程优化
	go func() {
		// 增加关注人数
		err := l.svcCtx.FollowModel.AddCacheFollowingCountHash(l.ctx, fromId)
		if err != nil {
			logx.Error(err)
		}
		// 增加粉丝数
		err = l.svcCtx.FollowModel.AddCacheFanCountHash(l.ctx, toId)
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
	// 发送kafka增加fromId-->toId关注关系到数据库, 异步
	msg := types.AddIsFollowMsg{
		FromId: fromId,
		ToId:   toId,
	}
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
	return nil
}
