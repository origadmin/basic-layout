package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	//"origadmin/basic-layout/toolkits/middlewares/logger"
	//"origadmin/basic-layout/toolkits/middlewares/metrics"
	//"origadmin/basic-layout/toolkits/middlewares/security"

	"github.com/origadmin/toolkits/runtime/config"
	"origadmin/basic-layout/internal/configs"
)

func LoadMiddlewares(name string, bootstrap *configs.Bootstrap, l log.Logger) ([]middleware.Middleware, error) {
	var middlewares []middleware.Middleware
	middlewares = append(middlewares, validate.Validator())
	mc := bootstrap.Middleware
	if mc == nil {
		return nil, nil
	}

	//if mc.Logger != nil && mc.Logger.Enabled {
	//	m, err := logger.Middleware(&logger.LoggerConfig{
	//		Name: mc.Logger.Name,
	//	}, l)
	//	if err != nil {
	//		return nil, err
	//	}
	//	middlewares = append(middlewares, m)
	//}
	//if v := mc.GetSecurity(); v != nil {
	//	m, err := security.Middleware(&security.SecurityConfig{
	//		AllowedMethodPaths: nil,
	//		Authorization:      nil,
	//		Casbin:             nil,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	middlewares = append(middlewares, m)
	//}
	//
	//if mc.Metrics != nil && mc.Metrics.Enabled {
	//	m, err := metrics.Middleware(metrics.SideServer, &metrics.MetricConfig{
	//		Name: mc.Metrics.Name,
	//	}, l)
	//	if err != nil {
	//		return nil, err
	//	}
	//	middlewares = append(middlewares, m)
	//}

	return middlewares, nil
}

// LoadGlobalMiddlewares Loading global middleware
func LoadGlobalMiddlewares(name string, conf *configs.Bootstrap, l log.Logger) ([]middleware.Middleware, error) {
	if conf.Middleware == nil {
		conf.Middleware = new(config.Middleware)
	}
	//if !conf.Middleware.RegisterAsGlobal {
	//	return nil, nil
	//}
	return LoadMiddlewares(name, conf, l)
}
