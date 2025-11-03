package server

import (
	"errors"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/google/wire"

	runtimeservice "github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/transport/grpc"
	"github.com/origadmin/runtime/service/transport/http"

	simplev1 "basic-layout/simple/simple_app/api/gen/go/simple/v1"
	"basic-layout/simple/simple_app/internal/service"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *transportv1.Servers, simple *service.SimpleService, logger log.Logger) (*runtimeservice.HTTPServer, error) {
	var opts = []runtimeservice.HTTPServerOption{
		runtimeservice.MiddlewareHTTP(
			recovery.Recovery(),
		),
	}
	var config *transportv1.Server
	for _, server := range c.GetConfigs() {
		if server.GetProtocol() == "http" {
			config = server
			break
		}
	}
	if config.GetHttp() == nil {
		return nil, errors.New("http server not found")
	}

	srv, err := http.NewServer(config.GetHttp(), &http.ServerOptions{
		ServiceOptions:    nil,
		HttpServerOptions: opts,
		CorsOptions:       nil,
	})
	if err != nil {
		return nil, err
	}
	helper := log.NewHelper(logger)
	simplev1.RegisterSimpleServiceHTTPServer(srv, simple)

	srv.WalkHandle(func(method, path string, handler stdhttp.HandlerFunc) {
		helper.Infof("HTTP %s %s", method, path)
	})
	return srv, nil
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *transportv1.Servers, simple *service.SimpleService, logger log.Logger) (*runtimeservice.GRPCServer, error) {
	var opts = []runtimeservice.GRPCServerOption{
		runtimeservice.MiddlewareGRPC(
			recovery.Recovery(),
		),
	}
	var config *transportv1.Server
	for _, server := range c.GetConfigs() {
		if server.GetProtocol() == "grpc" {
			config = server
			break
		}
	}
	if config.GetGrpc() == nil {
		return nil, errors.New("grpc server not found")
	}

	srv, err := grpc.NewServer(config.GetGrpc(), &grpc.ServerOptions{
		ServiceOptions:    nil,
		GrpcServerOptions: opts,
	})
	if err != nil {
		return nil, err
	}

	simplev1.RegisterSimpleServiceServer(srv, simple)
	return srv, nil
}
