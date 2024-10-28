package server

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bs *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	c := bs.Server
	if c.Grpc == nil {
		c.Grpc = new(configs.Server_GRPC)
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

	middlewares, err := bootstrap.LoadMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, grpc.Middleware(middlewares...))
	}

	naip, _ := netip.ParseAddrPort(bs.Server.Grpc.Addr)
	if bs.Server.Grpc.Endpoint == "" {
		bs.Server.Grpc.Endpoint = "grpc://" + bs.Server.Host + ":" + strconv.Itoa(int(naip.Port()))
	} else {
		prefix, suffix, ok := strings.Cut(bs.Server.Grpc.Endpoint, "://")
		if !ok {
			bs.Server.Grpc.Endpoint = "grpc://" + prefix
		} else {
			args := strings.SplitN(suffix, ":", 2)
			if len(args) == 2 {
				args[1] = strconv.Itoa(int(naip.Port()))
			} else if len(args) == 1 {
				args = append(args, strconv.Itoa(int(naip.Port())))
			} else {
				// unknown
				log.Infow("unknown http endpoint", bs.Server.Grpc.Endpoint)
			}
			bs.Server.Grpc.Endpoint = prefix + "://" + strings.Join(args, ":")
		}
	}

	log.Infof("Server.Grpc.Endpoint: %v", bs.Server.Grpc.Endpoint)
	ep, _ := url.Parse(bs.Server.Grpc.Endpoint)
	opts = append(opts, grpc.Endpoint(ep))
	srv := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServer(srv, greeter)
	return srv
}
