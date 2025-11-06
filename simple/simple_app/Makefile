# Makefile for the Basic Layout Example Project

# ============================================================================ #
#                              CONFIGURATION
# ============================================================================ #

# --------------------------- Basic Configuration ---------------------------- #
GOHOSTOS         ?= $(shell go env GOHOSTOS)
ENV              ?= dev
PROJECT_ORG      := OrigAdmin
APP_NAME         := simple_app
THIRD_PARTY_PATH := third_party
BUILT_BY         := $(PROJECT_ORG)

# ---------------------------- Git Information ----------------------------- #
GIT_COMMIT      := $(shell git rev-parse --short HEAD)
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD)

ifeq ($(GOHOSTOS), windows)
    SHELL          := powershell.exe
    .SHELLFLAGS    := -NoProfile -Command
    CMD_SEP        := ;
    GOWORKOFF      := $$env:GOWORK='off';
    RM_RF          := Remove-Item -Recurse -Force
    RM_F           := Remove-Item -Force
    CMD_EXISTS     := Get-Command -ErrorAction SilentlyContinue
    BUILD_OUTPUT_NAME := $(APP_NAME).exe
    GIT_HEAD_TAG   := $(shell git tag --points-at HEAD 2>$$null)
    BUILD_DATE     := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssK'")
    GIT_TREE_STATE := $(shell powershell -Command "if ((git status --porcelain)) { 'dirty' } else { 'clean' }")
else
    SHELL          := /bin/bash
    CMD_SEP        := &&
    GOWORKOFF      := GOWORK=off
    RM_RF          := rm -rf
    RM_F           := rm -f
    CMD_EXISTS     := command -v
    BUILD_OUTPUT_NAME := $(APP_NAME)
    GIT_HEAD_TAG   := $(shell git tag --points-at HEAD 2>/dev/null)
    BUILD_DATE     := $(shell TZ=Asia/Shanghai date +%FT%T%z)
    GIT_TREE_STATE := $(if $(shell git status --porcelain),dirty,clean)
endif

# Define GIT_TAG based on GIT_HEAD_TAG or GIT_COMMIT
ifeq ($(GOHOSTOS), windows)
    GIT_TAG        := $(shell powershell -Command "if ('${GIT_HEAD_TAG}') { '${GIT_HEAD_TAG}' } else { '${GIT_COMMIT}' }")
else
    GIT_TAG        := $(if $(GIT_HEAD_TAG),$(GIT_HEAD_TAG),$(GIT_COMMIT))
endif

# Define GIT_VERSION based on git describe and GIT_TREE_STATE
GIT_VERSION     := $(shell git describe --tags --always)
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

.PHONY: all init deps generate build test lint clean release run help

all: init deps generate build ## ✅ Run the full build process (init, deps, generate, build)

init: ## 🔧 Install necessary Go tools for development
	@echo "Installing Go tools..."
ifeq ($(GOHOSTOS), windows)
	@if (! ($(CMD_EXISTS) protoc)) { Write-Host "ERROR: protoc compiler not found. Please install it manually before proceeding."; Write-Host "  - For Windows: Download from https://github.com/protocolbuffers/protobuf/releases and add to PATH."; exit 1; }
else
	@if ! ($(CMD_EXISTS) protoc &> /dev/null); then \
		echo "ERROR: protoc compiler not found. Please install it manually before proceeding."; \
		echo "  - For macOS: brew install protobuf"; \
		echo "  - For Ubuntu/Debian: sudo apt-get install protobuf-compiler"; \
		exit 1; \
	fi
endif
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	@go install github.com/google/wire/cmd/wire@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@echo "Running go mod tidy to ensure tool dependencies are in go.mod..."
	@go mod tidy

deps: clean-deps ## 📦 Update and export third-party protobuf dependencies
	@echo "Updating buf dependencies..."
	@buf dep update
	@echo "Exporting protobuf dependencies to $(THIRD_PARTY_PATH)..."
	@buf export buf.build/bufbuild/protovalidate -o $(THIRD_PARTY_PATH)
	@buf export buf.build/protocolbuffers/wellknowntypes -o $(THIRD_PARTY_PATH)
	@buf export buf.build/googleapis/googleapis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/gnostic/gnostic -o $(THIRD_PARTY_PATH)
	@buf export buf.build/kratos/apis -o $(THIRD_PARTY_PATH)
	@buf export buf.build/origadmin/runtime -o $(THIRD_PARTY_PATH)


clean-deps: ## 🧹 Clean up third-party protobuf dependencies
	@echo "Cleaning up third-party protobuf dependencies..."
	@$(RM_RF) $(THIRD_PARTY_PATH)

generate: ## 🧬 Generate all code from Protobuf and run go generate (includes wire, openapi)
	@echo "Generating API Protobuf code using buf..."
	@buf generate
	@echo "Generating Internal Configs Protobuf code..."
	@protoc -I. -I./third_party --go_out=paths=source_relative:. --validate_out=paths=source_relative,lang=go:. configs/*.proto
	@echo "Running go generate for all packages (includes wire, etc.)..."
	@$(GOWORKOFF) go generate ./...
	@echo "Running go mod tidy after code generation..."
	@$(GOWORKOFF) go mod tidy

build: ## 🔨 Build the application
	@echo "Building $(APP_NAME) application..."
	@go build -ldflags "$(LDFLAGS)" -o ./dist/$(BUILD_OUTPUT_NAME) ./cmd/

test: ## 🧪 Run all Go tests
	@echo "Running all Go tests..."
	@go test ./...

lint: ## 🔍 Run golangci-lint (install if not present)
	@echo "Running golangci-lint..."
ifeq ($(GOHOSTOS), windows)
	@if (! ($(CMD_EXISTS) golangci-lint)) { Write-Host "golangci-lint not found, installing..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; } 
else
	@if ! ($(CMD_EXISTS) golangci-lint &> /dev/null); then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
endif
	@golangci-lint run ./...

clean: ## 🧹 Clean up generated files and build artifacts
	@echo "Cleaning up generated files and build artifacts..."
	@$(RM_RF) ./api/v1/gen
	@$(RM_RF) ./internal/mods/*/gen
	@$(RM_RF) ./dist
	@$(RM_RF) ./third_party
	@$(RM_F) ./go.sum # Clean go.sum to force re-download if needed

release: ## 🚀 Run GoReleaser to create releases
	@echo "Running GoReleaser..."
	@goreleaser release --clean

run: ## ▶️ Run the application
	@echo "Running $(APP_NAME) application..."
	@go run ./cmd/ -conf ./configs/bootstrap.yaml

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
