package main

import (
	"context"
	"flag"
	"fmt"
	"syscall"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	logger "github.com/origadmin/slog-kratos"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/protobuf/encoding/protojson"

	"origadmin/basic-layout/internal/bootstrap"
)

// go build -ldflags "-X main.Version=vx.y.z"
var (
	// Name is the name of the compiled software.
	name = "origadmin.service.v1.helloworld"
	// Version is the version of the compiled software.
	version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.DefaultFlags()
)

func init() {
	flag.StringVar(&flags.ConfigPath, "config", "resources", "config path, eg: -c config.toml")
	flag.StringVar(&flags.EnvPath, "env", "resources/env", "env path, eg: -e env.toml")
}

func main() {
	flag.Parse()
	flags.Name = name
	flags.Version = version
	flags.MetaData = make(map[string]string)
	l := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", flags.ID,
		"service.name", flags.Name,
		"service.version", flags.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	env, err := bootstrap.LoadEnv(flags.EnvPath)
	if err != nil {
		//log.Errorf("failed to load env: %s", err.Error())
		//os.Exit(1)
		log.Fatalf("failed to load env: %s", err.Error())
	}

	bs, err := bootstrap.FromLocal(name, flags.ConfigPath, env, l)
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	v, _ := protojson.Marshal(bs)
	fmt.Printf("show bootstrap config: %+v\n", string(v))
	ctx := context.Background()
	//info to ctx
	app, cleanup, err := buildInjectors(ctx, bs, l)
	if err != nil {
		log.Fatalf("failed to build injector: %s", err.Error())
	}
	defer cleanup()
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Fatalf("app stopped with error: %s", err.Error())
	}
}

func NewApp(ctx context.Context, injector *bootstrap.InjectorServer) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(flags.ID),
		kratos.Name(flags.Name),
		kratos.Version(flags.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(injector.Logger),
		kratos.Server(injector.ServerHTTP, injector.ServerGRPC),
	}
	if injector.Registrar != nil {
		opts = append(opts, kratos.Registrar(injector.Registrar))
	}

	return kratos.New(opts...)
}
