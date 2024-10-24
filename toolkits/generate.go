package toolkits

// uncomment this
////go:generate kratos proto client .

// generate protobuf without openapi.yaml
//go:generate protoc --proto_path=. --proto_path=../third_party --validate_out=lang=go,paths=source_relative:. --go_out=paths=source_relative:. ./middlewares/cors/*.proto
//go:generate protoc --proto_path=. --proto_path=../third_party --validate_out=lang=go,paths=source_relative:. --go_out=paths=source_relative:. ./middlewares/logger/*.proto
//go:generate protoc --proto_path=. --proto_path=../third_party --validate_out=lang=go,paths=source_relative:. --go_out=paths=source_relative:. ./middlewares/metrics/*.proto
//go:generate protoc --proto_path=. --proto_path=../third_party --validate_out=lang=go,paths=source_relative:. --go_out=paths=source_relative:. ./middlewares/security/*.proto
