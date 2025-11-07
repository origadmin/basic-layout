package data

import (
	"context"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/internal/features/user/data/ent"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/data/storage"
	"github.com/origadmin/runtime/interfaces"
	ifacestorage "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data encapsulates ent client and cache.
type Data struct {
	entClient *ent.Client
	cache     ifacestorage.Cache
	provider  ifacestorage.Provider
	config    interfaces.StructuredConfig
}

// NewData creates a new Data instance.
func NewData(rt *runtime.Runtime) (*Data, func(), error) {
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
}
