package svc

import (
	"github.com/xh-polaris/meowchat-collection-rpc/internal/config"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
)

type ServiceContext struct {
	Config   config.Config
	CatModel model.CatModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		CatModel: model.NewCatModel(c.Mongo.URL, c.Mongo.DB, c.Cache, c.Elasticsearch),
	}
}
