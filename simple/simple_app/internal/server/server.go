package server

import (
	"context"
	"errors"
	stdhttp "net/http"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/container"
	runtimeservice "github.com/origadmin/runtime/service"
	runtimetransport "github.com/origadmin/runtime/service/transport"
	"github.com/origadmin/runtime/service/transport/grpc"
	"github.com/origadmin/runtime/service/transport/http"

	simplev1 "basic-layout/simple/simple_app/api/gen/go/simple/v1"
	"basic-layout/simple/simple_app/internal/service"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

type serviceInstance struct {
	SimpleService *service.SimpleService
	logger        *log.Helper
}

func (s serviceInstance) Register(ctx context.Context, srv any) error {
	switch v := srv.(type) {
	case *runtimetransport.GRPCServer:
		simplev1.RegisterSimpleServiceServer(v, s.SimpleService)
	case *runtimetransport.HTTPServer:
		simplev1.RegisterSimpleServiceHTTPServer(v, s.SimpleService)
		err := v.WalkHandle(func(method, path string, handler stdhttp.HandlerFunc) {
			logger.Infof("HTTP %s %s", method, path)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewServer(rt *runtime.App, simple *service.SimpleService) ([]transport.Server, error) {
	provider, err := rt.Container().Middleware()
	if err != nil {
		return nil, err
	}
	middlewares, err := provider.ServerMiddlewares()
	if err != nil {
		return nil, err
	}
	serverConfigs, err := rt.StructuredConfig().DecodeServers()
	if err != nil {
		return nil, err
	}
	configs := serverConfigs.GetConfigs()
	var servers []transport.Server
	for _, config := range configs {
		srv, err := runtimeservice.NewServer(config,
			container.WithContainer(rt.Container()),
			grpc.WithServerMiddlewares(middlewares),
			http.WithServerMiddlewares(middlewares),
			runtimeservice.WithServerRegistrar(&serviceInstance{
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
func NewHTTPServer(c *transportv1.Servers, rt *runtime.App, simple *service.SimpleService, logger log.Logger) (*runtimetransport.HTTPServer, error) {
	middlewares := getMiddlewares(rt)
	var opts = []runtimetransport.HTTPServerOption{
		runtimetransport.MiddlewareHTTP(middlewares...),
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
		ServerOptions:     opts,
		CorsOptions:       nil,
		Registrar:         nil,
		ServerMiddlewares: nil,
		Context:           nil,
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
func NewGRPCServer(c *transportv1.Servers, rt *runtime.App, simple *service.SimpleService, logger log.Logger) (*runtimetransport.GRPCServer, error) {
	middlewares := getMiddlewares(rt)
	var opts = []runtimetransport.GRPCServerOption{
		runtimetransport.MiddlewareGRPC(middlewares...),
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
		ServerOptions:     opts,
		Context:           nil,
		Registrar:         nil,
		ServerMiddlewares: nil,
	})
	if err != nil {
		return nil, err
	}

	simplev1.RegisterSimpleServiceServer(srv, simple)
	return srv, nil
}

func getMiddlewares(rt *runtime.App) []middleware.Middleware {
	container := rt.Container()
	middlewaresMap, err := container.Middleware()
	if err != nil {
		logger.Fatalf("failed to get middlewares: %v", err)
	}
	middlewaresNames := middlewaresMap.Names()
	serverMiddlewares, err := middlewaresMap.ServerMiddlewares()
	if err != nil {
		return nil
	}

	// Build middleware slice in order
	middlewares := make([]middleware.Middleware, 0, len(middlewaresNames))
	for _, k := range middlewaresNames {
		middlewares = append(middlewares, serverMiddlewares[k])
	}
	return middlewares
}
