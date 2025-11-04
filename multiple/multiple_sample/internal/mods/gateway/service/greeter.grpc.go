/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"context"

	"github.com/google/wire"

	"basic-layout/multiple/multiple_sample/api/v1/gen/go/gateway"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/helloworld"
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/secondworld"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGatewayService)

// GatewayService implements the GatewayAPI service.
type GatewayService struct {
	gateway.UnimplementedGatewayAPIServer

	helloClient  helloworld.HelloGreeterAPIClient
	secondClient secondworld.SecondGreeterAPIClient
}

// NewGatewayService creates a new gateway service.
func NewGatewayService(hc helloworld.HelloGreeterAPIClient, sc secondworld.SecondGreeterAPIClient) *GatewayService {
	return &GatewayService{
		helloClient:  hc,
		secondClient: sc,
	}
}

// SayHello forwards the request to the helloworld service.
func (s *GatewayService) SayHello(ctx context.Context, in *helloworld.SayHelloRequest) (*helloworld.SayHelloResponse, error) {
	return s.helloClient.SayHello(ctx, in)
}

// SaySecond forwards the request to the secondworld service.
func (s *GatewayService) SaySecond(ctx context.Context, in *secondworld.SayHelloRequest) (*secondworld.SayHelloResponse, error) {
	// Note: The request for secondworld is also SayHelloRequest, as per the proto definition.
	return s.secondClient.SayHello(ctx, in)
}
