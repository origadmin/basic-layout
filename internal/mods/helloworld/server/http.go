/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	rtservice "github.com/origadmin/runtime/service"
	"origadmin/basic-layout/api/v1/gen/go/configs"
	"origadmin/basic-layout/api/v1/gen/go/helloworld" // Corrected import path
)

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(c *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) (*rtservice.HTTPServer,
	error) {
	var opts = []rtservice.HTTPServerOption{
		rtservice.MiddlewareHTTP(
			recovery.Recovery(),
		),
	}

	if service := c.GetServer().GetService(); service != nil {
		for _, srvConfig := range service.Servers { // Iterate through servers
			if srvConfig.Protocol == "http" && srvConfig.Http != nil {
				if srvConfig.Http.Network != "" {
					opts = append(opts, rtservice.NetworkHTTP(srvConfig.Http.Network))
				}
				if srvConfig.Http.Addr != "" {
					opts = append(opts, rtservice.AddressHTTP(srvConfig.Http.Addr))
				}
				if srvConfig.Http.Timeout != nil {
					opts = append(opts, rtservice.TimeoutHTTP(srvConfig.Http.Timeout.AsDuration()))
				}
				// Break after finding the first HTTP server config
				break
			}
		}
	}

	srv := rtservice.NewServerHTTP(opts...)
	helloworld.RegisterHelloGreeterAPIHTTPServer(srv, greeter)
	return srv, nil
}
