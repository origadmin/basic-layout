/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/api/v1/gen/go/gateway"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/order"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"github.com/origadmin/runtime/context"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGatewayService)

// GatewayService implements the GatewayAPI service.
type GatewayService struct {
	gateway.UnimplementedProxyGatewayAPIServer

	userClient  user.UserAPIClient
	orderClient order.OrderAPIClient
}

func (g GatewayService) GetOrder(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	return g.orderClient.GetOrder(ctx, request)
}

func (g GatewayService) GetUser(ctx context.Context, request *user.GetUserRequest) (*user.GetUserResponse, error) {
	return g.userClient.GetUser(ctx, request)
}

// NewGatewayService creates a new gateway service.
func NewGatewayService(uc user.UserAPIClient, oc order.OrderAPIClient) *GatewayService {
	return &GatewayService{
		userClient:  uc,
		orderClient: oc,
	}
}
