/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package client

import (
	"github.com/google/wire"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/service/transport/grpc"

	helloworld "origadmin/basic-layout/api/v1/gen/go/helloworld"
	secondworld "origadmin/basic-layout/api/v1/gen/go/secondworld"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(NewHelloworldClient, NewSecondworldClient)

// NewHelloworldClient creates a new HelloGreeterAPI client.
func NewHelloworldClient(c config.Config) helloworld.HelloGreeterAPIClient {
	return helloworld.NewHelloGreeterAPIClient(grpc.NewClient(c, "client.helloworld"))
}

// NewSecondworldClient creates a new SecondGreeterAPI client.
func NewSecondworldClient(c config.Config) secondworld.SecondGreeterAPIClient {
	return secondworld.NewSecondGreeterAPIClient(grpc.NewClient(c, "client.secondworld"))
}
