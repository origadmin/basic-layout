package server

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/origadmin/toolkits/runtime/config"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bs *configs.Bootstrap, greeter helloworld.GreeterAPIServer, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	c := bs.Service
	if c.Http == nil {
		c.Http = new(config.ServiceConfig_HTTP)
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
	//if c.Middleware == nil {
	//	c.Middleware = new(configs.Server_Middleware)
	//}
	middlewares, err := bootstrap.LoadMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, http.Middleware(middlewares...))
	}

	naip, _ := netip.ParseAddrPort(bs.Service.Http.Addr)
	if bs.Service.Http.Endpoint == "" {
		bs.Service.Http.Endpoint = "http://" + bs.Service.Host + ":" + strconv.Itoa(int(naip.Port()))
	} else {
		prefix, suffix, ok := strings.Cut(bs.Service.Http.Endpoint, "://")
		if !ok {
			bs.Service.Http.Endpoint = "http://" + prefix
		} else {
			args := strings.SplitN(suffix, ":", 2)
			if len(args) == 2 {
				args[1] = strconv.Itoa(int(naip.Port()))
			} else if len(args) == 1 {
				args = append(args, strconv.Itoa(int(naip.Port())))
			} else {
				// unknown
				log.Infow("unknown http endpoint", bs.Service.Http.Endpoint)
			}
			bs.Service.Http.Endpoint = prefix + "://" + strings.Join(args, ":")
		}
	}

	log.Infof("Server.Http.Endpoint: %v", bs.Service.Http.Endpoint)
	ep, _ := url.Parse(bs.Service.Http.Endpoint)
	opts = append(opts, http.Endpoint(ep))
	srv := http.NewServer(opts...)
	helloworld.RegisterGreeterAPIHTTPServer(srv, greeter)
	return srv
}
