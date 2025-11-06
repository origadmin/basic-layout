//go:build wireinject
// +build wireinject

/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// The build tag makes sure the stub is not built in the final build.
//go:generate go run -mod=mod github.com/google/wire/cmd/wire gen
package main

import (
	"context" // Added

	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/gateway/client"
	"basic-layout/multiple/multiple_sample/internal/mods/gateway/server"
	"basic-layout/multiple/multiple_sample/internal/mods/gateway/service"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/interfaces"
)

// provideRuntimeConfig extracts the runtime.Config interface from the runtime instance.
func provideRuntimeConfig(rt *runtime.Runtime) interfaces.Config {
	// rt.Config() returns interfaces.Config, which should implement runtimeConfig.Config
	return rt.Config()
}

// provideContext provides a background context.
func provideContext() context.Context {
	return context.Background()
}

// providerSet for components provided by the runtime.
var runtimeProviderSet = wire.NewSet(
	provideLogger,
	provideConfig,
	provideRuntimeConfig,
	provideServerConfig,
	provideContext,
)

// provideLogger extracts the logger from the runtime instance.
func provideLogger(rt *runtime.Runtime) kratoslog.Logger {
	return rt.Logger()
}

// provideConfig extracts and decodes the bootstrap config from the runtime instance.
func provideConfig(rt *runtime.Runtime) (*configs.Bootstrap, error) {
	var bc configs.Bootstrap
	if err := rt.Config().Decode("", &bc); err != nil { // Changed Scan to Decode
		return nil, err
	}
	return &bc, nil
}

// provideServerConfig extracts the server config from the bootstrap config.
func provideServerConfig(bc *configs.Bootstrap) *transportv1.Servers {
	return bc.GetServers()
}

// NewKratosApp creates the final kratos.App from the runtime and transport servers.
func NewKratosApp(rt *runtime.Runtime, servers []transport.Server) *kratos.App {
	return rt.NewApp(servers)
}

// wireApp initializes the application using wire.
func wireApp(rt *runtime.Runtime) (*kratos.App, func(), error) {
	panic(wire.Build(
		runtimeProviderSet,
		server.ProviderSet,
		client.ProviderSet,
		service.ProviderSet,
		NewKratosApp,
	))
}
