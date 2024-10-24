// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package start

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/internal/mods/helloworld/biz"
	"origadmin/basic-layout/internal/mods/helloworld/dal"
	"origadmin/basic-layout/internal/mods/helloworld/service"
	"origadmin/basic-layout/internal/mods/server"
)

// Injectors from wire.go:

// buildInjectors init kratos application.
func buildInjectors(contextContext context.Context, configsBootstrap *configs.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	discovery := bootstrap.NewDiscovery(configsBootstrap, logger)
	ginsServer := server.NewGINSServer(configsBootstrap, logger)
	httpServer := server.NewHTTPServer(configsBootstrap, logger)
	database, cleanup, err := dal.NewDB(configsBootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterDao := dal.NewGreeterDal(database, logger)
	greeterClient := biz.NewGreeterClient(greeterDao, logger)
	greeterServer := service.NewGreeterServer(greeterClient)
	injectorClient := &mods.InjectorClient{
		Logger:        logger,
		Discovery:     discovery,
		Bootstrap:     configsBootstrap,
		ServerGINS:    ginsServer,
		ServerHTTP:    httpServer,
		GreeterServer: greeterServer,
	}
	app := NewApp(contextContext, injectorClient)
	return app, func() {
		cleanup()
	}, nil
}
