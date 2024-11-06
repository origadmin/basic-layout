package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/origadmin/toolkits/runtime/kratos"

	"origadmin/basic-layout/internal/configs"
)

// NewRegistrar creates a new registrar.
func NewRegistrar(bootstrap *configs.Bootstrap) registry.Registrar {
	if bootstrap == nil || bootstrap.Registry == nil {
		log.Infof("no registry config")
		return nil
	}
	registrar, err := kratos.NewRegistrar(bootstrap.Registry)
	if err != nil {
		log.Errorf("new registrar failed: %v", err)
		return nil
	}
	return registrar
}

func NewDiscovery(bootstrap *configs.Bootstrap) registry.Discovery {
	if bootstrap == nil || bootstrap.Registry == nil {
		log.Infof("no registry config")
		return nil
	}
	discovery, err := kratos.NewDiscovery(bootstrap.Registry)
	if err != nil {
		log.Errorf("new discovery failed: %v", err)
		return nil
	}
	return discovery
}

func NewDiscoveries(source *Config, serviceNames ...string) map[string]registry.Discovery {
	discoveries := make(map[string]registry.Discovery, len(serviceNames))
	for _, name := range serviceNames {
		bs, err := FromRemote(name, source)
		if err != nil {
			return nil
		}
		if bs == nil {
			continue
		}
		discovery := NewDiscovery(bs)
		if discovery == nil {
			continue
		}
		discoveries[name] = discovery
	}
	return discoveries
}

func NewRegistrars(source *Config, serviceNames ...string) map[string]registry.Registrar {
	registrars := make(map[string]registry.Registrar, len(serviceNames))
	for _, name := range serviceNames {
		bs, err := FromRemote(name, source)
		if err != nil {
			return nil
		}
		if bs == nil {
			continue
		}
		registrar := NewRegistrar(bs)
		if registrar == nil {
			continue
		}
		registrars[name] = registrar
	}
	return registrars

}
