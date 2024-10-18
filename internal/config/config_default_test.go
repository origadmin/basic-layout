package config

import (
	"reflect"
	"testing"
)

func TestSaveConf(t *testing.T) {
	type args struct {
		path string
		conf *Bootstrap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				path: "test.toml",
				conf: DefaultBootstrap,
			},
		},
		{
			name: "test",
			args: args{
				path: "test.yml",
				conf: DefaultBootstrap,
			},
		},
		{
			name: "test",
			args: args{
				path: "test.json",
				conf: DefaultBootstrap,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveConf(tt.args.path, tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("SaveConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConf(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Bootstrap
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				path: "test.toml",
			},
			want:    DefaultBootstrap,
			wantErr: false,
		},
		{
			name: "test",
			args: args{
				path: "test.json",
			},
			want:    DefaultBootstrap,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConf(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
