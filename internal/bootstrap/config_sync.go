/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"path"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/origadmin/runtime"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"

	"origadmin/basic-layout/internal/configs"
)

func SyncConfig(serviceName string, bs *configs.Bootstrap, output string) error {
	if output != "" {
		err := GenerateRemoteConfig(serviceName, bs, output)
		if err != nil {
			return err
		}
	}
	sourceConfig := bs.GetRegistry()
	if sourceConfig == nil {
		return errors.String("sourceConfig config is nil")
	}
	switch sourceConfig.Type {
	case "file":
		return nil
	case "consul":
		consulConfig := api.DefaultConfig()
		if sourceConfig.Consul != nil {
			consulConfig.Address = sourceConfig.Consul.Address
			consulConfig.Scheme = sourceConfig.Consul.Scheme
			client, err := api.NewClient(consulConfig)
			if err != nil {
				return errors.Wrap(err, "consul client error")
			}

			registries, err := LoadRegistries(&Config{
				Type: "consul",
				Consul: &configv1.SourceConfig_Consul{
					Address: sourceConfig.Consul.Address,
					Scheme:  sourceConfig.Consul.Scheme,
					Path:    RegistryPath("registries.json"),
				},
			})
			if err != nil {
				return errors.Wrap(err, "load registries error")
			}
			exist := false
			for _, registryConfig := range registries {
				if registryConfig.GetServiceName() == serviceName {
					exist = true
					break
				}
			}

			if exist {
				log.Infof("sync config to consul path: %s exist", RegistryPath("registries.json"))
				return nil
			}
			opt := protojson.MarshalOptions{
				EmitUnpopulated: true,
				Indent:          " ",
			}
			sourceConfig.ServiceName = serviceName
			registries = append(registries, sourceConfig)
			marshal, err := opt.Marshal(&configs.Registry{
				Registries: registries,
			})
			if err != nil {
				return errors.Wrap(err, "marshal config error")
			}
			path := RegistryPath("registries.json")
			if _, err := client.KV().Put(&api.KVPair{
				Key:   path, //path.Join("configs", name, "bootstrap.json"),
				Value: marshal,
			}, nil); err != nil {
				return errors.Wrap(err, "consul put error")
			}
			log.Infof("sync config to consul path: %s success", RegistryPath("registries.json"))
			return nil
		}
	case "etcd":
		return nil
	}
	return errors.Errorf("unsupported sourceConfig type: %s", sourceConfig.Type)
}

func ConfigPath(serviceName string, configName string) string {
	return path.Join("configs", serviceName, configName)
}

func RegistryPath(configName string) string {
	return path.Join("registry", configName)
}

func LoadRegistries(sourceConfig *configv1.SourceConfig) ([]*configv1.Registry, error) {
	cfg, err := runtime.NewConfig(sourceConfig)
	if err != nil {
		return nil, err
	}

	var registries configs.Registry
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	if err := cfg.Scan(&registries); err != nil {
		return nil, err
	}
	return registries.Registries, nil
}

func GenerateRemoteConfig(serviceName string, bs *configs.Bootstrap, file string) error {
	cfg := bs.Registry
	if cfg == nil {
		return errors.String("registry config is nil")
	}

	var src configv1.SourceConfig
	src.Type = cfg.Type
	switch cfg.Type {
	case "file":
		src.File = &configv1.SourceConfig_File{
			Path: "resources/configs/bootstrap.json",
			//Format: cfg.File.Format,
		}
	case "consul":
		src.Consul = &configv1.SourceConfig_Consul{
			Address: cfg.Consul.Address,
			Scheme:  cfg.Consul.Scheme,
		}
	}
	//if cfg.File != nil {
	//	src.File = &config.SourceConfig_File{
	//		Path: cfg.File.Path,
	//		//Format: cfg.File.Format,
	//	}
	//}
	//if cfg.Consul != nil {
	//	src.Consul = &config.SourceConfig_Consul{
	//		Address: cfg.Consul.Address,
	//		Scheme:  cfg.Consul.Scheme,
	//	}
	//}

	err := codec.EncodeToFile(file, &src)
	if err != nil {
		return errors.Wrap(err, "marshal config error")
	}
	return nil
}

func SyncRegistry(name string, bs *configs.Bootstrap) error {
	path, _ := filepath.Abs("resources/registries.json")

	if bs.GetSource() == nil {
		bs.Source = NewFileSource(path)
	}
	registries, err := LoadRegistries(bs.Source)
	if err == nil {
		exist := false
		for i := range registries {
			if registries[i].GetServiceName() == name {
				exist = true
				break
			}
		}
		if exist {
			log.Infof("sync config to consul path: %s exist", path)
			return nil
		}
	}

	registry := bs.GetRegistry()
	if registry == nil {
		return errors.String("registry config is nil")
	}
	registry.ServiceName = name
	registries = append(registries, registry)
	//opt := protojson.MarshalOptions{
	//	EmitUnpopulated: true,
	//	Indent:          " ",
	//}
	//marshal, err := opt.Marshal(&configs.Registry{
	//	Registries: registries,
	//})
	//if err != nil {
	//	log.Errorf("marshal config error: %v", err)
	//	return
	//}
	if err := codec.EncodeToFile(path, &configs.Registry{
		Registries: registries,
	}); err != nil {
		return err
	}
	log.Infof("sync config to) consul path: %s success", path)
	return nil
}
