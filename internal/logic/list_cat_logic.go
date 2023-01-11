package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/meowchat-collection-rpc/types/pb"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
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
	data, countx, err := l.svcCtx.CatModel.FindManyByCommunityId(l.ctx, in.CommunityId, in.Skip, in.Count)
	if err != nil {
		return nil, err
	}
	var cats []*pb.Cat
	for _, val := range data {
		cat := &pb.Cat{}
		err = copier.Copy(cat, val)
		if err != nil {
			return nil, err
		}
		cat.Id = val.ID.Hex()
		cat.CreateAt = val.CreateAt.Unix()
		cats = append(cats, cat)
	}
	return &pb.ListCatResp{Cats: cats, Count: countx}, nil
}
