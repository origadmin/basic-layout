package main

import (
	"flag"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	logger "github.com/origadmin/slog-kratos"
	"github.com/origadmin/toolkits/errors"
	_ "go.uber.org/automaxprocs"

	"origadmin/basic-layout/internal/bootstrap"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the Name of the compiled software.
	Name string = "origadmin.service.v1.config"
	// Version is the Version of the compiled software.
	Version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.BootFlags{}
	// remote is the remote of bootstrap flags.
	output = "resources"
)

func init() {
	flags = bootstrap.NewBootFlags(Name, Version)
	flag.StringVar(&flags.ConfigPath, "c", "resources", "config path, eg: -c config.toml")
	flag.StringVar(&output, "o", "", "output a bootstrap config from local config, eg: -o bootstrap.toml")
}

func main() {
	flag.Parse()

	flags.Metadata = make(map[string]string)
	l := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", flags.ID,
		"service.name", flags.ServiceName,
		"service.version", flags.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(l)
	flags.ConfigPath = filepath.Join(flags.WorkDir, flags.ConfigPath, "local.toml")
	//env, _ := bootstrap.LoadEnv(flags.EnvPath)
	bs, err := bootstrap.FromFlags(flags, l)
	if err != nil {
		panic(errors.WithStack(err))
	}
	log.Infof("bootstrap: %+v", bootstrap.PrintString(bs))
	err = bootstrap.SyncConfig(bs.GetServiceName(), bs, output)
	if err != nil {
		panic(errors.WithStack(err))
	}

	return
}
