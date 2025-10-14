/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package service

import (
	"context"

	"origadmin/basic-layout/api/v1/gen/go/secondworld" // Corrected import path
)

// GreeterHTTPService is a greeter service.
type GreeterHTTPService struct {
	secondworld.SecondGreeterAPIServer

	client secondworld.SecondGreeterAPIHTTPClient
}

func (g GreeterHTTPService) SayHello(ctx context.Context, request *secondworld.SayHelloRequest) (*secondworld.SayHelloResponse, error) {
	return g.client.SayHello(ctx, request)
}

func (g GreeterHTTPService) PostHello(ctx context.Context, request *secondworld.PostHelloRequest) (*secondworld.PostHelloResponse, error) {
	return g.client.PostHello(ctx, request)
}

func (g GreeterHTTPService) CreateGreeter(ctx context.Context, request *secondworld.CreateGreeterRequest) (*secondworld.CreateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) UpdateGreeter(ctx context.Context, request *secondworld.UpdateGreeterRequest) (*secondworld.UpdateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) DeleteGreeter(ctx context.Context, request *secondworld.DeleteGreeterRequest) (*secondworld.DeleteGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) GetGreeter(ctx context.Context, request *secondworld.GetGreeterRequest) (*secondworld.GetGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) ListGreeter(ctx context.Context, request *secondworld.ListGreeterRequest) (*secondworld.ListGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterHTTPService new a greeter service.
func NewGreeterHTTPService(client secondworld.SecondGreeterAPIHTTPClient) *GreeterHTTPService {
	return &GreeterHTTPService{client: client}
}

// NewGreeterHTTPServer new a greeter service.
func NewGreeterHTTPServer(client secondworld.SecondGreeterAPIHTTPClient) secondworld.SecondGreeterAPIServer {
	return &GreeterHTTPService{client: client}
}

var _ secondworld.SecondGreeterAPIServer = (*GreeterHTTPService)(nil)
