package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/utils/replacer"
)

type Discovery struct {
	Name    string `yaml:"name" toml:"name" json:"name"`
	Address string `yaml:"address" toml:"address" json:"address"`
}

type ConsulConfig struct {
	Address string `yaml:"address" toml:"address" json:"address"`
}

type Bootstrap struct {
	Name      string       `yaml:"name" toml:"name" json:"name"`
	Type      string       `yaml:"type" toml:"type" json:"type"`
	Consul    ConsulConfig `yaml:"consul" toml:"consul" json:"consul"`
	Discovery Discovery    `yaml:"discovery" toml:"discovery" json:"discovery"`
}

func (c *Bootstrap) Setup() {
	return
}

// Load Loads configuration files in various formats from a directory,
// and parses them into a struct.
func Load(c *Bootstrap, path string) error {
	*c = DefaultBootstrap

	fullname, _ := filepath.Abs(path)
	info, err := os.Stat(fullname)
	if err != nil {
		return errors.Wrapf(err, "failed to state file %s", fullname)
	}
	if !info.IsDir() {
		if err := parseConfigFile(c, path); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", path)
		}
		return nil
	}
	if err := filepath.WalkDir(fullname, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", path)
		} else if d.IsDir() {
			return nil
		}
		if err := parseConfigFile(c, path); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", path)
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to walk config path %s", fullname)
	}

	return nil
}

func LoadWithEnv(c *Bootstrap, path string, envs map[string]string) error {
	*c = DefaultBootstrap
	fullname, _ := filepath.Abs(path)
	info, err := os.Stat(fullname)
	if err != nil {
		return errors.Wrapf(err, "failed to state file %s", fullname)
	}

	if !info.IsDir() {
		typo := codec.TypeFromExt(filepath.Ext(fullname))
		if typo == codec.UNKNOWN {
			return nil
		}

		file, err := os.ReadFile(fullname)
		if err != nil {
			return errors.Wrapf(err, "failed to read file %s", fullname)
		}
		if err := typo.Unmarshal(file, c); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", fullname)
		}

	}
	if err := filepath.WalkDir(fullname, func(path string, d os.DirEntry, err error) error {
		typo := codec.TypeFromExt(filepath.Ext(path))
		if typo == codec.UNKNOWN {
			return nil
		}
		file, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read file %s", path)
		}

		if err := typo.Unmarshal(file, c); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", path)
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to walk config path %s", fullname)
	}
	if err := replacer.ObjectReplacer(c, envs); err != nil {
		return errors.Wrap(err, "failed to parse config objects")
	}
	return nil
}

const (
	extNames = `.json,.toml,.yaml",.yml`
)

func parseConfigFile(c *Bootstrap, path string) error {
	ext := filepath.Ext(path)
	if ext == "" || !strings.Contains(extNames, ext) {
		return nil
	}
	err := codec.DecodeFromFile(path, c)
	if err != nil {
		return err
	}

	return nil
}

var DefaultBootstrap = Bootstrap{
	Name: "helloworld",
	Type: "consul",
	Consul: ConsulConfig{
		Address: "${consul_address}",
	},
	Discovery: Discovery{
		Name:    "helloworld",
		Address: "${discovery_address}",
	},
}
