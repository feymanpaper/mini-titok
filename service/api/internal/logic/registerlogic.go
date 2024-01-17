package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mini-titok/common/jwtx"
	"mini-titok/common/xcode"
	"mini-titok/service/api/internal/svc"
	"mini-titok/service/api/internal/types"
	userclient "mini-titok/service/user/rpc/userclient"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	res, err := l.svcCtx.UserRpc.Register(l.ctx, &userclient.RegisterRequest{
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
	return &types.RegisterResponse{
		BasicResponse: types.BasicResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId:      res.UserId,
		AccessToken: tokenString,
	}, nil
}
