/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"

	"origadmin/basic-layout/api/v1/gen/go/helloworld" // Corrected import path
	"origadmin/basic-layout/internal/mods/helloworld/dto"

	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	"github.com/origadmin/runtime/errors"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NewMessage(commonv1.ErrorReason_NOT_FOUND, "user not found")
)

// GreeterBiz is a Greeter use case.
type GreeterBiz struct {
	dao dto.GreeterRepo
	log *log.Helper
}

func (g GreeterBiz) SayHello(ctx context.Context, in *helloworld.SayHelloRequest, opts ...grpc.CallOption) (*helloworld.SayHelloResponse, error) {
	log.Infof("SayHello: %v data: %v", in.Id, in.Data.Name)
	return &helloworld.SayHelloResponse{
		Data: &dto.Greeter{
			Name: "hello " + in.Id,
		}}, nil
}

func (g GreeterBiz) PostHello(ctx context.Context, in *helloworld.PostHelloRequest, opts ...grpc.CallOption) (*helloworld.PostHelloResponse, error) {
	log.Infof("GreeterBiz.PostHello: %v", in.Data.Name)
	return &helloworld.PostHelloResponse{
		Data: &dto.Greeter{
			Name: "hello " + in.Data.Name,
		}}, nil
}

func (g GreeterBiz) CreateGreeter(ctx context.Context, in *helloworld.CreateGreeterRequest, opts ...grpc.CallOption) (*helloworld.CreateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) UpdateGreeter(ctx context.Context, in *helloworld.UpdateGreeterRequest, opts ...grpc.CallOption) (*helloworld.UpdateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) DeleteGreeter(ctx context.Context, in *helloworld.DeleteGreeterRequest, opts ...grpc.CallOption) (*helloworld.DeleteGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) GetGreeter(ctx context.Context, in *helloworld.GetGreeterRequest, opts ...grpc.CallOption) (*helloworld.GetGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) ListGreeter(ctx context.Context, in *helloworld.ListGreeterRequest, opts ...grpc.CallOption) (*helloworld.ListGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterBiz new a Greeter use case.
func NewGreeterBiz(repo dto.GreeterRepo, logger log.Logger) *GreeterBiz {
	return &GreeterBiz{dao: repo, log: log.NewHelper(logger)}
}

// NewGreeterClient new a Greeter use case.
func NewGreeterClient(repo dto.GreeterRepo, logger log.Logger) helloworld.HelloGreeterAPIClient {
	return &GreeterBiz{dao: repo, log: log.NewHelper(logger)}
}

var _ helloworld.HelloGreeterAPIClient = (*GreeterBiz)(nil)
