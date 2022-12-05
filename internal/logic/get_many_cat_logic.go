package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetManyCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetManyCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetManyCatLogic {
	return &GetManyCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetManyCatLogic) GetManyCat(in *pb.GetManyCatReq) (*pb.GetManyCatResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetManyCatResp{}, nil
}
