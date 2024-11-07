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
	"origadmin/basic-layout/internal/mods/helloworld/biz"
	"origadmin/basic-layout/internal/mods/helloworld/dal"
	"origadmin/basic-layout/internal/mods/helloworld/service"
	biz2 "origadmin/basic-layout/internal/mods/secondworld/biz"
	dal2 "origadmin/basic-layout/internal/mods/secondworld/dal"
	service2 "origadmin/basic-layout/internal/mods/secondworld/service"
	"origadmin/basic-layout/internal/mods/server"
)

// Injectors from wire.go:

// buildInjectors init kratos application.
func buildInjectors(contextContext context.Context, configsBootstrap *configs.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	discovery := bootstrap.NewDiscovery(configsBootstrap)
	ginsServer := server.NewGINSServer(configsBootstrap, logger)
	httpServer := server.NewHTTPServer(configsBootstrap, logger)
	database, cleanup, err := dal.NewDB(configsBootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterDao := dal.NewGreeterDal(database, logger)
	helloGreeterAPIClient := biz.NewGreeterClient(greeterDao, logger)
	helloGreeterAPIServer := service.NewGreeterServer(helloGreeterAPIClient)
	dalDatabase, cleanup2, err := dal2.NewDB(configsBootstrap, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	dtoGreeterDao := dal2.NewGreeterDal(dalDatabase, logger)
	secondGreeterAPIClient := biz2.NewGreeterClient(dtoGreeterDao, logger)
	secondGreeterAPIServer := service2.NewGreeterServer(secondGreeterAPIClient)
	injectorClient := &bootstrap.InjectorClient{
		Bootstrap:           configsBootstrap,
		Logger:              logger,
		Discovery:           discovery,
		ServerGINS:          ginsServer,
		ServerHTTP:          httpServer,
		HelloGreeterServer:  helloGreeterAPIServer,
		SecondGreeterServer: secondGreeterAPIServer,
	}
	app := NewApp(contextContext, injectorClient)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
