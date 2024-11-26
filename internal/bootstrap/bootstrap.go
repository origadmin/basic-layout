package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"github.com/origadmin/toolkits/runtime/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/api/v1/services/secondworld"
	"origadmin/basic-layout/internal/agent"
	"origadmin/basic-layout/internal/configs"
)

var (
	ProviderSet = wire.NewSet(
		NewRegistrar,
		NewDiscovery,
		NewDiscoveries,
		wire.Struct(new(InjectorServer), "*"),
		wire.Struct(new(InjectorClient), "*"),
	)
)

type InjectorClient struct {
	Config      *Config
	Bootstrap   *configs.Bootstrap
	Discoveries map[string]registry.Discovery
	Logger      log.Logger
	ServerGINS  *gins.Server
	ServerHTTP  *http.Server
}

type InjectorServer struct {
	Bootstrap  *configs.Bootstrap
	Logger     log.Logger
	Registrar  registry.Registrar
	ServerGRPC *grpc.Server
	ServerHTTP *http.Server
}

func InjectorGinServer(injector *InjectorClient) error {
	for name, discovery := range injector.Discoveries {
		switch name {
		case "origadmin.service.v1.helloworld":
			cli, err := agent.NewHelloGreeterAPIServer(injector.Bootstrap, discovery)
			if err != nil {
				return err
			}
			helloworld.RegisterHelloGreeterAPIGINSServer(injector.ServerGINS, cli)
			helloworld.RegisterHelloGreeterAPIHTTPServer(injector.ServerHTTP, cli)
		case "origadmin.service.v1.secondworld":
			cli, err := agent.NewSecondGreeterAPIServer(injector.Bootstrap, discovery)
			if err != nil {
				return err
			}
			secondworld.RegisterSecondGreeterAPIGINSServer(injector.ServerGINS, cli)
			secondworld.RegisterSecondGreeterAPIHTTPServer(injector.ServerHTTP, cli)
		}
	}
	return nil
}
