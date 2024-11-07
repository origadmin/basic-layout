package server

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/origadmin/toolkits/runtime/config"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

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
	c := bs.Service
	if c.Gins == nil {
		c.Gins = new(config.ServiceConfig_GINS)
	}
	if c.Gins.Network != "" {
		opts = append(opts, gins.Network(c.Gins.Network))
	}
	if c.Gins.Addr != "" {
		opts = append(opts, gins.Address(c.Gins.Addr))
	}
	if c.Gins.Timeout != nil {
		opts = append(opts, gins.Timeout(c.Gins.Timeout.AsDuration()))
	}
	//if c.Middleware == nil {
	//	c.Middleware = new(configs.Server_Middleware)
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

	log.Infof("Server.Gins.Endpoint: %v", bs.Service.Gins.Endpoint)
	ep, _ := url.Parse(bs.Service.Gins.Endpoint)
	opts = append(opts, gins.Endpoint(ep))
	srv := gins.NewServer(opts...)
	helloworld.RegisterHelloGreeterAPIGINSServer(srv, greeter)
	return srv
}
