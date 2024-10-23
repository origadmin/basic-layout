package bootstrap

import (
	"os"
	"path"
	"path/filepath"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/utils/replacer"

	"origadmin/basic-layout/internal/configs"
)

type FileSource struct {
	Path   string
	Format string
}

type ConsulSource struct {
	Address string
	Scheme  string
}

type EnvSource struct {
	Envs   map[string]string
	Files  []string
	Prefix []string
}

type SourceConfig struct {
	Type   string
	File   *FileSource
	Consul *ConsulSource
	Env    *EnvSource
}

// LoadConfig returns a new kratos config .
func LoadConfig(name string, cfg *SourceConfig, l log.Logger) (config.Config, error) {
	srcs, err := createConfigSource(name, cfg, l)
	if err != nil {
		return nil, err
	}

	return config.New(config.WithSource(srcs...)), nil
}

// createConfigSource creates the appropriate config source based on the type.
func createConfigSource(name string, cfg *SourceConfig, l log.Logger) ([]config.Source, error) {
	var source config.Source
	switch cfg.Type {
	case "file":
		if cfg.File == nil {
			return nil, errors.New("file config is nil")
		}
		source = file.NewSource(cfg.File.Path)
		if cfg.File.Format != "" {
			// todo
		}
	case "consul":
		if cfg.Consul == nil {
			return nil, errors.New("consul config is nil")
		}
		client, err := api.NewClient(&api.Config{
			Address: cfg.Consul.Address,
			Scheme:  cfg.Consul.Scheme,
		})
		if err != nil {
			return nil, errors.Wrap(err, "consul client error")
		}
		source, err = consul.New(client,
			consul.WithPath(consulConfigPath(name, "bootstrap.json")),
		)
		if err != nil {
			return nil, errors.Wrap(err, "consul source error")
		}
	default:
		return nil, errors.New("unsupported source type")
	}

	if source == nil {
		return nil, errors.New("source is nil")
	}
	srcs := []config.Source{source}
	if cfg.Env != nil {
		if cfg.Env.Files != nil {
			srcs = append(srcs, envf.NewSource(cfg.Env.Files, cfg.Env.Prefix...))
		}
		if cfg.Env.Envs != nil {
			srcs = append(srcs, envf.WithEnv(cfg.Env.Envs, cfg.Env.Prefix...))
		}
	}

	return srcs, nil
}

// NewFileSourceConfig returns a new SourceConfig for file-based configurations.
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
	cfg := NewFileSourceConfig(path)
	path, _ = filepath.Abs(path)
	stat, err := os.Stat(path)
	if err != nil {
		return cfg
	}

	if stat.IsDir() {
		cfg = loadConfigFromDir(path)
	} else {
		cfg = loadConfigFromFile(path)
	}

	if err := codec.DecodeFromFile(path, cfg); err != nil {
		return cfg
	}
	return cfg
}

// LoadEnv Loads configuration files in various formats from a directory,
func LoadEnv(path string) (map[string]string, error) {
	envs := make(map[string]string)
	if err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		} else if d.IsDir() {
			return nil
		}
		typo := codec.TypeFromExt(filepath.Ext(walkpath))
		if typo == codec.UNKNOWN {
			return nil
		}
		if err := codec.DecodeFromFile(walkpath, &envs); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", walkpath)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return envs, nil
}

func FromLocal(name, path string, envs map[string]string, l log.Logger) (*configs.Bootstrap, error) {
	path, _ = filepath.Abs(path)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to state file %s", path)
	}
	log.NewHelper(l).Infof("loading config from %s", path)

	var cfg *SourceConfig
	if stat.IsDir() {
		cfg = loadConfigFromDir(path)
	} else {
		cfg = loadConfigFromFile(path)
	}
	if cfg == nil {
		cfg = NewFileSourceConfig(path)
	}

	if envs != nil {
		cfg.Env = &EnvSource{
			Envs: envs,
		}
	}

	return LoadBootstrap(name, cfg, l)
}

// loadConfigFromDir loads configuration from a directory.
func loadConfigFromDir(path string) *SourceConfig {
	var cfg SourceConfig
	err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		} else if d.IsDir() {
			return nil
		}
		typo := codec.TypeFromExt(filepath.Ext(walkpath))
		if typo == codec.UNKNOWN {
			return nil
		}

		if err := codec.DecodeFromFile(walkpath, &cfg); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", walkpath)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return &cfg
}

// loadConfigFromFile loads configuration from a single file.
func loadConfigFromFile(path string) *SourceConfig {
	var cfg SourceConfig
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return nil
	}
	return &cfg
}

// LoadBootstrap Loads configuration files in various formats from a directory,
// and parses them into a struct.
func LoadBootstrap(name string, cfg *SourceConfig, l log.Logger) (*configs.Bootstrap, error) {
	source, err := LoadConfig(name, cfg, l)
	if err != nil {
		return nil, err
	}

	bs := configs.DefaultBootstrap
	if err := source.Load(); err != nil {
		return nil, errors.Wrap(err, "load config error")
	}
	if err := source.Scan(bs); err != nil {
		return nil, errors.Wrap(err, "scan config error")
	}
	if err := bs.ValidateAll(); err != nil {
		return nil, errors.Wrap(err, "validate config error")
	}
	return bs, nil
}

func ApplyEnv(content []byte, envs map[string]string) []byte {
	r := replacer.New(replacer.WithSeparator("="))
	return r.Replace(content, envs)
}

// consul path
func consulConfigPath(dir, name string) string {
	return path.Join("configs", dir, name)
}
