package server

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bootstrap *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	c := bootstrap.Server
	if c.Http == nil {
		c.Http = new(configs.Server_HTTP)
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
		c.Middleware = new(configs.Server_Middleware)
	}

	naip, _ := netip.ParseAddrPort(bootstrap.Server.Http.Addr)
	prefix, suffix, ok := strings.Cut(bootstrap.Server.Http.Endpoint, "://")
	if !ok {
		bootstrap.Server.Http.Endpoint = "http://" + prefix
	} else {
		args := strings.SplitN(suffix, ":", 2)
		if len(args) == 2 {
			args[1] = strconv.Itoa(int(naip.Port()))
		} else if len(args) == 1 {
			args = append(args, strconv.Itoa(int(naip.Port())))
		} else {
			// unknown
			log.NewHelper(l).Info("unknown http endpoint", bootstrap.Server.Http.Endpoint)
		}
		bootstrap.Server.Http.Endpoint = prefix + "://" + strings.Join(args, ":")
	}

	log.NewHelper(l).Infof("bootstrap.Server.Http.Endpoint: %v", bootstrap.Server.Http.Endpoint)
	ep, _ := url.Parse(bootstrap.Server.Http.Endpoint)
	opts = append(opts, http.Endpoint(ep))
	srv := http.NewServer(opts...)
	helloworld.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
