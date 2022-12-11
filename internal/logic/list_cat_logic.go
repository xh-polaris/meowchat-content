package logic

import (
	"context"

	. "github.com/xh-polaris/meowchat-collection-rpc/internal/common"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCatLogic {
	return &ListCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCatLogic) ListCat(in *pb.ListCatReq) (*pb.ListCatResp, error) {
	catList, err := l.svcCtx.CatModel.FindManyValidByCommunityIdValid(l.ctx, in.CommunityId, in.Skip, in.Count)
	if err != nil {
		return nil, err
	}
	var Cat []*pb.Cat
	for _, val := range catList {
		Cat = append(Cat, TransformPbCat(val))
	}
	return &pb.ListCatResp{Cats: Cat}, nil
}
