package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	_ "github.com/sqlite3ent/sqlite3"

	"basic-layout/simple/simple_app/internal/conf"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
)

var (
	Name    = "simple_app"
	Version = "1.0.0"
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "bootstrap.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// Get the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}
	log.Printf("current working directory: %s", wd)

	if !filepath.IsAbs(flagconf) {
		flagconf = filepath.Join("resources/configs/", flagconf)
	}

	// Create a new runtime application from the bootstrap configuration.
	// Bootstrap-specific options are now wrapped in runtime.WithBootstrapOptions.
	rt := runtime.New(Name, Version)
	err = rt.Load(flagconf, bootstrap.WithConfigTransformer(&conf.Config{}))
	if err != nil {
		log.Fatalf("failed to create runtime from bootstrap: %v", err)
	}
	defer rt.Config().Close()
	app, cleanup, err := wireApp(rt)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		log.Printf("failed to run app: %v", err)
		os.Exit(1)
	}
}
