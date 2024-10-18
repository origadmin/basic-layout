package server

import (
	"fmt"
	"net/netip"
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods/helloworld/conf"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bootstrap *conf.Bootstrap, greeter helloworld.GreeterServiceServer, logger log.Logger) *grpc.Server {
	c := bootstrap.Server
	if c.Grpc == nil {
		c.Grpc = new(conf.Server_GRPC)
	}
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	if c.Middleware == nil {
		c.Middleware = new(conf.Server_Middleware)
	}

	naip, _ := netip.ParseAddrPort(bootstrap.Server.Grpc.Addr)
	endpoint, _ := url.Parse(fmt.Sprintf("grpc://192.168.28.81:%d", naip.Port()))
	opts = append(opts, grpc.Endpoint(endpoint))
	srv := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServiceServer(srv, greeter)
	return srv
}
