package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"

	"origadmin/basic-layout/internal/configs"
)

type FileSource struct {
	Path string
}

type ConsulSource struct {
	Address string
	Scheme  string
}

type EnvSource struct {
	Files  []string
	Prefix []string
}

type SourceConfig struct {
	Type   string
	File   *FileSource
	Consul *ConsulSource
	Envs   *EnvSource
}

// NewSourceConfig returns a new kratos config .
func NewSourceConfig(cfg *SourceConfig, l log.Logger) (config.Config, error) {
	var source config.Source
	switch cfg.Type {
	case "file":
		source = file.NewSource(cfg.File.Path)
	case "consul":
		client, err := api.NewClient(&api.Config{
			Address: cfg.Consul.Address,
			Scheme:  cfg.Consul.Scheme,
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
	srcs := []config.Source{source}
	if cfg.Envs != nil && cfg.Envs.Files != nil {
		srcs = append(srcs, envf.NewSource(cfg.Envs.Files, cfg.Envs.Prefix...))
	}

	return config.New(config.WithSource(srcs...)), nil
}

// NewSourceConfigWithEnv returns a new kratos config .
func NewSourceConfigWithEnv(cfg *SourceConfig, envs map[string]string, l log.Logger) (config.Config, error) {
	var source config.Source
	switch cfg.Type {
	case "file":
		source = file.NewSource(cfg.File.Path)
	case "consul":
		client, err := api.NewClient(&api.Config{
			Address: cfg.Consul.Address,
			Scheme:  cfg.Consul.Scheme,
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
	srcs := []config.Source{source}

	var prefix []string
	if cfg.Envs != nil {
		prefix = append(prefix, cfg.Envs.Prefix...)
		srcs = append(srcs, envf.NewSource(cfg.Envs.Files, prefix...))
	}
	if envs != nil {
		srcs = append(srcs, envf.WithEnv(envs, prefix...))
	}

	return config.New(config.WithSource(srcs...)), nil
}

func NewFileSourceConfig(path string) *SourceConfig {
	return &SourceConfig{
		Type: "file",
		File: &FileSource{
			Path: path,
		},
	}
}

// LocalSourceConfig Loads configuration files in various formats from a directory,
// and parses them into a config.
func LocalSourceConfig(path string) *SourceConfig {
	var cfg SourceConfig
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return NewFileSourceConfig(path)
	}
	return &cfg
}

// LoadEnv Loads configuration files in various formats from a directory,
func LoadEnv(path string) (map[string]string, error) {
	envs := make(map[string]string)
	if err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", path)
		} else if d.IsDir() {
			return nil
		}
		if err := codec.DecodeFromFile(path, &envs); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return envs, nil
}

func FromLocal(path string, envs map[string]string, l log.Logger) (*configs.Bootstrap, error) {
	path, _ = filepath.Abs(path)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to state file %s", path)
	}

	if stat.IsDir() {
		var cfg SourceConfig
		err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return errors.Wrapf(err, "failed to get config file %s", path)
			} else if d.IsDir() {
				return nil
			}
			typo := codec.TypeFromExt(filepath.Ext(path))
			if typo == codec.UNKNOWN {
				return nil
			}

			if err := codec.DecodeFromFile(path, &cfg); err != nil {
				return errors.Wrapf(err, "failed to parse config file %s", path)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to walk config path %s", path)
		}
		return LoadBootstrap(&cfg, envs, l)
	}
	var cfg SourceConfig
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to parse config file %s", path)
	}

	return LoadBootstrap(&cfg, envs, l)
}

// LoadBootstrap Loads configuration files in various formats from a directory,
// and parses them into a struct.
func LoadBootstrap(cfg *SourceConfig, envs map[string]string, l log.Logger) (*configs.Bootstrap, error) {
	var source config.Config
	var err error
	if len(envs) == 0 {
		source, err = NewSourceConfig(cfg, l)
	} else {
		source, err = NewSourceConfigWithEnv(cfg, envs, l)
	}
	if err != nil {
		return nil, err
	}

	bs := configs.DefaultBootstrap
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
