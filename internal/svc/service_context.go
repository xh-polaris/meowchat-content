package svc

import (
	"github.com/xh-polaris/meowchat-collection-rpc/internal/config"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	CatModel model.CatModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.Datasource)
	return &ServiceContext{
		Config:   c,
		CatModel: model.NewCatModel(conn, c.CacheRedis),
	}
}
