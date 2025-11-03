# Makefile for the Basic Layout Example Project

# ============================================================================ #
#                              CONFIGURATION
# ============================================================================ #

# --------------------------- Basic Configuration ---------------------------- #
GOHOSTOS         ?= $(shell go env GOHOSTOS)
ENV              ?= dev
PROJECT_ORG      := OrigAdmin
THIRD_PARTY_PATH := third_party
BUILT_BY         := $(PROJECT_ORG)

# ---------------------------- Git Information ----------------------------- #
# Use /bin/bash for non-Windows, powershell for Windows for consistent command execution
ifeq ($(GOHOSTOS), windows)
    SHELL          := powershell.exe
    .SHELLFLAGS    := -NoProfile -Command
    # Git commands for Windows (PowerShell)
    GIT_COMMIT     := $(shell git rev-parse --short HEAD)
    GIT_BRANCH     := $(shell git rev-parse --abbrev-ref HEAD)
    GIT_VERSION    := $(shell git describe --tags --always)
    BUILD_DATE     := $(shell Get-Date -Format 'yyyy-MM-ddTHH:mm:ssK')
    GIT_TREE_STATE := $(shell if ((git status --porcelain)) { 'dirty' } else { 'clean' })
else
    SHELL          := /bin/bash
    # Git commands for Linux/macOS (Bash)
    GIT_COMMIT     := $(shell git rev-parse --short HEAD)
    GIT_BRANCH     := $(shell git rev-parse --abbrev-ref HEAD)
    GIT_VERSION    := $(shell git describe --tags --always)
    BUILD_DATE     := $(shell TZ=Asia/Shanghai date +%FT%T%z)
    GIT_TREE_STATE := $(if $(shell git status --porcelain),dirty,clean)
endif

# Append -dirty suffix if the working directory is not clean
ifneq ($(GIT_TREE_STATE), clean)
    GIT_VERSION := $(GIT_VERSION)-dirty
endif

# ----------------------------- Build Flags ------------------------------ #
# LDFLAGS for passing version information to main packages
# Note: The main package needs to define variables like main.version, main.commit, etc.
LDFLAGS := -X main.version=$(GIT_VERSION) \
           -X main.commit=$(GIT_COMMIT) \
           -X main.date=$(BUILD_DATE) \
           -X main.treeState=$(GIT_TREE_STATE) \
           -X main.builtBy=$(BUILT_BY)

ifeq ($(ENV), release)
    LDFLAGS += -s -w # Strip symbols and omit DWARF symbol table for release builds
endif

# ------------------------ Protobuf Configuration ------------------------ #
PROTO_API_PATH     := api
OPENAPI_DOCS_PATH  := resources/docs/openapi

# ============================================================================ #
#                           LIFECYCLE TARGETS
# ============================================================================ #

.PHONY: all init deps generate build test lint clean release run-helloworld run-secondworld run-gateway help

all: init deps generate build ## ✅ Run the full build process (init, deps, generate, build)

init: ## 🔧 Install necessary Go tools for development
	@echo "Installing Go tools..."
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	@go install github.com/google/wire/cmd/wire@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@echo "Running go mod tidy to ensure tool dependencies are in go.mod..."
	@go mod tidy

deps: ## 📦 Update and export third-party protobuf dependencies
	@echo "Updating buf dependencies..."
	@buf dep update
	@echo "Exporting protobuf dependencies to $(THIRD_PARTY_PATH)..."
	@buf export buf.build/bufbuild/protovalidate -o $(THIRD_PARTY_PATH)
	@buf export buf.build/protocolbuffers/wellknowntypes -o $(THIRD_PARTY_PATH)
	@buf export buf.build/googleapis/googleapis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/gnostic/gnostic -o $(THIRD_PARTY_PATH)
	@buf export buf.build/kratos/apis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/origadmin/runtime -o $(THIRD_PARTY_PATH)

generate: ## 🧬 Generate all code from Protobuf and run go generate (includes wire, openapi)
	@echo "Generating API Protobuf code using buf..."
	@buf generate
	@echo "Generating Internal Configs Protobuf code..."
	@protoc -I. -I./third_party --go_out=paths=source_relative:. --validate_out=paths=source_relative,lang=go:. internal/configs/*.proto
	@echo "Generating Wire code for dependency injection..."
	@cd cmd/helloworld && wire
	@cd cmd/secondworld && wire
	@cd cmd/gateway && wire
	@echo "Running go generate for other code generation tasks..."
	@go generate ./...
	@echo "Running go mod tidy after code generation..."
	@go mod tidy

build: ## 🔨 Build all services (helloworld, secondworld, gateway)
	@echo "Building helloworld service..."
	@go build -ldflags "$(LDFLAGS)" -o ./dist/helloworld ./cmd/helloworld
	@echo "Building secondworld service..."
	@go build -ldflags "$(LDFLAGS)" -o ./dist/secondworld ./cmd/secondworld
	@echo "Building gateway service..."
	@go build -ldflags "$(LDFLAGS)" -o ./dist/gateway ./cmd/gateway

test: ## 🧪 Run all Go tests
	@echo "Running all Go tests..."
	@go test ./...

lint: ## 🔍 Run golangci-lint (install if not present)
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ./...

clean: ## 🧹 Clean up generated files and build artifacts
	@echo "Cleaning up generated files and build artifacts..."
	@rm -rf ./api/v1/gen
	@rm -rf ./internal/mods/*/gen
	@rm -rf ./dist
	@rm -rf ./third_party
	@rm -f ./go.sum # Clean go.sum to force re-download if needed

release: ## 🚀 Run GoReleaser to create releases
	@echo "Running GoReleaser..."
	@goreleaser release --clean

run-helloworld: ## ▶️ Run the helloworld service
	@echo "Running helloworld service..."
	@go run ./cmd/helloworld/ -conf ./resources/configs/helloworld/bootstrap.yaml

run-secondworld: ## ▶️ Run the secondworld service
	@echo "Running secondworld service..."
	@go run ./cmd/secondworld/ -conf ./resources/configs/secondworld/bootstrap.yaml

run-gateway: ## ▶️ Run the gateway service
	@echo "Running gateway service..."
	@go run ./cmd/gateway/ -conf ./resources/configs/gateway/bootstrap.yaml

# ============================================================================ #
#                                     HELP
# ============================================================================ #

.PHONY: help
help: ## ✨ Show this help message
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Common Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  \033[36m%-22s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
