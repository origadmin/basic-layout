package data

import (
	"basic-layout/multiple/multiple_sample/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/order/data/ent"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/origadmin/runtime"
	runtime_data_database "github.com/origadmin/runtime/data/database/v1"
)

// ProviderSet is order module's data providers.
var ProviderSet = wire.NewSet(NewData, NewOrderRepo)

// Data .
type Data struct {
	db *ent.Client
}

// NewData .
func NewData(rt *runtime.Runtime, conf *configs.Bootstrap, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data/order"))

	// Get database config for this module
	dbConf, err := conf.GetData().GetDatabases().GetConfig("order_db")
	if err != nil {
		return nil, nil, err
	}

	// Create database client
	db, err := runtime_data_database.NewDatabase(rt, dbConf)
	if err != nil {
		return nil, nil, err
	}
	client := ent.NewClient(ent.Driver(db))

	cleanup := func() {
		helper.Info("closing the order module data resources")
		if err := client.Close(); err != nil {
			helper.Errorf("failed to close order db client: %v", err)
		}
	}
	return &Data{db: client}, cleanup, nil
}
