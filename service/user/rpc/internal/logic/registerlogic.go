package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"mini-titok/common/cryptx"
	"mini-titok/service/user/model"
	"mini-titok/service/user/rpc/userclient"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"mini-titok/service/user/rpc/internal/svc"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *userclient.RegisterRequest) (*userclient.RegisterResponse, error) {
	//判断该用户是否注册
	_, err := l.svcCtx.UserModel.FindOneByName(l.ctx, in.UserName)
	if err == nil {
		return nil, status.Error(100, "该用户已存在")
	}
	if err == model.ErrNotFound {
		newUser := model.User{
			Name:            in.UserName,
			Avatar:          "",
			BackgroundImage: "",
			Signature:       "",
			Password:        cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password),
			CreateTime:      time.Time{},
			UpdateTime:      time.Time{},
		}
		res, err := l.svcCtx.UserModel.Insert(l.ctx, &newUser)
		if err != nil {
			return nil, status.Error(500, err.Error())
		}
		newUser.Id, err = res.LastInsertId()
		if err != nil {
			return nil, status.Error(500, err.Error())
		}

		return &userclient.RegisterResponse{
			UserId: newUser.Id,
			Token:  "123",
		}, nil
	}
	return nil, status.Error(500, err.Error())
}
