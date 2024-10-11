package service

import (
	"context"

	"github.com/origadmin/basic-layout/api/v1/services/helloworld"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	helloworld.UnimplementedGreeterServiceServer

	uc *biz.GreeterBiz
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterBiz) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *helloworld.GreeterRequest) (*helloworld.GreeterReply, error) {
	return &helloworld.GreeterReply{Data: &helloworld.Greeter{
		Name: "hello " + in.Name,
	}}, nil
}
