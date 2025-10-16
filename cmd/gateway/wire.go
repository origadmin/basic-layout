//go:build wireinject
// +build wireinject

/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/interfaces"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/internal/mods/gateway/client"
	"origadmin/basic-layout/internal/mods/gateway/server"
	"origadmin/basic-layout/internal/mods/gateway/service"
)

// provideRuntimeConfig extracts the runtime.Config interface from the runtime instance.
func provideRuntimeConfig(rt *runtime.Runtime) interfaces.Config {
	// rt.Config() returns interfaces.Config, which should implement runtimeConfig.Config
	return rt.Config()
}

// providerSet for components provided by the runtime.
var runtimeProviderSet = wire.NewSet(
	provideLogger,
	provideConfig,
	provideRuntimeConfig, // Add this to provide the runtime.Config interface
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

// NewKratosApp creates the final kratos.App from the runtime and transport servers.
func NewKratosApp(rt *runtime.Runtime, hs *http.Server, gs *grpc.Server) *kratos.App {
	servers := []transport.Server{hs, gs}
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
