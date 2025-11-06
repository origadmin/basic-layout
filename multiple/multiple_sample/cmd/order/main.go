/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/google/uuid"

	_ "basic-layout/multiple/multiple_sample/helpers/configsource/oneof"
	"basic-layout/multiple/multiple_sample/internal/transformer"

	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	_ "github.com/origadmin/runtime/config/envsource"
	_ "github.com/origadmin/runtime/config/file"
)

var (
	// Name is the name of the compiled software.
	Name = "origadmin.service.v1.order"
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	// flagconf is the config flag.
	flagconf string
)

func init() {
	// The config path should be the directory containing configuration files.
	flag.StringVar(&flagconf, "conf", "bootstrap.yaml", "config path, eg: -conf bootstrap.yaml")
}

func main() {
	flag.Parse()

	// Log the config path for debugging
	log.Printf("Loading configuration from: %s\n", flagconf)

	if !filepath.IsAbs(flagconf) {
		flagconf = filepath.Join("resources/configs/order/", flagconf)
	}
	// Create AppInfo using the struct from the runtime package
	appInfo := &appv1.App{
		Id:      uuid.New().String(),
		Name:    Name,
		Version: Version,
	}

	// NewFromBootstrap handles config loading, logging, and container setup.
	rt, err := runtime.NewFromBootstrap(
		flagconf,
		//runtime.WithAppInfo(appInfo),
		bootstrap.WithConfigTransformer(transformer.New(appInfo)),
	)
	if err != nil {
		log.Fatalf("failed to create runtime: %v", err)
	}
	defer rt.Cleanup()

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
