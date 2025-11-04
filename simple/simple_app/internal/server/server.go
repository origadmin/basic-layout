package server

import (
	"context"
	"errors"
	stdhttp "net/http"
	"sort"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"

	"github.com/origadmin/runtime"
	runtimeservice "github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/transport/grpc"
	"github.com/origadmin/runtime/service/transport/http"

	simplev1 "basic-layout/simple/simple_app/api/gen/go/simple/v1"
	"basic-layout/simple/simple_app/internal/service"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

type serviceInstance struct {
	SimpleService *service.SimpleService
	logger        *log.Helper
}

func (s serviceInstance) Register(ctx context.Context, srv any) error {
	switch v := srv.(type) {
	case *runtimeservice.GRPCServer:
		simplev1.RegisterSimpleServiceServer(v, s.SimpleService)
	case *runtimeservice.HTTPServer:
		simplev1.RegisterSimpleServiceHTTPServer(v, s.SimpleService)
		v.WalkHandle(func(method, path string, handler stdhttp.HandlerFunc) {
			logger.Infof("HTTP %s %s", method, path)
		})
	}
	return nil
}

func NewServer(rt *runtime.Runtime, simple *service.SimpleService) ([]transport.Server, error) {
	serverConfigs, err := rt.StructuredConfig().DecodeServers()
	if err != nil {
		return nil, err
	}
	configs := serverConfigs.GetConfigs()
	var servers []transport.Server
	for _, config := range configs {
		srv, err := runtimeservice.NewServer(config,
			runtimeservice.WithContainer(rt.Container()),
			runtimeservice.WithRegistrar(&serviceInstance{
				logger:        log.NewHelper(rt.Logger()),
				SimpleService: simple,
			}))
		if err != nil {
			return nil, err
		}
		servers = append(servers, srv)
	}

	return servers, nil
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *transportv1.Servers, rt *runtime.Runtime, simple *service.SimpleService, logger log.Logger) (*runtimeservice.HTTPServer, error) {
	middlewares := getMiddlewares(rt)
	var opts = []runtimeservice.HTTPServerOption{
		runtimeservice.MiddlewareHTTP(middlewares...),
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
func NewGRPCServer(c *transportv1.Servers, rt *runtime.Runtime, simple *service.SimpleService, logger log.Logger) (*runtimeservice.GRPCServer, error) {
	middlewares := getMiddlewares(rt)
	var opts = []runtimeservice.GRPCServerOption{
		runtimeservice.MiddlewareGRPC(middlewares...),
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

func getMiddlewares(rt *runtime.Runtime) []middleware.Middleware {
	container := rt.Container()
	middlewaresMap := container.ServerMiddlewares()

	// Extract keys and sort them
	keys := make([]string, 0, len(middlewaresMap))
	for k := range middlewaresMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build middleware slice in order
	middlewares := make([]middleware.Middleware, 0, len(keys))
	for _, k := range keys {
		middlewares = append(middlewares, middlewaresMap[k])
	}
	return middlewares
}
