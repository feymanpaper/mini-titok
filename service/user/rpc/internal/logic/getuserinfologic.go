package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"mini-titok/common/xcode"
	"mini-titok/service/relation/rpc/relationclient"
	"mini-titok/service/user/model"
	"mini-titok/service/user/rpc/internal/svc"
	"mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *userclient.UserInfoRequest) (*userclient.UserInfoResponse, error) {
	res, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xcode.RecordNotFound
		}
		return nil, xcode.DBErr
	}

	relationCountResp, _ := l.svcCtx.RelationRpc.GetRelationCount(l.ctx, &relationclient.GetRelationCountRequest{UserId: in.UserId})
	userInfo := &userclient.UserInfo{
		Id:              res.Id,
		Name:            res.Name,
		FollowCount:     relationCountResp.FollowCount,
		FollowerCount:   relationCountResp.FanCount,
		IsFollow:        true,
		Avatar:          res.Avatar,
		BackgroundImage: res.BackgroundImage,
		Signature:       res.Signature,
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
	return &userclient.UserInfoResponse{
		UserInfo: userInfo,
	}, nil
}

// getUserListByIds 调用 userrpc mapreduce并发获取userInfoList
func (l *GetUserInfoLogic) GetUserListByIds(ctx context.Context, userIds []int64) ([]*userclient.UserInfo, error) {
	userInfoList, err := mr.MapReduce[int64, *userclient.UserInfo, []*userclient.UserInfo](func(source chan<- int64) {
		for _, aid := range userIds {
			if aid == -1 {
				continue
			}
			source <- aid
		}
	}, func(id int64, writer mr.Writer[*userclient.UserInfo], cancel func(error)) {
		userInfo, err2 := l.GetUserInfo(&userclient.UserInfoRequest{UserId: id})
		if err2 != nil {
			cancel(err2)
			return
		}
		writer.Write(userInfo.UserInfo)
	}, func(pipe <-chan *userclient.UserInfo, writer mr.Writer[[]*userclient.UserInfo], cancel func(error)) {
		var userInfoList []*userclient.UserInfo
		for userInfo := range pipe {
			userInfoList = append(userInfoList, userInfo)
		}
		writer.Write(userInfoList)
	})
	if err != nil {
		return nil, err
	}
	return userInfoList, nil
}
