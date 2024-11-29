/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Setting = func(*Option)

type Option struct {
	Logger log.Logger
}

func WithLogger(logger log.Logger) Setting {
	return func(o *Option) {
		o.Logger = logger
	}
}
