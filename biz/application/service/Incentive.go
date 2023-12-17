package service

import (
	"context"
	"strconv"
	"time"

	"github.com/google/wire"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util"
)

type IIncentiveService interface {
	GetMission(ctx context.Context, req *content.GetMissionReq) (*content.GetMissionResp, error)
	CheckIn(ctx context.Context, req *content.CheckInReq) (*content.CheckInResp, error)
}

type IncentiveService struct {
	Config *config.Config
	Redis  *redis.Redis
}

var IncentiveSet = wire.NewSet(
	wire.Struct(new(IncentiveService), "*"),
	wire.Bind(new(IIncentiveService), new(*IncentiveService)),
)

func (s *IncentiveService) GetMission(ctx context.Context, req *content.GetMissionReq) (*content.GetMissionResp, error) {
	resp := new(content.GetMissionResp)

	util.ParallelRun([]func(){
		func() {
			r, err := s.GetTimes(ctx, "contentTimes", "contentDate", req.GetUserId())
			if err != nil {
				return
			}
			resp.ContentTime = r
		},
		func() {
			r, err := s.GetTimes(ctx, "checkInTimes", "checkInDates", req.GetUserId())
			if err != nil {
				return
			}
			resp.SignInTime = r
		},
		func() {
			r, err := s.GetTimes(ctx, "likeTimes", "likeDates", req.GetUserId())
			if err != nil {
				return
			}
			resp.LikeTime = r
		},
		func() {
			r, err := s.GetTimes(ctx, "commentTimes", "commentDate", req.GetUserId())
			if err != nil {
				return
			}
			resp.CommentTime = r
		},
	})
	return resp, nil
}

func (s *IncentiveService) CheckIn(ctx context.Context, req *content.CheckInReq) (*content.CheckInResp, error) {

	res := new(content.CheckInResp)

	t, err := s.Redis.GetCtx(ctx, "checkInTimes"+req.UserId)
	if err != nil {
		return &content.CheckInResp{GetFish: false}, nil
	}
	r, err := s.Redis.GetCtx(ctx, "checkInDates"+req.UserId)
	if err != nil {
		return &content.CheckInResp{GetFish: false}, nil
	} else if r == "" {
		res.GetFish = true
		res.GetFishTimes = 1
		err = s.Redis.SetexCtx(ctx, "checkInTimes"+req.UserId, "1", 604800)
		if err != nil {
			return &content.CheckInResp{GetFish: false}, nil
		}
		err = s.Redis.SetexCtx(ctx, "checkInDates"+req.UserId, strconv.FormatInt(time.Now().Unix(), 10), 604800)
		if err != nil {
			return &content.CheckInResp{GetFish: false}, nil
		}
	} else {
		times, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return &content.CheckInResp{GetFish: false}, nil
		}
		date, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			return &content.CheckInResp{GetFish: false}, nil
		}
		lastTime := time.Unix(date, 0)
		err = s.Redis.SetexCtx(ctx, "checkInDates"+req.UserId, strconv.FormatInt(time.Now().Unix(), 10), 604800)
		if err != nil {
			return &content.CheckInResp{GetFish: false}, nil
		}
		if lastTime.Day() == time.Now().Day() && lastTime.Month() == time.Now().Month() && lastTime.Year() == time.Now().Year() {
			return &content.CheckInResp{GetFish: false}, nil
		}
		lastYear, lastWeek := lastTime.ISOWeek()
		nowYear, nowWeek := time.Now().ISOWeek()
		if lastWeek == nowWeek && lastYear == nowYear {
			res.GetFishTimes = times + 1
			err = s.Redis.SetexCtx(ctx, "checkInTimes"+req.UserId, strconv.FormatInt(times+1, 10), 604800)
			if err != nil {
				return &content.CheckInResp{GetFish: false}, nil
			}
			res.GetFish = true
		} else {
			err = s.Redis.SetexCtx(ctx, "checkInTimes"+req.UserId, "1", 604800)
			if err != nil {
				return &content.CheckInResp{GetFish: false}, nil
			}
			res.GetFish = true
			res.GetFishTimes = 1
		}
	}
	return res, nil
}

func (s *IncentiveService) GetTimes(ctx context.Context, reqTime string, reqDate string, userId string) (int64, error) {
	var res int64
	t, err := s.Redis.GetCtx(ctx, reqTime+userId)
	if err != nil {
		return 0, err
	}
	r, err := s.Redis.GetCtx(ctx, reqDate+userId)
	if err != nil {
		return 0, err
	} else if r == "" {
		res = 0
	} else {
		times, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, err
		}
		date, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			return 0, err
		}
		lastTime := time.Unix(date, 0)
		if lastTime.Day() == time.Now().Day() && lastTime.Month() == time.Now().Month() && lastTime.Year() == time.Now().Year() {
			res = times
		} else {
			res = 0
		}
	}
	return res, nil
}
