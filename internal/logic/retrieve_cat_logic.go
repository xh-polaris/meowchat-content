package logic

import (
	"context"
	"strconv"

	"github.com/xh-polaris/meowchat-collection-rpc/errorx"
	. "github.com/xh-polaris/meowchat-collection-rpc/internal/common"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

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
	id, err := strconv.ParseInt(in.CatId, 10, 64)
	if err != nil {
		return nil, err
	}
	cat, err := l.svcCtx.CatModel.FindOneValid(l.ctx, id)
	if err != nil {
		return nil, errorx.NoSuchCat
	}
	return &pb.RetrieveCatResp{Cat: TransformPbCat(cat)}, nil
}
