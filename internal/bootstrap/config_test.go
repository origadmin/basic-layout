package bootstrap

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/origadmin/toolkits/codec"
)

func TestNewFileSourceConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *SourceConfig
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				path: "resources/configs",
			},
			want: &SourceConfig{
				Type: "consul",
				File: &FileSource{
					Path: "resources/configs",
				},
				Consul: &ConsulSource{
					Address: "127.0.0.1:8500",
					Scheme:  "http",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &SourceConfig{
				Type: "consul",
				File: &FileSource{
					Path: "resources/configs",
				},
				Consul: &ConsulSource{
					Address: "127.0.0.1:8500",
					Scheme:  "http",
				},
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileSourceConfig() = %v, want %v", got, tt.want)
			}
			wd, _ := os.Getwd()
			path := filepath.Join(wd, "../../resources/local/config.toml")
			err := codec.EncodeToFile(path, got)
			if err != nil {
				t.Errorf("NewFileSourceConfig() = %v", err)
			}
		})
	}
}
