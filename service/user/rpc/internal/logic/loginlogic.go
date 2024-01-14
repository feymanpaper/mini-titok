package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"mini-titok/common/cryptx"
	"mini-titok/service/user/model"
	"mini-titok/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
	"mini-titok/service/user/rpc/internal/svc"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *userclient.LoginRequest) (*userclient.LoginResponse, error) {
	res, err := l.svcCtx.UserModel.FindOneByName(l.ctx, in.UserName)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "用户不存在")
		}
		return nil, status.Error(500, err.Error())
	}
	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	if password != res.Password {
		return nil, status.Error(100, "密码错误")
	}
	return &userclient.LoginResponse{
		UserId: res.Id,
		Token:  "123",
	}, nil
}
