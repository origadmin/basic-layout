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
)

// go build -ldflags "-X main.Version=vx.y.z -X main.Name=origadmin.service.v1.helloworld"
var (
	// Name is the Name of the compiled software.
	Name = "origadmin.service.v1.helloworld"
	// Version is the Version of the compiled software.
	Version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.BootFlags{}
)

func init() {
	flags = bootstrap.NewBootFlags(Name, Version)
	flag.StringVar(&flags.ConfigPath, "c", "resources", "config path, eg: -c config.toml")
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

	fmt.Printf("bootstrap flags: %+v\n", flags)
	bs, err := bootstrap.Load(flags, true)
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
		kratos.Name(flags.ServiceName),
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
