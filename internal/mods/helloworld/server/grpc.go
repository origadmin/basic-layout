/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	rtservice "github.com/origadmin/runtime/service"
	"origadmin/basic-layout/api/v1/gen/go/configs"
	"origadmin/basic-layout/api/v1/gen/go/helloworld"
)

// NewGRPCServer creates a new gRPC server.
func NewGRPCServer(c *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) (*rtservice.GRPCServer,
	error) {
	var opts = []rtservice.GRPCServerOption{
		rtservice.MiddlewareGRPC(
			recovery.Recovery(),
		),
	}

	if service := c.GetServer().GetService(); service != nil && service.Servers != nil {
		for _, srvConfig := range service.Servers { // Iterate through servers
			if srvConfig.Protocol == "grpc" && srvConfig.Grpc != nil { // Check for gRPC protocol and config
				if srvConfig.Grpc.Network != "" {
					opts = append(opts, rtservice.NetworkGRPC(srvConfig.Grpc.Network))
				}
				if srvConfig.Grpc.Addr != "" {
					opts = append(opts, rtservice.AddressGRPC(srvConfig.Grpc.Addr))
				}
				if srvConfig.Grpc.Timeout != nil {
					opts = append(opts, rtservice.TimeoutGRPC(srvConfig.Grpc.Timeout.AsDuration()))
				}
				// Break after finding the first gRPC server config
				break
			}
		}
	}

	srv := rtservice.NewServerGRPC(opts...)
	helloworld.RegisterHelloGreeterAPIServer(srv, greeter)
	return srv, nil
}
