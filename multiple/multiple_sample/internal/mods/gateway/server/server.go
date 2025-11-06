/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	gatewayAPI "basic-layout/multiple/multiple_sample/api/v1/gen/go/gateway"
	"basic-layout/multiple/multiple_sample/internal/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/gateway/service"

	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1" // Added import
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer, NewHTTPServer)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(logger log.Logger, c *configs.Bootstrap, gw *service.GatewayService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			// Add any HTTP middleware here if needed
		),
	}

	var httpServerConfig *transportv1.HttpServerConfig
	if c.Server != nil && c.Server.Service != nil {
		for _, s := range c.Server.Service.Servers {
			if s.GetProtocol() == "http" && s.GetHttp() != nil {
				httpServerConfig = s.GetHttp()
				break
			}
		}
	}

	if httpServerConfig != nil {
		if httpServerConfig.Addr != "" {
			opts = append(opts, http.Address(httpServerConfig.Addr))
		}
		if httpServerConfig.Timeout != nil {
			opts = append(opts, http.Timeout(httpServerConfig.Timeout.AsDuration()))
		}
	}
	srv := http.NewServer(opts...)
	gatewayAPI.RegisterGatewayAPIHTTPServer(srv, gw)
	return srv
}

// NewServer new a gRPC server.
func NewServer(logger log.Logger, c *configs.Bootstrap, gw *service.GatewayService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			// Add any gRPC middleware here if needed
		),
	}

	var grpcServerConfig *transportv1.GrpcServerConfig
	if c.Server != nil && c.Server.Service != nil {
		for _, s := range c.Server.Service.Servers {
			if s.GetProtocol() == "grpc" && s.GetGrpc() != nil {
				grpcServerConfig = s.GetGrpc()
				break
			}
		}
	}

	if grpcServerConfig != nil {
		if grpcServerConfig.Addr != "" {
			opts = append(opts, grpc.Address(grpcServerConfig.Addr))
		}
		if grpcServerConfig.Timeout != nil {
			opts = append(opts, grpc.Timeout(grpcServerConfig.Timeout.AsDuration()))
		}
	}
	srv := grpc.NewServer(opts...)
	gatewayAPI.RegisterGatewayAPIServer(srv, gw)
	return srv
}
