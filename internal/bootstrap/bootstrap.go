package bootstrap

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/internal/mods/helloworld/service"
)

var (
	ProviderSet = wire.NewSet(
		NewRegistrar,
		NewDiscovery,
		wire.Struct(new(InjectorServer), "*"),
		wire.Struct(new(InjectorClient), "*"),
	)
)

func InjectorGinServer(injector *InjectorClient) error {
	// Create route Filter: Filter instances whose version number is "2.0.0"
	filter := filter.Version("v1.0.0")
	// Create the Selector for the P2C load balancing algorithm and inject the route Filter
	selector.SetGlobalSelector(random.NewBuilder())
	//selector.SetGlobalSelector(wrr.NewBuilder())

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
		grpc.WithEndpoint("discovery:///origadmin.service.v1.helloworld"),
		grpc.WithDiscovery(injector.Discovery),
		grpc.WithNodeFilter(filter),
	)
	if err != nil {
		return err
	}
	gClient := helloworld.NewGreeterClient(conn)

	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///origadmin.service.v1.helloworld"),
		http.WithDiscovery(injector.Discovery),
		http.WithNodeFilter(filter),
	)
	if err != nil {
		return err
	}
	hClient := helloworld.NewGreeterHTTPClient(hConn)
	grpcClient := service.NewGreeterServer(gClient)
	httpClient := service.NewGreeterHTTPServer(hClient)
	// add _ to avoid unused
	_ = grpcClient
	_ = httpClient
	helloworld.RegisterGreeterGINServer(injector.ServerGINS, httpClient)
	helloworld.RegisterGreeterHTTPServer(injector.ServerHTTP, httpClient)
	return nil
}

type InjectorClient struct {
	Logger        log.Logger
	Discovery     registry.Discovery
	Bootstrap     *configs.Bootstrap
	ServerGINS    *gins.Server
	ServerHTTP    *http.Server
	GreeterServer helloworld.GreeterServer
}

type InjectorServer struct {
	Logger     log.Logger
	Registrar  registry.Registrar
	Bootstrap  *configs.Bootstrap
	ServerGRPC *grpc.Server
	ServerHTTP *http.Server
}
