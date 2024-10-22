package mods

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"origadmin/basic-layout/internal/configs"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Injector), "*"),
)

type Injector struct {
	Logger     log.Logger
	Registry   registry.Registrar
	Bootstrap  *configs.Bootstrap
	ServerGRPC *grpc.Server
	ServerHTTP *http.Server
	//ServerGINS *gins.Server
}

//func NewInjector(bootstrap *conf.Bootstrap, logger log.Logger) *Injector {
//	return &Injector{
//		Bootstrap: bootstrap,
//		Logger:    logger,
//	}
//}
