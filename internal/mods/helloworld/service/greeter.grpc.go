/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"context"

	helloworld "origadmin/basic-layout/api/v1/gen/go/helloworld" // Corrected import path
)

// GreeterService is a greeter service.
type GreeterService struct {
	helloworld.HelloGreeterAPIServer

	client helloworld.HelloGreeterAPIClient
}

func (g GreeterService) SayHello(ctx context.Context, request *helloworld.SayHelloRequest) (*helloworld.SayHelloResponse, error) {
	return g.client.SayHello(ctx, request)
}

func (g GreeterService) PostHello(ctx context.Context, request *helloworld.PostHelloRequest) (*helloworld.PostHelloResponse, error) {
	return g.client.PostHello(ctx, request)
}

func (g GreeterService) CreateGreeter(ctx context.Context, request *helloworld.CreateGreeterRequest) (*helloworld.CreateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) UpdateGreeter(ctx context.Context, request *helloworld.UpdateGreeterRequest) (*helloworld.UpdateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) DeleteGreeter(ctx context.Context, request *helloworld.DeleteGreeterRequest) (*helloworld.DeleteGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) GetGreeter(ctx context.Context, request *helloworld.GetGreeterRequest) (*helloworld.GetGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) ListGreeter(ctx context.Context, request *helloworld.ListGreeterRequest) (*helloworld.ListGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterService new a greeter service.
func NewGreeterService(client helloworld.HelloGreeterAPIClient) *GreeterService {
	return &GreeterService{client: client}
}

// NewGreeterServer new a greeter service.
func NewGreeterServer(client helloworld.HelloGreeterAPIClient) helloworld.HelloGreeterAPIServer {
	return &GreeterService{client: client}
}

var _ helloworld.HelloGreeterAPIServer = (*GreeterService)(nil)
