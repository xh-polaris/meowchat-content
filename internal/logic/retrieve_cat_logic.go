package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/meowchat-collection-rpc/errorx"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetrieveCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRetrieveCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrieveCatLogic {
	return &RetrieveCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RetrieveCatLogic) RetrieveCat(in *pb.RetrieveCatReq) (*pb.RetrieveCatResp, error) {
	data, err := l.svcCtx.CatModel.FindOne(l.ctx, in.CatId)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errorx.ErrNoSuchCat
	default:
		return nil, err
	}
	cat := &pb.Cat{}
	err = copier.Copy(cat, data)
	if err != nil {
		return nil, err
	}
	cat.Id = data.ID.Hex()
	cat.CreateAt = data.CreateAt.Unix()
	return &pb.RetrieveCatResp{Cat: cat}, nil
}
