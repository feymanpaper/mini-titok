package logic

import (
	"context"
	"mini-titok/common/jwtx"
	"mini-titok/service/user/rpc/userclient"
	"time"

	"mini-titok/service/user/api/internal/svc"
	"mini-titok/service/user/api/internal/types"

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
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, res.UserId)
	if err != nil {
		return nil, err
	}
	return &types.LoginResponse{
		StatusCode:  res.StatusCode,
		StatusMsg:   *res.StatusMsg,
		UserId:      res.UserId,
		AccessToken: accessToken,
	}, nil
}
