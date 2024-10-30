package bootloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/contrib/config/envf"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/runtime/bootstrap"
	"github.com/origadmin/toolkits/runtime/config"
	"github.com/origadmin/toolkits/runtime/kratos"
	"github.com/origadmin/toolkits/utils/replacer"

	"origadmin/basic-layout/internal/configs"
)

type BootFlags struct {
	bootstrap.Bootstrap
	bootstrap.Flags
}

func init() {
	kratos.RegistryConfig("file", NewFileConfig)
}
func NewBootFlags(serviceName, version string) BootFlags {
	return BootFlags{
		Bootstrap: bootstrap.DefaultBootstrap(),
		Flags:     bootstrap.NewFlags(serviceName, version),
	}
}

func Load(flags BootFlags, useEnv bool) (*configs.Bootstrap, error) {
	sourceConfig := LoadSourceFiles(flags.WorkDir, flags.ConfigPath)
	fmt.Printf("source: %+v\n", sourceConfig)
	switch sourceConfig.Type {
	case "file":
		//sourceConfig.File.Path
	case "consul":
		sourceConfig.Consul.Path = filepath.Join("configs", flags.ServiceName, "bootstrap.json")
	default:
		return nil, errors.New("invalid config type")
	}

	return LoadBootstrap(sourceConfig, nil)
}

// LoadSourceFiles Loads configuration files in various formats from a directory,
// and parses them into a config.
func LoadSourceFiles(wd, path string) *config.SourceConfig {
	if !filepath.IsAbs(path) {
		path = filepath.Join(wd, path)
	}
	path, _ = filepath.Abs(path)
	cfg := config.NewFileConfig(path)
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

func FromLocal(serviceName, path string, l log.Logger) (*configs.Bootstrap, error) {
	path, _ = filepath.Abs(path)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to state file %s", path)
	}
	log.NewHelper(l).Infof("loading config from %s", path)

	var cfg *config.SourceConfig
	if stat.IsDir() {
		cfg = loadSourceFromDir(path)
	} else {
		cfg = loadSourceFromFile(path)
	}
	if cfg == nil {
		cfg = config.NewFileConfig(path)
	}

	source, err := kratos.NewConfig(cfg)
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

// loadSourceFromDir loads configuration from a directory.
func loadSourceFromDir(path string) *config.SourceConfig {
	var cfg config.SourceConfig
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
func loadSourceFromFile(path string) *config.SourceConfig {
	var cfg config.SourceConfig
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return nil
	}
	return &cfg
}

// LoadBootConfig returns a new kratos config .
func LoadBootConfig(flags BootFlags, cfg *config.SourceConfig, l log.Logger) (config.Config, error) {
	//srcs, err := sourceFromConfig(serviceName, cfg, l)
	//if err != nil {
	//	return nil, err
	//}
	return kratos.NewConfig(cfg)
}

// LoadBootstrap Loads configuration files in various formats from a directory,
// and parses them into a struct.
func LoadBootstrap(cfg *config.SourceConfig, l log.Logger) (*configs.Bootstrap, error) {

	source, err := kratos.NewConfig(cfg)
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

func NewFileConfig(ccfg *config.SourceConfig, opts ...config.Option) (config.Config, error) {
	if ccfg.EnvArgs != nil {
		opts = append(opts, config.WithSource(file.NewSource(ccfg.File.Path), envf.WithEnv(ccfg.EnvArgs, ccfg.EnvPrefixes...)))
	} else {
		opts = append(opts, config.WithSource(file.NewSource(ccfg.File.Path)))
	}
	return config.New(opts...), nil
}
