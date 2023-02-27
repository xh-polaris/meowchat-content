package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListImageLogic {
	return &ListImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListImageLogic) ListImage(in *pb.ListImageReq) (*pb.ListImageResp, error) {
	res, err := l.svcCtx.ImageModel.ListImage(l.ctx, in.CatId, in.PrevId, in.Limit, in.Offset, in.Backward)
	if err != nil {
		return nil, err
	}
	// 如果是向前翻页且得到的数据小于Limit，说明向前翻页到了尽头，那么返回查询第一页的记录
	if in.Backward && len(res) < int(in.Limit) {
		res, err = l.svcCtx.ImageModel.ListImage(l.ctx, in.CatId, nil, in.Limit, in.Offset, false)
		if err != nil {
			return nil, err
		}
	}

	imageList := make([]*pb.Image, len(res))
	for i, image := range res {
		imageList[i] = &pb.Image{
			Id:    image.ID.Hex(),
			Url:   image.ImageUrl,
			CatId: image.CatId,
		}
	}

	return &pb.ListImageResp{Images: imageList}, nil
}
