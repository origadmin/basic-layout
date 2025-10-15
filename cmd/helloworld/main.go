/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"flag"
	"log"

	"github.com/google/uuid"

	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	"origadmin/basic-layout/internal/transformer"
)

var (
	// Name is the name of the compiled software.
	Name = "origadmin.service.v1.helloworld"
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	// flagconf is the config flag.
	flagconf string
)

func init() {
	// The config path should be the directory containing configuration files.
	flag.StringVar(&flagconf, "conf", "resources/configs/bootstrap.yaml", "config path, eg: -conf resources/configs/bootstrap.yaml")
}

func main() {
	flag.Parse()

	// Create AppInfo using the struct from the bootstrap package
	//appInfo := &runtime.AppInfo{
	//	ID:      uuid.New().String(), // Add ID field, using Name as ID
	//	Name:    Name,
	//	Version: Version,
	//}

	// NewFromBootstrap handles config loading, logging, and container setup.
	rt, cleanup, err := runtime.NewFromBootstrap(flagconf,
		bootstrap.WithConfigTransformer(transformer.New(&appv1.App{
			Id:      uuid.New().String(),
			Name:    Name,
			Version: Version,
		})))
	if err != nil {
		log.Fatalf("failed to create runtime: %v", err)
	}
	defer cleanup()

	// wireApp now takes the runtime instance and builds the kratos app.
	app, cleanupApp, err := wireApp(rt)
	if err != nil {
		log.Fatalf("failed to wire app: %v", err)
	}
	defer cleanupApp()

	// Run the application
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
