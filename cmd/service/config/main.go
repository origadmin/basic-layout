package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/protobuf/encoding/protojson"

	bootstrap "origadmin/basic-layout/internal/config"
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
		panic(errors.WithStack(err))
	}
	client, err := api.NewClient(&api.Config{
		Address: boot.Consul.Address,
	})
	if err != nil {
		panic(errors.WithStack(err))
	}

	fs := file.NewSource(flagconf)
	kvs, err := fs.Load()
	if err != nil {
		panic(errors.WithStack(err))
	}

	for _, kv := range kvs {
		fmt.Println("key:", kv.Key)
		typo := codec.TypeFromExt(filepath.Ext(kv.Key))
		if typo == codec.UNKNOWN {
			continue
		}
		fmt.Println("put key:", kv.Key)
		_, err := client.KV().Put(&api.KVPair{Key: "configs/" + kv.Key, Value: kv.Value}, nil)
		if err != nil {
			panic(errors.WithStack(err))
		}
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

	err = bc.ValidateAll()
	if err != nil {
		panic(errors.WithStack(err))
	}

	v, _ := protojson.Marshal(&bc)

	fmt.Printf("bootstrap config: %s\n", v)
}
