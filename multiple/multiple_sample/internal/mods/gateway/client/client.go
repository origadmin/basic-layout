/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package client

import (
	"context"
	"errors"

	"github.com/google/wire"

	"github.com/origadmin/runtime/service/transport/grpc" // Renamed for clarity

	"origadmin/basic-layout/api/v1/gen/go/helloworld"
	"origadmin/basic-layout/api/v1/gen/go/secondworld"
	"origadmin/basic-layout/internal/configs"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(NewHelloworldClient, NewSecondworldClient)

// NewHelloworldClient creates a new HelloGreeterAPI client.
func NewHelloworldClient(ctx context.Context, c *configs.Bootstrap) (helloworld.HelloGreeterAPIClient, error) {
	var clientConfig *transportv1.Client
	if c.Server != nil && c.Server.Service != nil {
		for _, cli := range c.Server.GetService().Clients {
			if cli.Name == "client.helloworld" {
				clientConfig = cli
				break
			}
		}
	}

	if clientConfig == nil {
		return nil, errors.New("client config not found: client.helloworld")
	}

	grpcConfig := clientConfig.GetGrpc()
	if grpcConfig == nil {
		return nil, errors.New("grpc client config not found: client.helloworld")
	}

	conn, err := grpc.NewClient(ctx, grpcConfig, &grpc.ClientOptions{}) // Changed nil to &grpc.ClientOptions{}
	if err != nil {
		return nil, err
	}
	return helloworld.NewHelloGreeterAPIClient(conn), nil
}

// NewSecondworldClient creates a new SecondGreeterAPI client.
func NewSecondworldClient(ctx context.Context, c *configs.Bootstrap) (secondworld.SecondGreeterAPIClient, error) {
	var clientConfig *transportv1.Client
	if c.Server != nil && c.Server.Service != nil {
		for _, cli := range c.Server.Service.Clients {
			if cli.Name == "client.secondworld" {
				clientConfig = cli
				break
			}
		}
	}

	if clientConfig == nil {
		return nil, errors.New("client config not found: client.secondworld")
	}

	grpcConfig := clientConfig.GetGrpc()
	if grpcConfig == nil {
		return nil, errors.New("grpc client config not found: client.secondworld")
	}

	conn, err := grpc.NewClient(ctx, grpcConfig, &grpc.ClientOptions{}) // Changed nil to &grpc.ClientOptions{}
	if err != nil {
		return nil, err
	}
	return secondworld.NewSecondGreeterAPIClient(conn), nil
}
