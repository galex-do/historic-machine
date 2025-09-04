# Historia ex Machina

A comprehensive web application for mapping historical events featuring interactive world map, timeline controls, BC/AD support, comprehensive tags, user authentication, and dataset management system.

## Features

- 🗺️ Interactive scalable world map with event markers
- 📅 Timeline filtering with BC/AD support and date ranges
- 🏛️ Ancient civilization navigation templates (3500 BC - 1000 BC)
- 🏷️ Comprehensive tagging system with color coding
- 👤 Role-based authentication (guest/user/editor/super)
- 📊 Professional admin panel with full CRUD operations
- 📁 Dataset import/export functionality
- 🔍 Advanced filtering and search capabilities

## Tech Stack

- **Backend**: Go with Gorilla Mux router
- **Frontend**: Vue.js 3 with Vite and Vue Router
- **Database**: PostgreSQL with PostGIS extensions
- **Authentication**: JWT tokens with bcrypt password hashing
- **Deployment**: Docker Compose with multi-container setup

## Quick Start

### Prerequisites
- Docker
- Docker Compose
- Make (optional, for convenience commands)

### Deployment with Docker Compose

1. **Clone and start services**:
```bash
make up
# Or manually: docker compose up -d
```

2. **Run database migrations**:
```bash
make migrate
# Or manually: docker compose exec backend goose -dir migrations postgres "postgres://postgres:password@db:5432/historical_events?sslmode=disable" up
```

3. **Create admin user** (choose one method):
```bash
# Method 1: Direct database insert (recommended)
docker compose exec db psql -U postgres -d historical_events -c \
  "INSERT INTO users (username, password_hash, access_level) VALUES 
   ('admin', crypt('your_secure_password', gen_salt('bf', 12)), 'super');"

# Method 2: Using registration API
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_secure_password","access_level":"super"}'

# Method 3: Interactive database session
make db-shell
# Then run SQL commands interactively
```

4. **Access the application**:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api
- **Database**: localhost:5432

### Development Mode

For development with hot reload:

1. **Start only the database**:
```bash
make dev
# Or: docker compose up -d db
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
| `make up` | Start all services with Docker Compose |
| `make down` | Stop all services |
| `make build` | Build Docker images |
| `make logs` | Show service logs |
| `make migrate` | Run database migrations |
| `make dev` | Start development environment (DB only) |
| `make admin-help` | Show admin user creation instructions |
| `make db-shell` | Connect to database shell |
| `make clean` | Clean up Docker resources |

## User Access Levels

- **guest**: View-only access to map and events
- **user**: Can view and create events
- **editor**: Can manage events and access admin tools
- **super**: Full system administration including user management

## Ancient Civilization Templates

Navigate through 5500+ years of history with 26 pre-configured templates across 6 major periods:

- **Dawn of Civilization** (3500-3000 BC): First Cities, Invention of Writing, Egyptian Unification
- **Early Dynastic Period** (3000-2350 BC): Egyptian Early Dynasties, Sumerian City-States, Great Pyramid Era
- **First Empires** (2350-2000 BC): Akkadian Empire, Egyptian Old Kingdom, Third Dynasty of Ur, Climate Crisis
- **Bronze Age Powers** (2000-1500 BC): Old Babylonian Period, Egyptian Middle Kingdom, Code of Hammurabi, Hyksos Period
- **Imperial Age** (1500-1200 BC): Egyptian New Kingdom, Hittite Empire, Hatshepsut, Thutmose III, Ramesses II, Battle of Kadesh
- **Bronze Age Collapse** (1200-1000 BC): Sea Peoples, Fall of Hittites, Civilizational Crisis, Rise of Iron Age

## API Endpoints

### Authentication
- `POST /api/auth/login` - User authentication
- `POST /api/auth/register` - User registration

### Events
- `GET /api/events` - List all events with filtering
- `POST /api/events` - Create new event
- `GET /api/events/{id}` - Get specific event
- `PUT /api/events/{id}` - Update event
- `DELETE /api/events/{id}` - Delete event

### Navigation Templates
- `GET /api/date-template-groups` - Get template groups
- `GET /api/date-templates` - Get all templates
- `GET /api/date-templates/group/{id}` - Get templates for group

### Dataset Management
- `GET /api/datasets` - List event datasets
- `POST /api/datasets/import` - Import event dataset
- `DELETE /api/datasets/{id}` - Delete dataset

### Tags
- `GET /api/tags` - List all tags
- `POST /api/tags` - Create new tag
- `GET /api/events/{id}/tags` - Get event tags
- `POST /api/events/{id}/tags` - Set event tags

## Project Structure

```
├── backend/                 # Go backend API
│   ├── internal/           # Internal packages
│   ├── migrations/         # Database migrations
│   ├── Dockerfile          # Backend container
│   └── main.go            # Application entry point
├── frontend/              # Vue.js frontend
│   ├── src/               # Source code
│   ├── Dockerfile         # Frontend container
│   └── nginx.conf         # Nginx configuration
├── datasets/              # Historical event datasets
├── docker-compose.yml     # Service orchestration
├── Makefile              # Build automation
└── README.md             # This documentation
```

## Database Schema

The application includes a comprehensive database schema with:

- **Events**: Core historical events with geo-temporal data
- **Users**: Authentication with role-based access control  
- **Tags**: Flexible tagging system with many-to-many relationships
- **Datasets**: Organization and tracking of imported event collections
- **Date Templates**: Pre-configured navigation periods for ancient civilizations
- **Views**: Optimized queries with computed display dates and tag aggregation

## Docker Services

- **db**: PostgreSQL with PostGIS extensions (port 5432)
- **backend**: Go API server (port 8080)
- **frontend**: Vue.js app served by Nginx (port 3000)

## Environment Variables

Docker Compose automatically configures:
- `DB_HOST=db`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=password`
- `DB_NAME=historical_events`
- `DB_SSL_MODE=disable`

## Contributing

1. Use snake_case naming convention for all functions and variables
2. Follow existing code patterns and project structure
3. Keep database changes in migrations
4. Test with Docker Compose deployment
5. Update documentation as needed

## License

MIT License - see LICENSE file for details.