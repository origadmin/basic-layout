package source

import (
	"path"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
	"github.com/origadmin/toolkits/errors"

	"origadmin/basic-layout/internal/configs"
)

func NewSource(path string, sourceConfig *configs.Consul) (config.Source, error) {
	client, err := api.NewClient(&api.Config{
		Address: sourceConfig.Address,
		Scheme:  sourceConfig.Scheme,
	})
	if err != nil {
		return nil, errors.Wrap(err, "consul client error")
	}
	source, err := consul.New(client,
		consul.WithPath(path),
	)
	if err != nil {
		return nil, errors.Wrap(err, "consul source error")
	}
	return source, nil
}

// ConfigPath path
func ConfigPath(dir, name string) string {
	return path.Join("configs", dir, name)
}
