/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/origadmin/runtime/service"

	"origadmin/basic-layout/api/v1/services/secondworld"
	"origadmin/basic-layout/internal/configs"
)

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(c *configs.Bootstrap, greeter secondworld.SecondGreeterAPIServer, l log.Logger) (*http.Server, error) {
	var opts = []http.ServerOption{
		service.MiddlewareHTTP(
			recovery.Recovery(),
		),
	}

	if c.Service != nil && c.Service.Server != nil && c.Service.Server.Http != nil {
		if c.Service.Server.Http.Network != "" {
			opts = append(opts, service.NetworkHTTP(c.Service.Server.Http.Network))
		}
		if c.Service.Server.Http.Addr != "" {
			opts = append(opts, service.AddressHTTP(c.Service.Server.Http.Addr))
		}
		if c.Service.Server.Http.Timeout != nil {
			opts = append(opts, service.TimeoutHTTP(c.Service.Server.Http.Timeout.AsDuration()))
		}
	}

	srv := service.NewServerHTTP(opts...)
	secondworld.RegisterSecondGreeterAPIHTTPServer(srv, greeter)
	return srv, nil
}
