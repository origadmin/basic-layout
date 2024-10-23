package service

import (
	"context"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

// GreeterHTTPService is a greeter service.
type GreeterHTTPService struct {
	helloworld.GreeterServer

	client helloworld.GreeterHTTPClient
}

func (s *GreeterHTTPService) CreateGreeter(ctx context.Context, request *helloworld.CreateGreeterRequest) (*helloworld.CreateGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterHTTPService) UpdateGreeter(ctx context.Context, request *helloworld.UpdateGreeterRequest) (*helloworld.UpdateGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterHTTPService) DeleteGreeter(ctx context.Context, request *helloworld.DeleteGreeterRequest) (*helloworld.DeleteGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterHTTPService) GetGreeter(ctx context.Context, request *helloworld.GetGreeterRequest) (*helloworld.GetGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterHTTPService) ListGreeter(ctx context.Context, request *helloworld.ListGreeterRequest) (*helloworld.ListGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterHTTPService new a greeter service.
func NewGreeterHTTPService(client helloworld.GreeterHTTPClient) *GreeterHTTPService {
	return &GreeterHTTPService{client: client}
}

// NewGreeterHTTPServer new a greeter service.
func NewGreeterHTTPServer(client helloworld.GreeterHTTPClient) helloworld.GreeterServer {
	return &GreeterHTTPService{client: client}
}

// SayHello implements helloworld.SayHello.
func (s *GreeterHTTPService) SayHello(ctx context.Context, in *helloworld.GreeterRequest) (*helloworld.GreeterReply, error) {
	return s.client.SayHello(ctx, in)
}

// PostHello implements helloworld.PostHello.
func (s *GreeterHTTPService) PostHello(ctx context.Context, in *helloworld.GreeterRequest) (*helloworld.GreeterReply, error) {
	return s.client.PostHello(ctx, in)
}

var _ helloworld.GreeterServer = (*GreeterHTTPService)(nil)
