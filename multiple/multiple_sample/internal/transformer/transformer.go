// Package transformer implements the functions, types, and interfaces for the module.
package transformer

import (
	"fmt"

	"basic-layout/multiple/multiple_sample/configs"
	datav1 "github.com/origadmin/runtime/api/gen/go/runtime/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
)

type Config struct {
	app       *appv1.App
	config    interfaces.Config
	bootstrap configs.Bootstrap
}

func (c *Config) Bootstrap() *configs.Bootstrap {
	return &c.bootstrap
}

func (c *Config) DecodeData() (*datav1.Data, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Config) DecodeDefaultDiscovery() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Config) DecodeServers() (*transportv1.Servers, error) {
	return &transportv1.Servers{
		Configs: c.bootstrap.GetService().GetServers(),
	}, nil
}

func (c *Config) DecodeClients() (*transportv1.Clients, error) {
	return &transportv1.Clients{
		Configs: c.bootstrap.GetService().GetClients(),
	}, nil
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

func (c *Config) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	return &discoveryv1.Discoveries{
		Configs: []*discoveryv1.Discovery{c.bootstrap.GetDiscovery()},
	}, nil
}

func (c *Config) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	middlewares := &middlewarev1.Middlewares{}
	middlewares.Configs = c.bootstrap.GetService().GetMiddlewares()
	return middlewares, nil
}

func (c *Config) Transform(config interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
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
		c.bootstrap.Service = &serverCfg
	}

	// If we still don't have a server config, create an empty one
	if c.bootstrap.Service == nil {
		logger.Warn("No server configuration found, using defaults")
		c.bootstrap.Service = &configs.ServiceConfig{
			Servers: []*transportv1.Server{
				{
					Name: c.app.GetName(),
				},
			}, // Version field doesn't exist in servicev1.Service
		}
	}

	// Ensure service name is set
	if c.bootstrap.Service.Servers == nil {
		c.bootstrap.Service.Servers = []*transportv1.Server{
			{
				Name: c.app.GetName(),
			},
		}
	}

	// Use app name if service name is not set
	if c.bootstrap.Service.Servers[0].Name == "" && c.app != nil {
		c.bootstrap.Service.Servers[0].Name = c.app.GetName()
	}

	logger.Infof("Service name: %s",
		c.bootstrap.Service.Servers[0].Name,
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
	return func(config interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
		return New(cfg).Transform(config, sc)
	}
}

func Transform(config interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
	return New(nil).Transform(config, sc)
}

var _ bootstrap.ConfigTransformer = (*Config)(nil)
