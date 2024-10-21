package conf

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/origadmin/toolkits/codec"
	"google.golang.org/protobuf/types/known/durationpb"
)

func SaveConf(path string, conf *Bootstrap) error {
	typo := codec.TypeFromExt(filepath.Ext(path))
	if typo == codec.UNKNOWN {
		return fmt.Errorf("unknown file type: %s", path)
	}
	err := codec.EncodeToFile(path, conf)
	if err != nil {
		return err
	}
	return nil
	//return os.WriteFile(path, marshal, os.ModePerm)
}

func LoadConf(path string) (*Bootstrap, error) {
	typo := codec.TypeFromExt(filepath.Ext(path))
	if typo == codec.UNKNOWN {
		return nil, fmt.Errorf("unknown file type: %s", path)
	}
	var b Bootstrap
	err := codec.DecodeFromFile(path, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

var DefaultBootstrap = &Bootstrap{
	ServiceName: "helloworld",
	Version:     "v1.0.0",
	CryptoType:  "argon2",
	Server: &Server{
		Gins: &Server_GINS{
			Network:         "tcp",
			Addr:            "${gins_address=0.0.0.0:7000}",
			UseTls:          false,
			CertFile:        "",
			KeyFile:         "",
			Timeout:         durationpb.New(3 * time.Minute),
			ShutdownTimeout: durationpb.New(3 * time.Minute),
			ReadTimeout:     durationpb.New(3 * time.Minute),
			WriteTimeout:    durationpb.New(3 * time.Minute),
			IdleTimeout:     durationpb.New(3 * time.Minute),
		},
		Http: &Server_HTTP{
			Network:         "tcp",
			Addr:            "${http_address=0.0.0.0:8000}",
			UseTls:          false,
			CertFile:        "",
			KeyFile:         "",
			Timeout:         durationpb.New(3 * time.Minute),
			ShutdownTimeout: durationpb.New(3 * time.Minute),
			ReadTimeout:     durationpb.New(3 * time.Minute),
			WriteTimeout:    durationpb.New(3 * time.Minute),
			IdleTimeout:     durationpb.New(3 * time.Minute),
		},
		Grpc: &Server_GRPC{
			Network:         "tcp",
			Addr:            "${grpc_address=0.0.0.0:9000}",
			UseTls:          false,
			CertFile:        "",
			KeyFile:         "",
			Timeout:         durationpb.New(3 * time.Minute),
			ShutdownTimeout: durationpb.New(3 * time.Minute),
			ReadTimeout:     durationpb.New(3 * time.Minute),
			WriteTimeout:    durationpb.New(3 * time.Minute),
			IdleTimeout:     durationpb.New(3 * time.Minute),
		},
		Middleware: &Server_Middleware{
			Cors: &Server_Middleware_Cors{
				Enabled:                false,
				AllowAllOrigins:        false,
				AllowOrigins:           nil,
				AllowMethods:           nil,
				AllowHeaders:           nil,
				AllowCredentials:       false,
				ExposeHeaders:          nil,
				MaxAge:                 0,
				AllowWildcard:          false,
				AllowBrowserExtensions: false,
				AllowWebSockets:        false,
				AllowFiles:             false,
			},
			Metrics: &Server_Middleware_Metrics{
				Enabled: false,
				Name:    "metrics",
			},
			Traces: &Server_Middleware_Traces{
				Enabled: false,
				Name:    "traces",
			},
			Logger: &Server_Middleware_Logger{
				Enabled: false,
				Name:    "logger",
			},
		},
	},
	Data: &Data{
		Database: &Data_Database{
			Driver: "mysql",
			Source: "dsn",
		},
		Redis: &Data_Redis{
			Network:      "tcp",
			Addr:         "${redis_address=127.0.0.1:6379}",
			ReadTimeout:  durationpb.New(3 * time.Minute),
			WriteTimeout: durationpb.New(3 * time.Minute),
		},
	},
	Discovery: &Discovery{
		Type: "${discovery_type:consul}",
		Consul: &Discovery_Consul{
			Address: "${consul_address=127.0.0.1:8500}",
			//Scheme:              "",
			//Token:               "",
			//Datacenter:          "",
			//Tag:                 "",
			//HealthCheckInterval: "",
			//HealthCheckTimeout:  "",
			HealthCheck: true,
			HeartBeat:   true,
		},
		Etcd: &Discovery_Etcd{
			Endpoints: "${etcd_address=127.0.0.1:2379}",
		},
	},
	Settings: &Settings{
		CryptoType: "argon2",
	},
}
