/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package client

import (
	"github.com/google/wire"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/service/transport/grpc"

	pb "origadmin/basic-layout/internal/configs"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(NewHelloworldClient, NewSecondworldClient)

// NewHelloworldClient creates a new Greeter client.
func NewHelloworldClient(c config.Config) pb.GreeterClient {
	return pb.NewGreeterClient(grpc.NewClient(c, "client.helloworld"))
}

// NewSecondworldClient creates a new Second client.
func NewSecondworldClient(c config.Config) pb.SecondClient {
	return pb.NewSecondClient(grpc.NewClient(c, "client.secondworld"))
}
