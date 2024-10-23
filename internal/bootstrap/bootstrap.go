package bootstrap

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/internal/mods/helloworld/service"
)

var (
	ProviderSet = wire.NewSet(NewRegistrar, NewDiscovery)
)

func InjectorGinServer(injector *mods.InjectorClient) error {
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
	_ = grpcClient
	_ = httpClient
	helloworld.RegisterGreeterGINServer(injector.ServerGINS, httpClient)
	return nil
}
