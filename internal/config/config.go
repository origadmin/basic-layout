package config

import (
	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
)

type Config struct {
	Path   string
	Type   string
	Consul struct {
		Address string
		Schema  string
	}
	Env string
}

// NewSourceConfig returns a new kratos config .
func NewSourceConfig(cfg *Config, l log.Logger) (config.Config, error) {
	var source config.Source
	switch cfg.Type {
	case "file":
		source = file.NewSource(cfg.Path)
	case "consul":
		client, err := api.NewClient(&api.Config{
			Address: cfg.Consul.Address,
		})
		if err != nil {
			return nil, errors.Wrap(err, "consul client error")
		}
		source, err = consul.New(client,
			consul.WithPath("configs/bootstrap.json"),
		)
		if err != nil {
			return nil, errors.Wrap(err, "consul source error")
		}
	}
	if source == nil {
		return nil, errors.New("source is nil")
	}

	return config.New(
		config.WithSource(source, envf.NewSource([]string{cfg.Env})),
	), nil
}

func LoadBootstrap(cfg *Config, l log.Logger) (*Bootstrap, error) {
	source, err := NewSourceConfig(cfg, l)
	if err != nil {
		return nil, errors.Wrap(err, "new source config error")
	}
	bs := DefaultBootstrap
	if err := source.Load(); err != nil {
		return nil, errors.Wrap(err, "load config error")
	}
	if err := source.Scan(&bs); err != nil {
		return nil, errors.Wrap(err, "scan config error")
	}
	if err := bs.ValidateAll(); err != nil {
		return nil, errors.Wrap(err, "validate config error")
	}
	return bs, nil
}
