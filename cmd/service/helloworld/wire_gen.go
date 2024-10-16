// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"origadmin/basic-layout/internal/config"
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/internal/mods/helloworld/biz"
	"origadmin/basic-layout/internal/mods/helloworld/conf"
	"origadmin/basic-layout/internal/mods/helloworld/dal"
	"origadmin/basic-layout/internal/mods/helloworld/server"
	"origadmin/basic-layout/internal/mods/helloworld/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// buildInjectors init kratos application.
func buildInjectors(contextContext context.Context, bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	registrar := bootstrap.NewRegistrar(bootstrap, logger)
	database, cleanup, err := dal.NewDB(bootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterDao := dal.NewGreeterDal(database, logger)
	greeterServiceClient := biz.NewGreeterClient(greeterDao, logger)
	greeterServiceServer := service.NewGreeterServer(greeterServiceClient)
	grpcServer := server.NewGRPCServer(bootstrap, greeterServiceServer, logger)
	httpServer := server.NewHTTPServer(bootstrap, greeterServiceServer, logger)
	injector := &mods.Injector{
		Logger:     logger,
		Registry:   registrar,
		Bootstrap:  bootstrap,
		ServerGRPC: grpcServer,
		ServerHTTP: httpServer,
	}
	app := NewApp(contextContext, injector)
	return app, func() {
		cleanup()
	}, nil
}
