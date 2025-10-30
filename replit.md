# Historical Events Mapping Application

## Overview
A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data. The project aims to provide a rich, interactive experience for exploring history, with capabilities for managing events, tags, and datasets, all within a localized and performant environment. This application has significant market potential for educational platforms, historical research, and general public engagement with historical data.

## User Preferences
- Use snake_case naming convention everywhere for elements and functions
- Develop professional, high-end code with proper patterns and templates
- Avoid duplication - refactor similar logic into reusable templates/functions
- Simple, everyday language for communication

## System Architecture

### UI/UX Decisions
- **Interactive World Map**: Utilizes Leaflet for scalable and interactive map visualization with event markers, clustering, and narrative flow visualization.
- **Admin Panel**: Professional table-based interface for managing events, tags, users, and datasets.
- **Localization**: Full internationalization (English/Russian) for all UI elements, historical period templates, and event data. Reactive locale switching.
- **Filter Panel**: Kibana-style removable chip interface for tag filtering with session storage persistence. Features a collapsible tag filter panel with a toggle button and chronological template sorting.
- **Pagination**: Minimalistic pagination moved to the header for constant visibility.
- **Date Handling**: BC/AD era selectors in forms, localized date formatting, and timezone-safe displays, with auto-expansion of date ranges for template groups.
- **Interactive Tags**: Clickable tags in event popups and expandable tags on event cards/modals for enhanced filtering and UX.
- **Date Template Selector**: Two-pane hierarchical selector for improved space efficiency and UX.

### Technical Implementations
- **Backend (Go)**:
    - **Framework**: Gorilla Mux HTTP router.
    - **API Pattern**: RESTful endpoints.
    - **Security**: Editor-level authentication for CRUD, Nginx hardening (XSS protection, CSP), role-based access control.
    - **Performance**: Nginx optimizations (worker processes, epoll, sendfile, gzip, caching), server-side OSM tile caching.
    - **Containerization**: Multi-stage Docker builds, Docker Compose setup.
- **Frontend (Vue.js)**:
    - **Framework**: Vue.js 3 with Vite build system.
    - **Routing**: Vue Router 4 for SPA navigation.
    - **State Management**: Session storage for filter conditions and map state persistence.
    - **Event Clustering**: Co-located events clustering with count badges.
    - **Reactive Data Fetching**: Event and template data re-fetches automatically on locale changes.
    - **Tag Filtering Logic**: AND logic for multi-tag filtering.

### Feature Specifications
- **Interactive World Map**: Displays event markers with hover details, click-to-add functionality, and narrative flow visualization between filtered events.
- **Timeline Filtering**: Date range selection (FROM/TO fields) and historical period templates with chronological sorting.
- **Event Management**: CRUD operations for events via an admin panel, including BC date editing.
- **Lens System**: Categorization system for different event types.
- **Tag Management**: CRUD for tags, with cross-page navigation to filter events by tag.
- **User Management**: Admin interface for managing user accounts and access levels (Guest, User, Editor, Admin, Super).
- **Dataset Management**: Import and creation of localized event datasets (e.g., Phoenician Mediterranean Empire, Assyrian Empire, Hurrian-Mitanni Kingdom, Thutmose III Campaigns).
- **Authentication**: Role-based access control, session persistence.

### System Design Choices
- **Database**: PostgreSQL with Goose migration tool.
- **API**: CORS enabled for frontend integration.
- **Deployment**: Docker Compose for orchestrating backend, frontend, and database services.

## External Dependencies

### Backend Dependencies
- `github.com/gorilla/mux`: HTTP router and URL matcher.
- `PostgreSQL driver`: Database connectivity.
- `Goose`: Database migration tool.

### Frontend Dependencies
- `vue@3.5.20`: Progressive JavaScript framework.
- `vite@7.1.3`: Build tool and dev server.
- `@vitejs/plugin-vue@6.0.1`: Vue plugin for Vite.
- `vue-router@4.5.1`: Official router for Vue.js SPA navigation.
- `leaflet@1.9.4`: Interactive map library for world map visualization.
- `vue-leaflet@0.1.0`: Vue.js wrapper for Leaflet maps.

### Other Integrations
- **OpenStreetMap**: Used for map tiles, with server-side caching via Nginx.