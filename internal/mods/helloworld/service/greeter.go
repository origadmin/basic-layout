package service

import (
	"context"
	"fmt"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

// GreeterService is a greeter service.
type GreeterService struct {
	helloworld.UnimplementedGreeterServiceServer

	//uc     *biz.GreeterBiz
	client helloworld.GreeterServiceClient
}

// NewGreeterService new a greeter service.
func NewGreeterService(client helloworld.GreeterServiceClient) *GreeterService {
	return &GreeterService{client: client}
}

// NewGreeterServer new a greeter service.
func NewGreeterServer(client helloworld.GreeterServiceClient) helloworld.GreeterServiceServer {
	return &GreeterService{client: client}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *helloworld.GreeterRequest) (*helloworld.GreeterReply, error) {
	fmt.Println("SayHello", in.Name)
	return s.client.SayHello(ctx, in)
}

var _ helloworld.GreeterServiceServer = (*GreeterService)(nil)
