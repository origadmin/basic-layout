/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"github.com/google/wire"
	helloworld "origadmin/basic-layout/api/v1/gen/go/helloworld" // Import helloworld package
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService, // Provides *GreeterService
	wire.Bind(new(helloworld.HelloGreeterAPIServer), new(*GreeterService)), // Binds *GreeterService to HelloGreeterAPIServer interface
)
