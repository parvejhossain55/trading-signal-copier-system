MAIN:=./
TARGET:=main
SERVER_CMD:=./${TARGET} serve-rest
GEN_JWT:=./${TARGET} gen-jwt
PROTOC_DEST:=./
PROTOC_FLAGS:=--go_out=${PROTOC_DEST} --go_opt=paths=source_relative --go-grpc_out=${PROTOC_DEST} --go-grpc_opt=paths=source_relative
PROTO_FILES:=$(shell find ./proto -name "*.proto" -type f)

# example migration create command -
# migrate create -ext sql -seq -dir migrations create-some-table

build-proto:
	@echo "Generating protobuf files..."
	protoc ${PROTOC_FLAGS} ${PROTO_FILES}
	@echo "Protobuf generation completed!"

clean-proto:
	@echo "Cleaning generated protobuf files..."
	find ./proto -name "*.pb.go" -type f -delete
	@echo "Clean completed!"

proto: clean-proto build-proto

run:
	${SERVER_CMD}

tidy:
	go mod tidy
	

install-proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

install-dev-deps:
	go install github.com/air-verse/air@latest

install-deps:
	go mod download

install-mockgen:
	go install go.uber.org/mock/mockgen@latest

prepare: install-proto-deps install-dev-deps install-deps tidy proto

dev: prepare
	air serve-rest

jwt: prepare
	air gen-jwt


build: install-deps
	go build -o ${TARGET} ${MAIN}

start: build run-server

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
	@echo "  build-proto    - Generate protobuf files from .proto definitions"
	@echo "  clean-proto    - Remove generated protobuf files"
	@echo "  proto          - Clean and rebuild protobuf files"
	@echo "  list-proto     - List all proto files found"
	@echo "  install-deps   - Install Go dependencies"
	@echo "  install-proto-deps - Install protobuf compiler and plugins"
	@echo "  prepare        - Install all dependencies and generate proto files"
	@echo "  dev            - Start development server with hot reload"
	@echo "  build          - Build the application"
	@echo "  run-server     - Run the server"
	@echo "  start          - Build and run the server"
	@echo "  tidy           - Run go mod tidy"
	@echo "  messaging-example - Run the messaging integration example"
	@echo "  infra-up       - Start Kafka, NATS, and other infrastructure services"
	@echo "  infra-down     - Stop infrastructure services"
	@echo "  infra-logs     - Show infrastructure service logs"
	@echo "  help           - Show this help message"