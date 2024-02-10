package logic

import (
	"context"
	"mini-titok/service/relation/rpc/relationclient"
	"mini-titok/service/user/rpc/internal/svc"
	"mini-titok/service/user/rpc/userclient"
	"sort"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*GetUserInfoLogic
}

func NewGetFollowUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowUserListLogic {
	return &GetFollowUserListLogic{
		ctx:              ctx,
		svcCtx:           svcCtx,
		Logger:           logx.WithContext(ctx),
		GetUserInfoLogic: NewGetUserInfoLogic(ctx, svcCtx),
	}
}

func (l *GetFollowUserListLogic) GetFollowUserList(in *userclient.GetFollowUserListRequest) (*userclient.GetFollowUserListResponse, error) {
	followIdListResp, err := l.svcCtx.RelationRpc.FollowList(l.ctx, &relationclient.FollowListRequest{
		UserId:   in.UserId,
		Cursor:   in.Cursor,
		PageSize: in.PageSize,
		EndId:    in.EndId,
	})
	if err != nil {
		return nil, err
	}
	followIdPairList := followIdListResp.FollowIdPairList
	idToTime := make(map[int64]int64)
	for _, pair := range followIdPairList {
		idToTime[pair.UserId] = pair.CreateTime
	}
	if len(followIdPairList) > int(in.PageSize) {
		followIdPairList = followIdPairList[:int(in.PageSize)]
	}

	useridList := make([]int64, len(followIdPairList))
	for i, followInfo := range followIdPairList {
		useridList[i] = followInfo.UserId
	}
	userList, err := l.GetUserListByIds(l.ctx, useridList)
	if err != nil {
		return nil, err
	}
	sort.Slice(userList, func(i, j int) bool {
		return idToTime[userList[i].Id] > idToTime[userList[j].Id]
	})
	resFollowUserListResp := &userclient.GetFollowUserListResponse{
		FollowUserList: userList,
		IsEnd:          followIdListResp.IsEnd,
		Cursor:         followIdListResp.Cursor,
		EndId:          followIdListResp.EndId,
	}
	return resFollowUserListResp, nil
}
