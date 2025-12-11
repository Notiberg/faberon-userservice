.PHONY: help build run test clean clean-all docker-build docker-up docker-down docker-restart docker-logs docker-clean docker-prune migrate-up migrate-down db-reset fixtures

# Variables
APP_NAME=smc-userservice
DOCKER_COMPOSE=docker-compose
GO=go

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application binary"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts and logs"
	@echo "  make clean-all      - Clean everything (artifacts, logs, Docker volumes, images)"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start all services (app, db, monitoring)"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-restart - Restart all services"
	@echo "  make docker-logs    - Show logs from all containers"
	@echo "  make docker-logs-app - Show logs from app container only"
	@echo "  make docker-clean   - Stop services and remove volumes"
	@echo "  make docker-prune   - Remove all project Docker images"
	@echo ""
	@echo "Database commands:"
	@echo "  make migrate-up     - Apply database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make db-reset       - Reset database (down volumes + up)"
	@echo "  make fixtures       - Load test fixtures into database"

# Build commands
build:
	@echo "Building application..."
	@$(GO) build -o bin/$(APP_NAME) ./cmd/main.go
	@echo "Build complete: bin/$(APP_NAME)"

run:
	@echo "Running application locally..."
	@$(GO) run cmd/main.go

test:
	@echo "Running tests..."
	@$(GO) test ./... -v

clean:
	@echo "Cleaning build artifacts and logs..."
	@rm -rf bin/
	@rm -rf logs/*.log
	@rm -rf docker/
	@echo "Clean complete"

clean-all: clean docker-clean docker-prune
	@echo "Full cleanup complete"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@$(DOCKER_COMPOSE) build userservice
	@echo "Docker image built successfully"

docker-up:
	@echo "Starting all services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Services started. Access:"
	@echo "  - App:        http://localhost:8080"
	@echo "  - Metrics:    http://localhost:8080/metrics"
	@echo "  - Prometheus: http://localhost:9091"
	@echo "  - Grafana:    http://localhost:3001 (admin/admin)"
	@echo "  - PostgreSQL: localhost:5435"

docker-down:
	@echo "Stopping all services..."
	@$(DOCKER_COMPOSE) down
	@echo "Services stopped"

docker-restart:
	@echo "Restarting all services..."
	@$(DOCKER_COMPOSE) restart
	@echo "Services restarted"

docker-logs:
	@$(DOCKER_COMPOSE) logs -f

docker-logs-app:
	@$(DOCKER_COMPOSE) logs -f userservice

docker-clean:
	@echo "Stopping services and removing volumes..."
	@$(DOCKER_COMPOSE) down -v
	@echo "Docker volumes removed"

docker-prune:
	@echo "Removing project Docker images..."
	@docker images | grep smc-userservice | awk '{print $$3}' | xargs -r docker rmi -f || true
	@echo "Docker images removed"

# Database commands
migrate-up:
	@echo "Applying database migrations..."
	@$(DOCKER_COMPOSE) up -d postgres
	@sleep 3
	@$(DOCKER_COMPOSE) up migrate
	@echo "Migrations applied"

migrate-down:
	@echo "Rolling back database migrations..."
	@$(DOCKER_COMPOSE) run --rm migrate -path /migrations -database "postgres://postgres:postgres@postgres:5432/smc_userservice?sslmode=disable" down
	@echo "Migrations rolled back"

db-reset:
	@echo "Resetting database..."
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d postgres
	@sleep 5
	@$(DOCKER_COMPOSE) up migrate
	@echo "Database reset complete"

fixtures:
	@echo "Loading test fixtures into database..."
	@$(DOCKER_COMPOSE) up -d postgres
	@sleep 2
	@docker exec -i smc-userservice-db psql -U postgres -d smc_userservice < migrations/fixtures/001_test_users.sql
	@echo "Fixtures loaded successfully"
	@echo ""
	@echo "Loaded data:"
	@echo "  - 11 users (7 clients, 3 managers, 1 without car)"
	@echo "  - 7 cars with selected flags"
	@echo ""
	@echo "Verify with: docker exec -it smc-userservice-db psql -U postgres -d smc_userservice -c 'SELECT COUNT(*) FROM users;'"

# Development helpers
dev:
	@echo "Starting development environment..."
	@$(DOCKER_COMPOSE) up -d postgres prometheus grafana
	@echo "Infrastructure started. Run 'make run' to start the app locally"

install:
	@echo "Installing Go dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "Dependencies installed"
