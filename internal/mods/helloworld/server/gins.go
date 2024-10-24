package server

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

// NewGINSServer new a gin server.
func NewGINSServer(bs *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
		),
	}
	c := bs.Server
	if c.Gins == nil {
		c.Gins = new(configs.Server_GINS)
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
	if c.Middleware == nil {
		c.Middleware = new(configs.Server_Middleware)
	}
	middlewares, err := bootstrap.LoadMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, gins.Middleware(middlewares...))
	}

	if l != nil {
		opts = append(opts, gins.WithLogger(log.With(l, "module", "gins")))
	}

	naip, _ := netip.ParseAddrPort(bs.Server.Gins.Addr)
	prefix, suffix, ok := strings.Cut(bs.Server.Gins.Endpoint, "://")
	if !ok {
		bs.Server.Gins.Endpoint = "http://" + prefix
	} else {
		args := strings.SplitN(suffix, ":", 2)
		if len(args) == 2 {
			args[1] = strconv.Itoa(int(naip.Port()))
		} else if len(args) == 1 {
			args = append(args, strconv.Itoa(int(naip.Port())))
		} else {
			// unknown
			log.Infow("unknown http endpoint", bs.Server.Gins.Endpoint)
		}
		bs.Server.Gins.Endpoint = prefix + "://" + strings.Join(args, ":")
	}

	log.Infof("bs.Server.Gins.Endpoint: %v", bs.Server.Gins.Endpoint)
	ep, _ := url.Parse(bs.Server.Gins.Endpoint)
	opts = append(opts, gins.Endpoint(ep))
	srv := gins.NewServer(opts...)
	helloworld.RegisterGreeterGINServer(srv, greeter)
	return srv
}
