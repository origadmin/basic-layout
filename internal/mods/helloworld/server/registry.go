package server

import (
	registryconsul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"origadmin/basic-layout/internal/mods/helloworld/conf"
)

// NewRegistrar creates a new registrar.
func NewRegistrar(bootstrap *conf.Bootstrap, l log.Logger) registry.Registrar {
	c := api.DefaultConfig()
	d := bootstrap.Discovery
	var reg registry.Registrar
	switch d.Type {
	case "none":
		return nil
	case "consul":
		cfg := d.GetConsul()
		c.Address = cfg.Address
		c.Scheme = cfg.Scheme
		c.Token = cfg.Token
		c.Datacenter = cfg.Datacenter
		//c.Tag = d.Consul.Tag
		//c.HealthCheckInterval = d.Consul.HealthCheckInterval
		//c.HealthCheckTimeout = d.Consul.HealthCheckTimeout
		cli, err := api.NewClient(c)
		if err != nil {
			panic(errors.Wrap(err, "consul client"))
		}
		reg = registryconsul.New(
			cli,
			registryconsul.WithHeartbeat(cfg.HeartBeat),
			registryconsul.WithHealthCheck(cfg.HealthCheck),
		)
		log.Infof("consul: %s", bootstrap.Discovery.Consul.Address)
	default:
		panic(errors.Errorf("unknown discovery type: %s", d.Type))
	}

	return reg
}
