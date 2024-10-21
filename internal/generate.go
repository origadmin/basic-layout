package internal

//generate config proto file
//go:generate protoc -I. -I../third_party -I../toolkits --go_out=paths=source_relative:./ --go-http_out=paths=source_relative:./ --go-grpc_out=paths=source_relative:./ --validate_out=lang=go:../  ./config/*.proto
