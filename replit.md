# Historical Events Mapping Application

## Overview

A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data.

## Recent Changes (August 29, 2025)

- Set up basic project structure with backend and frontend folders
- Created Go REST API backend with Gorilla Mux router
- Implemented Vue.js frontend with Vite build system
- Configured workflows: Backend on port 8080, Frontend on port 5000
- Basic event API endpoints for CRUD operations
- Initial UI with timeline controls and event display

## User Preferences

- Use snake_case naming convention everywhere for elements and functions
- Develop professional, high-end code with proper patterns and templates
- Avoid duplication - refactor similar logic into reusable templates/functions
- Simple, everyday language for communication

## Project Architecture

### Backend Architecture (Go)
- **Framework**: Gorilla Mux HTTP router
- **Port**: 8080
- **Database**: PostgreSQL (to be configured)
- **Migrations**: Goose migration tool (to be implemented)
- **API Pattern**: RESTful endpoints
- **CORS**: Enabled for frontend integration

### Frontend Architecture (Vue.js)
- **Framework**: Vue.js 3 with Vite build system
- **Port**: 5000 (configured for host 0.0.0.0 with allowedHosts: true)
- **State Management**: To be implemented with Pinia
- **Routing**: To be implemented with Vue Router
- **Map Library**: To be integrated (Leaflet or similar)

### Database Schema (Planned)
- **Events Table**: id, name, description, latitude, longitude, event_date, lens_type
- **Lenses Table**: id, name, description, color
- **Migration Tool**: Goose for Go

### Key Features
1. **Interactive World Map**: Scalable map with event markers
2. **Timeline Filtering**: Date range selection (FROM/TO fields)
3. **Event Management**: Click to add, hover for details
4. **Lens System**: Categorization system for different event types
5. **Modal Interface**: Add/edit event forms

## External Dependencies

### Backend Dependencies
- `github.com/gorilla/mux`: HTTP router and URL matcher
- PostgreSQL driver (to be added)
- Goose migration tool (to be added)

### Frontend Dependencies
- `vue@3.5.20`: Progressive JavaScript framework
- `vite@7.1.3`: Build tool and dev server
- `@vitejs/plugin-vue@6.0.1`: Vue plugin for Vite
- Map library (to be added - Leaflet recommended)

### Containerization (Planned)
- Docker Compose setup with backend, frontend, and PostgreSQL services
- Individual Dockerfiles for each service
- Makefile for easy deployment and migration management

## Development Workflow

### Current Status
- ✅ Basic project structure created
- ✅ Go backend with API endpoints running on port 8080
- ✅ Vue frontend with Vite dev server running on port 5000
- ✅ CORS configured for frontend-backend communication
- ✅ Basic event data structure implemented

### Next Steps
1. Set up PostgreSQL database and Goose migrations
2. Integrate world map library (Leaflet)
3. Implement timeline filtering functionality
4. Add event creation modal
5. Implement lens system
6. Docker containerization
7. Production deployment setup