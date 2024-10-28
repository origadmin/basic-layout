package server

import (
	stdhttp "net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

func NewGinHTTPServer(bs *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	c := bs.Server
	if c.Http == nil {
		c.Http = new(configs.Server_HTTP)
	}
	if c.Gins.Network != "" {
		opts = append(opts, http.Network(c.Gins.Network))
	}
	if c.Gins.Addr != "" {
		opts = append(opts, http.Address(c.Gins.Addr))
	}
	if c.Gins.Timeout != nil {
		opts = append(opts, http.Timeout(c.Gins.Timeout.AsDuration()))
	}
	if c.Middleware == nil {
		c.Middleware = new(configs.Server_Middleware)
	}
	middlewares, err := bootstrap.LoadMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, http.Middleware(middlewares...))
	}

	naip, _ := netip.ParseAddrPort(bs.Server.Gins.Addr)
	if bs.Server.Gins.Endpoint == "" {
		bs.Server.Gins.Endpoint = "http://" + bs.Server.Host + ":" + strconv.Itoa(int(naip.Port()))
	} else {
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
	}

	log.Infof("Server.GinHttp.Endpoint: %v", bs.Server.Gins.Endpoint)
	ep, _ := url.Parse(bs.Server.Gins.Endpoint)
	opts = append(opts, http.Endpoint(ep))
	srv := http.NewServer(opts...)
	engine := gin.New()

	srv.Server = &stdhttp.Server{
		Addr:         bs.Server.Gins.Addr,
		Handler:      engine.Handler(),
		ReadTimeout:  bs.Server.Gins.ReadTimeout.AsDuration(),
		WriteTimeout: bs.Server.Gins.WriteTimeout.AsDuration(),
		IdleTimeout:  bs.Server.Gins.IdleTimeout.AsDuration(),
	}
	helloworld.RegisterGreeterGINServer(engine, greeter)
	return srv
}
