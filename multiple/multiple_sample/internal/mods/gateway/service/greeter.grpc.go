/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"context"

	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/api/v1/gen/go/gateway"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/order"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGatewayService)

// GatewayService implements the GatewayAPI service.
type GatewayService struct {
	gateway.UnimplementedGatewayAPIServer

	userClient  user.UserAPIClient
	orderClient order.OrderAPIClient
}

// NewGatewayService creates a new gateway service.
func NewGatewayService(uc user.UserAPIClient, oc order.OrderAPIClient) *GatewayService {
	return &GatewayService{
		userClient:  uc,
		orderClient: oc,
	}
}

// SayUser forwards the request to the user service.
func (s *GatewayService) SayUser(ctx context.Context, in *user.SayUserRequest) (*user.SayUserResponse, error) {
	return s.userClient.SayUser(ctx, in)
}

// SayOrder forwards the request to the order service.
func (s *GatewayService) SayOrder(ctx context.Context, in *order.SayOrderRequest) (*order.SayOrderResponse, error) {
	// Note: The request for order is also SayOrderRequest, as per the proto definition.
	return s.orderClient.SayOrder(ctx, in)
}
