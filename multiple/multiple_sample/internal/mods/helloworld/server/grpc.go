/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"basic-layout/multiple/multiple_sample/api/v1/gen/go/helloworld"
	"basic-layout/multiple/multiple_sample/internal/configs"

	rtservice "github.com/origadmin/runtime/service"
)

// NewGRPCServer creates a new gRPC server for the helloworld service.
// It initializes the gRPC server with the provided configuration and sets up the necessary middleware.
func NewGRPCServer(bootstrap *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) (*rtservice.GRPCServer, error) {
	logger := log.NewHelper(log.With(l, "module", "helloworld/grpc"))
	logger.Info("Initializing gRPC server for helloworld servers")
	var opts = []rtservice.GRPCServerOption{
		rtservice.MiddlewareGRPC(
			recovery.Recovery(),
		),
	}

	if servers := bootstrap.GetService().GetServers(); servers != nil {
		logger.Debugf("Processing gRPC server configurations, total_servers: %d", len(servers))
		for _, srvConfig := range servers {
			logger.Debugf("Processing server configuration, protocol: %s", srvConfig.Protocol)

			if srvConfig.Protocol == "grpc" && srvConfig.Grpc != nil {
				if srvConfig.Grpc.Network != "" {
					opts = append(opts, rtservice.NetworkGRPC(srvConfig.Grpc.Network))
					logger.Debugf("Setting gRPC server network to, network: %s", srvConfig.Grpc.Network)
				}
				if srvConfig.Grpc.Addr != "" {
					opts = append(opts, rtservice.AddressGRPC(srvConfig.Grpc.Addr))
					logger.Debugf("Setting gRPC server address to, address: %s", srvConfig.Grpc.Addr)
				}
				if srvConfig.Grpc.Timeout != nil {
					opts = append(opts, rtservice.TimeoutGRPC(srvConfig.Grpc.Timeout.AsDuration()))
					logger.Debugf("Setting gRPC server timeout to, timeout: %s", srvConfig.Grpc.Timeout.AsDuration())
				}
				// Break after finding the first gRPC server config
				break
			}
			logger.Infof("gRPC server initialized successfully, servers: %s, methods: %v",
				srvConfig.GetName(), []string{
					"/helloworld.v1.HelloGreeter/SayHello",
				})
		}
	}

	srv := rtservice.NewServerGRPC(opts...)
	helloworld.RegisterHelloGreeterAPIServer(srv, greeter)

	return srv, nil
}
