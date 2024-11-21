package endpoint

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		scheme   string
		host     string
		addr     string
		endpoint string
		want     string
	}{
		{
			scheme:   "http",
			host:     "example.com",
			addr:     "127.0.0.1:8080",
			endpoint: "",
			want:     "http://example.com:8080",
		},
		{
			scheme:   "https",
			host:     "example.com",
			addr:     "127.0.0.1:8443",
			endpoint: "https://example.com:8080",
			want:     "https://example.com:8443",
		},
		{
			scheme:   "http",
			host:     "example.com",
			addr:     "127.0.0.1:80",
			endpoint: "example.com",
			want:     "http://example.com:80",
		},
		{
			scheme:   "http",
			host:     "example.com",
			addr:     "127.0.0.1:8081",
			endpoint: "http://example.com",
			want:     "http://example.com:8081",
		},
		{
			scheme:   "http",
			host:     "example.com",
			addr:     "127.0.0.1:8082",
			endpoint: "http://example.com:8080/path",
			want:     "http://example.com:8082",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Parse(tt.scheme, tt.host, tt.addr, tt.endpoint)
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
