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
	data := make([]*model.Image, len(in.Image))
	for i := 0; i < len(data); i++ {
		data[i] = &model.Image{
			CatId:    in.Image[i].CatId,
			ImageUrl: in.Image[i].Url,
		}
	}
	err := l.svcCtx.ImageModel.InsertMany(l.ctx, data)
	if err != nil {
		return nil, err
	}
	id := make([]string, len(data))
	for i := 0; i < len(data); i++ {
		id[i] = data[i].ID.Hex()
	}
	return &pb.CreateImageResp{ImageId: id}, nil
}
