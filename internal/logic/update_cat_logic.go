package logic

import (
	"context"

	. "github.com/xh-polaris/meowchat-collection-rpc/internal/common"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCatLogic {
	return &UpdateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCatLogic) UpdateCat(in *pb.UpdateCatReq) (*pb.UpdateCatResp, error) {
	err := l.svcCtx.CatModel.Update(l.ctx, TransformModelCat(in.Cat))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCatResp{}, nil
}
