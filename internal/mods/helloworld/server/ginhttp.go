package server

import (
	"fmt"
	stdhttp "net/http"
	"net/netip"
	"net/url"

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
	endpoint, _ := url.Parse(fmt.Sprintf("http://192.168.28.60:%d", naip.Port()))
	opts = append(opts, http.Endpoint(endpoint))
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
