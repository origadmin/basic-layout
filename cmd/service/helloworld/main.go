package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/hashicorp/consul/api"
	logger "github.com/origadmin/slog-kratos"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/idgen"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/protobuf/encoding/protojson"

	bootstrap "origadmin/basic-layout/internal/config"
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/internal/mods/helloworld/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
	flagenv  string
	flagport string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "resources/configs", "config path, eg: -conf config.toml")
	flag.StringVar(&flagenv, "env", "resources/env", "env path, eg: -env env.toml")
	flag.StringVar(&flagport, "port", "9000", "service port")
}

func main() {
	flag.Parse()

	var env map[string]string
	if err := codec.DecodeTOMLFile(filepath.Join(flagenv, "env.toml"), &env); err != nil {
		panic(err)
	}
	fmt.Printf("load env: %s\n", env)

	flagconf, _ = filepath.Abs(flagconf)
	fmt.Println("load config at:", flagconf)

	var boot bootstrap.Bootstrap
	err := bootstrap.LoadWithEnv(&boot, flagconf, env)
	if err != nil {
		panic(err)
	}
	client, err := api.NewClient(&api.Config{
		Address: boot.Consul.Address,
	})
	if err != nil {
		panic(errors.WithStack(err))
	}

	source, err := consul.New(client,
		consul.WithPath("configs/bootstrap.json"),
	)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c := config.New(
		config.WithSource(source, envf.WithEnv(env)),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(errors.WithStack(err))
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(errors.WithStack(err))
	}

	logger := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", bc.ServiceName,
		"service.version", bc.Version,
		"trace.id", idgen.GenID(),
		"span.id", tracing.SpanID(),
	)

	err = bc.ValidateAll()
	if err != nil {
		panic(errors.WithStack(err))
	}

	v, _ := protojson.Marshal(&bc)
	fmt.Printf("show bootstrap config: %+v\n", string(v))
	ctx := context.Background()
	app, cleanup, err := buildInjectors(ctx, &bc, logger)
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
	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(injector.Bootstrap.ServiceName),
		kratos.Version(injector.Bootstrap.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(injector.Logger),
		//kratos.Server(hs, gs, gss),
		kratos.Server(injector.ServerGINS, injector.ServerHTTP, injector.ServerGRPC),
	}
	ep1, _ := injector.ServerGINS.Endpoint()
	ep2, _ := injector.ServerHTTP.Endpoint()
	ep3, _ := injector.ServerGRPC.Endpoint()
	fmt.Println("endpoint:", ep1, ep2, ep3)

	if injector.Registry != nil {
		opts = append(opts, kratos.Registrar(injector.Registry))
	}

	return kratos.New(opts...)
}
