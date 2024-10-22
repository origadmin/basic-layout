package toolkits

// uncomment this
////go:generate kratos proto client .

////go:generate protoc --proto_path=. --proto_path=../third_party --go_out=paths=source_relative:./ --go-http_out=paths=source_relative:./ --go-grpc_out=paths=source_relative:./ ./middlewares/cors/*.proto
////go:generate protoc --proto_path=. --proto_path=../third_party --go_out=paths=source_relative:./ --go-http_out=paths=source_relative:./ --go-grpc_out=paths=source_relative:./ ./middlewares/logger/*.proto
////go:generate protoc --proto_path=. --proto_path=../third_party --go_out=paths=source_relative:./ --go-http_out=paths=source_relative:./ --go-grpc_out=paths=source_relative:./ ./middlewares/metrics/*.proto
////go:generate protoc --proto_path=. --proto_path=../third_party --go_out=paths=source_relative:./ --go-http_out=paths=source_relative:./ --go-grpc_out=paths=source_relative:./ ./middlewares/security/*.proto
