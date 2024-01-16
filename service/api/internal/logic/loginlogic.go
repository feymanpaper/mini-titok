package logic

import (
	"context"
	"mini-titok/common/jwtx"
	"mini-titok/common/xcode"
	"mini-titok/service/user/api/internal/svc"
	"mini-titok/service/user/api/internal/types"
	"mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginRequest{
		UserName: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	tokenString, err := jwtx.CreateToken(res.UserId, l.svcCtx.Config.JwtAuth.AccessSecret, l.svcCtx.Config.JwtAuth.AccessExpire)
	if err != nil {
		return nil, xcode.FailGenerateJwt
	}

	return &types.LoginResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId:      res.UserId,
		AccessToken: tokenString,
	}, nil
}
