package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/types/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCatLogic {
	return &DeleteCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCatLogic) DeleteCat(in *pb.DeleteCatReq) (*pb.DeleteCatResp, error) {
	err := l.svcCtx.CatModel.Delete(l.ctx, in.CatId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCatResp{}, nil
}
