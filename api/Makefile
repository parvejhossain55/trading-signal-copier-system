# ====== App metadata ======
MAIN := ./
TARGET := main
SERVER_CMD := ./$(TARGET) serve-rest
GEN_JWT := ./$(TARGET) gen-jwt

# ====== Protobuf ======
PROTOC_DEST := ./
PROTO_FILES := $(shell find ./proto -name "*.proto" -type f)

GO_BIN := $(shell go env GOBIN)
ifeq ($(GO_BIN),)
  GO_BIN := $(shell go env GOPATH)/bin
endif
export PATH := $(PATH):$(GO_BIN)

# ====== Targets ======

## Generate protobufs (robust: explicit plugin paths + timestamp mapping)
build-proto:
	@echo "Generating protobuf files..."
	PATH="$(PATH):$(GO_BIN)" \
	protoc \
	  --plugin=protoc-gen-go=$(GO_BIN)/protoc-gen-go \
	  --plugin=protoc-gen-go-grpc=$(GO_BIN)/protoc-gen-go-grpc \
	  --go_out=$(PROTOC_DEST) \
	  --go_opt=paths=source_relative \
	  --go_opt=Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb \
	  --go-grpc_out=$(PROTOC_DEST) \
	  --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILES)
	@echo "Protobuf generation completed!"

clean-proto:
	@echo "Cleaning generated protobuf files..."
	find ./proto -name "*.pb.go" -type f -delete
	@echo "Clean completed!"

proto: clean-proto build-proto

## Build/run app
build: install-deps
	go build -o $(TARGET) $(MAIN)

run:
	$(SERVER_CMD)

start: build run

## Dependencies / tooling
tidy:
	go mod tidy

install-deps:
	go mod download

install-proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

install-dev-deps:
	go install github.com/air-verse/air@latest

install-mockgen:
	go install go.uber.org/mock/mockgen@latest

prepare: install-proto-deps install-dev-deps install-deps tidy proto

## Dev (hot reload via Air)
dev: prepare
	$(GO_BIN)/air -c .air.toml

## JWT command (no hot reload)
jwt: build
	$(GEN_JWT)

## Utilities
list-proto:
	@echo "Found proto files:"
	@find ./proto -name "*.proto" -type f | sort

messaging-example:
	@echo "Running messaging example..."
	go run scripts/messaging_example.go

infra-up:
	@echo "Starting infrastructure services..."
	cd ../../infra && docker-compose up -d

infra-down:
	@echo "Stopping infrastructure services..."
	cd ../../infra && docker-compose down

infra-logs:
	@echo "Showing infrastructure logs..."
	cd ../../infra && docker-compose logs -f

help:
	@echo "Available commands:"
	@echo "  build-proto, clean-proto, proto, list-proto"
	@echo "  install-deps, install-proto-deps, install-dev-deps, tidy"
	@echo "  build, run, dev, start, jwt"
	@echo "  messaging-example, infra-up, infra-down, infra-logs, help"

.PHONY: build-proto clean-proto proto build run start tidy \
        install-deps install-proto-deps install-dev-deps install-mockgen \
        prepare dev jwt list-proto messaging-example infra-up infra-down \
        infra-logs help
