# RPC protocol usage

## Installation

* protoc **(protoc)**
  See https://grpc.io/docs/protoc-installation/

Otherwise, Can download [Protocol Buffers](https://github.com/protocolbuffers/protobuf/releases)
from https://github.com/protocolbuffers/protobuf/releases

```bash
> protoc --version
# libprotoc 27.2
```

* protoc-gen-go **(protoc-gen-go)**

```bash
> go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

> protoc-gen-go --version
# protoc-gen-go v1.34.2
```

* generate protobuf *

```bash
> protoc -I. --go_out=paths=source_relative:. rpc.proto
> protoc -I. --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. rpc.proto
> protoc -I. --go_out=paths=source_relative:. --go-errors_out=paths=source_relative:. rpc.proto
> protoc -I. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. rpc.proto
```