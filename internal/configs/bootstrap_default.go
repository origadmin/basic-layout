package configs

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/runtime/config"
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
	var bs Bootstrap
	err := codec.DecodeFromFile(path, &bs)
	if err != nil {
		return nil, err
	}
	return &bs, nil
}

func DefaultBootstrap() *Bootstrap {
	return &Bootstrap{
		ServiceName: "origadmin.service.v1.demo",
		Version:     "v1.0.0",
		CryptoType:  "argon2",
		Service: &config.ServiceConfig{
			Entry: &config.ServiceConfig_Entry{
				Network: "tcp",
				Addr:    "0.0.0.0:8000",
				Timeout: durationpb.New(3 * time.Minute),
			},
			Gins: &config.ServiceConfig_GINS{
				Network:         "tcp",
				Addr:            "${gins_address=0.0.0.0:8100}",
				UseTls:          false,
				CertFile:        "",
				KeyFile:         "",
				Timeout:         durationpb.New(3 * time.Minute),
				ShutdownTimeout: durationpb.New(3 * time.Minute),
				ReadTimeout:     durationpb.New(3 * time.Minute),
				WriteTimeout:    durationpb.New(3 * time.Minute),
				IdleTimeout:     durationpb.New(3 * time.Minute),
			},
			Http: &config.ServiceConfig_HTTP{
				Network:         "tcp",
				Addr:            "${http_address=0.0.0.0:8200}",
				UseTls:          false,
				CertFile:        "",
				KeyFile:         "",
				Timeout:         durationpb.New(3 * time.Minute),
				ShutdownTimeout: durationpb.New(3 * time.Minute),
				ReadTimeout:     durationpb.New(3 * time.Minute),
				WriteTimeout:    durationpb.New(3 * time.Minute),
				IdleTimeout:     durationpb.New(3 * time.Minute),
			},
			Grpc: &config.ServiceConfig_GRPC{
				Network:         "tcp",
				Addr:            "${grpc_address=0.0.0.0:8300}",
				UseTls:          false,
				CertFile:        "",
				KeyFile:         "",
				Timeout:         durationpb.New(3 * time.Minute),
				ShutdownTimeout: durationpb.New(3 * time.Minute),
				ReadTimeout:     durationpb.New(3 * time.Minute),
				WriteTimeout:    durationpb.New(3 * time.Minute),
				IdleTimeout:     durationpb.New(3 * time.Minute),
			},
			Host: "${host=127.0.0.1}",
		},
		Data: &config.DataConfig{},
		Source: &config.SourceConfig{
			Type: "file",
			Consul: &config.SourceConfig_Consul{
				Address: "${consul_address=127.0.0.1:8500}",
				Scheme:  "http",
			},
		},
		Settings: &Settings{
			CryptoType: "argon2",
		},
	}
}
