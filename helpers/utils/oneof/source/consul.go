package source

import (
	"path"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/hashicorp/consul/api"

	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/runtime/config"
)

func NewSource(name string, sourceConfig *config.SourceConfig_Consul) (config.Source, error) {
	client, err := api.NewClient(&api.Config{
		Address: sourceConfig.GetAddress(),
		Scheme:  sourceConfig.GetScheme(),
		Token:   sourceConfig.GetToken(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "consul client error")
	}
	source, err := consul.New(client,
		consul.WithPath(name),
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
