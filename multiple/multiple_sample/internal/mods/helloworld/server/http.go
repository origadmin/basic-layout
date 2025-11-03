/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"origadmin/basic-layout/api/v1/gen/go/helloworld" // Corrected import path
	"origadmin/basic-layout/internal/configs"

	rtservice "github.com/origadmin/runtime/service"
)

// NewHTTPServer creates a new HTTP server for the helloworld service.
// It initializes the server with the provided configuration and sets up the necessary middleware.
func NewHTTPServer(bootstrap *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) (*rtservice.HTTPServer, error) {
	logger := log.NewHelper(log.With(l, "module", "helloworld/http"))
	logger.Info("Initializing HTTP server for helloworld service")
	var opts = []rtservice.HTTPServerOption{
		rtservice.MiddlewareHTTP(
			recovery.Recovery(),
		),
	}

	if service := bootstrap.GetServer().GetService(); service != nil {
		logger.Debug("Processing server configurations", "total_servers", len(service.Servers))
		for _, srvConfig := range service.Servers {
			logger.Debug("Processing server configuration", "protocol", srvConfig.Protocol)
			if srvConfig.Protocol == "http" && srvConfig.Http != nil {
				if srvConfig.Http.Network != "" {
					opts = append(opts, rtservice.NetworkHTTP(srvConfig.Http.Network))
					logger.Debug("HTTP server network set", "network", srvConfig.Http.Network)
				}
				if srvConfig.Http.Addr != "" {
					opts = append(opts, rtservice.AddressHTTP(srvConfig.Http.Addr))
					logger.Debug("HTTP server address set", "addr", srvConfig.Http.Addr)
				}
				if srvConfig.Http.Timeout != nil {
					opts = append(opts, rtservice.TimeoutHTTP(srvConfig.Http.Timeout.AsDuration()))
					logger.Debug("HTTP server timeout set", "timeout", srvConfig.Http.Timeout.AsDuration())
				}
				// Break after finding the first HTTP server config
				break
			}
		}
	}

	srv := rtservice.NewServerHTTP(opts...)
	helloworld.RegisterHelloGreeterAPIHTTPServer(srv, greeter)
	srv.WalkRoute(func(route rtservice.RouteInfoHTTP) error {
		logger.Debugf("Registered HTTP route: %s %s", route.Method, route.Path)
		return nil
	})
	logger.Infof("HTTP server initialized successfully, service: %s, endpoints: /v1/helloworld/*",
		bootstrap.GetServer().GetService().GetName())
	return srv, nil
}
