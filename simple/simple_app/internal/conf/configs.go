// Package configs implements the functions, types, and interfaces for the module.
package conf

import (
	"cmp"

	"basic-layout/simple/simple_app/configs"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	cachev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/cache/v1"
	databasev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/database/v1"
	"github.com/origadmin/runtime/api/gen/go/runtime/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/runtime/trace/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/runtime/interfaces"
)

type Config struct {
	bootstrap configs.Bootstrap
}

func (c *Config) DecodeApp() (*appv1.App, error) {
	return c.bootstrap.GetApp(), nil
}

func (c *Config) DecodeData() (*datav1.Data, error) {
	return c.bootstrap.GetData(), nil
}

func (c *Config) DecodeDefaultDiscovery() (string, error) {
	return cmp.Or(
		c.bootstrap.GetDiscoveries().GetDefault(),
		c.bootstrap.GetDiscoveries().GetActive(),
		interfaces.GlobalDefaultKey), nil
}

func (c *Config) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	return c.bootstrap.GetDiscoveries(), nil
}

func (c *Config) DecodeLogger() (*loggerv1.Logger, error) {
	return c.bootstrap.GetLogger(), nil
}

func (c *Config) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	return c.bootstrap.GetMiddlewares(), nil
}

func (c *Config) DecodeServers() (*transportv1.Servers, error) {
	return c.bootstrap.GetServers(), nil
}

func (c *Config) DecodeClients() (*transportv1.Clients, error) {
	return c.bootstrap.GetClients(), nil
}

func (c *Config) Transform(config interfaces.Config, _ interfaces.StructuredConfig) (interfaces.StructuredConfig,
	error) {
	err := config.Decode("", &c.bootstrap)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func DefaultApp() *appv1.App {
	return &appv1.App{
		Id:      "test-app-id",
		Name:    "test-app-name",
		Version: "test-app-version",
		Env:     "dev",
	}
}

func DefaultServers() *transportv1.Servers {
	return &transportv1.Servers{
		Configs: []*transportv1.Server{
			{
				Name:     "grpc_server",
				Protocol: "grpc",
				Grpc: &grpcv1.Server{
					Network:     "tcp",
					Addr:        "0.0.0.0:9090",
					Middlewares: []string{"recovery", "logger"},
				},
			},
			{
				Name:     "http_server",
				Protocol: "http",
				Http: &httpv1.Server{
					Network:     "tcp",
					Addr:        "0.0.0.0:8080",
					Middlewares: []string{"recovery", "logger"},
				},
			},
		},
	}
}

func DefaultDiscoveries() *discoveryv1.Discoveries {
	return &discoveryv1.Discoveries{
		Configs: []*discoveryv1.Discovery{},
	}
}

func DefaultData() *datav1.Data {
	return &datav1.Data{
		Databases: &datav1.Databases{
			Configs: []*databasev1.DatabaseConfig{
				{
					Name:    "default",
					Dialect: "sqlite3",
					Source:  "file:./test.db?cache=shared&mode=memory&_fk=1",
				},
			},
		},
		Caches: &datav1.Caches{
			Configs: []*cachev1.CacheConfig{
				{
					Name:   "default",
					Driver: "memory",
				},
			},
		},
	}
}

func DefaultMiddlewares() *middlewarev1.Middlewares {
	return &middlewarev1.Middlewares{
		Configs: []*middlewarev1.Middleware{
			{
				Name:     "recovery",
				Type:     "recovery",
				Recovery: &middlewarev1.Recovery{},
			},
			{
				Name:    "logging",
				Type:    "logging",
				Logging: &middlewarev1.Logging{},
			},
		},
	}
}

func DefaultTrace() *tracev1.Trace {
	return &tracev1.Trace{
		Name:        "jaeger",
		Endpoint:    "localhost:6831",
		ServiceName: "test-app-name",
	}
}
