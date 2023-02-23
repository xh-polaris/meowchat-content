package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageLogic {
	return &CreateImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateImageLogic) CreateImage(in *pb.CreateImageReq) (*pb.CreateImageResp, error) {
	data := model.Image{
		CatId:    in.CatId,
		ImageUrl: in.ImageUrl,
	}
	err := l.svcCtx.ImageModel.Insert(l.ctx, &data)
	if err != nil {
		return nil, err
	}
	return &pb.CreateImageResp{ImageId: data.ID.Hex()}, nil
}
