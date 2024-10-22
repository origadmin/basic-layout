package server

import (
	"fmt"
	"net/netip"
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
)

// NewGINSServer new a gin server.
func NewGINSServer(bootstrap *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
		),
	}
	c := bootstrap.Server
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

	if l != nil {
		opts = append(opts, gins.WithLogger(log.With(l, "module", "gins")))
	}
	naip, _ := netip.ParseAddrPort(bootstrap.Server.Gins.Addr)
	endpoint, _ := url.Parse(fmt.Sprintf("http://192.168.28.60:%d", naip.Port()))
	opts = append(opts, gins.Endpoint(endpoint))
	srv := gins.NewServer(opts...)
	helloworld.RegisterGreeterGINServer(srv, greeter)
	return srv
}
