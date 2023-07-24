//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/google/wire"

	"github.com/xh-polaris/meowchat-collection/biz/adaptor"
)

func NewCollectionServerImpl() (*adaptor.CollectionServerImpl, error) {
	wire.Build(
		wire.Struct(new(adaptor.CollectionServerImpl), "*"),
		AllProvider,
	)
	return nil, nil
}
