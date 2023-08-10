package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/gopkg/pagination/esp"
	"github.com/xh-polaris/gopkg/pagination/mongop"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/plan"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/convertor"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IPlanService interface {
	ListPlan(ctx context.Context, req *content.ListPlanReq) (*content.ListPlanResp, error)
	CountPlan(ctx context.Context, req *content.CountPlanReq) (*content.CountPlanResp, error)
	RetrievePlan(ctx context.Context, req *content.RetrievePlanReq) (*content.RetrievePlanResp, error)
	CreatePlan(ctx context.Context, req *content.CreatePlanReq) (*content.CreatePlanResp, error)
	UpdatePlan(ctx context.Context, req *content.UpdatePlanReq) (*content.UpdatePlanResp, error)
	DeletePlan(ctx context.Context, req *content.DeletePlanReq) (*content.DeletePlanResp, error)
}

type PlanService struct {
	PlanMongoMapper plan.IMongoMapper
	PlanEsMapper    plan.IEsMapper
}

var PlanSet = wire.NewSet(
	wire.Struct(new(PlanService), "*"),
	wire.Bind(new(IPlanService), new(*PlanService)),
)

func (s *PlanService) ListPlan(ctx context.Context, req *content.ListPlanReq) (*content.ListPlanResp, error) {
	resp := new(content.ListPlanResp)
	var plans []*plan.Plan
	var total int64
	var err error

	filter := convertor.ParsePlanFilter(req.FilterOptions)
	p := convertor.ParsePagination(req.PaginationOptions)
	if req.SearchOptions == nil {
		plans, total, err = s.PlanMongoMapper.FindManyAndCount(ctx, filter, p, mongop.IdCursorType)
		print(plans)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			plans, total, err = s.PlanEsMapper.Search(ctx, convertor.ConvertPlanAllFieldsSearchQuery(o), filter, p, esp.ScoreCursorType)
		case *content.SearchOptions_MultiFieldsKey:
			plans, total, err = s.PlanEsMapper.Search(ctx, convertor.ConvertPlanMultiFieldsSearchQuery(o), filter, p, esp.ScoreCursorType)
		}
		if err != nil {
			return nil, err
		}
	}

	resp.Total = total
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Plans = make([]*content.Plan, 0, len(plans))
	for _, Plan_ := range plans {
		resp.Plans = append(resp.Plans, convertor.ConvertPlan(Plan_))
	}

	return resp, nil
}

func (s *PlanService) CountPlan(ctx context.Context, req *content.CountPlanReq) (*content.CountPlanResp, error) {
	resp := new(content.CountPlanResp)
	var err error
	filter := convertor.ParsePlanFilter(req.FilterOptions)
	if req.SearchOptions == nil {
		resp.Total, err = s.PlanMongoMapper.Count(ctx, filter)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			resp.Total, err = s.PlanEsMapper.CountWithQuery(ctx, convertor.ConvertPlanAllFieldsSearchQuery(o), filter)
		case *content.SearchOptions_MultiFieldsKey:
			resp.Total, err = s.PlanEsMapper.CountWithQuery(ctx, convertor.ConvertPlanMultiFieldsSearchQuery(o), filter)
		}
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *PlanService) RetrievePlan(ctx context.Context, req *content.RetrievePlanReq) (*content.RetrievePlanResp, error) {
	data, err := s.PlanMongoMapper.FindOne(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	m := convertor.ConvertPlan(data)
	return &content.RetrievePlanResp{Plan: m}, nil
}

func (s *PlanService) CreatePlan(ctx context.Context, req *content.CreatePlanReq) (*content.CreatePlanResp, error) {
	m := req.Plan
	data := &plan.Plan{
		CatId:        m.CatId,
		PlanType:     m.PlanType,
		StartTime:    time.Unix(m.StartTime, 0),
		EndTime:      time.Unix(m.EndTime, 0),
		Description:  m.Description,
		ImageUrls:    m.ImageUrls,
		Name:         m.Name,
		InitiatorIds: m.InitiatorIds,
	}

	err := s.PlanMongoMapper.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	return &content.CreatePlanResp{PlanId: data.ID.Hex()}, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, req *content.UpdatePlanReq) (*content.UpdatePlanResp, error) {
	m := req.Plan
	PlanId, err := primitive.ObjectIDFromHex(m.Id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	err = s.PlanMongoMapper.Update(ctx, &plan.Plan{
		ID:           PlanId,
		CatId:        m.CatId,
		PlanType:     m.PlanType,
		StartTime:    time.Unix(m.StartTime, 0),
		EndTime:      time.Unix(m.EndTime, 0),
		Description:  m.Description,
		ImageUrls:    m.ImageUrls,
		InitiatorIds: m.InitiatorIds,
	})
	if err != nil {
		return nil, err
	}

	return &content.UpdatePlanResp{}, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, req *content.DeletePlanReq) (*content.DeletePlanResp, error) {
	err := s.PlanMongoMapper.Delete(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &content.DeletePlanResp{}, nil
}