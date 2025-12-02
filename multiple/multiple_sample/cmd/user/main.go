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
	Name = "origadmin.service.v1.user"
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	// flagconf is the config flag.
	flagconf string

	// workdir is a flag to indicate whether to use the working directory as the config path.
	//workdir bool
)

func init() {
	// The config path should be the directory containing configuration files.
	flag.StringVar(&flagconf, "conf", "bootstrap.yaml", "config path, eg: -conf bootstrap.yaml")
	//flag.BoolVar(&workdir, "workdir", false, "use working directory as config path")
}

func main() {
	flag.Parse()

	// Log the config path for debugging
	log.Printf("Loading configuration from: %s\n", flagconf)

	// Create app info
	appInfo := &appv1.App{
		Id:      uuid.New().String(),
		Name:    Name,
		Version: Version,
	}

	// Log app info
	log.Printf("Starting %s %s (ID: %s)\n", appInfo.Name, appInfo.Version, appInfo.Id)
	// If workdir is set, use the working directory as the config path.

	if filepath.IsAbs(flagconf) {
		flagconf = filepath.Join("resources/configs/user/", flagconf)
	}
	//	directory = filepath.Dir(flagconf)
	//	flagconf = filepath.Base(flagconf)
	//}
	//log.Printf("Using config directory: %s\n", directory)
	// NewFromBootstrap handles config loading, logging, and container setup.
	rt := runtime.New(Name, Version)
	err := rt.Load(flagconf, bootstrap.WithConfigTransformer(conf.New(
		appInfo))) // Use just the filename since we changed the working directory
	if err != nil {
		log.Fatalf("failed to create runtime: %v", err)
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
