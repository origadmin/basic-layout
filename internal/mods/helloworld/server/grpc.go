package server

import (
	"fmt"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bootstrap *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *grpc.Server {
	c := bootstrap.Server
	if c.Grpc == nil {
		c.Grpc = new(configs.Server_GRPC)
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
		c.Middleware = new(configs.Server_Middleware)
	}

	naip, _ := netip.ParseAddrPort(bootstrap.Server.Grpc.Addr)
	prefix, suffix, ok := strings.Cut(bootstrap.Server.Grpc.Endpoint, "://")
	if !ok {
		bootstrap.Server.Grpc.Endpoint = "grpc://" + prefix
	} else {
		args := strings.SplitN(suffix, ":", 2)
		if len(args) == 2 {
			args[1] = strconv.Itoa(int(naip.Port()))
		} else if len(args) == 1 {
			args = append(args, strconv.Itoa(int(naip.Port())))
		} else {
			// unknown
			log.NewHelper(l).Info("unknown grpc endpoint", bootstrap.Server.Grpc.Endpoint)
		}
		bootstrap.Server.Grpc.Endpoint = prefix + "://" + strings.Join(args, ":")
	}
	fmt.Println("bootstrap.Server.Grpc.Endpoint", bootstrap.Server.Grpc.Endpoint)
	ep, _ := url.Parse(bootstrap.Server.Grpc.Endpoint)
	opts = append(opts, grpc.Endpoint(ep))
	srv := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServer(srv, greeter)
	return srv
}
