/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package dal

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/internal/configs"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewGreeterDal)

// Database .
type Database struct {
	// TODO wrapped database client
}

// NewDB .
func NewDB(bootstrap *configs.Bootstrap, logger log.Logger) (*Database, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Database{}, cleanup, nil
}
