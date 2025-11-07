# 1. Simple Application (`simple/simple_app`)

This project demonstrates a classic, layered architecture for a single, self-contained microservice. It's an ideal
starting point for most new services.

### Directory Structure

The structure is organized based on the "Clean Architecture" and "Dependency Inversion" principles to ensure high
cohesion and low coupling.

```
simple_app
├── api/            # Protobuf definitions for the application's API.
├── cmd/            # Main application entry point, responsible for initialization and startup.
├── internal/       # All private application code. Go prevents other projects from importing this.
│   ├── biz/        # Business Logic Layer: Defines business models and use cases (interfaces and structs).
│   ├── conf/       # Configuration structures and loading logic for internal use.
│   │   └── pb/     # Protobuf definitions for internal configurations (e.g., bootstrap.proto) and their generated Go code.
│   ├── data/       # Data Access Layer: Implements the interfaces defined in `biz`. Handles all interactions with databases, caches, etc.
│   ├── server/     # Server Layer: Initializes and configures transport servers (e.g., HTTP, gRPC).
│   └── service/    # Service Layer: Implements the API services defined in Protobuf. Acts as a bridge between the transport layer and the business logic layer.
├── resources/      # Contains various assets and non-code resources used by the application.
│   ├── configs/    # Application configuration files (e.g., bootstrap.yaml, conf.yaml) for each service.
│   └── api-docs/   # Generated OpenAPI documentation.
├── Makefile        # Provides common development commands (build, generate, run, etc.).
├── go.mod          # Go module definition.
└── go.sum          # Go module checksums.
```
