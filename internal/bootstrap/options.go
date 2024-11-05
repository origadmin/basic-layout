package bootstrap

type Type string

const (
	Default    Type = "none"
	File       Type = "file"
	Nacos      Type = "nacos"
	Consul     Type = "consul"
	Etcd       Type = "etcd"
	Apollo     Type = "apollo"
	Kubernetes Type = "kubernetes"
	Polaris    Type = "polaris"
)

func (t Type) String() string {
	return string(t)
}
