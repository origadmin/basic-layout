//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject && GOWORK
// +build !wireinject,GOWORK

// The build tag makes sure the stub is not built in the final build.
package main
