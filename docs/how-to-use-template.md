### 📖 How to Use These Templates

This guide explains how to use the `simple` and `multiple` layout templates to start your projects.

### Choosing the Right Starting Point

Before you begin, it's important to choose the template that best fits your needs.

*   **For a Single Microservice**: If you are building a single, independent service (e.g., a `user-service`), the **`simple`** layout is the recommended starting point. It provides a clean, minimal foundation for one service. To create a microservices architecture, simply create a new project from the `simple` template for each service.

*   **For a Monorepo**: If you prefer to manage multiple services within a single repository (a "monorepo"), the **`multiple`** layout is the right choice. It provides a structure for co-locating services like `user`, `order`, and `gateway` while sharing common configurations.

The following sections detail how to use the `multiple` template, including how to evolve it into separate microservices if your project's needs change over time.

---

### Path A: Creating a New Monorepo Project

Use this path if you want to start a new project that will be maintained within a single repository (monorepo) using the `multiple` template.

1.  **Copy the Template**: Copy the `examples/basic-layout/multiple/multiple_sample` directory and rename it to your project's name, e.g., `my-awesome-project`.

2.  **Initialize Your Module (Choose Your Naming Strategy)**: Navigate into your new project directory. How you name your module affects how it can be used.

    *   **Option 1: Public, Go-Gettable Module (Recommended)**
        If you plan to host your code on GitHub and want it to be a standard, fetchable Go module, use a full URL path.
        ```sh
        cd my-awesome-project
        go mod edit -module github.com/your-org/my-awesome-project
        ```

    *   **Option 2: Local-Only Module**
        If this is a private project and you don't intend for it to be fetched via `go get`, you can use a simple, non-URL name. This is exactly how the `basic-layout/multiple/multiple_sample` module itself is named.
        ```sh
        cd my-awesome-project
        go mod edit -module my-awesome-project
        ```

3.  **Synchronize Dependencies**: Run `go mod tidy`. This will download all required remote dependencies.
    ```sh
    go mod tidy
    ```

4.  **Update Import Paths**: Perform a global search-and-replace across your project to replace the old module prefix `basic-layout/multiple/multiple_sample` with the new one you chose in Step 2.

    **Important Note**: This not only includes Go file `import` paths but may also include references to the module path within configuration files such as `.goreleaser.yaml`, `buf.yaml`, and `resources/configs/*.yaml`. Please ensure a comprehensive replacement.

    *   **Windows (PowerShell) Example**:
        ```powershell
        Get-ChildItem -Path . -Recurse -Include *.go,*.yaml,*.yml,*.mod,*.toml | ForEach-Object {
            (Get-Content $_.FullName) | ForEach-Object {
                $_ -replace "basic-layout/multiple/multiple_sample", "github.com/your-org/my-awesome-project"
            } | Set-Content $_.FullName
        }
        ```
    *   **Linux/macOS (Bash) Example**:
        ```bash
        grep -rl "basic-layout/multiple/multiple_sample" . | xargs sed -i '' 's|basic-layout/multiple/multiple_sample|github.com/your-org/my-awesome-project|g'
        ```

5.  **Start Developing**: You can now customize the project by modifying or removing the example services (`user`, `order`, `gateway`).

---

### Path B: Advanced - Migrating a Monorepo to Independent Microservices

If you have started with a monorepo using the `multiple` template and now need to split it into separate, independent microservice repositories, this guide outlines the process. This is a common evolution for growing applications.

Let's demonstrate by splitting the `multiple_sample` monorepo into two separate projects: `user-service` and `gateway-service`.

##### Step 1: Create the `user-service`

1.  **Copy & Rename**: Copy `basic-layout/multiple/multiple_sample` to a new directory named `user-service`.
2.  **Prune the Project (Crucial: Prune Carefully)**: Delete all files and directories not related to the `user` service. To ensure correctness, it is recommended to adopt a 'keep what's needed, delete the rest' strategy:
    *   **Keep** `cmd/user/`
    *   **Keep** `internal/mods/user/`
    *   **Keep** `api/v1/proto/user/` (if `user.proto` is in this path)
    *   **Keep** `resources/configs/user/`
    *   **Delete** `cmd/gateway/` and `cmd/order/`.
    *   **Delete** `internal/mods/gateway/` and `internal/mods/order/`.
    *   **Delete** `api/v1/proto/gateway/` and `api/v1/proto/order/` (if these proto files are no longer needed).
    *   **Delete** `resources/configs/gateway/` and `resources/configs/order/`.
    *   **Check and delete** `.goreleaser.yaml` build configurations unrelated to the `user` service.
    *   **Check and delete** `buf.gen.yaml` and `buf.yaml` generation or module definitions unrelated to the `user` service.
    *   **Check and delete** `Makefile` commands unrelated to the `user` service.
    *   **Tip**: For larger projects, scripting this process can be considered, but manual inspection and deletion are crucial for ensuring a lean project.
3.  **Configure the Module**:
    *   `cd user-service`
    *   `go mod edit -module github.com/your-org/user-service`
    *   `go mod tidy`
    *   Update `.goreleaser.yaml` to only build the `user` binary.
    *   Globally replace `basic-layout/multiple/multiple_sample` with your new module name.

##### Step 2: Create the `gateway-service`

1.  **Copy & Rename**: Copy `basic-layout/multiple/multiple_sample` again to a new directory named `gateway-service`.
2.  **Prune the Project (Crucial: Prune Carefully)**: This time, delete business logic modules unrelated to the gateway. To ensure correctness, it is recommended to adopt a 'keep what's needed, delete the rest' strategy:
    *   **Keep** `cmd/gateway/`
    *   **Keep** `internal/mods/gateway/`
    *   **Keep** `api/v1/proto/gateway/` (if `gateway.proto` is in this path)
    *   **Keep** `resources/configs/gateway/`
    *   **Delete** `cmd/user/` and `cmd/order/`.
    *   **Delete** `internal/mods/user/` and `internal/mods/order/`.
    *   **Delete** `api/v1/proto/user/` and `api/v1/proto/order/` (if these proto files are no longer needed).
    *   **Delete** `resources/configs/user/` and `resources/configs/order/`.
    *   **Check and delete** `.goreleaser.yaml` build configurations unrelated to the `gateway` service.
    *   **Check and delete** `buf.gen.yaml` and `buf.yaml` generation or module definitions unrelated to the `gateway` service.
    *   **Check and delete** `Makefile` commands unrelated to the `gateway` service.
    *   **Tip**: For larger projects, scripting this process can be considered, but manual inspection and deletion are crucial for ensuring a lean project.
3.  **Configure the Module**:
    *   `cd gateway-service`
    *   `go mod edit -module github.com/your-org/gateway-service`
    *   `go mod tidy`
    *   Update `.goreleaser.yaml` to only build the `gateway` binary.
    *   Globally replace `basic-layout/multiple/multiple_sample` with your new module name.

##### Step 3: The Critical Change - Network Communication

Now that they are separate projects, the gateway can no longer call the `user` and `order` services in-process. It must call them over the network.

1.  **Configure the Clients**: In the `gateway-service` project, you must configure the gRPC clients to connect to the `user-service` and `order-service`. Modify the gateway's `resources/configs/gateway/conf.yaml` (or `bootstrap.yaml` if clients are configured there):

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

2.  **Update Dependency Injection**: The gateway's `wire.go` needs to be updated to use this new configuration to create the `user` and `order` clients, instead of relying on local providers. This typically involves creating `NewUserClient` and `NewOrderClient` functions that read from the configuration.

This process transforms the monorepo template into a distributed system, which is the natural evolution for many growing applications.
