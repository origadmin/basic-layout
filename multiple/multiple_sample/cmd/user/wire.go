//go:build wireinject && !GOWORK
// +build wireinject,!GOWORK

/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// The build tag makes sure the stub is not built in the final build.
//go:generate go run -mod=mod github.com/google/wire/cmd/wire gen
package main

import (
	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"

	confpb "basic-layout/multiple/multiple_sample/internal/conf/pb"
	"basic-layout/multiple/multiple_sample/internal/features/user/biz"
	"basic-layout/multiple/multiple_sample/internal/features/user/data"
	"basic-layout/multiple/multiple_sample/internal/features/user/server"
	"basic-layout/multiple/multiple_sample/internal/features/user/service"
	"github.com/origadmin/runtime"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
)

// provideLogger extracts the logger from the runtime instance.
func provideLogger(rt *runtime.Runtime) kratoslog.Logger {
	return rt.Logger()
}

// provideConfig extracts and decodes the bootstrap config from the runtime instance.
func provideConfig(rt *runtime.Runtime) (*confpb.Bootstrap, error) {
	var bc confpb.Bootstrap
	if err := rt.Config().Decode("", &bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// provideServerConfig extracts the server config from the bootstrap config.
func provideServerConfig(bc *confpb.Bootstrap) *transportv1.Servers {
	return bc.GetServers()
}

// provideDataConfig extracts the data config from the bootstrap config.
func provideDataConfig(bc *confpb.Bootstrap) *datav1.Data {
	return bc.GetData()
}

// providerSet for components provided by the runtime.
var runtimeProviderSet = wire.NewSet(
	provideLogger,
	provideConfig,
	provideServerConfig,
	provideDataConfig,
)

// NewKratosApp creates the final kratos.App from the runtime and transport servers.
func NewKratosApp(rt *runtime.Runtime, servers []transport.Server) *kratos.App {
	return rt.NewApp(servers)
}

// wireApp injects providers to initialize the application.
func wireApp(rt *runtime.Runtime) (*kratos.App, func(), error) {
	panic(wire.Build(
		runtimeProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		biz.ProviderSet,
		data.ProviderSet,
		NewKratosApp,
	))
}
