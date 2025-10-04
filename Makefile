# Signal Copier Project Makefile

.PHONY: help build up down logs clean restart migrate test

# Default target
help: ## Show this help message
	@echo "Signal Copier Project Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker Compose commands
build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs for all services
	docker-compose logs -f

logs-api: ## Show logs for API service only
	docker-compose logs -f api

logs-web: ## Show logs for Web service only
	docker-compose logs -f web

logs-db: ## Show logs for PostgreSQL service only
	docker-compose logs -f postgres

# Development commands
dev: ## Start services in development mode (with live reload)
	docker-compose up

dev-api: ## Start only API service in development mode
	docker-compose up api postgres redis kafka nats

dev-web: ## Start only Web service in development mode
	docker-compose up web

# Database commands
migrate: ## Run database migrations
	docker-compose exec api /app/main migrate

# Testing commands
test: ## Run all tests
	docker-compose exec api go test ./...

test-unit: ## Run unit tests
	docker-compose exec api go test ./tests/unit/...

test-integration: ## Run integration tests
	docker-compose exec api go test ./tests/integration/...

test-e2e: ## Run end-to-end tests
	docker-compose exec api go test ./tests/e2e/...

# Cleanup commands
clean: ## Remove all containers, networks, and volumes
	docker-compose down -v --remove-orphans
	docker system prune -f

clean-images: ## Remove all project images
	docker-compose down --rmi all

# Status commands
status: ## Show status of all services
	docker-compose ps

# Health check commands
health: ## Check health of all services
	@echo "Checking service health..."
	@echo "API Health:"
	@curl -s http://localhost:9090/health || echo "API not responding"
	@echo ""
	@echo "Web Health:"
	@curl -s http://localhost:3000 || echo "Web not responding"
	@echo ""
	@echo "PostgreSQL Health:"
	@docker-compose exec postgres pg_isready -U postgres || echo "PostgreSQL not responding"
	@echo ""
	@echo "Redis Health:"
	@docker-compose exec redis redis-cli ping || echo "Redis not responding"

# Quick start
start: build up migrate ## Build, start services, and run migrations
	@echo "Project started successfully!"
	@echo "API: http://localhost:9090"
	@echo "Web: http://localhost:3000"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Kafka: localhost:9092"
	@echo "NATS: localhost:4222"
