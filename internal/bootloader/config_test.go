package bootloader

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/runtime/config"
)

func TestNewFileSourceConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *config.SourceConfig
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				path: "resources/configs",
			},
			want: &config.SourceConfig{
				Type: "consul",
				File: &config.SourceConfig_File{
					Path: "resources/configs",
				},
				Consul: &config.SourceConfig_Consul{
					Address: "127.0.0.1:8500",
					Scheme:  "http",
				},
				Etcd: &config.SourceConfig_ETCD{
					Endpoints: []string{"127.0.0.1:2379"},
				},
				EnvArgs: map[string]string{
					"env": "dev",
				},
				EnvPrefixes: []string{"env"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &config.SourceConfig{
				Type: "consul",
				File: &config.SourceConfig_File{
					Path: "resources/configs",
				},
				Consul: &config.SourceConfig_Consul{
					Address: "127.0.0.1:8500",
					Scheme:  "http",
				},
				EnvArgs: map[string]string{
					"env": "dev",
				},
				EnvPrefixes: []string{"env"},
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileSourceConfig() = %v, want %v", got, tt.want)
			}
			wd, _ := os.Getwd()
			path := filepath.Join(wd, "../../resources/local2.toml")
			err := codec.EncodeToFile(path, got)
			if err != nil {
				t.Errorf("NewFileSourceConfig() = %v", err)
			}
		})
	}
}
