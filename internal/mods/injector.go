package mods

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(InjectorServer), "*"),
	wire.Struct(new(InjectorClient), "*"),
)

type InjectorClient struct {
	Logger        log.Logger
	Discovery     registry.Discovery
	Bootstrap     *configs.Bootstrap
	ServerGINS    *gins.Server
	GreeterServer helloworld.GreeterServer
}

type InjectorServer struct {
	Logger     log.Logger
	Registrar  registry.Registrar
	Bootstrap  *configs.Bootstrap
	ServerGRPC *grpc.Server
	ServerHTTP *http.Server
}
