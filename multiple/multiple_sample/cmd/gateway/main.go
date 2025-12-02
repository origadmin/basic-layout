/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/google/uuid"

	"basic-layout/multiple/multiple_sample/internal/conf"
	_ "basic-layout/multiple/multiple_sample/internal/helpers/configsource/oneof"
	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/bootstrap"
)

var (
	// Name is the name of the compiled software.
	Name = "origadmin.service.v1.gateway"
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
		flagconf = filepath.Join("resources/configs/gateway/", flagconf)
	}
	// Create AppInfo using the struct from the runtime package
	appInfo := &appv1.App{
		Id:      uuid.New().String(),
		Name:    Name,
		Version: Version,
	}

	// NewFromBootstrap handles config loading, logging, and container setup.
	rt := runtime.New(Name, Version)
	err := rt.Load(
		flagconf,
		bootstrap.WithConfigTransformer(conf.New(appInfo)),
	)
	if err != nil {
		return
	}
	defer rt.Config().Close()

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
