/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"

	"github.com/origadmin/contrib/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

// NewGINSServer new a gin server.
func NewGINSServer(bs *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
			metadata.Server(),
		),
	}
	cfg := bs.GetService().GetGins()
	if cfg == nil {
		cfg = bootstrap.DefaultServiceGins()
	}
	if cfg.Network != "" {
		opts = append(opts, gins.Network(cfg.Network))
	}
	if cfg.Addr != "" {
		opts = append(opts, gins.Address(cfg.Addr))
	}
	if cfg.Timeout != nil {
		opts = append(opts, gins.Timeout(cfg.Timeout.AsDuration()))
	}
	//if c.Build == nil {
	//	c.Build = new(configs.Server_Middleware)
	//}
	middlewares, err := bootstrap.LoadMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, gins.Middleware(middlewares...))
	}

	if l != nil {
		opts = append(opts, gins.WithLogger(log.With(l, "module", "gins")))
	}

	naip, _ := netip.ParseAddrPort(bs.Service.Gins.Addr)
	if bs.Service.Gins.Endpoint == "" {
		bs.Service.Gins.Endpoint = "http://" + bs.Service.Host + ":" + strconv.Itoa(int(naip.Port()))
	} else {
		prefix, suffix, ok := strings.Cut(bs.Service.Gins.Endpoint, "://")
		if !ok {
			bs.Service.Gins.Endpoint = "http://" + prefix
		} else {
			args := strings.SplitN(suffix, ":", 2)
			if len(args) == 2 {
				args[1] = strconv.Itoa(int(naip.Port()))
			} else if len(args) == 1 {
				args = append(args, strconv.Itoa(int(naip.Port())))
			} else {
				// unknown
				log.Infow("unknown http endpoint", bs.Service.Gins.Endpoint)
			}
			bs.Service.Gins.Endpoint = prefix + "://" + strings.Join(args, ":")
		}
	}

	log.Infof("Server.GINS.EndpointURL: %v", bs.Service.Gins.Endpoint)
	ep, _ := url.Parse(bs.Service.Gins.Endpoint)
	opts = append(opts, gins.Endpoint(ep))
	srv := gins.NewServer(opts...)
	helloworld.RegisterHelloGreeterAPIGINSServer(srv, greeter)
	return srv
}
