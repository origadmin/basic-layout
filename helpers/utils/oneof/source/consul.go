/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package source

import (
	"path"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/hashicorp/consul/api"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/runtime/source/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/toolkits/errors"
)

func NewSource(sourceConfig *sourcev1.SourceConfig, opts ...options.Option) (config.KSource, error) {
	consulConfig := sourceConfig.GetConsul()
	if consulConfig == nil {
		return nil, errors.New("consul config is nil")
	}
	client, err := api.NewClient(&api.Config{
		Address: consulConfig.GetAddress(),
		Scheme:  consulConfig.GetScheme(),
		Token:   consulConfig.GetToken(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "consul client error")
	}
	source, err := consul.New(client,
		consul.WithPath(consulConfig.Path),
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
