package provider

import (
	"github.com/google/wire"
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/config"

	"github.com/xh-polaris/meowchat-collection/biz/application/service"
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/mapper"
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)

var ApplicationSet = wire.NewSet(
	service.CatSet,
	service.ImageSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	MapperSet,
)

var MapperSet = wire.NewSet(
	mapper.CatSet,
	mapper.ImageSet,
)
