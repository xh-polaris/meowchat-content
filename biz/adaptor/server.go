package adaptor

import (
	"context"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"

	"github.com/xh-polaris/meowchat-content/biz/application/service"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
)

type ContentServerImpl struct {
	*config.Config
	CatService    service.ICatService
	ImageService  service.IImageService
	MomentService service.IMomentService
	PostService   service.IPostService
}

func (s *ContentServerImpl) SearchCat(ctx context.Context, req *content.SearchCatReq) (res *content.SearchCatResp, err error) {
	return s.CatService.SearchCat(ctx, req)
}

func (s *ContentServerImpl) ListCat(ctx context.Context, req *content.ListCatReq) (res *content.ListCatResp, err error) {
	return s.CatService.ListCat(ctx, req)
}

func (s *ContentServerImpl) RetrieveCat(ctx context.Context, req *content.RetrieveCatReq) (res *content.RetrieveCatResp, err error) {
	return s.CatService.RetrieveCat(ctx, req)
}

func (s *ContentServerImpl) CreateCat(ctx context.Context, req *content.CreateCatReq) (res *content.CreateCatResp, err error) {
	return s.CatService.CreateCat(ctx, req)
}

func (s *ContentServerImpl) UpdateCat(ctx context.Context, req *content.UpdateCatReq) (res *content.UpdateCatResp, err error) {
	return s.CatService.UpdateCat(ctx, req)
}

func (s *ContentServerImpl) DeleteCat(ctx context.Context, req *content.DeleteCatReq) (res *content.DeleteCatResp, err error) {
	return s.CatService.DeleteCat(ctx, req)
}

func (s *ContentServerImpl) CreateImage(ctx context.Context, req *content.CreateImageReq) (res *content.CreateImageResp, err error) {
	return s.ImageService.CreateImage(ctx, req)
}

func (s *ContentServerImpl) DeleteImage(ctx context.Context, req *content.DeleteImageReq) (res *content.DeleteImageResp, err error) {
	return s.ImageService.DeleteImage(ctx, req)
}

func (s *ContentServerImpl) ListImage(ctx context.Context, req *content.ListImageReq) (res *content.ListImageResp, err error) {
	return s.ImageService.ListImage(ctx, req)
}

func (s *ContentServerImpl) ListMoment(ctx context.Context, req *content.ListMomentReq) (res *content.ListMomentResp, err error) {
	return s.MomentService.ListMoment(ctx, req)
}

func (s *ContentServerImpl) CountMoment(ctx context.Context, req *content.CountMomentReq) (res *content.CountMomentResp, err error) {
	return s.MomentService.CountMoment(ctx, req)
}

func (s *ContentServerImpl) RetrieveMoment(ctx context.Context, req *content.RetrieveMomentReq) (res *content.RetrieveMomentResp, err error) {
	return s.MomentService.RetrieveMoment(ctx, req)
}

func (s *ContentServerImpl) CreateMoment(ctx context.Context, req *content.CreateMomentReq) (res *content.CreateMomentResp, err error) {
	return s.MomentService.CreateMoment(ctx, req)
}

func (s *ContentServerImpl) UpdateMoment(ctx context.Context, req *content.UpdateMomentReq) (res *content.UpdateMomentResp, err error) {
	return s.MomentService.UpdateMoment(ctx, req)
}

func (s *ContentServerImpl) DeleteMoment(ctx context.Context, req *content.DeleteMomentReq) (res *content.DeleteMomentResp, err error) {
	return s.MomentService.DeleteMoment(ctx, req)
}

func (s *ContentServerImpl) CreatePost(ctx context.Context, req *content.CreatePostReq) (res *content.CreatePostResp, err error) {
	return s.PostService.CreatePost(ctx, req)
}

func (s *ContentServerImpl) RetrievePost(ctx context.Context, req *content.RetrievePostReq) (res *content.RetrievePostResp, err error) {
	return s.PostService.RetrievePost(ctx, req)
}

func (s *ContentServerImpl) UpdatePost(ctx context.Context, req *content.UpdatePostReq) (res *content.UpdatePostResp, err error) {
	return s.PostService.UpdatePost(ctx, req)
}

func (s *ContentServerImpl) DeletePost(ctx context.Context, req *content.DeletePostReq) (res *content.DeletePostResp, err error) {
	return s.PostService.DeletePost(ctx, req)
}

func (s *ContentServerImpl) ListPost(ctx context.Context, req *content.ListPostReq) (res *content.ListPostResp, err error) {
	return s.PostService.ListPost(ctx, req)
}

func (s *ContentServerImpl) CountPost(ctx context.Context, req *content.CountPostReq) (res *content.CountPostResp, err error) {
	return s.PostService.CountPost(ctx, req)
}

func (s *ContentServerImpl) SetOfficial(ctx context.Context, req *content.SetOfficialReq) (res *content.SetOfficialResp, err error) {
	return s.PostService.SetOfficial(ctx, req)
}
