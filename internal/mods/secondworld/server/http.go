/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	servicehttp "github.com/origadmin/runtime/service/http"

	"origadmin/basic-layout/api/v1/services/secondworld"
	"origadmin/basic-layout/internal/configs"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bs *configs.Bootstrap, greeter secondworld.SecondGreeterAPIServer, l log.Logger) *http.Server {
	srv := servicehttp.NewServer(bs.GetService())
	secondworld.RegisterSecondGreeterAPIHTTPServer(srv, greeter)
	return srv
}

func RegisterHTTPServer(srv *http.Server, greeter secondworld.SecondGreeterAPIServer) {
	secondworld.RegisterSecondGreeterAPIHTTPServer(srv, greeter)
}
