package main

import (
	"net"

	"github.com/xh-polaris/meowchat-collection/biz/infrastructure/util/log"
	"github.com/xh-polaris/meowchat-collection/provider"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/collection/collection"
)

func main() {
	s, err := provider.NewCollectionServerImpl()
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", s.ListenOn)
	if err != nil {
		panic(err)
	}
	svr := collection.NewServer(
		s,
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: s.Name}),
	)

	err = svr.Run()

	if err != nil {
		log.Error(err.Error())
	}
}
