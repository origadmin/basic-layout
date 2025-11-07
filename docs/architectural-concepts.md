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

*   **DTO (Data Transfer Object)**



  The term DTO is best used to describe the objects used for transferring data across process or network boundaries. In

  this architecture, the Protobuf-generated request/response structs (e.g., `simplev1.SayHelloRequest`) in the `service`

  layer are perfect examples of **API DTOs**. The `service` layer's responsibility includes converting these API DTOs

  into the `biz` layer's DOs.



  However, a dedicated `dto/` directory (e.g., `internal/mods/user/dto/`) can also be used for **Internal DTOs** or

  **View Models**. These are data structures used for:

    *   Transferring data between different layers *within* the application, especially if the Domain Objects (DOs) are

      too rich or the API DTOs are too specific for internal use.

    *   Aggregating data from multiple DOs into a single response structure before converting to an API DTO.

    *   Representing specific views or projections of data that don't directly map to a single DO or API DTO.



  The `data` layer also plays a crucial role in DTO-like conversions, specifically mapping **Persistence Objects (POs)**

  (e.g., `ent` models) to Domain Objects (DOs) and vice-versa.



  While **DTO** is the industry-standard term, you might occasionally see them referred to as **API Models**, **Payloads

  **, or **Request/Response Models**. The key is not the name, but its role: a data structure dedicated to data transfer,

  separate from the internal business logic models.

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

* **Where to Place Custom Implementations, Tools, or Helpers**

  When adding your own custom code, it's important to place it in a location that aligns with Go's best practices and
  the project's architectural principles, especially regarding encapsulation.

    1. **Module-Specific Helpers (`internal/mods/<module_name>/helpers/`)**:
        * **Purpose**: For helper functions or tools that are specific to a particular service module (e.g., `user`,
          `order`, `gateway`) and are not part of its core business logic (`biz`) or data access (`data`).
        * **Example**: Data validation utilities, specific data transformation functions, or small, reusable components
          used only within that module.
        * **Rationale**: Keeps module-specific utilities organized and prevents them from polluting the global namespace
          or being accidentally imported by other modules.

    2. **Project-Level Helpers (`internal/helpers/`)**:
        * **Purpose**: For helper functions or tools that are used across multiple service modules within the
          `multiple_sample` project.
        * **Example**: Custom error handling, common logging utilities, or shared data formatting functions.
        * **Rationale**: Centralizes reusable helper code that is internal to the entire project, ensuring it's not
          exposed externally.

    3. **Standalone Tools (`tools/`)**:
        * **Purpose**: For standalone development tools, code generators, or scripts that are used for project-related
          tasks but are not part of the application's runtime code.
        * **Example**: Custom code generation scripts, database migration tools, or build automation scripts.
        * **Rationale**: Separates development-time utilities from the application's core source code, making it clear
          they are not meant to be compiled into the final binary.

    4. **Directly within `biz/` or `data/`**:
        * **Purpose**: If the custom code is tightly coupled with the business logic or data access layer of a specific
          module.
        * **Example**: A custom validation rule for a domain model would go in `internal/mods/<module_name>/biz/`. A
          specific database query builder would go in `internal/mods/<module_name>/data/`.
        * **Rationale**: Maintains high cohesion and keeps related code together. Ensure that `biz` remains pure and
          free of external dependencies.

  **Why use `internal/`?**
  The `internal` directory enforces encapsulation in Go. Any package within an `internal` directory can only be imported
  by code within the `internal` directory's direct parent or its ancestors. This prevents external projects from
  accidentally depending on your internal implementation details, promoting clear API boundaries and reducing future
  refactoring risks.

* **On Naming Protobuf Configuration Directories**

  The naming of directories containing Protobuf definitions for configurations can sometimes be ambiguous, especially
  when actual configuration value files (like YAML) exist elsewhere. Here's a breakdown of options and recommendations,
  distinguishing between internal and public API configurations:

  **1. Internal Configuration Protobuf Definitions**

  These are Protobuf definitions (`.proto` files) that define the schema for configurations used *internally* within
  your module and are not intended to be part of your module's public API.

  **Recommended Placement and Naming for Internal Configurations:**

  To ensure proper encapsulation, clarity, and adherence to Go best practices, internal configuration Protobuf definitions
  and their generated Go code should always be placed under the `internal/` directory. The specific subdirectory depends
  on whether a dedicated conversion layer is needed.

    *   **Scenario 1: Default/Direct Usage (No explicit conversion layer needed)**
        This is the primary recommendation when the Protobuf-generated configuration types can be used directly by
        `internal/conf` or other internal packages without significant adaptation.

        *   **Strongly Recommended Structure**: `internal/confpb/`
            *   **Purpose**: This is the preferred location for internal `.proto` files (e.g., `bootstrap.proto`) and their
              generated Go code (e.g., `bootstrap.pb.go`). The `pb` suffix clearly indicates Protobuf content, and
              `internal/` enforces encapsulation.
            *   **Go Package Name**: Typically `confpb`.
            *   **Role of `internal/conf/`**: In this scenario, `internal/conf/` (package `conf`) would import directly
              from `internal/confpb/` and use the `confpb` types for loading and handling configuration.
            *   **Benefits**: Simple, direct, and maintains clear separation of concerns.

    *   **Scenario 2: With Conversion Layer (Explicit adaptation to runtime interfaces needed)**
        This scenario applies when `internal/conf` (or a similar package) needs to perform significant transformations
        or adaptations from the raw Protobuf types to application-specific runtime interfaces.

        *   **Recommended Structure**: `internal/conf/pb/`
            *   **Purpose**: This subdirectory is the preferred location for internal `.proto` files and their generated
              Go code when `internal/conf` acts as a dedicated conversion layer. Nesting `pb` under `internal/conf/`
              consolidates all configuration-related artifacts for this specific conversion process.
            *   **Go Package Name**: Typically `confpb` (e.g., `option go_package = "your_module_path/internal/conf/pb;confpb";`).
            *   **Role of `internal/conf/`**: In this scenario, `internal/conf/` (package `conf`) would import from
              `internal/conf/pb/` and perform the necessary transformations or adaptations from the `conf_pb` types to
              application-specific runtime interfaces. This aligns perfectly with the "Protobuf + Conversion" pattern.
            *   **Benefits**: Consolidates configuration logic for conversion, enforces encapsulation, and avoids redundant
              top-level `config` directories.

    *   **Less Recommended Naming Conventions**
        *   `schema/config/` or `proto/config/`: These are generally less recommended due to deeper directory nesting and
          potential naming conflicts with other project modules or tools that might use `schema` or `proto` for different
          purposes.
        
**2. Public API Configuration Protobuf Definitions**

These are Protobuf definitions that define configuration schemas intended to be part of your module's *public API*,
meaning other services or clients might directly depend on them. These are typically managed by tools like `buf`.

    * **Recommended Placement**: Within your project's designated public API Protobuf root, e.g.,
      `<YOUR_API_PROTO_ROOT>/config/` or `<YOUR_API_PROTO_ROOT>/configpb/`.
        * The exact path `<YOUR_API_PROTO_ROOT>` (e.g., `api/v1/proto`) is project-specific and should be consistently
          applied across all public API Protobuf definitions.
        * **Benefits**: Allows for unified management and publishing of public API schemas.
