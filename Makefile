GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	TOOLKITS_PROTO_FILES=$(shell $(Git_Bash) -c "find toolkits -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	TOOLKITS_PROTO_FILES=$(shell find toolkits -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest


.PHONY: tools
# generate tools proto or use ./toolkits/generate.go
tools:
	protoc --proto_path=./internal \
		--proto_path=./third_party \
		--proto_path=./toolkits \
		--go_out=paths=source_relative:./toolkits \
		--validate_out=lang=go,paths=source_relative:./toolkits \
		$(TOOLKITS_PROTO_FILES)

.PHONY: config
# generate internal proto or use ./internal/generate.go
config: 
	protoc --proto_path=./internal \
		--proto_path=./third_party \
		--proto_path=./toolkits \
		--go_out=paths=source_relative:./internal \
		--validate_out=lang=go:. \
		$(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto or use ./api/generate.go
api:
#	protoc --proto_path=./api \
#	       --proto_path=./third_party \
# 	       --go_out=paths=source_relative:./api \
# 	       --go-http_out=paths=source_relative:./api \
# 	       --go-grpc_out=paths=source_relative:./api \
#	       --openapi_out=fq_schema_naming=true,default_response=false:. \
#	       $(API_PROTO_FILES)
	protoc --proto_path=./api \
		--proto_path=./third_party \
		--proto_path=./toolkits \
		--go_out=./api \
		--go-http_out=./api \
		--go-grpc_out=./api \
		--validate_out=lang=go:./api \
		--openapi_out=fq_schema_naming=true,default_response=false:. \
		$(API_PROTO_FILES)

.PHONY: pre
# pre
pre:
	goreleaser build --single-target --clean --snapshot

.PHONY: build
# build
build:
	mkdir -p dist/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./dist/ ./...

.PHONY: release
# release
release:
	goreleaser release --clean

#.PHONY: server
## server used generate a service at first
#server:
#	kratos proto server -t ./internal/mods/helloworld/service ./api/v1/protos/helloworld/greeter.proto
#
#.PHONY: client
## client used when proto file is in the same directory
#client:
#	kratos proto client ./api

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make tools;
	make api;
	make config;
	make generate;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
