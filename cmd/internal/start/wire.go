//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package start

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
	//helloworldbiz "origadmin/basic-layout/internal/mods/helloworld/biz"
	//helloworlddal "origadmin/basic-layout/internal/mods/helloworld/dal"
	//helloworldservice "origadmin/basic-layout/internal/mods/helloworld/service"
	//secondworldbiz "origadmin/basic-layout/internal/mods/secondworld/biz"
	//secondworlddal "origadmin/basic-layout/internal/mods/secondworld/dal"
	//secondworldservice "origadmin/basic-layout/internal/mods/secondworld/service"
	"origadmin/basic-layout/internal/mods/server"
)

// buildInjectors init kratos application.
func buildInjectors(context.Context, *bootstrap.Config, *configs.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		bootstrap.ProviderSet,
		server.ProviderSet,
		//helloworlddal.ProviderSet,
		//helloworldbiz.ProviderSet,
		//helloworldservice.ProviderSet,
		//secondworlddal.ProviderSet,
		//secondworldbiz.ProviderSet,
		//secondworldservice.ProviderSet,
		NewApp))
}
