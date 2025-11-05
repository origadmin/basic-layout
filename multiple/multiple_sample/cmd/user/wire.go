//go:build wireinject && !GOWORK
// +build wireinject,!GOWORK

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/user/biz"
	"basic-layout/multiple/multiple_sample/internal/mods/user/data"
	"basic-layout/multiple/multiple_sample/internal/mods/user/server"
	"basic-layout/multiple/multiple_sample/internal/mods/user/service"
	"github.com/origadmin/runtime"
)

// NewKratosApp creates the final kratos.App from the runtime and transport servers.
func NewKratosApp(rt *runtime.Runtime, hs *http.Server, gs *grpc.Server) *kratos.App {
	servers := []transport.Server{hs, gs}
	return rt.NewApp(servers)
}

// wireApp injects providers to initialize the application.
func wireApp(rt *runtime.Runtime, conf *configs.Bootstrap, logger kratoslog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// The order of ProviderSets matters.
		// Layers should be provided from the outside in: server -> service -> biz -> data.
		wire.FieldsOf(new(*configs.Bootstrap), "Service"),
		wire.FieldsOf(new(*configs.ServiceConfig), "Servers"),
		server.ProviderSet,
		service.ProviderSet,
		biz.ProviderSet,
		data.ProviderSet,
		NewKratosApp,
	))
}
