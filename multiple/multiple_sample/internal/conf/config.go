// Package conf implements the functions, types, and interfaces for the module.
package conf

import (
	"fmt"

	"github.com/goexts/generic/cmp"

	confpb "basic-layout/multiple/multiple_sample/internal/conf/pb"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
)

type Config struct {
	config    interfaces.Config
	bootstrap confpb.Bootstrap
}

func (c *Config) DecodeData() (*datav1.Data, error) {
	return c.bootstrap.GetData(), nil
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
	return &c.bootstrap
}

func (c *Config) Bootstrap() *confpb.Bootstrap {
	return &c.bootstrap
}

func (c *Config) DecodeDefaultDiscovery() (string, error) {
	return cmp.Or(c.bootstrap.GetDiscoveries().GetActive(), c.bootstrap.GetDiscoveries().GetDefault(), "default"), nil
}

func (c *Config) DecodeServers() (*transportv1.Servers, error) {
	return c.bootstrap.GetServers(), nil
}

func (c *Config) DecodeClients() (*transportv1.Clients, error) {
	return c.bootstrap.GetClients(), nil
}

func (c *Config) DecodeLogger() (*loggerv1.Logger, error) {
	return c.bootstrap.GetLogger(), nil
}

func (c *Config) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	return c.bootstrap.GetDiscoveries(), nil
}

func (c *Config) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	return c.bootstrap.GetMiddlewares(), nil
}

func (c *Config) Transform(config interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
	c.config = config
	logger := log.NewHelper(log.DefaultLogger)

	// Try to decode the entire config first
	if err := config.Decode("", &c.bootstrap); err != nil {
		return nil, fmt.Errorf("failed to decode bootstrap config: %w", err)
	}

	// Log the final configuration for debugging
	logger.Debugf("Final bootstrap config: %+v", &c.bootstrap)

	return c, nil
}

func New() *Config {
	return &Config{}
}

func TransformAfter(cfg *appv1.App) bootstrap.ConfigTransformFunc {
	return func(config interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
		return New().Transform(config, sc)
	}
}

func Transform(config interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
	return New().Transform(config, sc)
}

var _ bootstrap.ConfigTransformer = (*Config)(nil)
