package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteImageLogic {
	return &DeleteImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteImageLogic) DeleteImage(in *pb.DeleteImageReq) (*pb.DeleteImageResp, error) {

	err := l.svcCtx.ImageModel.Delete(l.ctx, in.ImageId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteImageResp{}, nil
}
