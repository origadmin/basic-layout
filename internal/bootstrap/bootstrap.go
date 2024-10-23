package bootstrap

import (
	"github.com/google/wire"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods"
)

var (
	ProviderSet = wire.NewSet(NewRegistrar)
)

func InjectorGinServer(injector *mods.InjectorClient) {
	helloworld.RegisterGreeterGINServer(injector.ServerGINS, injector.GreeterServer)

}
