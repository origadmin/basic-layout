package config

type Type string

const (
	Default    Type = "none"
	LocalFile  Type = "file"
	Nacos      Type = "nacos"
	Consul     Type = "consul"
	Etcd       Type = "etcd"
	Apollo     Type = "apollo"
	Kubernetes Type = "kubernetes"
	Polaris    Type = "polaris"
)
