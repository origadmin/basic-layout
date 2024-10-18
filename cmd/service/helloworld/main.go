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
	"github.com/origadmin/toolkits/errors"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/protobuf/encoding/protojson"

	"origadmin/basic-layout/internal/config"
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/toolkits/service"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name    string = "origadmin.service.v1.helloworld"
	Version string = "v1.0.0"
	cfg     string
	//// 配置，启动链路追踪
	//url := "http://192.168.0.103:14268/api/traces"
	//Name = "kratos.service.student"
	//id = "kratos.id.student.1"
)

func init() {
	flag.StringVar(&cfg, "config", "resources/config", "config path, eg: -c config.toml")
}

func main() {
	flag.Parse()
	inf := service.NewAppInfo("", Name, Version).SetMetadata(map[string]string{})
	logger := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", inf.InstanceID(),
		"service.name", inf.Name(),
		"service.version", inf.Version(),
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	bs, err := config.LoadBootstrap(&config.Config{}, logger)
	if err != nil {
		return
	}

	v, _ := protojson.Marshal(bs)
	fmt.Printf("show bootstrap config: %+v\n", string(v))
	ctx := context.Background()
	//info to ctx
	app, cleanup, err := buildInjectors(ctx, bs, logger)
	if err != nil {
		panic(errors.WithStack(err))
	}
	defer cleanup()
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(errors.WithStack(err))
	}
	//c.Watch("configs/bootstrap.json", func(key string, value config.Value) {
	//	if key != "configs/bootstrap.json" {
	//		return
	//	}
	//	err := value.Scan(&bc)
	//	if err != nil {
	//		logger.Log(log.LevelError, "error", err)
	//		return
	//	}
	//})
}

func NewApp(ctx context.Context, injector *mods.Injector) *kratos.App {
	info := service.NewAppInfo(id, Name, Version).SetMetadata(map[string]string{})
	opts := []kratos.Option{
		kratos.ID(info.InstanceID()),
		kratos.Name(info.Name()),
		kratos.Version(info.Version()),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(injector.Logger),
		//kratos.Server(hs, gs, gss),
		kratos.Server(injector.ServerHTTP, injector.ServerGRPC),
	}
	//ep1, _ := injector.ServerGINS.Endpoint()
	ep2, _ := injector.ServerHTTP.Endpoint()
	ep3, _ := injector.ServerGRPC.Endpoint()
	fmt.Println("endpoint:", ep2, ep3)

	if injector.Registry != nil {
		opts = append(opts, kratos.Registrar(injector.Registry))
	}

	return kratos.New(opts...)
}
