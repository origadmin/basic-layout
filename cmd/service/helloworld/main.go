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
	_ "go.uber.org/automaxprocs"

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

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "resources/configs", "config path, eg: -conf config.yaml")
}

func NewApp(ctx context.Context, config *conf.Server, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	var r registry.Registrar
	// example one: consul
	switch config.Discovery.GetType() {
	case "consul":
		cfg := config.Discovery.GetConsul()
		if cfg == nil {
			break
		}
		client, err := api.NewClient(&api.Config{
			Address: cfg.Address,
		})
		if err != nil {
			break
		}
		r = registryconsul.New(client)
	}

	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	}
	if r != nil {
		opts = append(opts, kratos.Registrar(r))
	}

	return kratos.New(opts...)
}

func main() {
	flag.Parse()

	flagconf, _ = filepath.Abs(flagconf)
	fmt.Println("load config at:", flagconf)

	client, err := api.NewClient(&api.Config{
		Address: "host:8500",
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
		_, err := client.KV().Put(&api.KVPair{Key: "configs/" + kv.Key, Value: kv.Value}, nil)
		if err != nil {
			panic(err)
		}
	}

	//consul.WithPath(testPath)
	source, err := consul.New(client, consul.WithPath("configs/bootstrap.toml"))
	if err != nil {
		panic(err)
	}
	c := config.New(
		//config.WithSource(file.NewSource(flagconf), source),
		config.WithSource(source),
		//config.WithResolveActualTypes(true),
		config.WithDecoder(codec.SourceDecoder),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	//v := c.Value("configs/bootstrap.toml")
	//vs, err := v.String()
	//if err != nil {
	//	panic(err)
	//}
	//typo := codec.SupportTypeFromExt(filepath.Ext("configs/bootstrap.toml"))
	//if typo == codec.UNKNOWN {
	//	panic("unknown file type")
	//}
	//if err := typo.Unmarshal([]byte(vs), &bc); err != nil {
	//	panic(err)
	//}
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//c.Watch("configs/bootstrap.toml", func(key string, value config.Value) {
	//	defer wg.Done()
	//	vs, err := value.String()
	//	if err != nil {
	//		return
	//	}
	//	typo := codec.SupportTypeFromExt(filepath.Ext(key))
	//	if typo == codec.UNKNOWN {
	//		return
	//	}
	//	if err := typo.Unmarshal([]byte(vs), &bc); err != nil {
	//		return
	//	}
	//})
	//wg.Wait()
	logger := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", bc.ServiceName,
		"service.version", bc.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	fmt.Printf("show bootstrap config: %+v\n", bc)
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
