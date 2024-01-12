package logic

import (
	"context"
	"mini-titok/service/user/api/internal/svc"
	"mini-titok/service/user/api/internal/types"
	userclient "mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
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
	return &types.RegisterResponse{
		StatusCode:  res.StatusCode,
		StatusMsg:   *res.StatusMsg,
		UserId:      res.UserId,
		AccessToken: res.Token,
	}, nil
}
