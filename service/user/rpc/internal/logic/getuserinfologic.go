package logic

import (
	"context"
	"mini-titok/common/xcode"
	"mini-titok/service/user/model"
	"mini-titok/service/user/rpc/userclient"

	"mini-titok/service/user/rpc/internal/svc"

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
	logx.Info(res)
	userInfo := &userclient.UserInfo{
		Id:              res.Id,
		Name:            res.Name,
		FollowCount:     &res.FollowCount,
		FollowerCount:   &res.FollowerCount,
		IsFollow:        true,
		Avatar:          &res.Avatar,
		BackgroundImage: &res.BackgroundImage,
		Signature:       &res.Signature,
		TotalFavorited:  &res.TotalFavorited,
		WorkCount:       &res.WorkCount,
		FavoriteCount:   &res.FavoriteCount,
	}
	return &userclient.UserInfoResponse{
		UserInfo: userInfo,
	}, nil
}
