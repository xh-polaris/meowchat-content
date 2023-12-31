package provider

import (
	"github.com/google/wire"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/config"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/cat"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/donate"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/fish"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/image"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/moment"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/plan"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mapper/post"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/mq"
	"github.com/xh-polaris/meowchat-content/biz/infrastructure/stores/redis"

	"github.com/xh-polaris/meowchat-content/biz/application/service"
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)

var ApplicationSet = wire.NewSet(
	service.CatSet,
	service.ImageSet,
	service.MomentSet,
	service.PostSet,
	service.PlanSet,
	service.IncentiveSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedis,
	mq.NewMqProducer,
	MapperSet,
)

var MapperSet = wire.NewSet(
	cat.NewMongoMapper,
	cat.NewEsMapper,
	image.NewMongoMapper,
	moment.NewMongoMapper,
	moment.NewEsMapper,
	post.NewMongoMapper,
	post.NewEsMapper,
	plan.NewMongoMapper,
	plan.NewEsMapper,
	fish.NewMongoMapper,
	donate.NewMongoMapper,
)
