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
	"origadmin/basic-layout/internal/mods/helloworld/biz"
	"origadmin/basic-layout/internal/mods/helloworld/dal"
	"origadmin/basic-layout/internal/mods/helloworld/server"
	"origadmin/basic-layout/internal/mods/helloworld/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// buildInjectors init kratos application.
func buildInjectors(contextContext context.Context, configsBootstrap *configs.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	registrar := bootstrap.NewRegistrar(configsBootstrap)
	database, cleanup, err := dal.NewDB(configsBootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterDao := dal.NewGreeterDal(database, logger)
	greeterAPIClient := biz.NewGreeterClient(greeterDao, logger)
	greeterAPIServer := service.NewGreeterServer(greeterAPIClient)
	grpcServer := server.NewGRPCServer(configsBootstrap, greeterAPIServer, logger)
	httpServer := server.NewHTTPServer(configsBootstrap, greeterAPIServer, logger)
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
