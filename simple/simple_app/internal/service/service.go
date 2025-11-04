package service

import (
	"context"

	"github.com/google/wire"

	simplev1 "basic-layout/simple/simple_app/api/gen/go/simple/v1"
	"basic-layout/simple/simple_app/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewSimpleService)

// SimpleService is a simple service.
type SimpleService struct {
	simplev1.UnimplementedSimpleServiceServer

	uc *biz.SimpleUsecase
}

// NewSimpleService new a simple service.
func NewSimpleService(uc *biz.SimpleUsecase) *SimpleService {
	return &SimpleService{uc: uc}
}

// SayHello implements SimpleServiceServer.
func (s *SimpleService) SayHello(ctx context.Context, in *simplev1.SayHelloRequest) (*simplev1.SayHelloResponse, error) {
	return s.uc.SayHello(ctx, &biz.Simple{Name: in.Name})
}
