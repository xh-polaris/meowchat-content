package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/moment"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/convertor"
	"github.com/xh-polaris/paginator-go/esp"
	"github.com/xh-polaris/paginator-go/mongop"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMomentService interface {
	ListMoment(ctx context.Context, req *content.ListMomentReq) (*content.ListMomentResp, error)
	CountMoment(ctx context.Context, req *content.CountMomentReq) (*content.CountMomentResp, error)
	RetrieveMoment(ctx context.Context, req *content.RetrieveMomentReq) (*content.RetrieveMomentResp, error)
	CreateMoment(ctx context.Context, req *content.CreateMomentReq) (*content.CreateMomentResp, error)
	UpdateMoment(ctx context.Context, req *content.UpdateMomentReq) (*content.UpdateMomentResp, error)
	DeleteMoment(ctx context.Context, req *content.DeleteMomentReq) (*content.DeleteMomentResp, error)
}

type MomentService struct {
	MomentMongoMapper moment.IMongoMapper
	MomentEsMapper    moment.IEsMapper
}

var MomentSet = wire.NewSet(
	wire.Struct(new(MomentService), "*"),
	wire.Bind(new(IMomentService), new(*MomentService)),
)

func (s *MomentService) ListMoment(ctx context.Context, req *content.ListMomentReq) (*content.ListMomentResp, error) {
	resp := new(content.ListMomentResp)
	var moments []*moment.Moment
	var total int64
	var err error

	filter := convertor.ParseMomentFilter(req.FilterOptions)
	p := convertor.ParsePagination(req.PaginationOptions)
	if req.SearchOptions == nil {
		moments, total, err = s.MomentMongoMapper.FindManyAndCount(ctx, filter, p, &mongop.IdSorter{})
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			moments, total, err = s.MomentEsMapper.Search(ctx, convertor.ConvertMomentAllFieldsSearchQuery(o), filter, p, &esp.ScoreSorter{})
		case *content.SearchOptions_MultiFieldsKey:
			moments, total, err = s.MomentEsMapper.Search(ctx, convertor.ConvertMomentMultiFieldsSearchQuery(o), filter, p, &esp.ScoreSorter{})
		}
		if err != nil {
			return nil, err
		}
	}

	resp.Total = total
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Moments = make([]*content.Moment, 0, len(moments))
	for _, moment_ := range moments {
		resp.Moments = append(resp.Moments, convertor.ConvertMoment(moment_))
	}

	return resp, nil
}

func (s *MomentService) CountMoment(ctx context.Context, req *content.CountMomentReq) (*content.CountMomentResp, error) {
	resp := new(content.CountMomentResp)
	var err error
	filter := convertor.ParseMomentFilter(req.FilterOptions)
	if req.SearchOptions == nil {
		resp.Total, err = s.MomentMongoMapper.Count(ctx, filter)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			resp.Total, err = s.MomentEsMapper.CountWithQuery(ctx, convertor.ConvertMomentAllFieldsSearchQuery(o), filter)
		case *content.SearchOptions_MultiFieldsKey:
			resp.Total, err = s.MomentEsMapper.CountWithQuery(ctx, convertor.ConvertMomentMultiFieldsSearchQuery(o), filter)
		}
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *MomentService) RetrieveMoment(ctx context.Context, req *content.RetrieveMomentReq) (*content.RetrieveMomentResp, error) {
	data, err := s.MomentMongoMapper.FindOne(ctx, req.MomentId)
	if err != nil {
		return nil, err
	}
	m := convertor.ConvertMoment(data)
	return &content.RetrieveMomentResp{Moment: m}, nil
}

func (s *MomentService) CreateMoment(ctx context.Context, req *content.CreateMomentReq) (*content.CreateMomentResp, error) {
	m := req.Moment
	data := &moment.Moment{
		Photos:      m.Photos,
		Title:       m.Title,
		Text:        m.Text,
		CommunityId: m.CommunityId,
		UserId:      m.UserId,
		CatId:       m.CatId,
	}

	err := s.MomentMongoMapper.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	return &content.CreateMomentResp{MomentId: data.ID.Hex()}, nil
}

func (s *MomentService) UpdateMoment(ctx context.Context, req *content.UpdateMomentReq) (*content.UpdateMomentResp, error) {
	m := req.Moment
	momentId, err := primitive.ObjectIDFromHex(m.Id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	err = s.MomentMongoMapper.Update(ctx, &moment.Moment{
		ID:          momentId,
		CatId:       m.CatId,
		CommunityId: m.CommunityId,
		Photos:      m.Photos,
		Title:       m.Title,
		Text:        m.Text,
		UserId:      m.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &content.UpdateMomentResp{}, nil
}

func (s *MomentService) DeleteMoment(ctx context.Context, req *content.DeleteMomentReq) (*content.DeleteMomentResp, error) {
	err := s.MomentMongoMapper.Delete(ctx, req.MomentId)
	if err != nil {
		return nil, err
	}
	return &content.DeleteMomentResp{}, nil
}
