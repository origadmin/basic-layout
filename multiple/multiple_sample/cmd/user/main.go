/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"flag"
	"log"

	_ "github.com/sqlite3ent/sqlite3"

	"github.com/joho/godotenv"

	"basic-layout/multiple/multiple_sample/internal/bootstrap"
	"basic-layout/multiple/multiple_sample/internal/conf"
	_ "basic-layout/multiple/multiple_sample/internal/helpers/configsource/oneof"
	"github.com/origadmin/runtime"
	runtimebootstrap "github.com/origadmin/runtime/bootstrap"
)

var (
	// Name is the name of the compiled software.
	Name = "origadmin.service.v1.user"
	// Version is the version of the compiled software.
	Version = "v1.0.0"

	// flagconf is the config flag.
	flagconf string
)

func init() {
	// The config path should be the directory containing configuration files.
	// The default is empty, so we can detect if the user has provided it.
	flag.StringVar(&flagconf, "conf", "", "config path, eg: -conf bootstrap.yaml")
}

func main() {
	// Load .env file for local development from resources directory.
	// It's safe to ignore the error, as the file may not exist in production.
	_ = godotenv.Load("resources/.env.user")

	flag.Parse()

	confPath := bootstrap.FindConfPath(flagconf)
	if confPath == "" {
		log.Fatalf("Could not find configuration file. Searched -conf flag, executable path, and development path.")
	}

	// Log the config path for debugging
	log.Printf("Loading configuration from: %s\n", confPath)

	// NewFromBootstrap handles config loading, logging, and container setup.
	rt := runtime.New(Name, Version)
	err := rt.Load(confPath, runtimebootstrap.WithConfigTransformer(conf.New()))
	if err != nil {
		log.Fatalf("failed to create runtime: %v", err)
	}
	defer rt.Config().Close()
	log.Printf("Starting %s %s (ID: %s)\n", rt.AppInfo().Name(), rt.AppInfo().Version(), rt.AppInfo().ID())

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
