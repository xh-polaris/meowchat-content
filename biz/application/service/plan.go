package service

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/xh-polaris/gopkg/pagination/esp"
	"github.com/xh-polaris/gopkg/pagination/mongop"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/donate"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/fish"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/plan"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/convertor"
)

type IPlanService interface {
	ListPlan(ctx context.Context, req *content.ListPlanReq) (*content.ListPlanResp, error)
	CountPlan(ctx context.Context, req *content.CountPlanReq) (*content.CountPlanResp, error)
	RetrievePlan(ctx context.Context, req *content.RetrievePlanReq) (*content.RetrievePlanResp, error)
	CreatePlan(ctx context.Context, req *content.CreatePlanReq) (*content.CreatePlanResp, error)
	UpdatePlan(ctx context.Context, req *content.UpdatePlanReq) (*content.UpdatePlanResp, error)
	DeletePlan(ctx context.Context, req *content.DeletePlanReq) (*content.DeletePlanResp, error)
	DonateFish(ctx context.Context, req *content.DonateFishReq) (*content.DonateFishResp, error)
	AddUserFish(ctx context.Context, req *content.AddUserFishReq) (*content.AddUserFishResp, error)
	ListFishByPlan(ctx context.Context, req *content.ListFishByPlanReq) (*content.ListFishByPlanResp, error)
	ListDonateByUser(ctx context.Context, req *content.ListDonateByUserReq) (*content.ListDonateByUserResp, error)
	RetrieveUserFish(ctx context.Context, req *content.RetrieveUserFishReq) (*content.RetrieveUserFishResp, error)
	CountDonateByPlan(ctx context.Context, req *content.CountDonateByPlanReq) (*content.CountDonateByPlanResp, error)
	CountDonateByUser(ctx context.Context, req *content.CountDonateByUserReq) (*content.CountDonateByUserResp, error)
}

type PlanService struct {
	PlanMongoMapper   plan.IMongoMapper
	PlanEsMapper      plan.IEsMapper
	DonateMongoMapper donate.IMongoMapper
	FishMongoMapper   fish.IMongoMapper
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
		CatId:       m.CatId,
		CommunityId: m.CommunityId,
		PlanType:    m.PlanType,
		StartTime:   time.Unix(m.StartTime, 0),
		EndTime:     time.Unix(m.EndTime, 0),
		Description: m.Description,
		ImageUrls:   m.ImageUrls,
		Name:        m.Name,
		InitiatorId: m.InitiatorId,
		CoverUrl:    m.CoverUrl,
		Instruction: m.Instruction,
		Summary:     m.Summary,
		PlanState:   m.PlanState,
		MaxFish:     m.MaxFish,
		NowFish:     0,
	}

	err := s.PlanMongoMapper.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	//发送使用url信息
	//var urls = make([]url.URL, len(m.ImageUrls))
	//for i := 0; i < len(m.ImageUrls); i++ {
	//	sendUrl, _ := url.Parse(m.ImageUrls[i])
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

	return &content.CreatePlanResp{PlanId: data.ID.Hex()}, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, req *content.UpdatePlanReq) (*content.UpdatePlanResp, error) {
	m := req.Plan
	PlanId, err := primitive.ObjectIDFromHex(m.Id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	err = s.PlanMongoMapper.Update(ctx, &plan.Plan{
		ID:          PlanId,
		CatId:       m.CatId,
		PlanType:    m.PlanType,
		StartTime:   time.Unix(m.StartTime, 0),
		EndTime:     time.Unix(m.EndTime, 0),
		Description: m.Description,
		ImageUrls:   m.ImageUrls,
		InitiatorId: m.InitiatorId,
		CoverUrl:    m.CoverUrl,
		Instruction: m.Instruction,
		Summary:     m.Summary,
		PlanState:   m.PlanState,
		NowFish:     m.NowFish,
		MaxFish:     m.MaxFish,
	})
	if err != nil {
		return nil, err
	}

	//发送使用url信息
	//var urls = make([]url.URL, len(m.ImageUrls))
	//for i := 0; i < len(m.ImageUrls); i++ {
	//	sendUrl, _ := url.Parse(m.ImageUrls[i])
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

	return &content.UpdatePlanResp{}, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, req *content.DeletePlanReq) (*content.DeletePlanResp, error) {
	err := s.PlanMongoMapper.Delete(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &content.DeletePlanResp{}, nil
}

func (s *PlanService) DonateFish(ctx context.Context, req *content.DonateFishReq) (*content.DonateFishResp, error) {
	if req.GetFish() <= 0 {
		return nil, consts.ErrDonateInvalid
	}

	dbClient, err := s.FishMongoMapper.StartClient(ctx)
	if err != nil {
		return nil, err
	}

	err = dbClient.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		r, err := s.RetrieveUserFish(ctx, &content.RetrieveUserFishReq{UserId: req.UserId})
		if err != nil {
			return err
		}
		if req.GetFish() > r.Fish {
			return consts.ErrFishNotEnough
		}
		err = s.FishMongoMapper.Add(sessionContext, req.UserId, -req.Fish)
		if err != nil {
			return err
		}
		data, err := s.DonateMongoMapper.FindOneById(sessionContext, req.PlanId, req.UserId)
		switch err {
		case nil:
			data.FishNum += req.Fish
			err := s.DonateMongoMapper.Update(sessionContext, data)
			if err != nil {
				err2 := sessionContext.AbortTransaction(sessionContext)
				if err2 != nil {
					return err2
				}
				return err
			}
		case monc.ErrNotFound:
			err = s.DonateMongoMapper.Insert(sessionContext, &donate.Donate{
				UserId:  req.UserId,
				PlanId:  req.PlanId,
				FishNum: req.Fish,
			})
			if err != nil {
				err2 := sessionContext.AbortTransaction(sessionContext)
				if err2 != nil {
					return err2
				}
				return err
			}
		default:
			err2 := sessionContext.AbortTransaction(sessionContext)
			if err2 != nil {
				return err2
			}
			return err
		}
		p, err := s.PlanMongoMapper.FindOne(ctx, req.PlanId)
		if p.NowFish+req.Fish > p.MaxFish {
			err2 := sessionContext.AbortTransaction(sessionContext)
			if err2 != nil {
				return err2
			}
			return consts.ErrDonateOverFlow
		} else if p.NowFish+req.Fish == p.MaxFish {
			p.PlanState = 2
			p.NowFish = p.MaxFish
			err := s.PlanMongoMapper.Update(ctx, p)
			if err != nil {
				err2 := sessionContext.AbortTransaction(sessionContext)
				if err2 != nil {
					return err2
				}
				return err
			}
		} else {
			p.NowFish += req.Fish
			err := s.PlanMongoMapper.Update(ctx, p)
			if err != nil {
				err2 := sessionContext.AbortTransaction(sessionContext)
				if err2 != nil {
					return err2
				}
				return err
			}
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			err2 := sessionContext.AbortTransaction(sessionContext)
			if err2 != nil {
				return err2
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &content.DonateFishResp{}, nil
}

func (s *PlanService) AddUserFish(ctx context.Context, req *content.AddUserFishReq) (*content.AddUserFishResp, error) {
	data, err := s.FishMongoMapper.FindOne(ctx, req.UserId)
	switch err {
	case nil:
		data.FishNum = data.FishNum + req.Fish
		err = s.FishMongoMapper.Update(ctx, data)
		if err != nil {
			return nil, err
		}
		return &content.AddUserFishResp{}, nil
	case consts.ErrNotFound:
		oid, err := primitive.ObjectIDFromHex(req.UserId)
		if err != nil {
			return nil, err
		}
		err = s.FishMongoMapper.Insert(ctx, &fish.Fish{
			UserId:  oid,
			FishNum: req.Fish,
		})
		if err != nil {
			return nil, err
		}
		return &content.AddUserFishResp{}, nil
	default:
		return nil, err
	}
}

func (s *PlanService) ListFishByPlan(ctx context.Context, req *content.ListFishByPlanReq) (*content.ListFishByPlanResp, error) {
	resp := new(content.ListFishByPlanResp)

	p := convertor.ParsePagination(req.PaginationOptions)
	data, total, err := s.DonateMongoMapper.FindManyAndCountByPlanId(ctx, req.PlanId, p)
	if err != nil {
		return nil, err
	}
	fishMap := make(map[string]int64, len(data))
	userIds := make([]string, 0, len(data))
	for _, value := range data {
		fishMap[value.UserId] = value.FishNum
		userIds = append(userIds, value.UserId)
	}
	resp.Total = total
	resp.FishMap = fishMap
	resp.UserIds = userIds
	return resp, nil
}

func (s *PlanService) ListDonateByUser(ctx context.Context, req *content.ListDonateByUserReq) (*content.ListDonateByUserResp, error) {
	resp := new(content.ListDonateByUserResp)

	p := convertor.ParsePagination(req.PaginationOptions)
	data, total, err := s.DonateMongoMapper.FindManyAndCountByUserId(ctx, req.UserId, p, mongop.IdCursorType)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Donations = make([]*content.Donation, 0)
	for _, v := range data {
		temp := &content.Plan{}
		temp.Id = v.PlanId
		plan_, err := s.PlanMongoMapper.FindOne(ctx, v.PlanId)
		if err == nil {
			temp = convertor.ConvertPlan(plan_)
		}
		resp.Donations = append(resp.Donations, &content.Donation{
			Plan:       temp,
			DonateTime: v.CreateAt.Unix(),
			DonateNum:  v.FishNum,
			DonateId:   v.ID.Hex(),
		})
	}
	return resp, nil
}

func (s *PlanService) RetrieveUserFish(ctx context.Context, req *content.RetrieveUserFishReq) (*content.RetrieveUserFishResp, error) {
	data, err := s.FishMongoMapper.FindOne(ctx, req.UserId)
	switch err {
	case nil:
		return &content.RetrieveUserFishResp{Fish: data.FishNum}, nil
	case consts.ErrNotFound:
		return &content.RetrieveUserFishResp{Fish: 0}, nil
	default:
		return nil, err
	}
}

func (s *PlanService) CountDonateByPlan(ctx context.Context, req *content.CountDonateByPlanReq) (*content.CountDonateByPlanResp, error) {
	total, err := s.DonateMongoMapper.CountByPlanId(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &content.CountDonateByPlanResp{
		Total: total,
	}, nil
}

func (s *PlanService) CountDonateByUser(ctx context.Context, req *content.CountDonateByUserReq) (*content.CountDonateByUserResp, error) {
	total, err := s.DonateMongoMapper.CountByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &content.CountDonateByUserResp{
		Total: total,
	}, nil
}
