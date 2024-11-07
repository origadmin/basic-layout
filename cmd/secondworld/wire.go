//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/internal/mods/secondworld/biz"
	"origadmin/basic-layout/internal/mods/secondworld/dal"
	"origadmin/basic-layout/internal/mods/secondworld/server"
	"origadmin/basic-layout/internal/mods/secondworld/service"
)

// buildInjectors init kratos application.
func buildInjectors(context.Context, *configs.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(bootstrap.ProviderSet, server.ProviderSet, dal.ProviderSet, biz.ProviderSet, service.ProviderSet, NewApp))
}
