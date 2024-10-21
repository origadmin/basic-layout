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
	// Name is the name of the compiled software.
	Name string = "origadmin.service.v1.config"
	// Version is the version of the compiled software.
	Version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.DefaultFlags()
)

func init() {
	flag.StringVar(&flags.ConfigPath, "config", "resources/local", "config path, eg: -c config.toml")
}

func main() {
	flag.Parse()

	flags.Name = Name
	flags.Version = Version
	flags.MetaData = make(map[string]string)
	logger := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", flags.IID(),
		"service.name", flags.Name,
		"service.version", flags.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	bs, err := bootstrap.FromLocal(flags.ConfigPath, nil, logger)
	if err != nil {
		panic(errors.WithStack(err))
	}

	_ = bs
	//todo
	//marshal, err := protojson.Marshal(bs)
	//if err != nil {
	//	panic(errors.WithStack(err))
	//}
	//
	//client, err := api.NewClient(&api.Config{
	//	Address: cfg.Consul.Address,
	//	Scheme:  cfg.Consul.Scheme,
	//})
	//if err != nil {
	//	return nil, errors.Wrap(err, "consul client error")
	//}
	//source, err = consul.New(client,
	//	consul.WithPath("configs/bootstrap.json"),
	//)
	//if err != nil {
	//	return nil, errors.Wrap(err, "consul source error")
	//}

	//for _, kv := range kvs {
	//	fmt.Println("key:", kv.Key)
	//	typo := codec.TypeFromExt(filepath.Ext(kv.Key))
	//	if typo == codec.UNKNOWN {
	//		continue
	//	}
	//	fmt.Println("put key:", kv.Key)
	//	_, err := client.KV().Put(&api.KVPair{Key: "configs/" + "bootstrap.json", Value: marshal}, nil)
	//	if err != nil {
	//		panic(errors.WithStack(err))
	//	}
	//}
	//
	//source, err := consul.New(client,
	//	consul.WithPath("configs/bootstrap.json"),
	//)
	//if err != nil {
	//	panic(errors.WithStack(err))
	//}
	//c := config.New(
	//	config.WithSource(source, envf.WithEnv(env)),
	//)
	//defer c.Close()
	//if err := c.Load(); err != nil {
	//	panic(errors.WithStack(err))
	//}
	//
	//var bc conf.Bootstrap
	//if err := c.Scan(&bc); err != nil {
	//	panic(errors.WithStack(err))
	//}
	//
	//err = bc.ValidateAll()
	//if err != nil {
	//	panic(errors.WithStack(err))
	//}
	//
	//v, _ := protojson.Marshal(&bc)

	//fmt.Printf("bootstrap config: %s\n", v)
}
