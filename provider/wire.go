//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/google/wire"

	"github.com/xh-polaris/meowchat-content/biz/adaptor"
)

func NewContentServerImpl() (*adaptor.ContentServerImpl, error) {
	wire.Build(
		wire.Struct(new(adaptor.ContentServerImpl), "*"),
		AllProvider,
	)
	return nil, nil
}
