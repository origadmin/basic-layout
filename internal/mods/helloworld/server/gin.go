package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/mods/helloworld/conf"
)

// NewGINServer new a gin server.
func NewGINServer(bootstrap *conf.Bootstrap, greeter helloworld.GreeterServiceServer, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
		),
	}
	c := bootstrap.Server
	if c.Gins == nil {
		c.Gins = new(conf.Server_GINS)
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
		c.Middleware = new(conf.Server_Middleware)
	}

	opts = append(opts)
	srv := gins.NewServer(opts...)
	helloworld.RegisterGreeterServiceGINServer(srv, greeter)
	return srv
}
