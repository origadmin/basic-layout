package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	registryconsul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
	logger "github.com/origadmin/slog-kratos"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/idgen"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/origadmin/basic-layout/internal/bootstrap"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/conf"
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
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "resources/configs", "config path, eg: -conf config.toml")
	flag.StringVar(&flagenv, "env", "resources/env", "env path, eg: -env env.toml")
}

func NewApp(ctx context.Context, cfg *conf.Server, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name("helloworld"),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	}

	// endpoint1, _ := url.Parse("http://192.168.28.81:8000")
	// endpoint2, _ := url.Parse("grpc://192.168.28.81:9000")
	// kratos.Endpoint(endpoint1, endpoint2)

	var reg registry.Registrar

	// example one: consul
	switch cfg.Discovery.GetType() {
	case "consul":
		cfg := cfg.Discovery.GetConsul()
		if cfg == nil {
			break
		}
		client, err := api.NewClient(&api.Config{
			Address: cfg.Address,
		})
		if err != nil {
			break
		}
		//endpoint, err := url.Parse(cfg.Address)
		//if err != nil {
		//	break
		//}
		//opts = append(opts, kratos.Endpoint(endpoint))
		reg = registryconsul.New(
			client,
			registryconsul.WithHealthCheck(true),
		)
	}

	if reg != nil {
		opts = append(opts, kratos.Registrar(reg))
	}

	return kratos.New(opts...)
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

	var boot bootstrap.Config
	err := bootstrap.LoadWithEnv(&boot, flagconf, env)
	if err != nil {
		panic(err)
	}
	client, err := api.NewClient(&api.Config{
		Address: boot.Consul.Address,
	})
	if err != nil {
		panic(err)
	}

	fs := file.NewSource(flagconf)
	kvs, err := fs.Load()
	if err != nil {
		panic(err)
	}

	for _, kv := range kvs {
		fmt.Println("key:", kv.Key)
		typo := codec.TypeFromExt(filepath.Ext(kv.Key))
		if typo == codec.UNKNOWN {
			continue
		}
		_, err := client.KV().Put(&api.KVPair{Key: "configs/" + kv.Key, Value: kv.Value}, nil)
		if err != nil {
			panic(err)
		}
	}

	source, err := consul.New(client,
		consul.WithPath("configs/bootstrap.json"),
	)
	if err != nil {
		panic(err)
	}
	c := config.New(
		config.WithSource(source),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
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
		panic(err)
	}
	v, _ := protojson.Marshal(&bc)

	fmt.Printf("show bootstrap config: %+v\n", string(v))
	ctx := context.Background()

	app, cleanup, err := buildApp(ctx, bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
