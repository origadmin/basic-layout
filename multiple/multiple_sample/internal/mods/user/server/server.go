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
	"basic-layout/multiple/multiple_sample/internal/mods/user/service"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *transportv1.HttpServer, userService *service.UserService, logger log.Logger) (*http.Server, error) {
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
func NewGRPCServer(cfg *transportv1.GrpcServer, userService *service.UserService, logger log.Logger) (*grpc.Server, error) {
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
