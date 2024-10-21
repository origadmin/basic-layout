package server

import (
	"fmt"
	"net/netip"
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/conf"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bootstrap *conf.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *http.Server {
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

	naip, _ := netip.ParseAddrPort(bootstrap.Server.Http.Addr)
	endpoint, _ := url.Parse(fmt.Sprintf("http://192.168.28.81:%d", naip.Port()))
	opts = append(opts, http.Endpoint(endpoint))
	srv := http.NewServer(opts...)
	helloworld.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
