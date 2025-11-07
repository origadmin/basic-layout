/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	gatewayAPI "basic-layout/multiple/multiple_sample/api/v1/gen/go/gateway"
	"basic-layout/multiple/multiple_sample/internal/features/gateway/service"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/config/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1" // Added import
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServers)

func NewServers(cfg *transportv1.Servers, userService *service.GatewayService, logger log.Logger) ([]transport.Server,
	error) {
	servers := cfg.GetConfigs()
	var transportServers []transport.Server
	for _, server := range servers {
		switch server.GetProtocol() {
		case "http":
			srv, err := NewHTTPServer(server.GetHttp(), userService, logger)
			if err != nil {
				return nil, err
			}
			transportServers = append(transportServers, srv)
		case "grpc":
			srv, err := NewGRPCServer(server.GetGrpc(), userService, logger)
			if err != nil {
				return nil, err
			}
			transportServers = append(transportServers, srv)
		default:
			return nil, errors.New("protocol is not supported")
		}
	}
	return transportServers, nil
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(httpServerConfig *httpv1.Server, gw *service.GatewayService, logger log.Logger) (*http.Server,
	error) {
	var opts = []http.ServerOption{
		http.Middleware(
			// Add any HTTP middleware here if needed
		),
	}
	if httpServerConfig == nil {
		return nil, errors.New("http server config is nil")
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
	gatewayAPI.RegisterProxyGatewayAPIHTTPServer(srv, gw)
	return srv, nil
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(grpcServerConfig *grpcv1.Server, gw *service.GatewayService, logger log.Logger) (*grpc.Server,
	error) {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			// Add any gRPC middleware here if needed
		),
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
	gatewayAPI.RegisterProxyGatewayAPIServer(srv, gw)
	return srv, nil
}
