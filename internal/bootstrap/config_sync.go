package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
	"google.golang.org/protobuf/encoding/protojson"

	"origadmin/basic-layout/internal/configs"
)

func SyncConfig(name string, bs *configs.Bootstrap, envs map[string]string, example string) error {
	cfg := bs.Config
	if cfg != nil && example != "" {
		var src SourceConfig
		src.Type = cfg.Type
		if cfg.File != nil {
			src.File = &FileSource{
				Path:   cfg.File.Path,
				Format: cfg.File.Format,
			}
		}
		if cfg.Consul != nil {
			src.Consul = &ConsulSource{
				Address: cfg.Consul.Address,
				Scheme:  cfg.Consul.Scheme,
			}
		}
		typo := codec.TypeFromExt(filepath.Ext(example))
		marshal, err := typo.Marshal(src)
		if err != nil {
			return errors.Wrap(err, "marshal config error")
		}
		marshal = ApplyEnv(marshal, envs)
		if err := os.WriteFile(example+".example", marshal, os.ModePerm); err != nil {
			return err
		}
	}
	switch bs.Config.Type {
	case "file":
		return nil
	case "consul":
		consulConfig := api.DefaultConfig()
		if bs.Config != nil && bs.Config.Consul != nil {
			consulConfig.Address = bs.Config.Consul.Address
			consulConfig.Scheme = bs.Config.Consul.Scheme
			client, err := api.NewClient(consulConfig)
			if err != nil {
				return errors.Wrap(err, "consul client error")
			}
			opt := protojson.MarshalOptions{
				Indent: " ",
			}
			marshal, err := opt.Marshal(bs)
			if err != nil {
				return errors.Wrap(err, "marshal config error")
			}
			if _, err := client.KV().Put(&api.KVPair{
				Key:   consulConfigPath(name, "bootstrap.json"), //path.Join("configs", name, "bootstrap.json"),
				Value: marshal,
			}, nil); err != nil {
				return errors.Wrap(err, "consul put error")
			}
			return nil
		}
	case "etcd":
		return nil
	}
	return errors.New("not support config type")
}
