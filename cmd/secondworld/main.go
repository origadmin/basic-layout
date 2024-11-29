/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"context"
	"flag"
	"syscall"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	_ "github.com/origadmin/contrib/consul/config"
	_ "github.com/origadmin/contrib/consul/registry"
	logger "github.com/origadmin/slog-kratos"

	"origadmin/basic-layout/internal/bootstrap"
)

// go build -ldflags "-X main.Version=vx.y.z -X main.Name=origadmin.service.v1.secondworld"
var (
	// Name is the Name of the compiled software.
	Name = "origadmin.service.v1.secondworld"
	// Version is the Version of the compiled software.
	Version = "v1.0.0"
	// boot are the bootstrap boot.
	boot = &bootstrap.Bootstrap{}
)

func init() {
	boot.SetFlags(Name, Version)
	flag.StringVar(&boot.ConfigPath, "c", "resources", "config path, eg: -c config.toml")
}

func main() {
	flag.Parse()

	l := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", boot.ID(),
		"service.name", boot.ServiceName(),
		"service.version", boot.Version(),
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(l)
	log.Infof("bootstrap flags: %+v\n", boot)
	bs, err := bootstrap.FromFlags(boot, bootstrap.WithLogger(l))
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}
	//if err := bootstrap.SyncRegistry(boot.ServiceName(), bs); err != nil {
	//	log.Warnf("sync config to consul error: %s", err.Error())
	//	return
	//}
	log.Infof("bootstrap config: %+v\n", bootstrap.PrintString(bs))
	ctx := context.Background()
	//info to ctx
	app, cleanup, err := buildInjectors(ctx, bs, l)
	if err != nil {
		log.Fatalf("failed to build injector: %s", err.Error())
	}
	defer cleanup()
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Fatalf("app stopped with error: %s", err.Error())
	}
}

func NewApp(ctx context.Context, injector *bootstrap.InjectorServer) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(boot.ID()),
		kratos.Name(boot.ServiceName()),
		kratos.Version(boot.Version()),
		kratos.Metadata(boot.Metadata()),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(injector.Logger),
		kratos.Server(injector.ServerHTTP, injector.ServerGRPC),
	}
	if injector.Registrar != nil {
		opts = append(opts, kratos.Registrar(injector.Registrar))
	}

	return kratos.New(opts...)
}
