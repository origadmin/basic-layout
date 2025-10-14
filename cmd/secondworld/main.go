/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"flag"
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
)

var (
	// Name is the name of the compiled software.
	Name = "origadmin.service.v1.secondworld"
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	// flagconf is the config flag.
	flagconf string
)

func init() {
	// The config path should be the directory containing configuration files.
	flag.StringVar(&flagconf, "conf", "resources", "config path, eg: -conf resources")
}

func main() {
	flag.Parse()

	// Create AppInfo
	appInfo := &interfaces.AppInfo{
		Name:    Name,
		Version: Version,
	}

	// NewFromBootstrap handles config loading, logging, and container setup.
	rt, cleanup, err := runtime.NewFromBootstrap(flagconf, bootstrap.WithAppInfo(appInfo))
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
