/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"

	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	"github.com/origadmin/runtime/errors"
	"origadmin/basic-layout/api/v1/gen/go/secondworld"
	"origadmin/basic-layout/internal/mods/secondworld/dto"
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

func (g GreeterBiz) SayHello(ctx context.Context, in *secondworld.SayHelloRequest, opts ...grpc.CallOption) (*secondworld.SayHelloResponse, error) {
	log.Infof("SayHello: %v data: %v", in.Id, in.Data.Name)
	return &secondworld.SayHelloResponse{
		Data: in.Data,
	}, nil
}

func (g GreeterBiz) PostHello(ctx context.Context, in *secondworld.PostHelloRequest, opts ...grpc.CallOption) (*secondworld.PostHelloResponse, error) {
	log.Infof("GreeterBiz.PostHello: %v", in.Data.Name)
	return &secondworld.PostHelloResponse{
		Data: &secondworld.GreeterData{
			Id:   in.Data.Id,
			Name: "hello " + in.Data.Name,
		},
	}, nil
}

func (g GreeterBiz) CreateGreeter(ctx context.Context, in *secondworld.CreateGreeterRequest, opts ...grpc.CallOption) (*secondworld.CreateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) UpdateGreeter(ctx context.Context, in *secondworld.UpdateGreeterRequest, opts ...grpc.CallOption) (*secondworld.UpdateGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) DeleteGreeter(ctx context.Context, in *secondworld.DeleteGreeterRequest, opts ...grpc.CallOption) (*secondworld.DeleteGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) GetGreeter(ctx context.Context, in *secondworld.GetGreeterRequest, opts ...grpc.CallOption) (*secondworld.GetGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GreeterBiz) ListGreeter(ctx context.Context, in *secondworld.ListGreeterRequest, opts ...grpc.CallOption) (*secondworld.ListGreeterResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterBiz new a Greeter use case.
func NewGreeterBiz(repo dto.GreeterRepo, logger log.Logger) *GreeterBiz {
	return &GreeterBiz{dao: repo, log: log.NewHelper(logger)}
}

// NewGreeterClient new a Greeter use case.
func NewGreeterClient(repo dto.GreeterRepo, logger log.Logger) secondworld.SecondGreeterAPIClient {
	return &GreeterBiz{dao: repo, log: log.NewHelper(logger)}
}

var _ secondworld.SecondGreeterAPIClient = (*GreeterBiz)(nil)
