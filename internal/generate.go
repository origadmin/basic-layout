package api

//generate helloworld proto file
// if you want to generate the client code to the same directory, please use the following command
////go:generate kratos proto client -p=../toolkits -p=../third_party ./conf/*.proto

//=paths=source_relative:. outputs to the same directory with the proto file

// uncomment this line to generate the client code to the same directory
//go:generate protoc -I. -I../third_party -I../toolkits --go_out=paths=source_relative:../internal --validate_out=lang=go:../  ./configs/*.proto
