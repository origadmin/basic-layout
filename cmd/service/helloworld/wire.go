//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/origadmin/basic-layout/internal/mods/helloworld/biz"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/conf"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/dal"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/server"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

//go:generate go run -mod=mod --tags wireinject github.com/google/wire/cmd/wire

// buildApp init kratos application.
func buildApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, dal.ProviderSet, biz.ProviderSet, service.ProviderSet, NewApp))
}
