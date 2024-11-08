package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Setting = func(*Option)

type Option struct {
	Logger log.Logger
}
