package service

import (
	"context"
	"fmt"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

// GreeterService is a greeter service.
type GreeterService struct {
	helloworld.GreeterServiceServer

	//uc     *biz.GreeterBiz
	client helloworld.GreeterServiceClient
}

func (s *GreeterService) CreateGreeter(ctx context.Context, request *helloworld.CreateGreeterRequest) (*helloworld.CreateGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterService) UpdateGreeter(ctx context.Context, request *helloworld.UpdateGreeterRequest) (*helloworld.UpdateGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterService) DeleteGreeter(ctx context.Context, request *helloworld.DeleteGreeterRequest) (*helloworld.DeleteGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterService) GetGreeter(ctx context.Context, request *helloworld.GetGreeterRequest) (*helloworld.GetGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *GreeterService) ListGreeter(ctx context.Context, request *helloworld.ListGreeterRequest) (*helloworld.ListGreeterReply, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterService new a greeter service.
func NewGreeterService(client helloworld.GreeterServiceClient) *GreeterService {
	return &GreeterService{client: client}
}

// NewGreeterServer new a greeter service.
func NewGreeterServer(client helloworld.GreeterServiceClient) helloworld.GreeterServiceServer {
	return &GreeterService{client: client}
}

// SayHello implements helloworld.SayHello.
func (s *GreeterService) SayHello(ctx context.Context, in *helloworld.GreeterRequest) (*helloworld.GreeterReply, error) {
	fmt.Println("GreeterService.SayHello", in.Name)
	return s.client.SayHello(ctx, in)
}

// PostHello implements helloworld.PostHello.
func (s *GreeterService) PostHello(ctx context.Context, in *helloworld.GreeterRequest) (*helloworld.GreeterReply, error) {
	fmt.Println("GreeterService.PostHello", in.Data.Name)
	return s.client.PostHello(ctx, in)
}

var _ helloworld.GreeterServiceServer = (*GreeterService)(nil)
