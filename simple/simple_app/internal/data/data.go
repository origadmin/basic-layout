package data

import (
	"cmp"
	"context"
	"errors"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/wire"

	"basic-layout/simple/simple_app/internal/biz"
	"basic-layout/simple/simple_app/internal/data/entity/ent"
	"github.com/origadmin/runtime"

	"github.com/origadmin/runtime/interfaces"
	ifacestorage "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/storage"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSimpleRepo)

// Data encapsulates ent client and cache.
type Data struct {
	entClient *ent.Client
	dbs       map[string]*entsql.Driver
	cache     ifacestorage.Cache
}

var ErrNoDatabaseConfig = errors.New("no database config found")

// NewData creates a new Data instance.
func NewData(rt *runtime.Runtime) (*Data, func(), error) {
	logHelper := log.NewHelper(rt.Logger())

	dataConfig, err := rt.StructuredConfig().DecodeData()
	if err != nil {
		return nil, nil, err
	}

	dbs := make(map[string]*entsql.Driver, len(dataConfig.GetDatabases().GetConfigs()))
	for _, config := range dataConfig.GetDatabases().GetConfigs() {
		database, err := storage.OpenDatabase(config)
		if err != nil {
			return nil, nil, err
		}
		key := cmp.Or(config.GetName(), config.GetDialect())
		dbs[key] = entsql.OpenDB(config.GetDialect(), database)
	}

	var activeDB *entsql.Driver
	// Determine the active database connection.
	// Priority: active > default > first in map.
	if config := dataConfig.GetDatabases(); config != nil {
		switch len(dbs) {
		case 0:
			return nil, nil, ErrNoDatabaseConfig
		case 1:
			for key := range dbs {
				activeDB = dbs[key]
				break
			}
		default:
			defaultKey := cmp.Or(config.GetActive(), config.GetDefault(), interfaces.GlobalDefaultKey)
			activeDB = dbs[defaultKey]
		}
	}
	if activeDB == nil {
		return nil, nil, ErrNoDatabaseConfig
	}

	entClient := ent.NewClient(ent.Driver(activeDB))
	// Run the auto migration tool.
	if err := entClient.Schema.Create(context.Background()); err != nil {
		logHelper.Fatalf("failed creating schema resources: %v", err)
	}

	caches := make(map[string]ifacestorage.Cache)
	for _, config := range dataConfig.GetCaches().GetConfigs() {
		cache, err := storage.New(config)
		if err != nil {
			return nil, nil, err
		}
		key := cmp.Or(config.GetName(), config.GetDriver())
		caches[key] = cache
	}

	var cache ifacestorage.Cache
	// Determine the active cache.
	// Priority: default > first in map.
	if len(caches) == 0 {
		return nil, nil, errors.New("no cache config found")
	} else {
		defaultKey := cmp.Or(dataConfig.GetCaches().GetDefault(), interfaces.GlobalDefaultKey)
		cache = caches[defaultKey]
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
		dbs:       dbs,
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
