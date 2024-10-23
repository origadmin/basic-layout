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
	"origadmin/basic-layout/internal/configs"
)

func NewGinHTTPServer(bootstrap *configs.Bootstrap, greeter helloworld.GreeterServer, l log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	c := bootstrap.Server
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

	naip, _ := netip.ParseAddrPort(bootstrap.Server.Gins.Addr)
	prefix, suffix, ok := strings.Cut(bootstrap.Server.Gins.Endpoint, "://")
	if !ok {
		bootstrap.Server.Gins.Endpoint = "http://" + prefix
	} else {
		args := strings.SplitN(suffix, ":", 2)
		if len(args) == 2 {
			args[1] = strconv.Itoa(int(naip.Port()))
		} else if len(args) == 1 {
			args = append(args, strconv.Itoa(int(naip.Port())))
		} else {
			// unknown
			log.NewHelper(l).Info("unknown http endpoint", bootstrap.Server.Gins.Endpoint)
		}
		bootstrap.Server.Gins.Endpoint = prefix + "://" + strings.Join(args, ":")
	}

	log.NewHelper(l).Infof("bootstrap.Server.Gins.Endpoint: %v", bootstrap.Server.Gins.Endpoint)
	ep, _ := url.Parse(bootstrap.Server.Gins.Endpoint)
	opts = append(opts, http.Endpoint(ep))
	srv := http.NewServer(opts...)
	engine := gin.New()

	srv.Server = &stdhttp.Server{
		Addr:         bootstrap.Server.Gins.Addr,
		Handler:      engine.Handler(),
		ReadTimeout:  bootstrap.Server.Gins.ReadTimeout.AsDuration(),
		WriteTimeout: bootstrap.Server.Gins.WriteTimeout.AsDuration(),
		IdleTimeout:  bootstrap.Server.Gins.IdleTimeout.AsDuration(),
	}
	helloworld.RegisterGreeterGINServer(engine, greeter)
	return srv
}
