package api

// for linux or macos you can use make generate
// for windows you can use make this generate

//generate helloworld proto file
// if you want to generate the client code to the same directory, please use the following command
////go:generate kratos proto client -p=../third_party .

//=paths=source_relative:. outputs to the same directory with the proto file
////go:generate protoc -I. -I../third_party --go_out=. --go-http_out=. --go-grpc_out=. --validate_out=lang=go:. --go-gins_out=. ./v1/protos/helloworld/*.proto

// generate *.pb.go
//go:generate protoc -I. -I../third_party --go_out=../api ./v1/protos/helloworld/*.proto
// generate *_http.pb.go
//go:generate protoc -I. -I../third_party --go-http_out=../api ./v1/protos/helloworld/*.proto
// generate *_grpc.pb.go
//go:generate protoc -I. -I../third_party --go-grpc_out=../api ./v1/protos/helloworld/*.proto
// generate *_gins.pb.go
//go:generate protoc -I. -I../third_party --go-gins_out=../api ./v1/protos/helloworld/*.proto
// generate *.pb.validate.go
//go:generate protoc -I. -I../third_party --validate_out=lang=go:../api ./v1/protos/helloworld/*.proto
//
//go:generate protoc -I. -I../third_party --openapi_out=naming=proto,fq_schema_naming=true,default_response=false:../api ./v1/protos/helloworld/*.proto
// generate a greeter server template
////go:generate kratos proto server -t ../internal/service v1/protos/helloworld/greeter.proto
