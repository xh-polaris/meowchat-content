package adaptor

import (
	"context"

	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"

	"github.com/xh-polaris/meowchat-content/biz/application/service"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
)

type ContentServerImpl struct {
	*config.Config
	CatService       service.ICatService
	ImageService     service.IImageService
	MomentService    service.IMomentService
	PostService      service.IPostService
	PlanService      service.IPlanService
	IncentiveService service.IIncentiveService
}

func (s *ContentServerImpl) CheckIn(ctx context.Context, req *content.CheckInReq) (res *content.CheckInResp, err error) {
	return s.IncentiveService.CheckIn(ctx, req)
}

func (s *ContentServerImpl) GetMission(ctx context.Context, req *content.GetMissionReq) (res *content.GetMissionResp, err error) {
	return s.IncentiveService.GetMission(ctx, req)
}

func (s *ContentServerImpl) CountDonateByPlan(ctx context.Context, req *content.CountDonateByPlanReq) (res *content.CountDonateByPlanResp, err error) {
	return s.PlanService.CountDonateByPlan(ctx, req)
}

func (s *ContentServerImpl) CountDonateByUser(ctx context.Context, req *content.CountDonateByUserReq) (res *content.CountDonateByUserResp, err error) {
	return s.PlanService.CountDonateByUser(ctx, req)
}

func (s *ContentServerImpl) DonateFish(ctx context.Context, req *content.DonateFishReq) (res *content.DonateFishResp, err error) {
	return s.PlanService.DonateFish(ctx, req)
}

func (s *ContentServerImpl) AddUserFish(ctx context.Context, req *content.AddUserFishReq) (res *content.AddUserFishResp, err error) {
	return s.PlanService.AddUserFish(ctx, req)
}

func (s *ContentServerImpl) ListFishByPlan(ctx context.Context, req *content.ListFishByPlanReq) (res *content.ListFishByPlanResp, err error) {
	return s.PlanService.ListFishByPlan(ctx, req)
}

func (s *ContentServerImpl) ListDonateByUser(ctx context.Context, req *content.ListDonateByUserReq) (res *content.ListDonateByUserResp, err error) {
	return s.PlanService.ListDonateByUser(ctx, req)
}

func (s *ContentServerImpl) RetrieveUserFish(ctx context.Context, req *content.RetrieveUserFishReq) (res *content.RetrieveUserFishResp, err error) {
	return s.PlanService.RetrieveUserFish(ctx, req)
}

func (s *ContentServerImpl) ListPlan(ctx context.Context, req *content.ListPlanReq) (res *content.ListPlanResp, err error) {
	return s.PlanService.ListPlan(ctx, req)
}

func (s *ContentServerImpl) CountPlan(ctx context.Context, req *content.CountPlanReq) (res *content.CountPlanResp, err error) {
	return s.PlanService.CountPlan(ctx, req)
}

func (s *ContentServerImpl) RetrievePlan(ctx context.Context, req *content.RetrievePlanReq) (res *content.RetrievePlanResp, err error) {
	return s.PlanService.RetrievePlan(ctx, req)
}

func (s *ContentServerImpl) CreatePlan(ctx context.Context, req *content.CreatePlanReq) (res *content.CreatePlanResp, err error) {
	return s.PlanService.CreatePlan(ctx, req)
}

func (s *ContentServerImpl) UpdatePlan(ctx context.Context, req *content.UpdatePlanReq) (res *content.UpdatePlanResp, err error) {
	return s.PlanService.UpdatePlan(ctx, req)
}

func (s *ContentServerImpl) DeletePlan(ctx context.Context, req *content.DeletePlanReq) (res *content.DeletePlanResp, err error) {
	return s.PlanService.DeletePlan(ctx, req)
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
