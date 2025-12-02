// Package confpb implements the functions, types, and interfaces for the module.
package conf

import (
	"cmp"
	"fmt"

	confpb "basic-layout/simple/simple_app/internal/conf/pb"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
	"github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/config/trace/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/config/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/interfaces"
)

type Config struct {
	bootstrap        confpb.Bootstrap
	defaultDiscovery string
}

func (c *Config) DecodeCaches() (*datav1.Caches, error) {
	return c.bootstrap.GetData().GetCaches(), nil
}

func (c *Config) DecodeDatabases() (*datav1.Databases, error) {
	return c.bootstrap.GetData().GetDatabases(), nil
}

func (c *Config) DecodeObjectStores() (*datav1.ObjectStores, error) {
	return c.bootstrap.GetData().GetObjectStores(), nil
}

func (c *Config) DecodedConfig() any {
	//TODO implement me
	panic("implement me")
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

	// Determine the default discovery key with a clear priority: active > default > global default.
	defaultKey := cmp.Or(
		c.bootstrap.GetDiscoveries().GetActive(),
		c.bootstrap.GetDiscoveries().GetDefault(),
		interfaces.GlobalDefaultKey)

	discoveries := c.bootstrap.GetDiscoveries()
	if discoveries != nil {
		confpb := discoveries.GetConfigs()
		// If there's only one discovery config, it should be the default, regardless of the key.
		switch {
		case len(confpb) == 1:
			c.defaultDiscovery = confpb[0].GetName()
		case len(confpb) > 1:
			for _, discovery := range confpb {
				if discovery.GetName() == defaultKey {
					c.defaultDiscovery = defaultKey
					break
				}
			}
			if c.defaultDiscovery == "" {
				return nil, fmt.Errorf("failed to determine default discovery: key '%s' not found among multiple configurations", defaultKey)
			}
		}
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
