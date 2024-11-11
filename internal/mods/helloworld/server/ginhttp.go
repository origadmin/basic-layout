package server

import (
	stdhttp "net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/origadmin/toolkits/runtime/config"

	"origadmin/basic-layout/api/v1/services/helloworld"
	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/configs"
)

func NewGinHTTPServer(bs *configs.Bootstrap, greeter helloworld.HelloGreeterAPIServer, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			metadata.Server(),
		),
	}
	c := bs.Service
	if c.Http == nil {
		c.Http = new(config.ServiceConfig_HTTP)
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
	//if c.Middleware == nil {
	//	c.Middleware = new(configs.Server_Middleware)
	//}
	middlewares, err := bootstrap.LoadMiddlewares(bs.GetServiceName(), bs, l)
	if err == nil && len(middlewares) > 0 {
		opts = append(opts, http.Middleware(middlewares...))
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

	log.Infof("Server.GINHTTP.Endpoint: %v", bs.Service.Gins.Endpoint)
	ep, _ := url.Parse(bs.Service.Gins.Endpoint)
	opts = append(opts, http.Endpoint(ep))
	srv := http.NewServer(opts...)
	engine := gin.New()

	srv.Server = &stdhttp.Server{
		Addr:         bs.Service.Gins.Addr,
		Handler:      engine.Handler(),
		ReadTimeout:  bs.Service.Gins.ReadTimeout.AsDuration(),
		WriteTimeout: bs.Service.Gins.WriteTimeout.AsDuration(),
		IdleTimeout:  bs.Service.Gins.IdleTimeout.AsDuration(),
	}
	helloworld.RegisterHelloGreeterAPIGINSServer(engine, greeter)
	return srv
}
