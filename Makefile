# Historia ex Machina
.PHONY: help build up down logs clean migrate dev admin-help

# Default target
help:
	@echo "Available commands:"
	@echo "  make build       - Build all Docker images"
	@echo "  make up          - Start all services with Docker Compose"
	@echo "  make down        - Stop all services"
	@echo "  make logs        - Show logs from all services"
	@echo "  make clean       - Remove all containers and images"
	@echo "  make migrate     - Run database migrations"
	@echo "  make dev         - Start development environment (DB only)"
	@echo "  make admin-help  - Show admin user creation instructions"

# Docker operations
build:
	@echo "Building Docker images..."
	docker compose build

up:
	@echo "Starting services..."
	docker compose up --build
	@echo "Services started!"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Backend API: http://localhost:8080/api"
	@echo ""
	@echo "Run 'make migrate' to set up the database"
	@echo "Run 'make admin-help' for admin user creation"

down:
	@echo "Stopping services..."
	docker compose down

logs:
	@echo "Showing logs..."
	docker compose logs -f

# Migration operations
migrate:
	@echo "Running database migrations..."
	docker compose run --rm migrate -dir /migrations postgres "postgres://postgres:password@db:5432/historical_events?sslmode=disable" up
	@echo "Database migration completed!"

# Development environment
dev:
	@echo "Starting development environment..."
	docker compose up -d db
	@echo "Database started on localhost:5432"
	@echo "Run backend and frontend locally for development."

# Admin user creation help
admin-help:
	@echo "Create admin user (choose one method):"
	@echo ""
	@echo "Method 1: Using PostgreSQL directly"
	@echo "  docker compose exec db psql -U postgres -d historical_events -c \\"
	@echo "    \"INSERT INTO users (username, password_hash, access_level) VALUES" 
	@echo "     ('admin', crypt('your_password', gen_salt('bf', 12)), 'super');\""
	@echo ""
	@echo "Method 2: Using registration API"
	@echo "  curl -X POST http://localhost:8080/api/auth/register \\"
	@echo "    -H \"Content-Type: application/json\" \\"
	@echo "    -d '{\"username\":\"admin\",\"password\":\"your_password\"}'"
	@echo ""
	@echo "Method 3: Interactive psql session"
	@echo "  make db-shell"
	@echo "  Then run: UPDATE users SET access_level='super' WHERE username='admin';"

# Database shell access
db-shell:
	@echo "Connecting to database..."
	docker compose exec db psql -U postgres -d historical_events

# Cleanup operations
clean:
	@echo "Cleaning up Docker resources..."
	docker compose down -v --remove-orphans
	docker system prune -f
