// Package transformer implements the functions, types, and interfaces for the module.
package transformer

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	"github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"origadmin/basic-layout/api/v1/gen/go/configs"
)

type Config struct {
	app       *appv1.App
	config    interfaces.Config
	bootstrap configs.Bootstrap
}

func (c *Config) Load() error {
	return c.config.Load()
}

func (c *Config) Decode(key string, value any) error {
	return c.config.Decode(key, value)
}

func (c *Config) Raw() any {
	return c.config.Raw()
}

func (c *Config) Close() error {
	return c.config.Close()
}

func (c *Config) DecodeApp() (*appv1.App, error) {
	return c.app, nil
}

func (c *Config) DecodeLogger() (*loggerv1.Logger, error) {
	return &loggerv1.Logger{}, nil
}

func (c *Config) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	discoveries := make(map[string]*discoveryv1.Discovery)
	v := c.bootstrap.GetDiscovery()
	if v != nil {
		discoveries[v.Name] = v
	}
	return discoveries, nil
}

func (c *Config) DecodeMiddleware() (*middlewarev1.Middlewares, error) {
	middlewares := &middlewarev1.Middlewares{}
	middlewares.Middlewares = c.bootstrap.GetServer().GetService().GetMiddlewares()
	return middlewares, nil
}

func (c *Config) Transform(config interfaces.Config) (interfaces.StructuredConfig, error) {
	c.config = config
	if err := config.Decode("", &c.bootstrap); err != nil {
		return nil, err
	}
	return c, nil
}

func New(app *appv1.App) *Config {
	if app == nil {
		app = &appv1.App{}
	}
	return &Config{
		app: app,
	}
}

func TransformAfter(cfg *appv1.App) bootstrap.ConfigTransformFunc {
	return func(config interfaces.Config) (interfaces.StructuredConfig, error) {
		return New(cfg).Transform(config)
	}
}

func Transform(config interfaces.Config) (interfaces.StructuredConfig, error) {
	return New(nil).Transform(config)
}

var _ bootstrap.ConfigTransformer = (*Config)(nil)
