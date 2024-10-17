package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods/helloworld/conf"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bootstrap *conf.Bootstrap, greeter helloworld.GreeterServiceServer, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	c := bootstrap.Server
	if c.Http == nil {
		c.Http = new(conf.Server_HTTP)
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	if c.Middleware == nil {
		c.Middleware = new(conf.Server_Middleware)
	}

	srv := http.NewServer(opts...)
	helloworld.RegisterGreeterServiceHTTPServer(srv, greeter)
	return srv
}
