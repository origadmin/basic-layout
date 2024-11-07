package bootstrap

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/origadmin/toolkits/errors"
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

type InjectorClient struct {
	Bootstrap     *configs.Bootstrap
	Logger        log.Logger
	Discovery     registry.Discovery
	ServerGINS    *gins.Server
	ServerHTTP    *http.Server
	GreeterServer helloworld.GreeterAPIServer
}

type InjectorServer struct {
	Bootstrap  *configs.Bootstrap
	Logger     log.Logger
	Registrar  registry.Registrar
	ServerGRPC *grpc.Server
	ServerHTTP *http.Server
}

func InjectorGinServer(injector *InjectorClient) error {
	// Create route Filter: Filter instances whose version number is "2.0.0"
	filter := filter.Version("v1.0.0")
	// Create the Selector for the P2C load balancing algorithm and inject the route Filter
	selector.SetGlobalSelector(random.NewBuilder())
	//selector.SetGlobalSelector(wrr.NewBuilder())

	serviceName := "origadmin.service.v1.helloworld"
	discovery := injector.Discovery
	if discovery == nil {
		return errors.String("discovery is nil")
	}
	//if discovery, ok := injector.Discoveries[serviceName]; ok {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
		),
		grpc.WithEndpoint("discovery:///"+serviceName),
		grpc.WithDiscovery(discovery),
		grpc.WithNodeFilter(filter),
	)
	if err != nil {
		return err
	}
	gClient := helloworld.NewGreeterAPIClient(conn)
	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
		),
		http.WithEndpoint("discovery:///"+serviceName),
		http.WithDiscovery(discovery),
		http.WithNodeFilter(filter),
	)
	if err != nil {
		return err
	}
	hClient := helloworld.NewGreeterAPIHTTPClient(hConn)

	var client helloworld.GreeterAPIServer
	if entry := injector.Bootstrap.GetEntry(); entry != nil && entry.Scheme == "http" {
		client = service.NewGreeterHTTPServer(hClient)
	} else {
		client = service.NewGreeterServer(gClient)
	}
	//grpcClient := service.NewGreeterServer(gClient)
	//httpClient := service.NewGreeterHTTPServer(hClient)
	//// add _ to avoid unused
	//_ = grpcClient
	//_ = httpClient
	helloworld.RegisterGreeterAPIGINSServer(injector.ServerGINS, client)
	helloworld.RegisterGreeterAPIHTTPServer(injector.ServerHTTP, client)
	//}

	return nil
}
