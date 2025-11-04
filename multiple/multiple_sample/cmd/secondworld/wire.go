//go:build wireinject
// +build wireinject

/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/internal/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/secondworld/biz"
	"basic-layout/multiple/multiple_sample/internal/mods/secondworld/dal"
	"basic-layout/multiple/multiple_sample/internal/mods/secondworld/server"
	"basic-layout/multiple/multiple_sample/internal/mods/secondworld/service"
	"github.com/origadmin/runtime"
)

// providerSet for components provided by the runtime.
var runtimeProviderSet = wire.NewSet(
	provideLogger,
	provideConfig,
)

// provideLogger extracts the logger from the runtime instance.
func provideLogger(rt *runtime.Runtime) kratoslog.Logger {
	return rt.Logger()
}

// provideConfig extracts and decodes the bootstrap config from the runtime instance.
func provideConfig(rt *runtime.Runtime) (*configs.Bootstrap, error) {
	var bc configs.Bootstrap
	if err := rt.Config().Decode("", &bc); err != nil { // Added missing key argument ""
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
		dal.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		NewKratosApp,
	))
}
