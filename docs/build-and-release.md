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
