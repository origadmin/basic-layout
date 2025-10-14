/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/service"
	"origadmin/basic-layout/api/v1/gen/go/secondworld"
	"origadmin/basic-layout/internal/configs"
)

// NewGRPCServer creates a new gRPC server.
func NewGRPCServer(c *configs.Bootstrap, greeter secondworld.SecondGreeterAPIServer, l log.Logger) (*grpc.Server, error) {
	var opts = []grpc.ServerOption{
		service.MiddlewareGRPC(
			recovery.Recovery(),
		),
	}

	if c.Service != nil && c.Service.Servers != nil {
		for _, srvConfig := range c.Service.Servers { // Iterate through servers
			if srvConfig.Protocol == "grpc" && srvConfig.Grpc != nil { // Check for gRPC protocol and config
				if srvConfig.Grpc.Network != "" {
					opts = append(opts, service.NetworkGRPC(srvConfig.Grpc.Network))
				}
				if srvConfig.Grpc.Addr != "" {
					opts = append(opts, service.AddressGRPC(srvConfig.Grpc.Addr))
				}
				if srvConfig.Grpc.Timeout != nil {
					opts = append(opts, service.TimeoutGRPC(srvConfig.Grpc.Timeout.AsDuration()))
				}
				// Break after finding the first gRPC server config
				break
			}
		}
	}

	srv := service.NewServerGRPC(opts...)
	secondworld.RegisterSecondGreeterAPIServer(srv, greeter)
	return srv, nil
}
