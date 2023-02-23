package svc

import (
	"github.com/xh-polaris/meowchat-collection-rpc/internal/config"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
)

type ServiceContext struct {
	Config     config.Config
	CatModel   model.CatModel
	ImageModel model.ImageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		CatModel:   model.NewCatModel(c.Mongo.URL, c.Mongo.DB, c.Cache, c.Elasticsearch),
		ImageModel: model.NewImageModel(c.Mongo.URL, c.Mongo.DB, c.Cache),
	}
}
