//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"

	"origadmin/basic-layout/internal/mods/helloworld/biz"
	"origadmin/basic-layout/internal/mods/helloworld/conf"
	"origadmin/basic-layout/internal/mods/helloworld/dal"
	"origadmin/basic-layout/internal/mods/helloworld/server"
	"origadmin/basic-layout/internal/mods/helloworld/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

//go:generate go run -mod=mod --tags wireinject github.com/google/wire/cmd/wire

// buildInjectors init kratos application.
func buildInjectors(context.Context, *conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, dal.ProviderSet, biz.ProviderSet, service.ProviderSet, NewApp))
}
