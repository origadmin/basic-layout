/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"

	//"origadmin/basic-layout/toolkits/middlewares/logger"
	//"origadmin/basic-layout/toolkits/middlewares/metrics"
	//"origadmin/basic-layout/toolkits/middlewares/security"

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
	//	m, err := logger.BuildMiddleware(&logger.LoggerConfig{
	//		Name: mc.Logger.Name,
	//	}, l)
	//	if err != nil {
	//		return nil, err
	//	}
	//	middlewares = append(middlewares, m)
	//}
	//if v := mc.GetSecurity(); v != nil {
	//	m, err := security.BuildMiddleware(&security.SecurityConfig{
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
	//	m, err := metrics.BuildMiddleware(metrics.SideServer, &metrics.MetricConfig{
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
		conf.Middleware = new(configv1.Middleware)
	}
	//if !conf.BuildMiddleware.RegisterAsGlobal {
	//	return nil, nil
	//}
	return LoadMiddlewares(name, conf, l)
}
