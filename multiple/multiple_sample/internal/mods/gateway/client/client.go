/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package client

import (
	"context"
	"errors"

	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/api/v1/gen/go/order"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"basic-layout/multiple/multiple_sample/configs"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/service/transport/grpc"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(NewUserClient, NewOrderClient)

// NewUserClient creates a new UserAPI client.
func NewUserClient(ctx context.Context, c *configs.Bootstrap) (user.UserAPIClient, error) {
	var clientConfig *transportv1.Client
	if c.GetClients().GetConfigs() != nil {
		for _, cli := range c.GetClients().GetConfigs() {
			if cli.Name == "client.user" {
				clientConfig = cli
				break
			}
		}
	}

	if clientConfig == nil {
		return nil, errors.New("client config not found: client.user")
	}

	grpcConfig := clientConfig.GetGrpc()
	if grpcConfig == nil {
		return nil, errors.New("grpc client config not found: client.user")
	}

	conn, err := grpc.NewClient(ctx, grpcConfig, &grpc.ClientOptions{}) // Changed nil to &grpc.ClientOptions{}
	if err != nil {
		return nil, err
	}
	return user.NewUserAPIClient(conn), nil
}

// NewOrderClient creates a new OrderAPI client.
func NewOrderClient(ctx context.Context, c *configs.Bootstrap) (order.OrderAPIClient, error) {
	var clientConfig *transportv1.Client
	if c.GetClients().GetConfigs() != nil {
		for _, cli := range c.GetClients().GetConfigs() {
			if cli.Name == "client.order" {
				clientConfig = cli
				break
			}
		}
	}

	if clientConfig == nil {
		return nil, errors.New("client config not found: client.order")
	}

	grpcConfig := clientConfig.GetGrpc()
	if grpcConfig == nil {
		return nil, errors.New("grpc client config not found: client.order")
	}

	conn, err := grpc.NewClient(ctx, grpcConfig, &grpc.ClientOptions{}) // Changed nil to &grpc.ClientOptions{}
	if err != nil {
		return nil, err
	}
	return order.NewOrderAPIClient(conn), nil
}
