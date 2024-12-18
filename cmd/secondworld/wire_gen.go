// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/internal/mods/secondworld/biz"
	"origadmin/basic-layout/internal/mods/secondworld/dal"
	"origadmin/basic-layout/internal/mods/secondworld/server"
	"origadmin/basic-layout/internal/mods/secondworld/service"
)

import (
	_ "github.com/origadmin/contrib/consul/config"
	_ "github.com/origadmin/contrib/consul/registry"
)

// Injectors from wire.go:

// buildInjectors init kratos application.
func buildInjectors(contextContext context.Context, configsBootstrap *configs.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	registrar := bootstrap.NewRegistrar(configsBootstrap)
	database, cleanup, err := dal.NewDB(configsBootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := dal.NewGreeterDal(database, logger)
	secondGreeterAPIClient := biz.NewGreeterClient(greeterRepo, logger)
	secondGreeterAPIServer := service.NewGreeterServer(secondGreeterAPIClient)
	grpcServer := server.NewGRPCServer(configsBootstrap, secondGreeterAPIServer, logger)
	httpServer := server.NewHTTPServer(configsBootstrap, secondGreeterAPIServer, logger)
	injectorServer := &bootstrap.InjectorServer{
		Bootstrap:  configsBootstrap,
		Logger:     logger,
		Registrar:  registrar,
		ServerGRPC: grpcServer,
		ServerHTTP: httpServer,
	}
	app := NewApp(contextContext, injectorServer)
	return app, func() {
		cleanup()
	}, nil
}
