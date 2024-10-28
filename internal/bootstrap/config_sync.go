package bootstrap

import (
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
	"google.golang.org/protobuf/encoding/protojson"

	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/toolkits/oneof/source"
)

func SyncConfig(serviceName string, bs *configs.Bootstrap, output string) error {
	if output != "" {
		err := GenerateRemoteConfig(serviceName, bs, output)
		if err != nil {
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
				EmitUnpopulated: true,
				Indent:          " ",
			}
			marshal, err := opt.Marshal(bs)
			if err != nil {
				return errors.Wrap(err, "marshal config error")
			}
			if _, err := client.KV().Put(&api.KVPair{
				Key:   source.ConfigPath(serviceName, "bootstrap.json"), //path.Join("configs", name, "bootstrap.json"),
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

func GenerateRemoteConfig(serviceName string, bs *configs.Bootstrap, file string) error {
	cfg := bs.Config
	if cfg == nil {
		return errors.New("config is nil")
	}

	var src SourceConfig
	src.Type = cfg.Type
	if cfg.File != nil {
		src.File = FileSource{
			Path:   cfg.File.Path,
			Format: cfg.File.Format,
		}
	}
	if cfg.Consul != nil {
		src.Consul = ConsulSource{
			Address: cfg.Consul.Address,
			Scheme:  cfg.Consul.Scheme,
		}
	}

	err := codec.EncodeToFile(file, &src)
	if err != nil {
		return errors.Wrap(err, "marshal config error")
	}
	return nil
}
