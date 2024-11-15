package client

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"github.com/origadmin/toolkits/runtime/transport/gins"
	"origadmin/basic-layout/helpers/errors"
	"origadmin/basic-layout/internal/configs"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGINSServer, NewHTTPServer)

func NewGINSServer(bootstrap *configs.Bootstrap, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
		),
		gins.ErrorEncoder(errors.GinErrorEncoder),
	}

	if bootstrap.GetEntry() == nil {
		bootstrap.Entry = new(configs.Bootstrap_Entry)
	}
	if bootstrap.Entry.Network != "" {
		opts = append(opts, gins.Network(bootstrap.Entry.Network))
	}
	if bootstrap.Entry.Addr != "" {
		opts = append(opts, gins.Address(bootstrap.Entry.Addr))
	}
	if bootstrap.Entry.Timeout != nil {
		opts = append(opts, gins.Timeout(bootstrap.Entry.Timeout.AsDuration()))
	}
	//if c.Middleware == nil {
	//	c.Middleware = new(configs.Server_Middleware)
	//}

	//middlewares, err := bootstrap.LoadGlobalMiddlewares(bootstrap.GetServiceName(), bootstrap, l)
	//if err == nil && len(middlewares) > 0 {
	//	opts = append(opts, gins.Middleware(middlewares...))
	//}

	if l != nil {
		opts = append(opts, gins.WithLogger(log.With(l, "module", "gins")))
	}

	srv := gins.NewServer(opts...)
	return srv
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bootstrap *configs.Bootstrap, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.ErrorEncoder(errors.HttpErrorEncoder),
	}

	bootstrap.GetEntry()
	if bootstrap.GetEntry() == nil {
		bootstrap.Entry = new(configs.Bootstrap_Entry)
	}
	if bootstrap.Entry.Network != "" {
		opts = append(opts, http.Network(bootstrap.Entry.Network))
	}
	if bootstrap.Entry.Addr != "" {
		opts = append(opts, http.Address(bootstrap.Entry.Addr))
	}
	if bootstrap.Entry.Timeout != nil {
		opts = append(opts, http.Timeout(bootstrap.Entry.Timeout.AsDuration()))
	}
	//if bootstrap.Middleware == nil {
	//	bootstrap.Middleware = new(configs.Server_Middleware)
	//}
	//middlewares, err := bootstrap.LoadGlobalMiddlewares(bootstrap.GetServiceName(), bootstrap, l)
	//if err == nil && len(middlewares) > 0 {
	//	opts = append(opts, http.Middleware(middlewares...))
	//}

	srv := http.NewServer(opts...)
	return srv
}
