/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	userv1 "basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"basic-layout/multiple/multiple_sample/configs"
	"basic-layout/multiple/multiple_sample/internal/mods/user/service"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/config/transport/http/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *configs.ServiceConfig, userService *service.UserService, logger log.Logger) (*http.Server, error) {
	servers := cfg.GetServers()
	var httpCfg *httpv1.Server
	for _, server := range servers {
		if server.GetProtocol() == "http" {
			httpCfg = server.GetHttp()
			break
		}
	}
	if httpCfg == nil {
		return nil, nil
	}

	var opts []http.ServerOption
	if httpCfg.GetAddr() != "" {
		opts = append(opts, http.Address(httpCfg.GetAddr()))
	}
	if httpCfg.GetTimeout() != nil {
		opts = append(opts, http.Timeout(httpCfg.GetTimeout().AsDuration()))
	}
	srv := http.NewServer(opts...)
	userv1.RegisterUserAPIHTTPServer(srv, userService)
	return srv, nil
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *configs.ServiceConfig, userService *service.UserService, logger log.Logger) (*grpc.Server, error) {
	servers := cfg.GetServers()
	var grpcCfg *grpcv1.Server
	for _, server := range servers {
		if server.GetProtocol() == "grpc" {
			grpcCfg = server.GetGrpc()
			break
		}
	}
	if grpcCfg == nil {
		return nil, nil
	}

	var opts []grpc.ServerOption
	if grpcCfg.GetAddr() != "" {
		opts = append(opts, grpc.Address(grpcCfg.GetAddr()))
	}
	if grpcCfg.GetTimeout() != nil {
		opts = append(opts, grpc.Timeout(grpcCfg.GetTimeout().AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	userv1.RegisterUserAPIServer(srv, userService)
	return srv, nil
}
