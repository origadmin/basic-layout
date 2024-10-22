package bootstrap

import (
	registryconsul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"origadmin/basic-layout/internal/configs"
)

// NewRegistrar creates a new registrar.
func NewRegistrar(bootstrap *configs.Bootstrap, l log.Logger) registry.Registrar {
	d := bootstrap.Discovery
	var reg registry.Registrar
	switch Type(d.Type) {
	case Default:
		return nil
	case Consul:
		cfg := d.GetConsul()
		if cfg == nil {
			return nil
		}
		c := api.DefaultConfig()
		c.Address = cfg.Address
		c.Scheme = cfg.Scheme
		c.Token = cfg.Token
		c.Datacenter = cfg.Datacenter
		//c.Tag = cfg.Consul.Tag
		//c.HealthCheckInterval = d.Consul.HealthCheckInterval
		//c.HealthCheckTimeout = d.Consul.HealthCheckTimeout
		cli, err := api.NewClient(c)
		if err != nil {
			log.Fatalf("consul client %+v", err)
		}
		reg = registryconsul.New(
			cli,
			registryconsul.WithHeartbeat(cfg.HeartBeat),
			registryconsul.WithHealthCheck(cfg.HealthCheck),
		)
		log.Infof("consul: %s", cfg.Address)
	default:
		panic(errors.Errorf("unknown discovery type: %s", d.Type))
	}

	return reg
}
