package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"strconv"

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
	id, err := strconv.ParseInt(in.CatId, 10, 64)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.CatModel.DeleteSoftly(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCatResp{}, nil
}
