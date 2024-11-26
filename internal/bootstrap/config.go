package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goexts/generic/settings"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/codec/json"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/runtime"
	"github.com/origadmin/toolkits/runtime/bootstrap"
	"github.com/origadmin/toolkits/runtime/config"
	configv1 "github.com/origadmin/toolkits/runtime/gen/go/config/v1"

	"origadmin/basic-layout/helpers/oneof/source"
	"origadmin/basic-layout/internal/configs"
)

type Bootstrap = bootstrap.Bootstrap

type Config = configv1.SourceConfig

func init() {
	//cfg := runtime.ConfigBuildFunc(NewFileConfig)
	runtime.RegisterConfigFunc("file", NewFileConfig)
}
func DefaultBootstrap() *Bootstrap {
	//boot := &Config{}
	//boot.Flags = bootstrap.NewFlags(serviceName, version)
	//return boot
	return bootstrap.DefaultBootstrap()
}

// LoadSourceFiles Loads configuration files in various formats from a directory,
// and parses them into a config.
func LoadSourceFiles(wd, path string) *Config {
	path = WorkPath(wd, path)
	cfg := NewFileSource(path)
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

func FromRemote(serviceName string, source *Config, ss ...Setting) (*configs.Bootstrap, error) {
	o := settings.ApplyOrZero(ss...)
	switch source.Type {
	case "consul":
		return FromConsul(serviceName, source, o)
	default:

	}
	return nil, errors.Errorf("unsupported config type: %s", source.Type)
}

func FromConsul(serviceName string, cfg *Config, option *Option) (*configs.Bootstrap, error) {
	if cfg.Consul == nil {
		return nil, errors.String("consul config is nil")
	}
	cfg.Consul.Path = source.ConfigPath(serviceName, "bootstrap.json")
	return LoadBootstrap(cfg, option)
}

func FromFlags(boot *Bootstrap, ss ...Setting) (*configs.Bootstrap, error) {
	o := settings.ApplyOrZero(ss...)
	path := WorkPath(boot.WorkDir, boot.ConfigPath)
	if path == "" {
		return nil, errors.String("config path is empty")
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to state file %s", path)
	}
	log.Infof("loading config from %s", path)

	var cfg *Config
	if stat.IsDir() {
		cfg = loadSourceFromDir(path)
	} else {
		cfg = loadSourceFromFile(path)
	}
	if cfg == nil {
		cfg = NewFileSource(path)
	}

	return LoadBootstrap(PathToSource(cfg, boot.ServiceName()), o)
}

func FromLocal(serviceName string, source *Config, option *Option) (*configs.Bootstrap, error) {
	if source.File == nil {
		return nil, errors.String("file config is nil")
	}

	path := WorkPath("", source.File.Path)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to state file %s", path)
	}
	log.Infof("loading config from %s", path)

	var cfg *Config
	if stat.IsDir() {
		cfg = loadSourceFromDir(path)
	} else {
		cfg = loadSourceFromFile(path)
	}
	if cfg == nil {
		cfg = NewFileSource(path)
	}
	return LoadBootstrap(PathToSource(cfg, serviceName), option)
}

func FromLocalPath(serviceName string, path string, ss ...Setting) (*configs.Bootstrap, error) {
	o := settings.ApplyOrZero(ss...)
	return FromLocal(serviceName, NewFileSource(path), o)
}

// loadSourceFromDir loads configuration from a directory.
func loadSourceFromDir(path string) *Config {
	var cfg Config
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
func loadSourceFromFile(path string) *Config {
	var cfg Config
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return nil
	}
	return &cfg
}

// LoadBootstrap Loads configuration files in various formats from a directory,
// and parses them into a struct.
func LoadBootstrap(cfg *Config, option *Option) (*configs.Bootstrap, error) {
	source, err := runtime.NewConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new kratos config error")
	}
	log.Infof("load bootstrap config from: %s, values: %+v", cfg.Type, PrintString(cfg))
	bs := DefaultConfig()
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

func NewFileConfig(ccfg *Config, opts ...config.Option) (config.Config, error) {
	if ccfg.EnvArgs != nil {
		opts = append(opts, config.WithSource(file.NewSource(ccfg.File.Path), envf.WithEnv(ccfg.EnvArgs, ccfg.EnvPrefixes...)))
	} else {
		opts = append(opts, config.WithSource(file.NewSource(ccfg.File.Path)))
	}
	return config.New(opts...), nil
}

func NewFileSource(path string) *Config {
	return &Config{
		Type: "file",
		File: &configv1.SourceConfig_File{
			Path: path,
		},
	}
}

func WorkPath(wd, path string) string {
	if wd != "" && !filepath.IsAbs(path) {
		path = filepath.Join(wd, path)
	}
	path, _ = filepath.Abs(path)
	return path
}

func PrintString(v any) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(bytes)
}

func PathToSource(cfg *Config, serviceName string) *Config {
	switch cfg.Type {
	case "file":
		cfg.File.Path = WorkPath("", cfg.File.Path)
	case "consul":
		cfg.Consul.Path = source.ConfigPath(serviceName, "bootstrap.json")
	//case "etcd":
	//	return source.ConfigPath(serviceName, "bootstrap.json")
	//default:
	//	return source.ConfigPath(serviceName, "bootstrap.json")
	default:

	}
	return cfg
}
