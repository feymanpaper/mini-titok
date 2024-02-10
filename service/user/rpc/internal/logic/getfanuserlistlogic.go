package logic

import (
	"context"
	"mini-titok/service/relation/rpc/relationclient"
	"sort"

	"mini-titok/service/user/rpc/internal/svc"
	"mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFanUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	*GetUserInfoLogic
}

func NewGetFanUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFanUserListLogic {
	return &GetFanUserListLogic{
		ctx:              ctx,
		svcCtx:           svcCtx,
		Logger:           logx.WithContext(ctx),
		GetUserInfoLogic: NewGetUserInfoLogic(ctx, svcCtx),
	}
}

func (l *GetFanUserListLogic) GetFanUserList(in *userclient.GetFanUserListRequest) (*userclient.GetFanUserListResponse, error) {
	fanListResp, err := l.svcCtx.RelationRpc.FansList(l.ctx, &relationclient.FansListRequest{
		UserId:   in.UserId,
		Cursor:   in.Cursor,
		PageSize: in.PageSize,
		EndId:    in.EndId,
	})
	if err != nil {
		logx.Errorf(err.Error())
		return nil, err
	}
	fanIdPairList := fanListResp.FollowIdPairList
	idToTime := make(map[int64]int64)
	for _, pair := range fanIdPairList {
		idToTime[pair.UserId] = pair.CreateTime
	}
	if len(fanIdPairList) > int(in.PageSize) {
		fanIdPairList = fanIdPairList[:int(in.PageSize)]
	}

	useridList := make([]int64, len(fanIdPairList))
	for i, followInfo := range fanIdPairList {
		useridList[i] = followInfo.UserId
	}
	userList, err := l.GetUserListByIds(l.ctx, useridList)
	if err != nil {
		logx.Errorf(err.Error())
		return nil, err
	}
	sort.Slice(userList, func(i, j int) bool {
		return idToTime[userList[i].Id] > idToTime[userList[j].Id]
	})
	resFanUserListResp := &userclient.GetFanUserListResponse{
		FanUserList: userList,
		IsEnd:       fanListResp.IsEnd,
		Cursor:      fanListResp.Cursor,
		EndId:       fanListResp.EndId,
	}
	return resFanUserListResp, nil
}
