/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	stdhttp "net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/origadmin/runtime/middleware"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/helpers/endpoint"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

func NewGinHTTPServer(bs *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) *http.Server {
	ms := middleware.NewServer(bs.GetMiddleware())
	var opts = []http.ServerOption{
		http.Middleware(ms...),
	}
	cfg := bs.GetService().GetGins()
	if cfg == nil {
		cfg = bootstrap.DefaultServiceGins()
	}

	if cfg.Network != "" {
		opts = append(opts, http.Network(cfg.Network))
	}
	if cfg.Addr != "" {
		opts = append(opts, http.Address(cfg.Addr))
	}
	if cfg.Timeout != nil {
		opts = append(opts, http.Timeout(cfg.Timeout.AsDuration()))
	}

	cfg.Endpoint = endpoint.Parse("http", bs.GetService().Host, cfg.Addr, cfg.Endpoint)
	log.Infof("Server.GinHttp.EndpointURL: %v", cfg.Endpoint)
	ep, _ := url.Parse(cfg.Endpoint)
	opts = append(opts, http.Endpoint(ep))
	srv := http.NewServer(opts...)
	engine := gin.New()

	srv.Server = &stdhttp.Server{
		Addr:         cfg.Addr,
		Handler:      engine.Handler(),
		ReadTimeout:  cfg.ReadTimeout.AsDuration(),
		WriteTimeout: cfg.WriteTimeout.AsDuration(),
		IdleTimeout:  cfg.IdleTimeout.AsDuration(),
	}

	helloworld.RegisterHelloGreeterAPIGINSServer(engine, greeter)
	return srv
}
