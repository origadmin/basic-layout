package service

import (
	"context"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

// GreeterHTTPService is a greeter service.
type GreeterHTTPService struct {
	helloworld.GreeterAPIServer

	client helloworld.GreeterAPIHTTPClient
}

func (g GreeterHTTPService) SayHello(ctx context.Context, request *helloworld.SayHelloRequest) (*helloworld.SayHelloResponse, error) {
	return g.client.SayHello(ctx, request)
}

func (g GreeterHTTPService) PostHello(ctx context.Context, request *helloworld.PostHelloRequest) (*helloworld.PostHelloResponse, error) {
	return g.client.PostHello(ctx, request)
}

func (g GreeterHTTPService) CreateGreeter(ctx context.Context, request *helloworld.CreateGreeterRequest) (*helloworld.CreateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) UpdateGreeter(ctx context.Context, request *helloworld.UpdateGreeterRequest) (*helloworld.UpdateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) DeleteGreeter(ctx context.Context, request *helloworld.DeleteGreeterRequest) (*helloworld.DeleteGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) GetGreeter(ctx context.Context, request *helloworld.GetGreeterRequest) (*helloworld.GetGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterHTTPService) ListGreeter(ctx context.Context, request *helloworld.ListGreeterRequest) (*helloworld.ListGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterHTTPService new a greeter service.
func NewGreeterHTTPService(client helloworld.GreeterAPIHTTPClient) *GreeterHTTPService {
	return &GreeterHTTPService{client: client}
}

// NewGreeterHTTPServer new a greeter service.
func NewGreeterHTTPServer(client helloworld.GreeterAPIHTTPClient) helloworld.GreeterAPIServer {
	return &GreeterHTTPService{client: client}
}

var _ helloworld.GreeterAPIServer = (*GreeterHTTPService)(nil)
