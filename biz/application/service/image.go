package service

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"

	imagemapper "github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/image"

	"github.com/google/wire"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
)

type IImageService interface {
	CreateImage(ctx context.Context, req *content.CreateImageReq) (*content.CreateImageResp, error)
	DeleteImage(ctx context.Context, req *content.DeleteImageReq) (*content.DeleteImageResp, error)
	ListImage(ctx context.Context, req *content.ListImageReq) (*content.ListImageResp, error)
}

type ImageService struct {
	ImageModel imagemapper.IMongoMapper
	MqProducer rocketmq.Producer
}

var ImageSet = wire.NewSet(
	wire.Struct(new(ImageService), "*"),
	wire.Bind(new(IImageService), new(*ImageService)),
)

func (s *ImageService) CreateImage(ctx context.Context, req *content.CreateImageReq) (*content.CreateImageResp, error) {
	data := make([]*imagemapper.Image, len(req.Images))
	for i := 0; i < len(data); i++ {
		data[i] = &imagemapper.Image{
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

	//发送使用url信息
	//var urls = make([]url.URL, len(data))
	//for i := 0; i < len(data); i++ {
	//	sendUrl, _ := url.Parse(data[i].ImageUrl)
	//	urls = append(urls, *sendUrl)
	//}
	//json, err := sonic.Marshal(urls)
	//if err != nil {
	//	return nil, err
	//}
	//msg := &mqprimitive.Message{
	//	Topic: "sts_used_url",
	//	Body:  json,
	//}
	//_, err = s.MqProducer.SendSync(ctx, msg)
	//if err != nil {
	//	return nil, err
	//}

	return &content.CreateImageResp{ImageIds: id}, nil
}

func (s *ImageService) DeleteImage(ctx context.Context, req *content.DeleteImageReq) (*content.DeleteImageResp, error) {
	err := s.ImageModel.Delete(ctx, req.ImageId)
	if err != nil {
		return nil, err
	}

	return &content.DeleteImageResp{}, nil
}

func (s *ImageService) ListImage(ctx context.Context, req *content.ListImageReq) (*content.ListImageResp, error) {
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

	imageList := make([]*content.Image, len(res))
	for i, image := range res {
		imageList[i] = &content.Image{
			Id:    image.ID.Hex(),
			Url:   image.ImageUrl,
			CatId: image.CatId,
		}
	}

	total, err := s.ImageModel.CountImage(ctx, req.CatId)
	if err != nil {
		return nil, err
	}
	return &content.ListImageResp{Images: imageList, Total: total}, nil
}
