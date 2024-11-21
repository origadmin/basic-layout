package bootstrap

import (
	"time"

	configv1 "github.com/origadmin/toolkits/runtime/gen/go/config/v1"
	"google.golang.org/protobuf/types/known/durationpb"

	"origadmin/basic-layout/internal/configs"
)

func DefaultConfig() *configs.Bootstrap {
	return &configs.Bootstrap{
		ServiceName: "origadmin.service.v1.demo",
		Version:     "v1.0.0",
		CryptoType:  "argon2",
		Entry: &configs.Bootstrap_Entry{
			Network: "tcp",
			Addr:    "0.0.0.0:8000",
			Timeout: durationpb.New(3 * time.Minute),
		},
		//Service: &config.Service{
		//	Gins: &config.Service_GINS{
		//		Network:         "tcp",
		//		Addr:            "${gins_address:0.0.0.0:8100}",
		//		UseTls:          false,
		//		CertFile:        "",
		//		KeyFile:         "",
		//		Timeout:         durationpb.New(3 * time.Minute),
		//		ShutdownTimeout: durationpb.New(3 * time.Minute),
		//		ReadTimeout:     durationpb.New(3 * time.Minute),
		//		WriteTimeout:    durationpb.New(3 * time.Minute),
		//		IdleTimeout:     durationpb.New(3 * time.Minute),
		//	},
		//	Http: &config.Service_HTTP{
		//		Network:         "tcp",
		//		Addr:            "${http_address:0.0.0.0:8200}",
		//		UseTls:          false,
		//		CertFile:        "",
		//		KeyFile:         "",
		//		Timeout:         durationpb.New(3 * time.Minute),
		//		ShutdownTimeout: durationpb.New(3 * time.Minute),
		//		ReadTimeout:     durationpb.New(3 * time.Minute),
		//		WriteTimeout:    durationpb.New(3 * time.Minute),
		//		IdleTimeout:     durationpb.New(3 * time.Minute),
		//	},
		//	Grpc: &config.Service_GRPC{
		//		Network:         "tcp",
		//		Addr:            "${grpc_address:0.0.0.0:8300}",
		//		UseTls:          false,
		//		CertFile:        "",
		//		KeyFile:         "",
		//		Timeout:         durationpb.New(3 * time.Minute),
		//		ShutdownTimeout: durationpb.New(3 * time.Minute),
		//		ReadTimeout:     durationpb.New(3 * time.Minute),
		//		WriteTimeout:    durationpb.New(3 * time.Minute),
		//		IdleTimeout:     durationpb.New(3 * time.Minute),
		//	},
		//	Host: "${host:127.0.0.1}",
		//},
		//Data: &config.Data{},
	}
}

func DefaultServiceGins() *configv1.Service_GINS {
	return &configv1.Service_GINS{
		Network:         "tcp",
		Addr:            "${gins_address:0.0.0.0:8100}",
		UseTls:          false,
		CertFile:        "",
		KeyFile:         "",
		Timeout:         durationpb.New(3 * time.Minute),
		ShutdownTimeout: durationpb.New(3 * time.Minute),
		ReadTimeout:     durationpb.New(3 * time.Minute),
		WriteTimeout:    durationpb.New(3 * time.Minute),
		IdleTimeout:     durationpb.New(3 * time.Minute),
	}
}
