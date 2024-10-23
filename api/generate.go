package api

//generate helloworld proto file
// if you want to generate the client code to the same directory, please use the following command
////go:generate kratos proto client -p=../third_party .

//=paths=source_relative:. outputs to the same directory with the proto file
////go:generate protoc -I. -I../third_party --go_out=. --go-http_out=. --go-grpc_out=. --validate_out=lang=go:. --go-gin_out=. ./v1/protos/helloworld/*.proto

// uncomment the following line to generate the client code to the same directory
////go:generate protoc -I. -I../third_party --go_out=../api ./v1/protos/helloworld/*.proto
////go:generate protoc -I. -I../third_party --go-http_out=../api ./v1/protos/helloworld/*.proto
////go:generate protoc -I. -I../third_party --go-grpc_out=../api ./v1/protos/helloworld/*.proto
//go:generate protoc -I. -I../third_party --go-gin_out=../api ./v1/protos/helloworld/*.proto
////go:generate protoc -I. -I../third_party --validate_out=lang=go:../api ./v1/protos/helloworld/*.proto
//
////go:generate protoc -I. -I../third_party --openapi_out=fq_schema_naming=true,default_response=false:. ./v1/protos/helloworld/*.proto
//
////go:generate kratos proto server -t ../internal/service v1/protos/helloworld/greeter.proto
