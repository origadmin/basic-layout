package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/basic-layout/api/v1/services/helloworld"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/conf"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http == nil {
		c.Http = new(conf.Server_HTTP)
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	if c.Middleware == nil {
		c.Middleware = new(conf.Middleware)
	}
	//var middlewares []middleware.Middleware
	//middlewares = append(middlewares, validate.Validator())
	//if c.Middleware.Logger.Enabled {
	//	m, err := logger.Middleware(logger.Config{
	//		Name: c.Middleware.Logger.Name,
	//	}, nil)
	//	if err != nil {
	//		l.Log(log.LevelError, "init logger middleware error", err)
	//	} else {
	//		middlewares = append(middlewares, m)
	//	}
	//
	//}
	//if c.Middleware.Traces.Enabled {
	//	m, err := traces.Middleware(traces.Config{
	//		Name: c.Middleware.Traces.Name,
	//	})
	//	if err != nil {
	//		l.Log(log.LevelError, "init traces middleware error", err)
	//	} else {
	//		middlewares = append(middlewares, m)
	//	}
	//}
	//
	//if c.Middleware.Metrics.Enabled {
	//	m, err := metrics.Middleware(metrics.Config{
	//		Name: c.Middleware.Traces.Name,
	//		Side: metrics.SideServer,
	//	})
	//	if err != nil {
	//		l.Log(log.LevelError, "init metrics middleware error", err)
	//	} else {
	//		middlewares = append(middlewares, m)
	//	}
	//}

	srv := http.NewServer(opts...)
	helloworld.RegisterGreeterServiceHTTPServer(srv, greeter)
	return srv
}
