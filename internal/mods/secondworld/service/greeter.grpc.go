package service

import (
	"context"

	"origadmin/basic-layout/api/v1/services/secondworld"
)

// GreeterService is a greeter service.
type GreeterService struct {
	secondworld.SecondGreeterAPIServer

	client secondworld.SecondGreeterAPIClient
}

func (g GreeterService) SayHello(ctx context.Context, request *secondworld.SayHelloRequest) (*secondworld.SayHelloResponse, error) {
	return g.client.SayHello(ctx, request)
}

func (g GreeterService) PostHello(ctx context.Context, request *secondworld.PostHelloRequest) (*secondworld.PostHelloResponse, error) {
	return g.client.PostHello(ctx, request)
}

func (g GreeterService) CreateGreeter(ctx context.Context, request *secondworld.CreateGreeterRequest) (*secondworld.CreateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) UpdateGreeter(ctx context.Context, request *secondworld.UpdateGreeterRequest) (*secondworld.UpdateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) DeleteGreeter(ctx context.Context, request *secondworld.DeleteGreeterRequest) (*secondworld.DeleteGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) GetGreeter(ctx context.Context, request *secondworld.GetGreeterRequest) (*secondworld.GetGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterService) ListGreeter(ctx context.Context, request *secondworld.ListGreeterRequest) (*secondworld.ListGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterService new a greeter service.
func NewGreeterService(client secondworld.SecondGreeterAPIClient) *GreeterService {
	return &GreeterService{client: client}
}

// NewGreeterServer new a greeter service.
func NewGreeterServer(client secondworld.SecondGreeterAPIClient) secondworld.SecondGreeterAPIServer {
	return &GreeterService{client: client}
}

var _ secondworld.SecondGreeterAPIServer = (*GreeterService)(nil)
