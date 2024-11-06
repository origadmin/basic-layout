package server

import (
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/origadmin/toolkits/runtime/config"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"

	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
	"origadmin/basic-layout/toolkits/errors"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGINSServer, NewHTTPServer)

func NewGINSServer(bs *configs.Bootstrap, l log.Logger) *gins.Server {
	var opts = []gins.ServerOption{
		gins.Middleware(
			recovery.Recovery(),
		),
		gins.ErrorEncoder(func(ctx *gins.Context, err error) {
			se := errors.ToHttpError(err)
			codec, _ := http.CodecForRequest(ctx.Request, "Accept")
			body, err := codec.Marshal(se)
			if err != nil {
				ctx.Writer.WriteHeader(nethttp.StatusInternalServerError)
				return
			}
			ctx.Writer.Header().Set("Content-Type", ctx.Request.Header.Get("Content-Type"))
			ctx.Writer.WriteHeader(int(se.Code))
			_, _ = ctx.Writer.Write(body)
		}),
	}
	c := bs.Service
	if c == nil {
		c = new(config.ServiceConfig)
	}
	if c.Entry == nil {
		c.Entry = new(config.ServiceConfig_Entry)
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
	//if c.Middleware == nil {
	//	c.Middleware = new(configs.Server_Middleware)
	//}

	middlewares, err := bootstrap.LoadGlobalMiddlewares(bs.GetServiceName(), bs, l)
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
		http.ErrorEncoder(encodeError),
	}
	c := bs.Service
	if c == nil {
		c = new(config.ServiceConfig)
	}
	if c.Entry == nil {
		c.Entry = new(config.ServiceConfig_Entry)
	}
	if c.Entry.Network != "" {
		opts = append(opts, http.Network(c.Entry.Network))
	}
	if c.Entry.Addr != "" {
		opts = append(opts, http.Address(c.Entry.Addr))
	}
	if c.Entry.Timeout != nil {
		opts = append(opts, http.Timeout(c.Entry.Timeout.AsDuration()))
	}
	//if c.Middleware == nil {
	//	c.Middleware = new(configs.Server_Middleware)
	//}
	middlewares, err := bootstrap.LoadGlobalMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, http.Middleware(middlewares...))
	}

	srv := http.NewServer(opts...)
	return srv
}

func encodeError(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.ToHttpError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}
