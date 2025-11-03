/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"origadmin/basic-layout/api/v1/gen/go/secondworld" // Corrected import path
	"origadmin/basic-layout/internal/configs"

	rtservice "github.com/origadmin/runtime/service"
)

// NewHTTPServer creates a new HTTP server for the secondworld service.
// It initializes the server with the provided configuration and sets up the necessary middleware.
func NewHTTPServer(bootstrap *configs.Bootstrap, greeter secondworld.SecondGreeterAPIServer, l log.Logger) (*rtservice.HTTPServer, error) {
	logger := log.NewHelper(log.With(l, "module", "secondworld/http"))
	logger.Info("Initializing HTTP server for secondworld service")
	var opts = []rtservice.HTTPServerOption{
		rtservice.MiddlewareHTTP(
			recovery.Recovery(),
		),
	}

	if service := bootstrap.GetServer().GetService(); service != nil {
		logger.Debugf("Processing server configurations, total_servers: %d", len(service.Servers))

		for _, srvConfig := range service.Servers {
			logger.Debugf("Processing server configuration, protocol: %s", srvConfig.Protocol)

			if srvConfig.Protocol == "http" && srvConfig.Http != nil {
				if srvConfig.Http.Network != "" {
					opts = append(opts, rtservice.NetworkHTTP(srvConfig.Http.Network))
					logger.Debugf("Setting HTTP server network to %s", srvConfig.Http.Network)
				}
				if srvConfig.Http.Addr != "" {
					opts = append(opts, rtservice.AddressHTTP(srvConfig.Http.Addr))
					logger.Debugf("Setting HTTP server address to %s", srvConfig.Http.Addr)
				}
				if srvConfig.Http.Timeout != nil {
					opts = append(opts, rtservice.TimeoutHTTP(srvConfig.Http.Timeout.AsDuration()))
					logger.Debugf("Setting HTTP server timeout to %s", srvConfig.Http.Timeout.AsDuration())
				}
				// Break after finding the first HTTP server config
				break
			}
		}
	}

	srv := rtservice.NewServerHTTP(opts...)
	secondworld.RegisterSecondGreeterAPIHTTPServer(srv, greeter)
	srv.WalkRoute(func(info rtservice.RouteInfoHTTP) error {
		logger.Debugf("Registered HTTP route: %s %s", info.Method, info.Path)
		return nil
	})
	logger.Infof("HTTP server initialized successfully, service: %s, endpoints: /v1/secondworld/*",
		bootstrap.GetServer().GetService().GetName())
	return srv, nil
}
