GO_PLUGIN_PKG=google.golang.org/protobuf/cmd/protoc-gen-go
GO_GRPC_PLUGIN_PKG=google.golang.org/grpc/cmd/protoc-gen-go-grpc
OPENAPI_PLUGIN_PKG=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
GATEWAY_PLUGIN_PKG=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

PROTO_DIR=api/proto
PROJECT=github.com/guodongq/jigsaw
BUILD_VERSION=v0.0.1
GIT_VERSION = $(shell git show -s --format='format:%h %aI')
BUILD_FLAGS=-ldflags '-X "$(PROJECT)/pkg/version.BuildVersion=$(BUILD_VERSION)" -X "$(PROJECT)/pkg/version.BuildGitHash=$(GIT_VERSION)"'

DEFAULT: all

plugin:
	go get $(GO_PLUGIN_PKG)
	go get $(GO_GRPC_PLUGIN_PKG)
	go get $(GATEWAY_PLUGIN_PKG)
	go get $(OPENAPI_PLUGIN_PKG)

vendor:
	go mod tidy

generate_files:
	@for file in $$(buf ls-files $(PROTO_DIR)); do \
		echo "files:$$file"; \
		buf generate --path $$file; \
	done;

generate_update:
	buf build -o lock.bin

generate_check:
	buf check breaking --against lock.bin

clean:
	rm -rf api/gen
	rm -rf api/spec

all: clean plugin generate_check generate_files generate_update vendor

designer:
	mkdir -p _output/designer
	rm -rf _output/designer/*
	go build -o _output/designer/designer $(BUILD_FLAGS) cmd/designer/main.go
	cp -rf hack/designer/* _output/designer
