package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/toolkits/runtime"
	"github.com/origadmin/toolkits/runtime/registry"

	"origadmin/basic-layout/internal/configs"
)

// NewRegistrar creates a new registrar.
func NewRegistrar(bootstrap *configs.Bootstrap) registry.Registrar {
	cfg := bootstrap.GetRegistry()
	if cfg == nil {
		log.Infof("no registry config")
		return nil
	}
	registrar, err := runtime.NewRegistrar(cfg)
	if err != nil {
		log.Errorf("new registrar failed: %v", err)
		return nil
	}
	return registrar
}

func NewDiscovery(bootstrap *configs.Bootstrap) registry.Discovery {
	cfg := bootstrap.GetRegistry()
	if cfg == nil {
		log.Infof("no registry config")
		return nil
	}
	discovery, err := runtime.NewDiscovery(cfg)
	if err != nil {
		log.Errorf("new discovery failed: %v", err)
		return nil
	}
	return discovery
}

func NewDiscoveries(bootstrap *configs.Bootstrap) (map[string]registry.Discovery, error) {
	registries := bootstrap.GetRegistries()
	discoveries := make(map[string]registry.Discovery, len(registries))
	for i := range registries {
		cfg, err := runtime.NewDiscovery(registries[i])
		if err != nil {
			return nil, err
		}
		discoveries[registries[i].ServiceName] = cfg
	}

	return discoveries, nil
}
