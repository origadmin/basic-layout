package v1

//generate helloworld proto file
//=paths=source_relative:. outputs to the same directory with the proto file
//go:generate protoc -I. -I../third_party --go_out=. --go-http_out=. --go-grpc_out=. --validate_out=lang=go:. --go-gin_out=. ./v1/protos/helloworld/*.proto
