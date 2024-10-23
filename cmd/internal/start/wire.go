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
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/internal/mods/helloworld/biz"
	"origadmin/basic-layout/internal/mods/helloworld/dal"
	"origadmin/basic-layout/internal/mods/helloworld/service"
	"origadmin/basic-layout/internal/mods/server"
)

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

// buildInjectors init kratos application.
func buildInjectors(context.Context, *configs.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(mods.ProviderSet, bootstrap.ProviderSet, server.ProviderSet, dal.ProviderSet, biz.ProviderSet, service.ProviderSet, NewApp))
}
