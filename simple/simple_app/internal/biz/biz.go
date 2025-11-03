package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	simplev1 "basic-layout/simple/simple_app/api/gen/go/simple/v1"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewSimpleUsecase)

// SimpleRepo is a Simple repo.
type SimpleRepo interface {
	Save(context.Context, *Simple) (*Simple, error)
	Update(context.Context, *Simple) (*Simple, error)
	FindByID(context.Context, int64) (*Simple, error)
	ListByName(context.Context, string) ([]*Simple, error)
	ListAll(context.Context) ([]*Simple, error)
}

// Simple is a Simple model.
type Simple struct {
	Name string
}

// SimpleUsecase is a Simple usecase.
type SimpleUsecase struct {
	repo SimpleRepo
	log  *log.Helper
}

// NewSimpleUsecase new a Simple usecase.
func NewSimpleUsecase(repo SimpleRepo, logger log.Logger) *SimpleUsecase {
	return &SimpleUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Greet creates a greeting, and returns the new greeting.
func (uc *SimpleUsecase) Greet(ctx context.Context, s *Simple) (*simplev1.SayHelloResponse, error) {
	uc.log.WithContext(ctx).Infof("Greet: %v", s.Name)
	return &simplev1.SayHelloResponse{Message: "Hello " + s.Name}, nil
}
