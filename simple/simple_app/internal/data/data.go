package data

import (
	"context"
	"errors"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/wire"

	"basic-layout/simple/simple_app/internal/biz"
	"basic-layout/simple/simple_app/internal/data/entity/ent"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/interfaces"

	"github.com/origadmin/runtime/data/storage"
	ifacestorage "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/log"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSimpleRepo)

// Data encapsulates ent client and cache.
type Data struct {
	entClient *ent.Client
	cache     ifacestorage.Cache
	provider  ifacestorage.Provider
	config    interfaces.StructuredConfig
}

var ErrNoDatabaseConfig = errors.New("no database config found")

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

// GetEntClient returns the ent client.
func (d *Data) GetEntClient() *ent.Client {
	return d.entClient
}

// GetCache returns the cache.
func (d *Data) GetCache() ifacestorage.Cache {
	return d.cache
}

type simpleRepo struct {
	data *Data
	log  *log.Helper
}

// NewSimpleRepo creates a new implementation of biz.SimpleRepo.
func NewSimpleRepo(rt *runtime.Runtime, data *Data) biz.SimpleRepo {
	return &simpleRepo{
		data: data,
		log:  log.NewHelper(rt.Logger()),
	}
}

func (r *simpleRepo) Save(ctx context.Context, s *biz.Simple) (*biz.Simple, error) {
	// Placeholder implementation
	r.log.WithContext(ctx).Infof("Saving Simple: %s", s.Name)
	return s, nil
}

func (r *simpleRepo) Update(ctx context.Context, s *biz.Simple) (*biz.Simple, error) {
	// Placeholder implementation
	return s, nil
}

func (r *simpleRepo) FindByID(ctx context.Context, id int64) (*biz.Simple, error) {
	// Placeholder implementation
	return &biz.Simple{Name: "found_by_id"}, nil
}

func (r *simpleRepo) ListByName(ctx context.Context, name string) ([]*biz.Simple, error) {
	// Placeholder implementation
	return []*biz.Simple{{Name: name}}, nil
}

func (r *simpleRepo) ListAll(ctx context.Context) ([]*biz.Simple, error) {
	// Placeholder implementation
	return []*biz.Simple{{Name: "all"}}, nil
}
