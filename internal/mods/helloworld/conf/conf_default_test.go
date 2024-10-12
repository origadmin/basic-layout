package conf

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
				conf: DefaultConf(),
			},
		},
		{
			name: "test",
			args: args{
				path: "test.json",
				conf: DefaultConf(),
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
			want:    DefaultConf(),
			wantErr: false,
		},
		{
			name: "test",
			args: args{
				path: "test.json",
			},
			want:    DefaultConf(),
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
