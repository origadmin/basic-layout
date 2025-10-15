/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"context"

	"github.com/google/wire"

	pb "origadmin/basic-layout/internal/configs"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGatewayService)

// GatewayService is a gateway service.
// It acts as a client to other backend services.
type GatewayService struct {
	pb.UnimplementedGatewayServer

	helloworld  pb.GreeterClient
	secondworld pb.SecondClient
}

// NewGatewayService new a gateway service.
func NewGatewayService(helloworld pb.GreeterClient, secondworld pb.SecondClient) *GatewayService {
	return &GatewayService{
		helloworld:  helloworld,
		secondworld: secondworld,
	}
}

// SayHello forwards the request to the helloworld service.
func (s *GatewayService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return s.helloworld.SayHello(ctx, in)
}

// SaySecond forwards the request to the secondworld service.
func (s *GatewayService) SaySecond(ctx context.Context, in *pb.SecondRequest) (*pb.SecondReply, error) {
	return s.secondworld.SaySecond(ctx, in)
}
