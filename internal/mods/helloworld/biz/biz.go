/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterBiz, NewGreeterClient)
