//go:build wireinject
// +build wireinject

/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// The build tag makes sure the stub is not built in the final build.
package start

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/client"
	"origadmin/basic-layout/internal/configs"
)

// buildInjectors init kratos application.
func buildInjectors(context.Context, *bootstrap.Config, *configs.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		bootstrap.ProviderSet,
		agent.ProviderSet,
		//helloworlddal.ProviderSet,
		//helloworldbiz.ProviderSet,
		//helloworldservice.ProviderSet,
		//secondworlddal.ProviderSet,
		//secondworldbiz.ProviderSet,
		//secondworldservice.ProviderSet,
		NewApp))
}
