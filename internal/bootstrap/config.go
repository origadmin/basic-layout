package bootstrap

import (
	"os"
	"path"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/utils/replacer"

	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/toolkits/oneof/source"
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
	Files  []string
	Envs   map[string]string
	Prefix []string
}

type SourceConfig struct {
	Type   string
	File   FileSource
	Consul ConsulSource
	Env    EnvSource
}

func Load(path string, serviceName string) (*configs.Bootstrap, error) {
	//var s SourceConfig
	//err := codec.DecodeFromFile(path, s)
	//if err != nil {
	//	return nil, err
	//}
	sourceConfig := LoadSourceFiles(path)
	if sourceConfig.Env.Files != nil {
		envs, err := LoadEnvFiles(sourceConfig.Env.Files...)
		if err != nil {
			return nil, err
		}
		if sourceConfig.Env.Envs != nil {
			for k, v := range envs {
				if _, has := sourceConfig.Env.Envs[k]; !has {
					envs[k] = v
				}
			}
		} else {
			sourceConfig.Env.Envs = envs
		}
	}
	return LoadBootstrap(serviceName, sourceConfig, nil)
}

// LoadSourceFiles Loads configuration files in various formats from a directory,
// and parses them into a config.
func LoadSourceFiles(path string) *SourceConfig {
	cfg := NewFileSourceConfig(path)
	path, _ = filepath.Abs(path)
	stat, err := os.Stat(path)
	if err != nil {
		return cfg
	}

	if stat.IsDir() {
		cfg = loadSourceFromDir(path)
	} else {
		cfg = loadSourceFromFile(path)
	}

	return cfg
}

// sourceFromConfig creates the appropriate config source based on the type.
func sourceFromConfig(name string, cfg *SourceConfig, l log.Logger) ([]config.Source, error) {
	var configSource config.Source
	var err error
	switch cfg.Type {
	case "file":
		if cfg.File.Path == "" {
			return nil, errors.New("file config is nil")
		}
		configSource = file.NewSource(cfg.File.Path)
		if cfg.File.Format != "" {
			// todo
		}
	case "consul":
		if cfg.Consul.Address == "" {
			return nil, errors.New("consul config is nil")
		}
		path := source.ConfigPath(name, "bootstrap.json")
		configSource, err = source.NewSource(path, &configs.Consul{
			Address: cfg.Consul.Address,
			Scheme:  cfg.Consul.Scheme,
		})
		if err != nil {
			return nil, errors.Wrap(err, "consul source error")
		}
	default:
		return nil, errors.New("unsupported source type")
	}

	if configSource == nil {
		return nil, errors.New("source is nil")
	}
	srcs := []config.Source{configSource}
	if cfg.Env.Envs != nil {
		srcs = append(srcs, envf.WithEnv(cfg.Env.Envs, cfg.Env.Prefix...))
	}

	return srcs, nil
}

// NewFileSourceConfig returns a new SourceConfig for file-based configurations.
func NewFileSourceConfig(path string) *SourceConfig {
	return &SourceConfig{
		Type: "file",
		File: FileSource{
			Path: path,
		},
	}
}

// LoadEnvFiles Loads configuration files in various formats from a directory,
func LoadEnvFiles(paths ...string) (map[string]string, error) {
	envs := make(map[string]string)
	for i := range paths {
		if err := filepath.WalkDir(paths[i], func(walkpath string, d os.DirEntry, err error) error {
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
		cfg = loadSourceFromDir(path)
	} else {
		cfg = loadSourceFromFile(path)
	}
	if cfg == nil {
		cfg = NewFileSourceConfig(path)
	}

	if envs != nil {
		cfg.Env = EnvSource{
			Envs: envs,
		}
	}

	return LoadBootstrap(name, cfg, l)
}

// loadSourceFromDir loads configuration from a directory.
func loadSourceFromDir(path string) *SourceConfig {
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

// loadSourceFromFile loads configuration from a single file.
func loadSourceFromFile(path string) *SourceConfig {
	var cfg SourceConfig
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return nil
	}
	return &cfg
}

// LoadBootConfig returns a new kratos config .
func LoadBootConfig(serviceName string, cfg *SourceConfig, l log.Logger) (config.Config, error) {
	srcs, err := sourceFromConfig(serviceName, cfg, l)
	if err != nil {
		return nil, err
	}

	return config.New(config.WithSource(srcs...)), nil
}

// LoadBootstrap Loads configuration files in various formats from a directory,
// and parses them into a struct.
func LoadBootstrap(serviceName string, cfg *SourceConfig, l log.Logger) (*configs.Bootstrap, error) {
	source, err := LoadBootConfig(serviceName, cfg, l)
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
