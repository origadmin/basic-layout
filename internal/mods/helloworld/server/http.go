/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/origadmin/runtime/service"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1" // Import transportv1

	helloworld "origadmin/basic-layout/api/v1/gen/go/helloworld" // Corrected import path
	"origadmin/basic-layout/internal/configs"
)

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(c *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) (*http.Server, error) {
	var opts = []http.ServerOption{
		service.MiddlewareHTTP(
			recovery.Recovery(),
		),
	}

	if c.Service != nil && c.Service.Servers != nil {
		for _, srvConfig := range c.Service.Servers { // Iterate through servers
			if srvConfig.Protocol == "http" && srvConfig.Http != nil { // Check for HTTP protocol and config
				if srvConfig.Http.Network != "" {
					opts = append(opts, service.NetworkHTTP(srvConfig.Http.Network))
				}
				if srvConfig.Http.Addr != "" {
					opts = append(opts, service.AddressHTTP(srvConfig.Http.Addr))
				}
				if srvConfig.Http.Timeout != nil {
					opts = append(opts, service.TimeoutHTTP(srvConfig.Http.Timeout.AsDuration()))
				}
				// Break after finding the first HTTP server config
				break
			}
		}
	}

	srv := service.NewServerHTTP(opts...)
	helloworld.RegisterHelloGreeterAPIHTTPServer(srv, greeter)
	return srv, nil
}
