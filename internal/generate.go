package api

//generate helloworld proto file
// if you want to generate the client code to the same directory, please use the following command
////go:generate kratos proto client -p=../toolkits -p=../third_party ./conf/*.proto

//=paths=source_relative:. outputs to the same directory with the proto file
//go:generate protoc -I. -I../third_party -I../toolkits --go_out=../ --go-http_out=../ --go-grpc_out=../ --validate_out=lang=go:../  ./conf/*.proto
