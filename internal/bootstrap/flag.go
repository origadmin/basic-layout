package bootstrap

import (
	"os"
)

type Flags struct {
	ID         string
	Name       string
	Version    string
	EnvPath    string
	ConfigPath string
	MetaData   map[string]string
}

// IID returns the instance id
func (f *Flags) IID() string {
	return f.ID + "." + f.Name
}

func (f *Flags) Setup() {
	f.ID, _ = os.Hostname()
}

func DefaultFlags() Flags {
	id, _ := os.Hostname()
	return Flags{
		ID: id,
	}
}
