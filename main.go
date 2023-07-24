package main

import (
	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/util/log"
	"github.com/xh-polaris/meowchat-collection/provider"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/collection/collection"
)

func main() {
	s, err := provider.NewCollectionServerImpl()
	if err != nil {
		panic(err)
	}
	svr := collection.NewServer(s)

	err = svr.Run()

	if err != nil {
		log.Error(err.Error())
	}
}
