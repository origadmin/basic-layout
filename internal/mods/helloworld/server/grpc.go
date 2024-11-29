/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	servicegrpc "github.com/origadmin/runtime/service/grpc"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/api/v1/services/secondworld"
	"origadmin/basic-layout/internal/configs"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bs *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) *grpc.Server {
	srv := servicegrpc.NewServer(bs.GetService())
	helloworld.RegisterHelloGreeterAPIServer(srv, greeter)
	return srv
}
func RegisterGRPCServer(srv *grpc.Server, greeter secondworld.SecondGreeterAPIServer) {
	secondworld.RegisterSecondGreeterAPIServer(srv, greeter)
}
