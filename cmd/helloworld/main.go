/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	_ "github.com/origadmin/runtime/config/envsource"
	_ "github.com/origadmin/runtime/config/file"
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

	// Get the absolute path to the config file
	configPath := flagconf
	if !filepath.IsAbs(configPath) {
		absPath, err := filepath.Abs(filepath.Join(".", configPath))
		if err != nil {
			log.Fatalf("failed to get absolute path for config: %v", err)
		}
		configPath = absPath
	}

	// Log the config path for debugging
	log.Printf("Loading configuration from: %s\n", configPath)

	// Verify config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	// Set working directory to the directory containing the config file
	configDir := filepath.Dir(configPath)
	if err := os.Chdir(configDir); err != nil {
		log.Fatalf("failed to change working directory to %s: %v", configDir, err)
	}
	log.Printf("Working directory set to: %s\n", configDir)

	// Create app info
	appInfo := &appv1.App{
		Id:      uuid.New().String(),
		Name:    Name,
		Version: Version,
	}

	// Log app info
	log.Printf("Starting %s %s (ID: %s)\n", appInfo.Name, appInfo.Version, appInfo.Id)

	// NewFromBootstrap handles config loading, logging, and container setup.
	rt, cleanup, err := runtime.NewFromBootstrap(
		filepath.Base(configPath), // Use just the filename since we changed the working directory
		bootstrap.WithConfigTransformer(transformer.New(appInfo)),
	)
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
