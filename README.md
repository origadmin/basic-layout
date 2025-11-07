# Basic Layout Examples

This directory contains example projects demonstrating different application layouts using the OrigAdmin framework.
These examples serve as practical guides and starting points for building your own applications.

## Projects

1. [**simple/simple_app**](./simple/simple_app): A standard single-module application.
2. [**multiple/multiple_sample**](./multiple/multiple_sample): A multi-module application within a single project
   boundary, featuring `user`, `order`, and `gateway` services.

---

## 1. Simple Application (`simple/simple_app`)

This project demonstrates a classic, layered architecture for a single, self-contained microservice. It's an ideal
starting point for most new services.

### Directory Structure

The structure is organized based on the "Clean Architecture" and "Dependency Inversion" principles to ensure high
cohesion and low coupling.

```
simple_app
├── api/            # Protobuf definitions for the application's API.
├── cmd/            # Main application entry point, responsible for initialization and startup.
├── configs/        # Protobuf definitions for internal configurations (e.g., bootstrap.proto) and their generated Go code.
├── internal/       # All private application code. Go prevents other projects from importing this.
│   ├── biz/        # Business Logic Layer: Defines business models and use cases (interfaces and structs).
│   ├── conf/       # Configuration structures and loading logic for internal use.
│   ├── data/       # Data Access Layer: Implements the interfaces defined in `biz`. Handles all interactions with databases, caches, etc.
│   ├── server/     # Server Layer: Initializes and configures transport servers (e.g., HTTP, gRPC).
│   └── service/    # Service Layer: Implements the API services defined in Protobuf. Acts as a bridge between the transport layer and the business logic layer.
├── resources/      # Contains various assets and non-code resources used by the application.
│   ├── configs/    # Application configuration files (e.g., bootstrap.yaml, conf.yaml) for each service.
│   └── api-docs/   # Generated OpenAPI documentation.
├── Makefile        # Provides common development commands (build, generate, run, etc.).
└── go.mod          # Go module definition.
```

### Key Architectural Concepts

* **Dependency Rule**: Dependencies flow inwards. `service` depends on `biz`, and `data` depends on `biz`. The core
  business logic in `biz` has no external dependencies.
* **Interface Segregation**: Interfaces (e.g., `SimpleRepo`) are defined in the `biz` layer, representing the contracts
  that the business logic needs.
* **Dependency Injection**: We use `google/wire` for dependency injection. The dependency graph is defined in
  `cmd/wire.go` and the generated code is in `cmd/wire_gen.go`.

### Architectural Philosophy & Best Practices

This section provides deeper insights into some of the design choices and suggests best practices for extending this
layout.

* **On Naming the Business Logic Layer (`biz`)**

  While this example uses `biz` (for "business"), you might encounter other names for the core logic layer in different
  projects. The choice often reflects a specific architectural philosophy:
    * **`biz`**: A common and straightforward name, clearly indicating it contains business logic. It's a great
      general-purpose choice.
    * **`domain`**: This name is often used in projects following **Domain-Driven Design (DDD)**. It implies a richer
      model, including Entities, Value Objects, Aggregates, and Domain Events. Use this when your business logic is
      complex and you want to model the domain explicitly.
    * **`usecase`**: This name comes from **Clean Architecture** and emphasizes the application's specific user-facing
      actions or "use cases" (e.g., `CreateUserUseCase`). It makes the application's capabilities very explicit.

* **On Naming the Data Access Layer (`data`)**

  Similar to the business layer, the data access layer also has several common naming conventions, although we recommend
  `data` for consistency within this framework.
    * **`data`**: The default choice in these examples. It's a general and widely understood term for the layer that
      handles data persistence and retrieval.
    * **`dal`** (Data Access Layer): A more traditional and explicit acronym that clearly states the layer's purpose.
    * **`repository`**: A great choice when strictly following Domain-Driven Design (DDD). This name emphasizes that the
      layer's primary role is to provide concrete implementations of the repository interfaces defined in the `biz`/
      `domain` layer.
    * **`persistence`**: This name strongly focuses on the storage aspect, making it a good fit when the layer's
      responsibility is strictly limited to database interactions.

  Ultimately, the name is less important than its role: to implement the data-handling interfaces required by the
  business layer and to isolate the core logic from the details of data storage.

* **Separating Data Models (DO vs. PO)**

  In this simple example, the `biz` layer's model (`biz.Simple`) and the `data` layer's model (e.g., an `ent.Simple`
  struct) might look identical. For simplicity, one might be tempted to use the data layer's model directly in the
  business layer. However, for more complex applications, it's a crucial best practice to keep them separate.

    * **Persistence Object (PO)**: This is the model in the `data` layer (e.g., `ent.Simple`). Its structure mirrors the
      database table and may contain database-specific tags or fields (`id`, `created_at`, ORM tags). Its sole purpose
      is data persistence.
    * **Domain Object (DO)**: This is the model in the `biz` layer (e.g., `biz.Simple`). It represents a concept in your
      business domain and should be "pure," containing only fields and methods relevant to the business logic.

  **Why separate them?**
    1. **Decoupling**: The `biz` layer remains completely independent of the database schema. You can change your
       database tables without affecting your core business logic.
    2. **Clarity**: Each model has a single, clear responsibility.
    3. **Flexibility**: A single business operation in the `biz` layer might need to compose data from multiple database
       tables. Having a distinct DO makes this aggregation clean and straightforward.

  To implement this, you would create conversion functions within the `data` layer to map between POs and DOs.

* **DTO (Data Transfer Object)**

  The term DTO is best used to describe the objects used for transferring data across process or network boundaries. In
  this architecture, the Protobuf-generated request/response structs (e.g., `simplev1.SayHelloRequest`) in the `service`
  layer are perfect examples of DTOs. The `service` layer's responsibility includes converting these DTOs into the `biz`
  layer's DOs.

  While **DTO** is the industry-standard term, you might occasionally see them referred to as **API Models**, **Payloads
  **, or **Request/Response Models**. The key is not the name, but its role: a data structure dedicated to API
  communication, separate from the internal business logic models.

* **On Gateway Design Patterns**

  The API Gateway is the front door to your system. Its design has a major impact on scalability and maintainability.
  Here are three common patterns, each with its own trade-offs:

    1. **Edge Gateway Pattern**:
        * **Concept**: The gateway defines its own, independent API contract (`.proto` file). It acts as a "boundary"
          that protects internal services from external changes.
        * **Implementation**: The gateway's `service` layer is responsible for converting the gateway's DTOs to the
          internal modules' DTOs before forwarding the request.
        * **Pros**: **Maximum Decoupling**. Provides a stable API for external clients. Allows for API aggregation
          and裁剪. Internal services can evolve freely without breaking external clients.
        * **Cons**: Higher implementation overhead due to the need for manual data transformation code.

    2. **Transparent Proxy Pattern**:
        * **Concept**: The gateway's `.proto` file directly imports and reuses the API contracts from internal modules.
        * **Implementation**: The gateway's `service` layer can directly forward requests with minimal or no data
          transformation.
        * **Pros**: **Simple and Fast to Implement**. Low code overhead.
        * **Cons**: **Tight Coupling**. Changes in internal module APIs immediately break the gateway's public API.
          Hinders independent evolution and versioning.

    3. **gRPC-Gateway (Reverse Proxy) Pattern**:
        * **Concept**: This pattern eliminates the need for a manually written gateway `service` or even a
          `gateway.proto` file for routing. It uses the `grpc-gateway` protoc plugin to auto-generate a reverse proxy.
        * **Implementation**: HTTP routing rules are defined directly within the downstream services' `.proto` files (
          e.g., `user.proto`, `order.proto`) using `google.api.http` annotations. The gateway application becomes a thin
          host that simply runs the generated proxy server and points it to the downstream gRPC services.
        * **Pros**: **Single Source of Truth**. The API contract (gRPC and REST) is defined in one place. The gateway is
          truly "unaware" of business logic, making it extremely low-maintenance.
        * **Cons**: Less flexibility for complex API aggregation or transformation logic at the gateway layer itself.
          All logic is pushed to the downstream services.

  **Recommendation**:
    * Start with the **gRPC-Gateway Pattern** for most use cases due to its simplicity and low maintenance.
    * Evolve to the **Edge Gateway Pattern** when you need to provide a stable, aggregated, or secured public API that
      differs significantly from your internal service APIs.
    * Use the **Transparent Proxy Pattern** sparingly, perhaps for internal-only gateways where coupling is less of a
      concern.

* **On the Role of `service` and `server` Layers**

  These two layers handle the application's external-facing concerns, completing the separation of business logic from
  transport details.

    * **Service Layer (`service`)**: This is the **API Implementation Layer**. It acts as the crucial bridge between the
      network and your business logic.
        * **Role**: To implement the service interfaces defined in your `.proto` files.
        * **Responsibilities**: It receives request DTOs (e.g., `user.CreateUserRequest`), converts them into the `biz`
          layer's DOs, calls the appropriate business use case, and then converts the resulting DOs back into response
          DTOs.
        * **Core Principle**: This layer should be kept "thin". Its primary job is translation and delegation, not
          implementing business rules.

    * **Server Layer (`server`)**: This is the **Transport Layer**, where the actual network servers (gRPC, HTTP) are
      configured and run.
        * **Role**: To manage the lifecycle of transport servers.
        * **Responsibilities**: It initializes servers, attaches middleware chains (like recovery or logging), and, most
          importantly, **registers** the `service` layer implementations onto the server instances (e.g.,
          `user.RegisterUserAPIServer(grpcSrv, userService)`).
        * **Core Principle**: This layer is concerned with the "how" of communication (protocols, ports, timeouts), not
          the "what" (the business operations).

* **Scaling for Complexity: Additional Layers**

  The current `biz`/`data`/`service` layout is an excellent foundation. As a project's complexity grows, you can
  introduce more layers to maintain clarity and separation of concerns.

    * **Application Layer (`application`)**: In very complex systems, you might introduce an `application` layer between
      `service` and `biz`.
        * **Role**: It acts as an orchestrator for complex use cases. A single `application` service call might
          coordinate multiple `biz`/`domain` objects, manage database transactions, handle authorization, and dispatch
          events.
        * **Benefit**: It keeps the `service` layer thin (only handling DTO conversion and transport concerns) and the
          `biz`/`domain` layer pure (only containing core business rules), while centralizing complex workflow logic.

    * **Infrastructure Layer (`infrastructure`)**: The current `data` layer is a form of infrastructure. In a larger
      project, you can make this more explicit by creating a dedicated `infrastructure` layer.
        * **Role**: This layer contains all the concrete implementations for interacting with the outside world (
          databases, message queues, caches, third-party APIs, etc.). It implements the interfaces defined by the `biz`/
          `domain` layer.
        * **Benefit**: It groups all external dependencies and their "glue code" in one place, making it clear what
          external systems the application depends on.

  A more complex project structure might look like this:

  ```
  internal/
  ├── application/    # Orchestrates use cases, handles transactions.
  ├── domain/         # The core business logic (formerly `biz`).
  ├── infrastructure/ # Implementations for external systems.
  │   ├── persistence/  # Database implementations (formerly `data`).
  │   ├── messaging/    # Message queue clients (e.g., Kafka, RabbitMQ).
  │   └── cache/        # Cache implementations (e.g., Redis).
  ├── service/        # API service implementations (thin layer).
  └── ...
  ```

### How to Run

Navigate to the `simple/simple_app` directory and use the `Makefile` commands:

```bash
# Generate code (protobuf, wire, etc.)
make generate

# Build the application binary
make build

# Run the application
make run
```

---

## 2. Multiple Module Application (`multiple/multiple_sample`)

A multi-module backend project example based on the [Kratos](https://go-kratos.dev/) framework, demonstrating how to
build and manage multiple services within a single Go module.

This project serves as a template for creating a standalone Go application with multiple interconnected services. It
depends on several public `origadmin` packages (like `runtime` and `toolkits`), which are managed as standard Go module
dependencies.

### ✨ Features

* **Kratos v2**: Built upon the latest Kratos framework.
* **Standard Go Modules**: Manages all dependencies through `go.mod`.
* **Multi-Module Monorepo**: Contains three independent services (`user`, `order`, `gateway`) within a single repository
  structure.
* **Microservice Ready**: Provides a clear path for splitting the monorepo into multiple independent microservice
  repositories.
* **API Gateway**: The `gateway` service acts as a unified entry point, routing requests to `user` and `order` services.
* **Code Generation & Publishing**: Integrated with `buf` and `GoReleaser`.
* **Makefile**: Provides convenient `make` commands for development workflows.

### 📂 Project Structure

```
.
├── api/                # Protobuf API definitions for all services (user, order, gateway)
├── cmd/                # Service entrypoints (main.go for user, order, gateway)
├── dist/               # Compiled service binaries and release artifacts (generated by 'make build' or 'goreleaser').
├── internal/           # Internal business logic, organized by module
│   └── mods/
│       ├── user/       # User service implementation
│       ├── order/      # Order service implementation
│       └── gateway/    # API Gateway implementation
├── resources/          # Contains various assets and non-code resources used by the application.
│   ├── configs/        # Application configuration files (e.g., bootstrap.yaml, conf.yaml) for each service.
│   │   ├── user/       # Specific configs for user service
│   │   ├── order/      # Specific configs for order service
│   │   ├── gateway/    # Specific configs for gateway service
│   │   ├── bootstrap.yaml.example  # Example bootstrap config for new projects
│   │   └── conf.yaml.example       # Example main config for new projects
│   ├── api-docs/       # Project-specific documentation, such as OpenAPI specifications or Buf configurations.
│   ├── env/            # Environment-specific configuration files or templates.
│   │   └── .env.example            # Example .env file for environment variable overrides
│   ├── release/        # Release-related scripts or templates.
│   ├── statics/        # (Optional) Static files like images, CSS, JavaScript for web interfaces.
│   ├── templates/      # (Optional) Template files for HTML, emails, etc.
│   ├── tests/          # (Optional) Test data or test-specific configurations.
│   ├── bin/            # (Optional) Pre-compiled binaries or helper executables.
│   └── data/           # (Optional) Data files, such as database migration scripts or seed data.
├── .goreleaser.yaml    # Automated release configuration
├── buf.gen.yaml        # Buf code generation configuration
├── buf.yaml            # Buf module definition
├── go.mod              # Go module file for the project
└── Makefile            # Developer helper commands
```

### 📖 How to Use This Template: From Monolith to Microservices

This guide provides two main paths for using this template.

#### Path A: Create a New Single-Repo Project

Use this path if you want to start a new project that will be maintained within a single repository.

1. **Copy the Template**: Copy the `examples/basic-layout/multiple/multiple_sample` directory and rename it to your
   project's name, e.g., `my-awesome-project`.

2. **Initialize Your Module (Choose Your Naming Strategy)**: Navigate into your new project directory. How you name your
   module affects how it can be used.

    * **Option 1: Public, Go-Gettable Module (Recommended)**
      If you plan to host your code on GitHub and want it to be a standard, fetchable Go module, use a full URL path.
      ```sh
      cd my-awesome-project
      go mod edit -module github.com/your-org/my-awesome-project
      ```

    * **Option 2: Local-Only Module**
      If this is a private project and you don't intend for it to be fetched via `go get`, you can use a simple, non-URL
      name. This is exactly how the `basic-layout/multiple/multiple_sample` module itself is named.
      ```sh
      cd my-awesome-project
      go mod edit -module my-awesome-project
      ```

3. **Synchronize Dependencies**: Run `go mod tidy`. This will download all required remote dependencies.
   ```sh
   go mod tidy
   ```

4. **Update Import Paths**: Perform a global search-and-replace across your project to replace the old module prefix
   `basic-layout/multiple/multiple_sample` with the new one you chose in Step 2.

5. **Start Developing**: You can now customize the project by modifying or removing the example services (`user`,
   `order`, `gateway`).

#### Path B: Split into Multiple Independent Microservice Repos

This is the advanced path for building a true microservices architecture where each service lives in its own repository.
Let's demonstrate by splitting `multiple_sample` into two separate projects: `user-service` and `gateway-service`.

##### Step 1: Create the `user-service`

1. **Copy & Rename**: Copy `basic-layout/multiple/multiple_sample` to a new directory named `user-service`.
2. **Prune the Project**: Delete all files and directories not related to the `user` service.
    * Delete `cmd/gateway/` and `cmd/order/`.
    * Delete `internal/mods/gateway/` and `internal/mods/order/`.
    * Delete `api/v1/proto/gateway/` and `api/v1/proto/order/`.
    * Delete `resources/configs/gateway/` and `resources/configs/order/`.
3. **Configure the Module**:
    * `cd user-service`
    * `go mod edit -module github.com/your-org/user-service`
    * `go mod tidy`
    * Update `.goreleaser.yaml` to only build the `user` binary.
    * Globally replace `basic-layout/multiple/multiple_sample` with your new module name.

##### Step 2: Create the `gateway-service`

1. **Copy & Rename**: Copy `basic-layout/multiple/multiple_sample` again to a new directory named `gateway-service`.
2. **Prune the Project**: This time, delete the business logic modules not related to the gateway.
    * Delete `cmd/user/` and `cmd/order/`.
    * Delete `internal/mods/user/` and `internal/mods/order/`.
    * Keep `api/v1/proto/gateway/`, but you may need to adjust its imports if the proto files for `user` and `order` are
      now in separate repositories.
    * Delete `resources/configs/user/` and `resources/configs/order/`.
3. **Configure the Module**:
    * `cd gateway-service`
    * `go mod edit -module github.com/your-org/gateway-service`
    * `go mod tidy`
    * Update `.goreleaser.yaml` to only build the `gateway` binary.
    * Globally replace `basic-layout/multiple/multiple_sample` with your new module name.

##### Step 3: The Critical Change - Network Communication

Now that they are separate projects, the gateway can no longer call the `user` and `order` services in-process. It must
call them over the network.

1. **Configure the Clients**: In the `gateway-service` project, you must configure the gRPC clients to connect to the
   `user-service` and `order-service`. Modify the gateway's `resources/configs/gateway/conf.yaml` (or `bootstrap.yaml`
   if clients are configured there):

   ```yaml
   server:
     # ... (server config) ...
   client: # Assuming client configurations are under a 'client' key
     # Add a configuration block for the user gRPC client
     user:
       protocol: grpc
       grpc:
         # The endpoint where user-service is running
         endpoint: "discovery:///user-service" # Using service discovery
         # Or a direct address for local testing:
         # endpoint: "localhost:9001"
         timeout: 2s
     # Add a configuration block for the order gRPC client
     order:
       protocol: grpc
       grpc:
         # The endpoint where order-service is running
         endpoint: "discovery:///order-service" # Using service discovery
         # Or a direct address for local testing:
         # endpoint: "localhost:9002"
         timeout: 2s
   ```

2. **Update Dependency Injection**: The gateway's `wire.go` needs to be updated to use this new configuration to create
   the `user` and `order` clients, instead of relying on local providers. This typically involves creating
   `NewUserClient` and `NewOrderClient` functions that read from the configuration.

This process transforms the monorepo template into a distributed system, which is the natural evolution for many growing
applications.

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
# Example: Get a user through the gateway
curl http://localhost:8000/api/v1/gateway/user/123

# Example: Get an order through the gateway
curl http://localhost:8000/api/v1/gateway/order/456
```

*(Note: The actual API paths might differ based on your protobuf definitions and Kratos HTTP rule configurations. Adjust
the curl commands accordingly.)*

### 📦 Build and Release

#### Using Makefile

The `Makefile` provides several useful commands:

* `make generate`: Generate code from Protobuf definitions and run `wire`.
* `make build`: Compile all services into the `dist/` directory.
* `make test`: Run all tests.
* `make lint`: Run the linter.

#### Using GoReleaser

This project uses `GoReleaser` for automated releases. The configuration is in `.goreleaser.yaml`.

To perform a local test release (this will not publish to GitHub):

```sh
# This will create binaries and archives in a 'dist' directory
goreleaser release --snapshot --clean
```

When a new version tag (e.g., `v1.2.0`) is pushed to the `main` branch of the repository, a GitHub Action workflow will
automatically trigger GoReleaser to build and create a new GitHub Release with all the compiled assets.
