// Package configs implements the functions, types, and interfaces for the module.
package configs

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	storagev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/storage/v1"
	"github.com/origadmin/runtime/api/gen/go/runtime/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	tracev1 "github.com/origadmin/runtime/api/gen/go/runtime/trace/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/grpc/v1"
	httpv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/http/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

func DefaultApp() *appv1.App {
	return &appv1.App{
		Id:      "test-app-id",
		Name:    "test-app-name",
		Version: "test-app-version",
		Env:     "dev",
	}
}

func DefaultServers() *transportv1.Servers {
	return &transportv1.Servers{
		Servers: []*transportv1.Server{
			{
				Name:     "grpc_server",
				Protocol: "grpc",
				Grpc: &grpcv1.Server{
					Network: "tcp",
					Addr:    "0.0.0.0:50051",
				},
			},
			{
				Name:     "http_server",
				Protocol: "http",
				Http: &httpv1.Server{
					Network: "tcp",
					Addr:    "0.0.0.0:8080",
				},
			},
		},
	}
}

func DefaultDiscoveries() *discoveryv1.Discoveries {
	return &discoveryv1.Discoveries{
		Discoveries: []*discoveryv1.Discovery{},
	}
}

func DefaultData() *datav1.Data {
	return &datav1.Data{
		Databases: &datav1.Databases{
			Configs: []*storagev1.DatabaseConfig{
				{
					Name:    "default",
					Dialect: "sqlite3",
					Source:  "file:./test.db?cache=shared&mode=memory&_fk=1",
				},
			},
		},
		Caches: &datav1.Caches{
			Configs: []*storagev1.CacheConfig{
				{
					Name:   "default",
					Driver: "memory",
				},
			},
		},
	}
}

func DefaultMiddlewares() *middlewarev1.Middlewares {
	return &middlewarev1.Middlewares{
		Middlewares: []*middlewarev1.Middleware{},
	}
}

func DefaultTrace() *tracev1.Trace {
	return &tracev1.Trace{
		Name:        "jaeger",
		Endpoint:    "localhost:6831",
		ServiceName: "test-app-name",
	}
}
