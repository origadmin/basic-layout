/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	gatewayAPI "origadmin/basic-layout/api/v1/gen/go/gateway"
	configs "origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/internal/mods/gateway/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(logger log.Logger, c *configs.Bootstrap, gw *service.GatewayService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			// Add any HTTP middleware here if needed
		),
	}
	if c.Server != nil && c.Server.Http != nil {
		if c.Server.Http.Addr != "" {
			opts = append(opts, http.Address(c.Server.Http.Addr))
		}
		if c.Server.Http.Timeout != nil {
			opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
		}
	}
	srv := http.NewServer(opts...)
	gatewayAPI.RegisterGatewayAPIHTTPServer(srv, gw)
	return srv
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(logger log.Logger, c *configs.Bootstrap, gw *service.GatewayService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			// Add any gRPC middleware here if needed
		),
	}
	if c.Server != nil && c.Server.Grpc != nil {
		if c.Server.Grpc.Addr != "" {
			opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
		}
		if c.Server.Grpc.Timeout != nil {
			opts = append(opts, grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()))
		}
	}
	srv := grpc.NewServer(opts...)
	gatewayAPI.RegisterGatewayAPIServer(srv, gw)
	return srv
}
