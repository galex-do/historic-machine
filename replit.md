# Historical Events Mapping Application

## Overview

A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data.

## Recent Changes (October 10, 2025)

- **Complete UI Localization**: Full internationalization support for all UI elements with EN/RU locales
  - Header navigation menu (Map, Admin, Events, Tags, Datasets, Users)
  - Authentication system (Login/Logout buttons, Welcome messages, Access level badges)
  - Login modal (Title, form labels, placeholders, button states)
  - Reactive locale switching with instant UI updates
  - Locale selector positioned at header far right with flag indicators
- **Localized Dataset System**: Complete ancient civilizations dataset with 125 events
  - All events include professional Russian translations (name_ru, description_ru)
  - Chronologically ordered from 3500 BC to 1000 BC
  - Dual-language support for event data
  - Locale-aware import/export workflow
- **Session Storage Implementation**: All filter conditions (dates, lens types, templates, sidebar state) persist across page reloads
- **Performance Optimization**: Events are filtered immediately on load instead of showing all events first 
- **State Management Fix**: Event editing preserves current filter state, zoom level, and map position
- **Era Selector Implementation**: Added BC/AD era selectors to all event forms (admin panel and map creation)
- **Database Type Conflict Resolution**: Fixed PostgreSQL type casting issues for latitude/longitude coordinates
- **Authentication Header Fix**: Resolved 401 errors by properly merging authentication headers in API requests
- **Admin Area Implementation**: Complete admin panel with navigation menu
- **Vue Router Integration**: Proper SPA routing between Map and Admin views  
- **Role-Based Access**: Admin panel restricted to editor/super access levels
- **Full CRUD Interface**: Professional table-based event management with proper era handling
- **Enhanced Authentication**: Added editor role support with proper permissions
- **Navigation Menu**: Header navigation tabs for authenticated users
- **Bug Fixes**: Resolved marker binding issues during auth state changes and BC/AD date preservation

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
- **Routing**: Vue Router 4 for SPA navigation (Map view, Admin panel)
- **Map Library**: Leaflet integration for interactive world map
- **Admin Panel**: Full CRUD table interface for event management

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
- `vue-router@4.5.1`: Official router for Vue.js SPA navigation
- `leaflet@1.9.4`: Interactive map library for world map visualization
- `vue-leaflet@0.1.0`: Vue.js wrapper for Leaflet maps

### Containerization (Planned)
- Docker Compose setup with backend, frontend, and PostgreSQL services
- Individual Dockerfiles for each service
- Makefile for easy deployment and migration management

## Development Workflow

### Current Status
- ✅ Complete project architecture with backend/frontend separation
- ✅ Go backend with full REST API on port 8080
- ✅ Vue.js SPA with routing on port 5000
- ✅ PostgreSQL database with event management
- ✅ Interactive world map with Leaflet integration
- ✅ Professional admin panel with full CRUD operations
- ✅ Role-based authentication system (guest/user/editor/super)
- ✅ Navigation menu and multi-page application structure

### Completed Features
1. ✅ Interactive world map with event markers
2. ✅ Timeline filtering and date controls  
3. ✅ Event creation and editing modals
4. ✅ Admin table interface for event management
5. ✅ Authentication with session persistence
6. ✅ Co-located events grouping in popups
7. ✅ Minimalist edit icons and user interface
8. ✅ Full internationalization (EN/RU) for UI and event data
9. ✅ Localized dataset system with dual-language support