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
- **Database**: PostgreSQL with comprehensive schema
- **Authentication**: JWT tokens with bcrypt password hashing
- **Migrations**: Goose migration tool

## Quick Start

### Prerequisites
- Go 1.19+
- Node.js 18+
- PostgreSQL database
- Environment variables: `DATABASE_URL`, `JWT_SECRET`

### Setup

1. **Install dependencies:**
```bash
# Backend
cd backend && go mod download

# Frontend  
cd frontend && npm install
```

2. **Database setup:**
```bash
# Run migrations to create complete schema
cd backend && goose postgres $DATABASE_URL up
```

3. **Create admin user** (choose one method):

**Method 1: Using PostgreSQL with bcrypt**
```bash
# Connect to database
psql $DATABASE_URL

# Create admin user (replace 'your_secure_password' with actual password)
INSERT INTO users (username, password_hash, access_level) VALUES 
('admin', crypt('your_secure_password', gen_salt('bf', 12)), 'super');
```

**Method 2: Using registration API**
```bash
# Start backend server first
cd backend && go run main.go

# Register admin user via API
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_secure_password", 
    "access_level": "super"
  }'
```

**Method 3: Generate hash manually**
```bash
# Generate bcrypt hash
htpasswd -bnBC 12 "" your_password | cut -d: -f2

# Insert with generated hash
psql $DATABASE_URL -c "INSERT INTO users (username, password_hash, access_level) VALUES ('admin', '\$2a\$12\$generated_hash', 'super');"
```

4. **Start the application:**
```bash
# Terminal 1: Backend
cd backend && go run main.go

# Terminal 2: Frontend  
cd frontend && npm run dev
```

5. **Access the application:**
- Frontend: http://localhost:5000
- Backend API: http://localhost:8080/api
- Login with your created admin credentials

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
│   │   ├── handlers/       # HTTP handlers
│   │   ├── services/       # Business logic
│   │   ├── models/         # Data models
│   │   └── database/       # Database layer
│   ├── migrations/         # Database migrations
│   └── main.go            # Application entry point
├── frontend/              # Vue.js frontend
│   ├── src/               # Source code
│   │   ├── components/    # Vue components
│   │   ├── views/         # Page views
│   │   ├── composables/   # Vue composables
│   │   └── utils/         # Utilities
│   └── public/            # Static assets
└── datasets/              # Historical event datasets
    └── ancient_civilizations_pre_1000bc.json
```

## Database Schema

The application includes a comprehensive database schema with:

- **Events**: Core historical events with geo-temporal data
- **Users**: Authentication with role-based access control
- **Tags**: Flexible tagging system with many-to-many relationships
- **Datasets**: Organization and tracking of imported event collections
- **Date Templates**: Pre-configured navigation periods
- **Views**: Optimized queries with computed display dates

## Contributing

1. Use snake_case naming convention for all functions and variables
2. Follow existing code patterns and project structure
3. Keep database changes in migrations
4. Update documentation as needed
5. Test authentication and permission levels

## Environment Variables

Required environment variables:
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT token signing
- `PORT` - Server port (default: 8080)

## License

MIT License - see LICENSE file for details.