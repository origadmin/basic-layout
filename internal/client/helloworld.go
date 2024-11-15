/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package client implements the functions, types, and interfaces for the module.
package client

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/runtime/registry"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
	helloworldservice "origadmin/basic-layout/internal/mods/helloworld/service"
)

const (
	// DefaultHelloWorldServiceName is the name of the service.
	DefaultHelloWorldServiceName = "origadmin.service.v1.helloworld"
)

func NewHelloGreeterAPIServer(bootstrap *configs.Bootstrap, discovery registry.Discovery) (helloworld.HelloGreeterAPIServer, error) {
	// Create route Filter: Filter instances whose version number is "2.0.0"
	filter := filter.Version("v1.0.0")

	// Create the Selector for the P2C load balancing algorithm and inject the route Filter
	selector.SetGlobalSelector(random.NewBuilder())
	//selector.SetGlobalSelector(wrr.NewBuilder())

	serviceName := DefaultHelloWorldServiceName
	//if registry.ServiceName != "" {
	//	serviceName = registry.ServiceName
	//}
	//
	//discovery, err := runtime.NewDiscovery(registry)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to create discovery")
	//}

	//if discovery, ok := injector.Discoveries[serviceName]; ok {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
		),
		grpc.WithEndpoint("discovery:///"+serviceName),
		grpc.WithDiscovery(discovery),
		grpc.WithNodeFilter(filter),
		grpc.WithPrintDiscoveryDebugLog(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client")
	}
	gClient := helloworld.NewHelloGreeterAPIClient(conn)
	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
		),
		http.WithEndpoint("discovery:///"+serviceName),
		http.WithDiscovery(discovery),
		http.WithNodeFilter(filter),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http client")
	}
	hClient := helloworld.NewHelloGreeterAPIHTTPClient(hConn)

	var client helloworld.HelloGreeterAPIServer
	if entry := bootstrap.GetEntry(); entry != nil && entry.Scheme == "http" {
		client = helloworldservice.NewGreeterHTTPServer(hClient)
	} else {
		client = helloworldservice.NewGreeterServer(gClient)
	}
	//grpcClient := service.NewGreeterServer(gClient)
	//httpClient := service.NewGreeterHTTPServer(hClient)
	//// add _ to avoid unused
	//_ = grpcClient
	//_ = httpClient
	//helloworld.RegisterHelloGreeterAPIGINSServer(injector.ServerGINS, client)
	//helloworld.RegisterHelloGreeterAPIHTTPServer(injector.ServerHTTP, client)
	//}

	return client, nil
}
