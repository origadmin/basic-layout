package bootstrap

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/utils/replacer"

	"origadmin/basic-layout/internal/config"
)

// Load Loads configuration files in various formats from a directory,
// and parses them into a struct.
func Load(c *config.Bootstrap, path string) error {
	*c = *config.DefaultBootstrap

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

func LoadWithEnv(c *config.Bootstrap, path string, envs map[string]string) error {
	*c = *config.DefaultBootstrap
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

func parseConfigFile(c *config.Bootstrap, path string) error {
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
