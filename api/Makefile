# ====== App metadata ======
MAIN := ./
TARGET := main
SERVER_CMD := ./$(TARGET) serve-rest
GEN_JWT := ./$(TARGET) gen-jwt
MIGRATE_CMD := echo "PostgreSQL migrations are automatic on startup"
SEED_CMD := ./$(TARGET) seed

# ====== Go environment ======
GO_BIN := $(shell go env GOBIN)
ifeq ($(GO_BIN),)
  GO_BIN := $(shell go env GOPATH)/bin
endif
export PATH := $(PATH):$(GO_BIN)

# ====== Targets ======

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

install-dev-deps:
	go install github.com/air-verse/air@latest

install-mockgen:
	go install go.uber.org/mock/mockgen@latest

prepare: install-dev-deps install-deps tidy

## Dev (hot reload via Air)
dev: prepare
	$(GO_BIN)/air -c .air.toml

## JWT command (no hot reload)
jwt: build
	$(GEN_JWT)

## Database commands
migrate: build
	$(MIGRATE_CMD)

seed: build
	$(SEED_CMD)

migrate-seed: migrate seed

help:
	@echo "Available commands:"
	@echo "  install-deps, install-dev-deps, install-mockgen, tidy"
	@echo "  build, run, dev, start, jwt, help"
	@echo "  migrate, seed, migrate-seed"

.PHONY: build run start tidy install-deps install-dev-deps install-mockgen \
        prepare dev jwt migrate seed migrate-seed help
