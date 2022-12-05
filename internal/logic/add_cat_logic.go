package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCatLogic {
	return &AddCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCatLogic) AddCat(in *pb.AddCatReq) (*pb.AddCatResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddCatResp{}, nil
}
