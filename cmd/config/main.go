package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	logger "github.com/origadmin/slog-kratos"
	"github.com/origadmin/toolkits/errors"
	_ "go.uber.org/automaxprocs"

	"origadmin/basic-layout/internal/bootstrap"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// name is the name of the compiled software.
	name string = "origadmin.service.v1.config"
	// version is the version of the compiled software.
	version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.DefaultFlags()
	// remote is the remote of bootstrap flags.
	remote = "resources/remote.toml"
)

func init() {
	flag.StringVar(&flags.ConfigPath, "c", "resources", "config path, eg: -c config.toml")
	flag.StringVar(&flags.EnvPath, "e", "resources/env", "env path, eg: -e env.toml")
	flag.StringVar(&remote, "r", "remote.toml", "export remote config, eg: -r remote.toml.")
}

func main() {
	flag.Parse()

	flags.Name = name
	flags.Version = version
	flags.MetaData = make(map[string]string)
	l := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", flags.IID(),
		"service.name", flags.Name,
		"service.version", flags.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	env, _ := bootstrap.LoadEnv(flags.EnvPath)
	bs, err := bootstrap.FromLocal("", flags.ConfigPath, env, l)
	if err != nil {
		panic(errors.WithStack(err))
	}

	err = bootstrap.SyncConfig(bs.GetServiceName(), bs, env, remote)
	if err != nil {
		panic(errors.WithStack(err))
	}

	return
}
