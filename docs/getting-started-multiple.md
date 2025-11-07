### 🚀 Getting Started

#### 1. Prerequisites

Ensure you have the following tools installed:

* Go (>= 1.24)
* [Buf](https://buf.build/docs/installation)
* [Protobuf (`protoc`)](https://grpc.io/docs/protoc-installation/)
* Go-related `protoc` plugins (installed via `go install`):
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
  go install github.com/envoyproxy/protoc-gen-validate@latest
  go install github.com/google/wire/cmd/wire@latest
  go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
  ```

#### 2. Generate Code

From your project directory, run the `make` command to generate all gRPC/HTTP code from the `.proto` files, and also run
`wire` for dependency injection.

```sh
make generate
```

#### 3. Run the Services

You can run each service in a separate terminal. Ensure a service discovery agent (like Consul) is running if you use
`discovery://` endpoints.

**Local Service Discovery Agent Setup (Optional, Recommended for Multi-Service Testing):**

If you use `discovery://` endpoints, a service discovery agent is required. Below is a simple example of starting Consul
using Docker Compose:

1. Create a `docker-compose.yaml` file (e.g., in the project root directory):
   ```yaml
   version: '3.8'
   services:
     consul:
       image: consul:1.10.0 # You can use the latest stable version
       container_name: consul
       ports:
         - "8500:8500" # UI Port
         - "8600:8600/udp" # DNS Port
       command: "agent -server -bootstrap-expect=1 -client=0.0.0.0 -ui -node=consul-server-1"
       healthcheck:
         test: ["CMD", "consul", "members"]
         interval: 10s
         timeout: 5s
         retries: 3
   ```
2. Run in the directory containing `docker-compose.yaml`:
   ```bash
   docker-compose up -d
   ```
   This will start a local Consul service discovery agent.

```sh
# Terminal 1: Start the user service
go run ./cmd/user/ -conf ./resources/configs/user/bootstrap.yaml

# Terminal 2: Start the order service
go run ./cmd/order/ -conf ./resources/configs/order/bootstrap.yaml

# Terminal 3: Start the gateway service
go run ./cmd/gateway/ -conf ./resources/configs/gateway/bootstrap.yaml
```

#### 4. Test the API

Send requests to the `user` and `order` services through the `gateway`.

```sh
# Example: Get a user through the gateway (Transparent Proxy Pattern)
curl http://localhost:8000/api/v1/proxy/user/123

# Example: Get an order through the gateway (Transparent Proxy Pattern)
curl http://localhost:8000/api/v1/proxy/order/456

# Example: Get a user through the gateway (Edge Gateway Pattern)
curl http://localhost:8000/api/v1/edge/user/123

# Example: Get an order through the gateway (Edge Gateway Pattern)
curl http://localhost:8000/api/v1/edge/order/456

# Example: Directly call the user service (if exposed)
curl http://localhost:9001/api/v1/user/123

# Example: Directly call the order service (if exposed)
curl http://localhost:9002/api/v1/order/456
```

*(Note: The actual API paths might differ based on your protobuf definitions and Kratos HTTP rule configurations. Adjust
the curl commands accordingly.)*
