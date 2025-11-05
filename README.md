# Basic Layout Examples

This directory contains example projects demonstrating different application layouts using the OrigAdmin framework.
These examples serve as practical guides and starting points for building your own applications.

## Projects

1. [**simple/simple_app**](./simple/simple_app): A standard single-module application.
2. [**multiple/multiple_sample**](./multiple/multiple_sample): A multi-module application within a single project
   boundary.

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
├── configs/        # Default configuration files (e.g., bootstrap.yaml).
├── internal/       # All private application code. Go prevents other projects from importing this.
│   ├── biz/        # Business Logic Layer: Defines business models and use cases (interfaces and structs).
│   ├── data/       # Data Access Layer: Implements the interfaces defined in `biz`. Handles all interactions with databases, caches, etc.
│   ├── server/     # Server Layer: Initializes and configures transport servers (e.g., HTTP, gRPC).
│   └── service/    # Service Layer: Implements the API services defined in Protobuf. Acts as a bridge between the transport layer and the business logic layer.
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

This project demonstrates a more complex scenario where multiple logical modules or services coexist within a single
application. This layout is useful for:

* Monolithic applications with clearly separated internal domains.
* Services that expose multiple, distinct APIs (e.g., an admin API and a public API).

### Directory Structure

The key difference is the introduction of a `mods` or `modules` directory within `internal`, where each sub-directory
represents a distinct functional module. Each module can have its own `biz`, `data`, and `service` layers.

```
multiple_sample
└── internal/
    └── mods/
        ├── user/         # User module
        │   ├── biz/
        │   ├── data/
        │   └── service/
        └── order/        # Order module
            ├── biz/
            ├── data/
            └── service/
```

This structure helps manage complexity by keeping the code for different business domains isolated, while still allowing
them to be compiled and deployed as a single unit.