package biz

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/basic-layout/api/v1/services/helloworld"
	"github.com/origadmin/basic-layout/internal/mods/helloworld/dto"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = helloworld.ErrorHTTP(helloworld.HelloWorldErrorReason_USER_NOT_FOUND, http.StatusNotFound, "user not found")
)

// GreeterBiz is a Greeter use case.
type GreeterBiz struct {
	repo dto.GreeterDao
	log  *log.Helper
}

// NewGreeterBiz new a Greeter use case.
func NewGreeterBiz(repo dto.GreeterDao, logger log.Logger) *GreeterBiz {
	return &GreeterBiz{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterBiz) CreateGreeter(ctx context.Context, g *dto.Greeter) (*dto.Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

func (uc *GreeterBiz) ListGreeter(ctx context.Context, g *dto.Greeter) ([]*dto.Greeter, error) {
	uc.log.WithContext(ctx).Infof("ListGreeter: %v", g.Hello)
	return uc.repo.ListByHello(ctx, g.Hello, &dto.GreeterQueryParam{})
}
