// Package transformer implements the functions, types, and interfaces for the module.
package transformer

import (
	"fmt"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	"github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	servicev1 "github.com/origadmin/runtime/api/gen/go/runtime/service/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	"origadmin/basic-layout/internal/configs"
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
	return &loggerv1.Logger{
		Default: true,
		Level:   "debug",
	}, nil
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
	logger := log.NewHelper(log.DefaultLogger)

	// Try to decode the entire config first
	if err := config.Decode("", &c.bootstrap); err != nil {
		logger.Errorf("Failed to decode bootstrap config: %v", err)
		// Try to decode just the server part
		var serverCfg configs.ServiceConfig
		if err := config.Decode("server", &serverCfg); err != nil {
			logger.Errorf("Failed to decode server config: %v", err)
			return nil, fmt.Errorf("failed to decode configuration: %v", err)
		}
		c.bootstrap.Server = &serverCfg
	}

	// If we still don't have a server config, create an empty one
	if c.bootstrap.Server == nil {
		logger.Warn("No server configuration found, using defaults")
		c.bootstrap.Server = &configs.ServiceConfig{
			Service: &servicev1.Service{
				Name: c.app.GetName(),
			}, // Version field doesn't exist in servicev1.Service
		}
	}

	// Ensure service name is set
	if c.bootstrap.Server.Service == nil {
		c.bootstrap.Server.Service = &servicev1.Service{}
	}

	// Use app name if service name is not set
	if c.bootstrap.Server.Service.Name == "" && c.app != nil {
		c.bootstrap.Server.Service.Name = c.app.GetName()
	}

	logger.Infof("Service name: %s",
		c.bootstrap.Server.Service.Name,
	)

	// Log the final configuration for debugging
	logger.Debugf("Final bootstrap config: %+v", &c.bootstrap)

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
