package server

import (
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/configs"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGINSServer, NewHTTPServer)

func NewGINSServer(bs *configs.Bootstrap, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
		),
	}
	c := bs.Server
	if c.Entry == nil {
		c.Entry = new(configs.Server_Entry)
	}
	if c.Entry.Network != "" {
		opts = append(opts, gins.Network(c.Entry.Network))
	}
	if c.Entry.Addr != "" {
		opts = append(opts, gins.Address(c.Entry.Addr))
	}
	if c.Entry.Timeout != nil {
		opts = append(opts, gins.Timeout(c.Entry.Timeout.AsDuration()))
	}
	if c.Middleware == nil {
		c.Middleware = new(configs.Server_Middleware)
	}

	middlewares, err := bootloader.LoadGlobalMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, gins.Middleware(middlewares...))
	}

	if l != nil {
		opts = append(opts, gins.WithLogger(log.With(l, "module", "gins")))
	}

	srv := gins.NewServer(opts...)
	return srv
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bs *configs.Bootstrap, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.ErrorEncoder(func(w http.ResponseWriter, r *http.Request, err error) {
			se := helloworld.ToHttpError(err)
			codec, _ := http.CodecForRequest(r, "Accept")
			body, err := codec.Marshal(se)
			if err != nil {
				w.WriteHeader(nethttp.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(int(se.Code))
			_, _ = w.Write(body)
		}),
	}
	c := bs.Server
	if c.Http == nil {
		c.Http = new(configs.Server_HTTP)
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(":8000"))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	if c.Middleware == nil {
		c.Middleware = new(configs.Server_Middleware)
	}
	middlewares, err := bootloader.LoadGlobalMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, http.Middleware(middlewares...))
	}

	srv := http.NewServer(opts...)
	return srv
}
