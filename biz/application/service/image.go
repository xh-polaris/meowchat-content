package service

import (
	"context"

	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/data/db"
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/mapper"

	"github.com/google/wire"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/collection"
)

type ImageService interface {
	CreateImage(ctx context.Context, req *collection.CreateImageReq) (*collection.CreateImageResp, error)
	DeleteImage(ctx context.Context, req *collection.DeleteImageReq) (*collection.DeleteImageResp, error)
	ListImage(ctx context.Context, req *collection.ListImageReq) (*collection.ListImageResp, error)
}

type ImageServiceImpl struct {
	ImageModel mapper.ImageModel
}

var ImageSet = wire.NewSet(
	wire.Struct(new(ImageServiceImpl), "*"),
	wire.Bind(new(ImageService), new(*ImageServiceImpl)),
)

func (s *ImageServiceImpl) CreateImage(ctx context.Context, req *collection.CreateImageReq) (*collection.CreateImageResp, error) {
	data := make([]*db.Image, len(req.Images))
	for i := 0; i < len(data); i++ {
		data[i] = &db.Image{
			CatId:    req.Images[i].CatId,
			ImageUrl: req.Images[i].Url,
		}
	}
	err := s.ImageModel.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	id := make([]string, len(data))
	for i := 0; i < len(data); i++ {
		id[i] = data[i].ID.Hex()
	}
	return &collection.CreateImageResp{ImageIds: id}, nil
}

func (s *ImageServiceImpl) DeleteImage(ctx context.Context, req *collection.DeleteImageReq) (*collection.DeleteImageResp, error) {
	err := s.ImageModel.Delete(ctx, req.ImageId)
	if err != nil {
		return nil, err
	}

	return &collection.DeleteImageResp{}, nil
}

func (s *ImageServiceImpl) ListImage(ctx context.Context, req *collection.ListImageReq) (*collection.ListImageResp, error) {
	res, err := s.ImageModel.ListImage(ctx, req.CatId, req.PrevId, req.Limit, req.Offset, req.Backward)
	if err != nil {
		return nil, err
	}
	// 如果是向前翻页且得到的数据小于Limit，说明向前翻页到了尽头，那么返回查询第一页的记录
	if req.Backward && len(res) < int(req.Limit) {
		res, err = s.ImageModel.ListImage(ctx, req.CatId, nil, req.Limit, req.Offset, false)
		if err != nil {
			return nil, err
		}
	}

	imageList := make([]*collection.Image, len(res))
	for i, image := range res {
		imageList[i] = &collection.Image{
			Id:    image.ID.Hex(),
			Url:   image.ImageUrl,
			CatId: image.CatId,
		}
	}

	total, err := s.ImageModel.CountImage(ctx, req.CatId)
	if err != nil {
		return nil, err
	}
	return &collection.ListImageResp{Images: imageList, Total: total}, nil
}
