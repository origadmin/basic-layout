package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/google/wire"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

func NewServer(bootstrap *configs.Bootstrap, l log.Logger) *gins.Server {
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

	opts = append(opts)
	srv := gins.NewServer(opts...)
	//helloworld.RegisterGreeterServiceGINServer(srv, greeter)
	//todo
	return srv
}
