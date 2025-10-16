/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	rtservice "github.com/origadmin/runtime/service"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/api/v1/gen/go/secondworld"
)

// NewGRPCServer creates a new gRPC server for the secondworld service.
// It initializes the gRPC server with the provided configuration and sets up the necessary middleware.
func NewGRPCServer(bootstrap *configs.Bootstrap, greeter secondworld.SecondGreeterAPIServer, l log.Logger) (*rtservice.ServerGRPC, error) {
	logger := log.NewHelper(log.With(l, "module", "secondworld/grpc"))
	logger.Info("Initializing gRPC server for secondworld service")
	var opts = []rtservice.GRPCServerOption{
		rtservice.MiddlewareGRPC(
			recovery.Recovery(),
		),
	}

	if service := bootstrap.GetServer().GetService(); service != nil {
		logger.Debugf("Processing gRPC server configurations, total_servers: %d", len(service.Servers))
		for _, srvConfig := range service.Servers {
			logger.Debugf("Processing server configuration, protocol: %s", srvConfig.Protocol)

			if srvConfig.Protocol == "grpc" && srvConfig.Grpc != nil {
				if srvConfig.Grpc.Network != "" {
					opts = append(opts, rtservice.NetworkGRPC(srvConfig.Grpc.Network))
					logger.Debugf("Setting gRPC server network to %s", srvConfig.Grpc.Network)
				}
				if srvConfig.Grpc.Addr != "" {
					opts = append(opts, rtservice.AddressGRPC(srvConfig.Grpc.Addr))
					logger.Debugf("Setting gRPC server address to %s", srvConfig.Grpc.Addr)
				}
				if srvConfig.Grpc.Timeout != nil {
					opts = append(opts, rtservice.TimeoutGRPC(srvConfig.Grpc.Timeout.AsDuration()))
					logger.Debugf("Setting gRPC server timeout to %s", srvConfig.Grpc.Timeout.AsDuration())
				}
				// Break after finding the first gRPC server config
				break
			}
		}
	}

	srv := rtservice.NewServerGRPC(opts...)
	secondworld.RegisterSecondGreeterAPIServer(srv, greeter)

	logger.Infof("gRPC server initialized successfully, service: %s, methods: %v",
		bootstrap.GetServer().GetService().GetName(), []string{
			"/secondworld.v1.SecondGreeter/SayHello",
		})
	return srv, nil
}
