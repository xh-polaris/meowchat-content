package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/errorx"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"time"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetImageByCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetImageByCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageByCatLogic {
	return &GetImageByCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetImageByCatLogic) GetImageByCat(in *pb.GetImageByCatReq) (*pb.GetImageByCatResp, error) {

	var lastId primitive.ObjectID
	if in.PrevId == nil {
		lastId = primitive.NewObjectIDFromTimestamp(time.Unix(math.MaxInt32, 0))
	} else {
		var err error
		lastId, err = primitive.ObjectIDFromHex(*in.PrevId)
		if err != nil {
			return nil, errorx.ErrInvalidId
		}
	}
	res, err := l.svcCtx.ImageModel.ListImageByCat(l.ctx, in.CatId, lastId, in.Limit)
	if err != nil {
		return nil, err
	}
	imageList := make([]string, len(res))
	for i := 0; i < len(res); i++ {
		imageList[i] = res[i].ImageUrl
	}
	if len(res) > 0 {
		lastId = res[len(res)-1].ID
	}

	return &pb.GetImageByCatResp{ImageUrl: imageList, LastId: lastId.Hex()}, nil
}
