/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/secondworld" // Import secondworld package

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService,                                                        // Provides *GreeterService
	wire.Bind(new(secondworld.SecondGreeterAPIServer), new(*GreeterService)), // Binds *GreeterService to SecondGreeterAPIServer interface
)
