// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v5.27.2
// source: v1/protos/helloworld/greeter.proto

package helloworld

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationGreeterCreateGreeter = "/api.v1.services.helloworld.Greeter/CreateGreeter"
const OperationGreeterDeleteGreeter = "/api.v1.services.helloworld.Greeter/DeleteGreeter"
const OperationGreeterGetGreeter = "/api.v1.services.helloworld.Greeter/GetGreeter"
const OperationGreeterListGreeter = "/api.v1.services.helloworld.Greeter/ListGreeter"
const OperationGreeterPostHello = "/api.v1.services.helloworld.Greeter/PostHello"
const OperationGreeterSayHello = "/api.v1.services.helloworld.Greeter/SayHello"
const OperationGreeterUpdateGreeter = "/api.v1.services.helloworld.Greeter/UpdateGreeter"

type GreeterHTTPServer interface {
	// CreateGreeter CreateGreeter creates a new Greeter
	CreateGreeter(context.Context, *CreateGreeterRequest) (*CreateGreeterReply, error)
	// DeleteGreeter DeleteGreeter deletes a Greeter
	DeleteGreeter(context.Context, *DeleteGreeterRequest) (*DeleteGreeterReply, error)
	// GetGreeter GetGreeter gets a Greeter
	GetGreeter(context.Context, *GetGreeterRequest) (*GetGreeterReply, error)
	// ListGreeter ListGreeter lists Greeters
	ListGreeter(context.Context, *ListGreeterRequest) (*ListGreeterReply, error)
	// PostHello PostHello is a post method
	PostHello(context.Context, *GreeterRequest) (*GreeterReply, error)
	// SayHello SayHello is a get method
	SayHello(context.Context, *GreeterRequest) (*GreeterReply, error)
	// UpdateGreeter UpdateGreeter updates a Greeter
	UpdateGreeter(context.Context, *UpdateGreeterRequest) (*UpdateGreeterReply, error)
}

func RegisterGreeterHTTPServer(s *http.Server, srv GreeterHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/greeter/{id}/hello", _Greeter_SayHello0_HTTP_Handler(srv))
	r.POST("/api/v1/greeter/{id}/hello", _Greeter_PostHello0_HTTP_Handler(srv))
	r.POST("/api/v1/greeter", _Greeter_CreateGreeter0_HTTP_Handler(srv))
	r.PUT("/api/v1/greeter/{id}", _Greeter_UpdateGreeter0_HTTP_Handler(srv))
	r.DELETE("/api/v1/greeter/{id}", _Greeter_DeleteGreeter0_HTTP_Handler(srv))
	r.GET("/api/v1/greeter/{id}", _Greeter_GetGreeter0_HTTP_Handler(srv))
	r.GET("/api/v1/greeter", _Greeter_ListGreeter0_HTTP_Handler(srv))
}

func _Greeter_SayHello0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GreeterRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterSayHello)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayHello(ctx, req.(*GreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GreeterReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_PostHello0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GreeterRequest
		if err := ctx.Bind(&in.Data); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterPostHello)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.PostHello(ctx, req.(*GreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GreeterReply)
		return ctx.Result(200, reply.Data)
	}
}

func _Greeter_CreateGreeter0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateGreeterRequest
		if err := ctx.Bind(&in.Data); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterCreateGreeter)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateGreeter(ctx, req.(*CreateGreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateGreeterReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_UpdateGreeter0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateGreeterRequest
		if err := ctx.Bind(&in.Data); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterUpdateGreeter)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateGreeter(ctx, req.(*UpdateGreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateGreeterReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_DeleteGreeter0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteGreeterRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterDeleteGreeter)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteGreeter(ctx, req.(*DeleteGreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteGreeterReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_GetGreeter0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetGreeterRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterGetGreeter)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetGreeter(ctx, req.(*GetGreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetGreeterReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_ListGreeter0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListGreeterRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGreeterListGreeter)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListGreeter(ctx, req.(*ListGreeterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListGreeterReply)
		return ctx.Result(200, reply)
	}
}

type GreeterHTTPClient interface {
	CreateGreeter(ctx context.Context, req *CreateGreeterRequest, opts ...http.CallOption) (rsp *CreateGreeterReply, err error)
	DeleteGreeter(ctx context.Context, req *DeleteGreeterRequest, opts ...http.CallOption) (rsp *DeleteGreeterReply, err error)
	GetGreeter(ctx context.Context, req *GetGreeterRequest, opts ...http.CallOption) (rsp *GetGreeterReply, err error)
	ListGreeter(ctx context.Context, req *ListGreeterRequest, opts ...http.CallOption) (rsp *ListGreeterReply, err error)
	PostHello(ctx context.Context, req *GreeterRequest, opts ...http.CallOption) (rsp *GreeterReply, err error)
	SayHello(ctx context.Context, req *GreeterRequest, opts ...http.CallOption) (rsp *GreeterReply, err error)
	UpdateGreeter(ctx context.Context, req *UpdateGreeterRequest, opts ...http.CallOption) (rsp *UpdateGreeterReply, err error)
}

type GreeterHTTPClientImpl struct {
	cc *http.Client
}

func NewGreeterHTTPClient(client *http.Client) GreeterHTTPClient {
	return &GreeterHTTPClientImpl{client}
}

func (c *GreeterHTTPClientImpl) CreateGreeter(ctx context.Context, in *CreateGreeterRequest, opts ...http.CallOption) (*CreateGreeterReply, error) {
	var out CreateGreeterReply
	pattern := "/api/v1/greeter"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGreeterCreateGreeter))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in.Data, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GreeterHTTPClientImpl) DeleteGreeter(ctx context.Context, in *DeleteGreeterRequest, opts ...http.CallOption) (*DeleteGreeterReply, error) {
	var out DeleteGreeterReply
	pattern := "/api/v1/greeter/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGreeterDeleteGreeter))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GreeterHTTPClientImpl) GetGreeter(ctx context.Context, in *GetGreeterRequest, opts ...http.CallOption) (*GetGreeterReply, error) {
	var out GetGreeterReply
	pattern := "/api/v1/greeter/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGreeterGetGreeter))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GreeterHTTPClientImpl) ListGreeter(ctx context.Context, in *ListGreeterRequest, opts ...http.CallOption) (*ListGreeterReply, error) {
	var out ListGreeterReply
	pattern := "/api/v1/greeter"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGreeterListGreeter))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GreeterHTTPClientImpl) PostHello(ctx context.Context, in *GreeterRequest, opts ...http.CallOption) (*GreeterReply, error) {
	var out GreeterReply
	pattern := "/api/v1/greeter/{id}/hello"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGreeterPostHello))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in.Data, &out.Data, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GreeterHTTPClientImpl) SayHello(ctx context.Context, in *GreeterRequest, opts ...http.CallOption) (*GreeterReply, error) {
	var out GreeterReply
	pattern := "/api/v1/greeter/{id}/hello"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGreeterSayHello))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GreeterHTTPClientImpl) UpdateGreeter(ctx context.Context, in *UpdateGreeterRequest, opts ...http.CallOption) (*UpdateGreeterReply, error) {
	var out UpdateGreeterReply
	pattern := "/api/v1/greeter/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGreeterUpdateGreeter))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in.Data, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
