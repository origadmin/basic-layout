/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/origadmin/runtime/service"

	"origadmin/basic-layout/api/v1/services/secondworld"
	"origadmin/basic-layout/internal/configs"
)

// NewGRPCServer creates a new gRPC server.
func NewGRPCServer(c *configs.Bootstrap, greeter secondworld.SecondGreeterAPIServer, l log.Logger) (*grpc.Server, error) {
	var opts = []grpc.ServerOption{
		service.MiddlewareGRPC(
			recovery.Recovery(),
		),
	}

	if c.Service != nil && c.Service.Server != nil && c.Service.Server.Grpc != nil {
		if c.Service.Server.Grpc.Network != "" {
			opts = append(opts, service.NetworkGRPC(c.Service.Server.Grpc.Network))
		}
		if c.Service.Server.Grpc.Addr != "" {
			opts = append(opts, service.AddressGRPC(c.Service.Server.Grpc.Addr))
		}
		if c.Service.Server.Grpc.Timeout != nil {
			opts = append(opts, service.TimeoutGRPC(c.Service.Server.Grpc.Timeout.AsDuration()))
		}
	}

	srv := service.NewServerGRPC(opts...)
	secondworld.RegisterSecondGreeterAPIServer(srv, greeter)
	return srv, nil
}
