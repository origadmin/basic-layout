package biz

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods/helloworld/dto"
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

func (uc *GreeterBiz) ListGreeter(ctx context.Context, in *helloworld.ListGreeterRequest, opts ...grpc.CallOption) (*helloworld.ListGreeterReply, error) {
	uc.log.WithContext(ctx).Infof("ListGreeter")
	return &helloworld.ListGreeterReply{}, nil
}

func (uc *GreeterBiz) SayHello(ctx context.Context, in *helloworld.GreeterRequest, opts ...grpc.CallOption) (*helloworld.GreeterReply, error) {
	uc.log.WithContext(ctx).Infof("SayHello: %v", in.Name)
	return &helloworld.GreeterReply{Data: &helloworld.Greeter{
		Name: "hello " + in.Name,
	}}, nil
}

func (uc *GreeterBiz) CreateGreeter(ctx context.Context, in *helloworld.CreateGreeterRequest, opts ...grpc.CallOption) (*helloworld.CreateGreeterReply, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", in.Data.Name)
	_, err := uc.repo.Save(ctx, &dto.Greeter{
		Name: in.Data.Name,
	})
	if err != nil {
		return nil, err
	}
	return &helloworld.CreateGreeterReply{}, nil
}

func (uc *GreeterBiz) UpdateGreeter(ctx context.Context, in *helloworld.UpdateGreeterRequest, opts ...grpc.CallOption) (*helloworld.UpdateGreeterReply, error) {
	uc.log.WithContext(ctx).Infof("UpdateGreeter: %v", in.Data.Name)
	_, err := uc.repo.Update(ctx, &dto.Greeter{
		Name: in.Data.Name,
	})
	if err != nil {
		return nil, err
	}
	return &helloworld.UpdateGreeterReply{}, nil
}

func (uc *GreeterBiz) DeleteGreeter(ctx context.Context, in *helloworld.DeleteGreeterRequest, opts ...grpc.CallOption) (*helloworld.DeleteGreeterReply, error) {
	uc.log.WithContext(ctx).Infof("DeleteGreeter: %v", in.Id)
	//err := uc.repo.Delete(ctx, in.Id)
	//if err != nil {
	//	return nil, err
	//}
	return &helloworld.DeleteGreeterReply{}, nil
}

func (uc *GreeterBiz) GetGreeter(ctx context.Context, in *helloworld.GetGreeterRequest, opts ...grpc.CallOption) (*helloworld.GetGreeterReply, error) {
	uc.log.WithContext(ctx).Infof("GetGreeter: %v", in.Id)
	return &helloworld.GetGreeterReply{}, nil
}

// NewGreeterBiz new a Greeter use case.
func NewGreeterBiz(repo dto.GreeterDao, logger log.Logger) *GreeterBiz {
	return &GreeterBiz{repo: repo, log: log.NewHelper(logger)}
}

// NewGreeterClient new a Greeter use case.
func NewGreeterClient(repo dto.GreeterDao, logger log.Logger) helloworld.GreeterServiceClient {
	return &GreeterBiz{repo: repo, log: log.NewHelper(logger)}
}

var _ helloworld.GreeterServiceClient = (*GreeterBiz)(nil)
