package logic

import (
	"context"

	"mini-titok/service/relation/rpc/internal/svc"
	"mini-titok/service/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActionLogic) Action(in *relationclient.ActionRequest) (*relationclient.ActionResponse, error) {
	// todo: add your logic here and delete this line

	return &relationclient.ActionResponse{}, nil
}
