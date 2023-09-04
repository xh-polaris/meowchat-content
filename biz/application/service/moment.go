package service

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	mqprimitive "github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/xh-polaris/gopkg/pagination/esp"
	"github.com/xh-polaris/gopkg/pagination/mongop"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/consts"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/moment"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/convertor"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"strconv"
	"time"
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
	Config            *config.Config
	MomentMongoMapper moment.IMongoMapper
	MomentEsMapper    moment.IEsMapper
	Redis             *redis.Redis
	MqProducer        rocketmq.Producer
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
		moments, total, err = s.MomentMongoMapper.FindManyAndCount(ctx, filter, p, mongop.IdCursorType)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := req.SearchOptions.Type.(type) {
		case *content.SearchOptions_AllFieldsKey:
			moments, total, err = s.MomentEsMapper.Search(ctx, convertor.ConvertMomentAllFieldsSearchQuery(o), filter, p, esp.ScoreCursorType)
		case *content.SearchOptions_MultiFieldsKey:
			moments, total, err = s.MomentEsMapper.Search(ctx, convertor.ConvertMomentMultiFieldsSearchQuery(o), filter, p, esp.ScoreCursorType)
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
	resp := new(content.CreateMomentResp)
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

	resp.MomentId = data.ID.Hex()

	//发送使用url信息
	var urls []url.URL
	for _, u := range m.Photos {
		sendUrl, _ := url.Parse(u)
		urls = append(urls, *sendUrl)
	}
	go s.SendDelayMessage(urls)

	//小鱼干奖励
	t, err := s.Redis.GetCtx(ctx, "contentTimes"+m.UserId)
	if err != nil {
		return resp, nil
	}
	r, err := s.Redis.GetCtx(ctx, "contentDate"+m.UserId)
	if err != nil {
		return resp, nil
	} else if r == "" {
		resp.GetFish = true
		resp.GetFishTimes = 1
		err = s.Redis.SetexCtx(ctx, "contentTimes"+m.UserId, "1", 86400)
		if err != nil {
			resp.GetFish = false
			return resp, nil
		}
		err = s.Redis.SetexCtx(ctx, "contentDate"+m.UserId, strconv.FormatInt(time.Now().Unix(), 10), 86400)
		if err != nil {
			resp.GetFish = false
			return resp, nil
		}
	} else {
		times, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return resp, nil
		}
		resp.GetFishTimes = times + 1
		date, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			return resp, nil
		}
		lastTime := time.Unix(date, 0)
		err = s.Redis.SetexCtx(ctx, "contentTimes"+m.UserId, strconv.FormatInt(times+1, 10), 86400)
		if err != nil {
			return resp, nil
		}
		err = s.Redis.SetexCtx(ctx, "contentDate"+m.UserId, strconv.FormatInt(time.Now().Unix(), 10), 86400)
		if err != nil {
			return resp, nil
		}
		if lastTime.Day() == time.Now().Day() && lastTime.Month() == time.Now().Month() && lastTime.Year() == time.Now().Year() {
			err = s.Redis.SetexCtx(ctx, "contentTimes"+m.UserId, strconv.FormatInt(times+1, 10), 86400)
			if err != nil {
				return resp, nil
			}
			if times >= s.Config.GetFishTimes {
				resp.GetFish = false
			} else {
				resp.GetFish = true
			}
		} else {
			err = s.Redis.SetexCtx(ctx, "contentTimes"+m.UserId, "1", 86400)
			if err != nil {
				return resp, nil
			}
			resp.GetFish = true
			resp.GetFishTimes = 1
		}
	}
	return resp, nil
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

	//发送使用url信息
	var urls []url.URL
	for _, u := range m.Photos {
		sendUrl, _ := url.Parse(u)
		urls = append(urls, *sendUrl)
	}
	go s.SendDelayMessage(urls)

	return &content.UpdateMomentResp{}, nil
}

func (s *MomentService) DeleteMoment(ctx context.Context, req *content.DeleteMomentReq) (*content.DeleteMomentResp, error) {
	err := s.MomentMongoMapper.Delete(ctx, req.MomentId)
	if err != nil {
		return nil, err
	}
	return &content.DeleteMomentResp{}, nil
}

func (s *MomentService) SendDelayMessage(message interface{}) {
	json, _ := sonic.Marshal(message)
	msg := &mqprimitive.Message{
		Topic: "sts_used_url",
		Body:  json,
	}

	res, err := s.MqProducer.SendSync(context.Background(), msg)
	if err != nil || res.Status != mqprimitive.SendOK {
		for i := 0; i < 2; i++ {
			res, err := s.MqProducer.SendSync(context.Background(), msg)
			if err == nil && res.Status == mqprimitive.SendOK {
				break
			}
		}
	}
}
