package internal

//generate config proto file
//go:generate protoc -I. -I../third_party --go_out=../ --go-http_out=../ --go-grpc_out=../ --validate_out=lang=go:../  ./mods/helloworld/conf/*.proto
