package main

import (
	"net"

	"github.com/xh-polaris/meowchat-content/biz/infrastructure/util/log"
	"github.com/xh-polaris/meowchat-content/provider"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	content "github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowchat/content/contentservice"
)

func main() {
	s, err := provider.NewContentServerImpl()
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", s.ListenOn)
	if err != nil {
		panic(err)
	}
	svr := content.NewServer(
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
