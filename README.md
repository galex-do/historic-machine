# Historical Events Mapping Application

A web application for mapping historical events on an interactive world map with timeline functionality.

## Features

- Interactive scalable world map
- Timeline filtering with date ranges (FROM/TO)
- Click-to-add events with geo-coordinates
- Event categorization with lens system
- Persistent data storage with PostgreSQL

## Tech Stack

- **Backend**: Go with Gorilla Mux router
- **Frontend**: Vue.js 3 with Vite
- **Database**: PostgreSQL 17
- **Containerization**: Docker & Docker Compose
- **Migrations**: Goose migration tool

## Quick Start

### Prerequisites
- Docker
- Docker Compose
- Make (optional, for convenience commands)

### Running the Application

1. **Clone and start services**:
   ```bash
   make up
   ```
   
   Or manually:
   ```bash
   docker-compose up -d
   ```

2. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

3. **Run database migrations**:
   ```bash
   make migrate
   ```

### Development Mode

For development with hot reload:

1. **Start only the database**:
   ```bash
   make dev
   ```

2. **Run backend locally**:
   ```bash
   cd backend
   go run main.go
   ```

3. **Run frontend locally**:
   ```bash
   cd frontend
   npm run dev
   ```

## Available Commands

| Command | Description |
|---------|-------------|
| `make up` | Start all services |
| `make down` | Stop all services |
| `make build` | Build Docker images |
| `make logs` | Show service logs |
| `make migrate` | Run database migrations |
| `make clean` | Clean up Docker resources |
| `make dev` | Start development environment |

## API Endpoints

- `GET /api/events` - Get all events
- `POST /api/events` - Create new event
- `GET /api/events/{id}` - Get specific event

## Project Structure

```
├── backend/              # Go REST API
│   ├── main.go          # Main application
│   ├── Dockerfile       # Backend container
│   └── migrations/      # Database migrations
├── frontend/            # Vue.js application
│   ├── src/            # Source files
│   ├── Dockerfile      # Frontend container
│   └── nginx.conf      # Nginx configuration
├── docker-compose.yml  # Service orchestration
├── Makefile           # Build automation
└── README.md          # This file
```

## Environment Variables

Backend environment variables:
- `DB_HOST` - Database host
- `DB_PORT` - Database port  
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSL_MODE` - SSL mode for database

## Contributing

1. Use snake_case naming convention for all functions and variables
2. Follow existing code patterns and structure
3. Add appropriate tests for new features
4. Update documentation as needed