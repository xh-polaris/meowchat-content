package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/meowchat-collection-rpc/types/pb"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCatLogic {
	return &SearchCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchCatLogic) SearchCat(in *pb.SearchCatReq) (*pb.SearchCatResp, error) {
	data, err := l.svcCtx.CatModel.Search(l.ctx, in.CommunityId, in.Keyword, in.Skip, in.Count)
	if err != nil {
		return nil, err
	}
	count, err := l.svcCtx.CatModel.SearchNumber(l.ctx, in.CommunityId, in.Keyword)
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
	return &pb.SearchCatResp{Cats: cats, Count: count}, nil
}
