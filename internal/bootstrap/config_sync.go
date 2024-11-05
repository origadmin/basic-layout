package bootstrap

import (
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/runtime/config"
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
	cfg := bs.Source
	switch cfg.Type {
	case "file":
		return nil
	case "consul":
		consulConfig := api.DefaultConfig()
		if bs.Source != nil && cfg.Consul != nil {
			consulConfig.Address = cfg.Consul.Address
			consulConfig.Scheme = cfg.Consul.Scheme
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
	cfg := bs.Source
	if cfg == nil {
		return errors.New("config is nil")
	}

	var src config.SourceConfig
	src.Type = cfg.Type
	if cfg.File != nil {
		src.File = &config.SourceConfig_File{
			Path: cfg.File.Path,
			//Format: cfg.File.Format,
		}
	}
	if cfg.Consul != nil {
		src.Consul = &config.SourceConfig_Consul{
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
