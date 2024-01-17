package logic

import (
	"context"
	"mini-titok/service/api/internal/svc"
	"mini-titok/service/api/internal/types"
	"mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	res, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: req.Id,
		Token:  req.Token,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		User: types.User{
			Id:              res.UserInfo.Id,
			Name:            res.UserInfo.Name,
			FollowCount:     *res.UserInfo.FollowCount,
			FollowerCount:   *res.UserInfo.FollowerCount,
			IsFollow:        res.UserInfo.IsFollow,
			Avatar:          *res.UserInfo.Avatar,
			BackgroundImage: *res.UserInfo.BackgroundImage,
			Signature:       *res.UserInfo.Signature,
			TotalFavorited:  *res.UserInfo.TotalFavorited,
			WorkCount:       *res.UserInfo.WorkCount,
			FavoriteCount:   *res.UserInfo.FavoriteCount,
		},
	}, nil

}
