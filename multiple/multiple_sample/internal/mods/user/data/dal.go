/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package data

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/user/data/ent"
	"github.com/origadmin/runtime/data/storage"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(bootstrap *configs.Bootstrap, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(rt.Logger())

	provider, err := storage.New(rt.StructuredConfig())
	if err != nil {
		return nil, nil, err
	}

	database, err := provider.DefaultDatabase()
	if err != nil {
		return nil, nil, err
	}

	activeDB := entsql.OpenDB(database.Dialect(), database.DB())
	entClient := ent.NewClient(ent.Driver(activeDB))
	// Run the auto migration tool.
	if err := entClient.Schema.Create(context.Background()); err != nil {
		logHelper.Fatalf("failed creating schema resources: %v", err)
	}

	cache, err := provider.DefaultCache()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		logHelper.Info("closing the data resources")
		if entClient != nil {
			if err := entClient.Close(); err != nil {
				logHelper.Errorf("failed to close ent client: %v", err)
			}
		}
	}

	return &Data{
		config:    rt.StructuredConfig(),
		provider:  provider,
		entClient: entClient,
		cache:     cache,
	}, cleanup, nil

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Database{}, cleanup, nil
}
