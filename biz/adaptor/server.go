package adaptor

import (
	"context"

	"github.com/xh-polaris/meowchat-collection/biz/application/service"
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/config"

	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/collection"
)

type CollectionServerImpl struct {
	*config.Config
	CatService   service.CatService
	ImageService service.ImageService
}

func (s *CollectionServerImpl) SearchCat(ctx context.Context, req *collection.SearchCatReq) (res *collection.SearchCatResp, err error) {
	return s.CatService.SearchCat(ctx, req)
}

func (s *CollectionServerImpl) ListCat(ctx context.Context, req *collection.ListCatReq) (res *collection.ListCatResp, err error) {
	return s.CatService.ListCat(ctx, req)
}

func (s *CollectionServerImpl) RetrieveCat(ctx context.Context, req *collection.RetrieveCatReq) (res *collection.RetrieveCatResp, err error) {
	return s.CatService.RetrieveCat(ctx, req)
}

func (s *CollectionServerImpl) CreateCat(ctx context.Context, req *collection.CreateCatReq) (res *collection.CreateCatResp, err error) {
	return s.CatService.CreateCat(ctx, req)
}

func (s *CollectionServerImpl) UpdateCat(ctx context.Context, req *collection.UpdateCatReq) (res *collection.UpdateCatResp, err error) {
	return s.CatService.UpdateCat(ctx, req)
}

func (s *CollectionServerImpl) DeleteCat(ctx context.Context, req *collection.DeleteCatReq) (res *collection.DeleteCatResp, err error) {
	return s.CatService.DeleteCat(ctx, req)
}

func (s *CollectionServerImpl) CreateImage(ctx context.Context, req *collection.CreateImageReq) (res *collection.CreateImageResp, err error) {
	return s.ImageService.CreateImage(ctx, req)
}

func (s *CollectionServerImpl) DeleteImage(ctx context.Context, req *collection.DeleteImageReq) (res *collection.DeleteImageResp, err error) {
	return s.ImageService.DeleteImage(ctx, req)
}

func (s *CollectionServerImpl) ListImage(ctx context.Context, req *collection.ListImageReq) (res *collection.ListImageResp, err error) {
	return s.ImageService.ListImage(ctx, req)
}
