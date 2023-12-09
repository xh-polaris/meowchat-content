// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	"github.com/xh-polaris/meowchat-content/biz/adaptor"
	"github.com/xh-polaris/meowchat-content/biz/application/service"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/cat"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/donate"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/fish"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/image"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/moment"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/plan"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/post"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/stores/redis"
)

// Injectors from wire.go:

func NewContentServerImpl() (*adaptor.ContentServerImpl, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	iMongoMapper := cat.NewMongoMapper(configConfig)
	iEsMapper := cat.NewEsMapper(configConfig)
	catService := &service.CatService{
		CatMongoMapper: iMongoMapper,
		CatEsMapper:    iEsMapper,
	}
	imageIMongoMapper := image.NewMongoMapper(configConfig)
	imageService := &service.ImageService{
		ImageModel: imageIMongoMapper,
	}
	momentIMongoMapper := moment.NewMongoMapper(configConfig)
	momentIEsMapper := moment.NewEsMapper(configConfig)
	redisRedis := redis.NewRedis(configConfig)
	momentService := &service.MomentService{
		Config:            configConfig,
		MomentMongoMapper: momentIMongoMapper,
		MomentEsMapper:    momentIEsMapper,
		ImageMapper:       imageIMongoMapper,
		Redis:             redisRedis,
	}
	postIMongoMapper := post.NewMongoMapper(configConfig)
	postIEsMapper := post.NewEsMapper(configConfig)
	postService := &service.PostService{
		Config:          configConfig,
		PostMongoMapper: postIMongoMapper,
		PostEsMapper:    postIEsMapper,
		Redis:           redisRedis,
	}
	planIMongoMapper := plan.NewMongoMapper(configConfig)
	planIEsMapper := plan.NewEsMapper(configConfig)
	donateIMongoMapper := donate.NewMongoMapper(configConfig)
	fishIMongoMapper := fish.NewMongoMapper(configConfig)
	planService := &service.PlanService{
		PlanMongoMapper:   planIMongoMapper,
		PlanEsMapper:      planIEsMapper,
		DonateMongoMapper: donateIMongoMapper,
		FishMongoMapper:   fishIMongoMapper,
	}
	contentServerImpl := &adaptor.ContentServerImpl{
		Config:        configConfig,
		CatService:    catService,
		ImageService:  imageService,
		MomentService: momentService,
		PostService:   postService,
		PlanService:   planService,
	}
	return contentServerImpl, nil
}
