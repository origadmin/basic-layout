package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/sqlite3ent/sqlite3"

	"basic-layout/simple/simple_app/internal/conf"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
)

var (
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

	// Create a new runtime application from the bootstrap configuration.
	rt, err := runtime.NewFromBootstrap(flagconf, bootstrap.WithConfigTransformer(&conf.Config{}))
	if err != nil {
		log.Fatalf("failed to create runtime from bootstrap: %v", err)
	}

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
