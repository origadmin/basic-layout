package bootstrap

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/validate"

	"origadmin/basic-layout/internal/config"
	"origadmin/basic-layout/toolkits/middlewares/logger"
	"origadmin/basic-layout/toolkits/middlewares/metrics"
	"origadmin/basic-layout/toolkits/middlewares/traces"
)

func LoadMiddlewares(name string, conf *config.Bootstrap) ([]middleware.Middleware, error) {
	var middlewares []middleware.Middleware
	middlewares = append(middlewares, validate.Validator())
	mc := conf.Middlewares
	if mc.Logger.Enabled {
		m, err := logger.Middleware(logger.Config{
			Name: mc.Logger.Name,
		}, nil)
		if err != nil {
			return nil, err
		}
		middlewares = append(middlewares, m)
	}
	if mc.Traces.Enabled {
		m, err := traces.Middleware(traces.Config{
			Name: mc.Traces.Name,
		})
		if err != nil {
			return nil, err
		}
		middlewares = append(middlewares, m)
	}

	if mc.Metrics.Enabled {
		m, err := metrics.Middleware(metrics.Config{
			Name: mc.Traces.Name,
			Side: metrics.SideServer,
		})
		if err != nil {
			return nil, err
		}
		middlewares = append(middlewares, m)
	}

	return middlewares, nil
}
