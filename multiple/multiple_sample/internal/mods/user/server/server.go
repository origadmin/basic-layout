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

	userv1 "basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"basic-layout/multiple/multiple_sample/internal/mods/user/service"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/config/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServers)

func NewServers(cfg *transportv1.Servers, userService *service.UserService, logger log.Logger) ([]transport.Server,
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
func NewHTTPServer(cfg *httpv1.Server, userService *service.UserService, logger log.Logger) (*http.Server, error) {
	if cfg == nil {
		return nil, errors.New("http config is nil")
	}

	var opts []http.ServerOption
	if cfg.GetAddr() != "" {
		opts = append(opts, http.Address(cfg.GetAddr()))
	}
	if cfg.GetTimeout() != nil {
		opts = append(opts, http.Timeout(cfg.GetTimeout().AsDuration()))
	}
	srv := http.NewServer(opts...)
	userv1.RegisterUserAPIHTTPServer(srv, userService)
	return srv, nil
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *grpcv1.Server, userService *service.UserService, logger log.Logger) (*grpc.Server, error) {
	if cfg == nil {
		return nil, errors.New("grpc config is nil")
	}

	var opts []grpc.ServerOption
	if cfg.GetAddr() != "" {
		opts = append(opts, grpc.Address(cfg.GetAddr()))
	}
	if cfg.GetTimeout() != nil {
		opts = append(opts, grpc.Timeout(cfg.GetTimeout().AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	userv1.RegisterUserAPIServer(srv, userService)
	return srv, nil
}
