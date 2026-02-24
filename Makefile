# Signal Copier Project Makefile

.PHONY: help build up down logs logs-api logs-db dev dev-api migrate test test-unit test-integration test-e2e clean clean-images status health start

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

logs-db: ## Show logs for Postgres service only
	docker-compose logs -f postgres

# Development commands
dev: ## Start services in development mode
	docker-compose up

dev-api: ## Start API and Postgres service
	docker-compose up api postgres

# Database commands
migrate: ## Run database migrations (automatic in API startup)
	docker-compose exec api /app/main migrate

seed: ## Seed the database
	docker-compose exec api /app/main seed

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
	@echo "Postgres Health:"
	@docker-compose exec postgres pg_isready -U admin -d copier_db || echo "Postgres not responding"

# Quick start
start: build up ## Build and start services
	@echo "Project started successfully!"
	@echo "API: http://localhost:9090"
	@echo "Health: http://localhost:9090/health"
