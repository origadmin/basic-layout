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

* **On Organizing Business Modules in a Monorepo**

  In a multi-module project (monorepo), structuring the internal code is critical for clarity and scalability. The goal is to clearly separate independent business modules from each other and from shared infrastructure code.

  **Recommended Structure: Separating Features from Infrastructure**

  The recommended approach is to separate components based on their architectural role. This leads to a more intuitive and scalable structure. The following is a conceptual illustration of this principle:

  ```
  internal/
  ├── features/          # <-- Houses the core business features/domains.
  │   ├── order/         # <-- The "order" business feature.
  │   └── user/          # <-- The "user" business feature.
  ├── gateway/           # <-- The API Gateway, as a distinct infrastructure/integration layer.
  └── helpers/           # <-- Shared helper utilities for the entire project.
  ```

  > For a complete and concrete example of this structure in practice, please refer to the project layout in the **[Multiple Module Application](./multiple-module-application.md)** document.

  **Why this structure is recommended:**
    *   **Clarity of Responsibility**: It makes the architectural intent clear. `features` contains the application's core business value, while `gateway` is a supporting infrastructure component responsible for routing, aggregation, and access control.
    *   **Eliminates Naming Conflicts**: We no longer need to find a single ambiguous term to describe both business logic (`user`) and infrastructure (`gateway`). The directory name `features` accurately describes its contents.
    *   **Scalability**: When adding a new business capability (e.g., `payment`), you simply add a new directory under `features`. If you add a new cross-cutting infrastructure component (e.g., an `eventbus`), it can be added at the top level of `internal`, parallel to `gateway`.

  **Alternative Naming Schemes for the Business Module Directory:**

  While `features` is the recommended default, other names are also viable depending on the project's philosophy:

    *   **`components`**: A great, neutral, and widely understood term. It refers to the sub-directories as independent parts of a larger system. This is an excellent alternative to `features` and avoids any potential ambiguity.
    *   **`domains`**: This is the preferred choice if your project strictly follows Domain-Driven Design (DDD). It signals that each subdirectory is a self-contained business domain. However, it may be less appropriate for non-domain components like a `gateway`.

  **A Note on Flatter Structures:**

  For very simple projects with only two or three modules, you might consider a flatter structure by placing all modules directly under `internal/`. While this reduces directory depth, it is **not recommended for most projects** because it mixes components with different architectural roles (business vs. infrastructure) at the same level. As the project grows, this can lead to a cluttered and less organized `internal` directory.

* **Project-Level vs. Feature-Level Configuration**

  In a monorepo, it's crucial to distinguish between shared project-level configurations and configurations that are specific to a single feature.

    *   **Project-Level Configuration (`internal/conf`)**: This directory is for configuration logic that is shared across the entire project. This includes the main configuration adapter that bridges raw config values (from YAML) to the runtime framework's expected interfaces, as well as Protobuf definitions for shared settings (e.g., logging, tracing).

    *   **Feature-Level Configuration (`internal/features/<feature_name>/conf`)**: When a specific feature (e.g., `user`) requires its own unique configuration that other features do not need to know about, it should be defined within that feature's own directory. This promotes encapsulation and high cohesion, keeping the feature self-contained.

* **On Structuring the Data Layer (with `ent`)**

  A well-organized `data` layer is crucial for maintainability, especially when using code generation tools like `ent`. The following structure is recommended for each feature:

  ```
  internal/features/user/
  └── data/
      ├── data.go      # 1. Data layer's public entrypoint and Wire ProviderSet.
      ├── user.go      # 2. Implementation of the biz.UserRepo interface.
      └── ent/         # 3. Code generated by ent; treated as a third-party library.
          ├── schema/
          └── ...
  ```

  **Key Responsibilities:**
    1.  **`data.go`**: This file acts as the public API of the data layer for that feature. It should contain the `NewData` function (to initialize database clients) and the `ProviderSet` for `google/wire`, which groups all repository implementations (`NewUserRepo`, etc.) for dependency injection.
    2.  **`user.go` (Repository Implementation)**: This file implements the `biz.UserRepo` interface. It contains the actual database logic, calling `ent` functions and performing the crucial conversion from `ent` models (Persistence Objects, POs) to `biz` models (Domain Objects, DOs).
    3.  **`ent/`**: This directory is managed entirely by the `ent` tool. It should not be manually edited. The rest of the application should not import from this directory directly; only the repository implementations (like `user.go`) are allowed to.

  **Advanced: Customizing `ent` with Templates**

  For advanced use cases, such as generating custom repository code, `ent` templates are invaluable. These templates are developer tools, not runtime code, and should be placed accordingly.

  **Recommended Structure for `ent` Templates:**
  ```
  <project_root>/
  └── tools/
      └── entc/              # Directory for custom ent code generation.
          ├── entc.go        # Go script to configure and run the code generation.
          └── templates/     # Directory for all custom .tpl template files.
              └── repository.tpl
  ```

  **How to Trigger Generation:**

  Instead of using fragile relative paths in `go:generate` comments, the best practice is to use the Go module path to run the `entc.go` script. Add the following comment to a file within the feature's `data` package (e.g., `internal/features/user/data/data.go`):

  ```go
  //go:generate go run <your_module_path>/tools/entc
  ```

  Replace `<your_module_path>` with your project's actual module path from `go.mod` (e.g., `basic-layout/multiple/multiple_sample`). This creates a robust, location-independent way to trigger your custom code generation.

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

  The term DTO is best used to describe the objects used for transferring data across process or network boundaries. In this architecture, the Protobuf-generated request/response structs (e.g., `simplev1.SayHelloRequest`) in the `service` layer are perfect examples of **API DTOs**. The `service` layer's responsibility includes converting these API DTOs into the `biz` layer's DOs.

  However, a dedicated `dto/` directory (e.g., `internal/features/user/dto/`) can also be used for **Internal DTOs** or **View Models**. These are data structures used for:

    *   Transferring data between different layers *within* the application.
    *   Aggregating data from multiple DOs into a single response structure.
    *   Representing specific views or projections of data that don't directly map to a single DO or API DTO.

  While **DTO** is the industry-standard term, you might occasionally see them referred to as **API Models**, **Payloads**, or **Request/Response Models**. The key is not the name, but its role: a data structure dedicated to data transfer, separate from the internal business logic models.

* **On Gateway Design Patterns**

  The API Gateway is the front door to your system. Its design has a major impact on scalability and maintainability. Here are three common patterns, each with its own trade-offs:

    1. **Edge Gateway Pattern**: The gateway defines its own, independent API contract and is responsible for transforming data before forwarding requests to internal features.
    2. **Transparent Proxy Pattern**: The gateway directly reuses the API contracts from internal features, forwarding requests with minimal transformation.
    3. **gRPC-Gateway (Reverse Proxy) Pattern**: This pattern uses a plugin to auto-generate a reverse proxy from `google.api.http` annotations in the `.proto` files of the downstream features.

  **Recommendation**:
    * Start with the **gRPC-Gateway Pattern** for simplicity.
    * Evolve to the **Edge Gateway Pattern** when you need a stable, aggregated public API.
    * Use the **Transparent Proxy Pattern** sparingly for internal-only gateways.

* **On the Role of `service` and `server` Layers**

  These two layers handle the application's external-facing concerns:

    * **Service Layer (`service`)**: The **API Implementation Layer**. It translates request DTOs into the `biz` layer's DOs, calls the business logic, and translates the results back into response DTOs. It should be kept "thin".

    * **Server Layer (`server`)**: The **Transport Layer**. It manages the lifecycle of network servers (gRPC, HTTP), attaches middleware, and registers the `service` implementations.

* **Scaling for Complexity: Additional Layers**

  As a project grows, you can introduce more layers:

    * **Application Layer (`application`)**: Placed between `service` and `biz`, it orchestrates complex use cases involving multiple domains, transactions, and events.

    * **Infrastructure Layer (`infrastructure`)**: An explicit layer containing all concrete implementations for interacting with the outside world (databases, message queues, etc.), implementing interfaces defined by the `biz`/`domain` layer.

* **Where to Place Custom Implementations, Tools, or Helpers**

  When adding your own custom code, it's important to place it in a location that aligns with the project's architectural principles:

    1. **Feature-Specific Helpers (`internal/features/<feature_name>/helpers/`)**: For helper functions or tools that are specific to a particular business feature (e.g., `user`, `order`).

    2. **Project-Level Helpers (`internal/helpers/`)**: For helper functions or tools that are used across multiple features within the project.

    3. **Standalone Tools (`tools/`)**: For standalone development tools, code generators, or scripts that are not part of the application's runtime code.

    4. **Directly within `biz/` or `data/`**: If the custom code is tightly coupled with the business logic or data access layer of a specific feature.

  **Why use `internal/`?**
  The `internal` directory enforces encapsulation in Go. Any package within an `internal` directory can only be imported by code within the `internal` directory's direct parent or its ancestors. This prevents external projects from accidentally depending on your internal implementation details, promoting clear API boundaries and reducing future refactoring risks.
