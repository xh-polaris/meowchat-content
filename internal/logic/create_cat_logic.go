package logic

import (
	"context"
	"strconv"

	. "github.com/xh-polaris/meowchat-collection-rpc/internal/common"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatLogic {
	return &CreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCatLogic) CreateCat(in *pb.CreateCatReq) (*pb.CreateCatResp, error) {
	res, err := l.svcCtx.CatModel.Insert(l.ctx, TransformModelCat(in.Cat))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &pb.CreateCatResp{CatId: strconv.FormatInt(id, 10)}, nil
}
