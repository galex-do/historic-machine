# Historical Events Mapping Application
.PHONY: help build up down restart logs clean migrate migrate-up migrate-down migrate-status

# Default target
help:
	@echo "Available commands:"
	@echo "  make build       - Build all Docker images"
	@echo "  make up          - Start all services with Docker Compose"
	@echo "  make down        - Stop all services"
	@echo "  make restart     - Restart all services"
	@echo "  make logs        - Show logs from all services"
	@echo "  make clean       - Remove all containers and images"
	@echo "  make migrate     - Run database migrations"
	@echo "  make migrate-up  - Run migrations up"
	@echo "  make migrate-down - Run migrations down"
	@echo "  make dev         - Start development environment"

# Docker operations
build:
	@echo "Building Docker images..."
	docker compose build

up:
	@echo "Starting services..."
	docker compose up
	@echo "Services started! Frontend: http://localhost:3000"

down:
	@echo "Stopping services..."
	docker compose down

restart:
	@echo "Restarting services..."
	docker compose restart

logs:
	@echo "Showing logs..."
	docker compose logs -f

# Migration operations
migrate: migrate-up

migrate-up:
	@echo "Running database migrations up..."
	docker compose exec backend ./goose -dir migrations postgres "postgres://postgres:password@db:5432/historical_events?sslmode=disable" up

migrate-down:
	@echo "Running database migrations down..."
	docker compose exec backend ./goose -dir migrations postgres "postgres://postgres:password@db:5432/historical_events?sslmode=disable" down

migrate-status:
	@echo "Checking migration status..."
	docker compose exec backend ./goose -dir migrations postgres "postgres://postgres:password@db:5432/historical_events?sslmode=disable" status

# Development environment
dev:
	@echo "Starting development environment..."
	docker compose up -d db
	@echo "Database started. Run backend and frontend separately for development."

# Cleanup operations
clean:
	@echo "Cleaning up Docker resources..."
	docker compose down -v --remove-orphans
	docker system prune -f
