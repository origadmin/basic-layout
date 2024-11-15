/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"reflect"
	"testing"

	"github.com/origadmin/toolkits/codec"
	"origadmin/basic-layout/internal/configs"
)

func TestSaveConf(t *testing.T) {
	type args struct {
		path string
		conf *configs.Bootstrap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "toml",
			args: args{
				path: "test.toml",
				conf: DefaultConfig(),
			},
		},
		{
			name: "yml",
			args: args{
				path: "test.yml",
				conf: DefaultConfig(),
			},
		},
		{
			name: "json",
			args: args{
				path: "test.json",
				conf: DefaultConfig(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveConf(tt.args.path, tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("SaveConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		//opt := protojson.MarshalOptions{
		//	EmitUnpopulated: true,
		//	Indent:          " ",
		//}
		//bs, _ := opt.Marshal(DefaultConfig())
		//_ = os.WriteFile("test.json", bs, os.ModePerm)
	}
}

func SaveConf(path string, conf *configs.Bootstrap) error {
	err := codec.EncodeToFile(path, conf)
	if err != nil {
		return err
	}
	return nil
}

func TestLoadConf(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *configs.Bootstrap
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				path: "test.toml",
			},
			want:    DefaultConfig(),
			wantErr: false,
		},
		{
			name: "test",
			args: args{
				path: "test.json",
			},
			want:    DefaultConfig(),
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

func LoadConf(path string) (*configs.Bootstrap, error) {
	var conf *configs.Bootstrap
	err := codec.DecodeFromFile(path, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
