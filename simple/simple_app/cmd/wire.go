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
	"github.com/google/wire"

	"github.com/origadmin/runtime"
	datav1 "github.com/origadmin/runtime/api/gen/go/runtime/data/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"

	"basic-layout/simple/simple_app/internal/biz"
	"basic-layout/simple/simple_app/internal/server"

	"basic-layout/simple/simple_app/internal/conf"
	"basic-layout/simple/simple_app/internal/data"
	"basic-layout/simple/simple_app/internal/service"
)

// providerSet for components provided by the runtime.
var runtimeProviderSet = wire.NewSet(
	provideLogger,
	provideConfig,
	provideServerConfig,
	provideDataConfig,
)

// provideLogger extracts the logger from the runtime instance.
func provideLogger(rt *runtime.Runtime) kratoslog.Logger {
	return rt.Logger()
}

// provideConfig extracts and decodes the bootstrap config from the runtime instance.
func provideConfig(rt *runtime.Runtime) (*conf.Bootstrap, error) {
	var bc conf.Bootstrap
	if err := rt.Config().Decode("", &bc); err != nil { // Changed Scan to Decode
		return nil, err
	}
	return &bc, nil
}

// provideServerConfig extracts the server config from the bootstrap config.
func provideServerConfig(bc *conf.Bootstrap) *transportv1.Servers {
	return bc.GetServers()
}

// provideDataConfig extracts the data config from the bootstrap config.
func provideDataConfig(bc *conf.Bootstrap) *datav1.Data {
	return bc.GetData()
}

// NewKratosApp creates the final kratos.App from the runtime and transport servers.
func NewKratosApp(rt *runtime.Runtime, servers []transport.Server) *kratos.App {
	return rt.NewApp(servers)
}

// wireApp initializes the application using wire.
func wireApp(rt *runtime.Runtime) (*kratos.App, func(), error) {
	panic(wire.Build(
		runtimeProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		NewKratosApp,
	))
}
