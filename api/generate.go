package v1

//generate helloworld proto file
//go:generate protoc -I. -I../third_party --go_out=. --go-http_out=. --go-grpc_out=. ./v1/protos/helloworld/*.proto
