package service

import "os"

type AppInfo interface {
	ID() string
	InstanceID() string
	Name() string
	Version() string
	Metadata() map[string]string
	SetID(string) AppInfo
	SetName(string) AppInfo
	SetVersion(string) AppInfo
	SetMetadata(map[string]string) AppInfo
}

type appInfo struct {
	id       string
	name     string
	version  string
	metadata map[string]string
}

func (s *appInfo) ID() string {
	return s.id
}

func (s *appInfo) InstanceID() string {
	return s.id + "." + s.name
}

func (s *appInfo) SetID(id string) AppInfo {
	s.id = id
	return s
}

func (s *appInfo) Name() string {
	return s.name
}

func (s *appInfo) SetName(name string) AppInfo {
	s.name = name
	return s
}

func (s *appInfo) Version() string {
	return s.version
}

func (s *appInfo) SetVersion(version string) AppInfo {
	s.version = version
	return s
}

func (s *appInfo) Metadata() map[string]string {
	return s.metadata
}

func (s *appInfo) SetMetadata(metadata map[string]string) AppInfo {
	s.metadata = metadata
	return s
}

func NewAppInfo(id, name, version string) AppInfo {
	if id == "" {
		id, _ = os.Hostname()
	}
	return &appInfo{
		id:       id,
		name:     name,
		version:  version,
		metadata: map[string]string{},
	}
}
