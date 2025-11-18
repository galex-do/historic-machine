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
- **Mobile-First Layout**: Events grid above the map on mobile, compact control header, hidden event cards and pagination, content-dependent sidebar height, and optimized map space.
- **Timeline View Modal**: Ultra-compact vertical timeline visualization with date grouping, expandable descriptions, BC/AD aware sorting, and accessibility features.
- **Compact Event List View**: Redesigned to a minimalistic, single-line layout for events, showing full descriptions inline and colored #tag format.
- **Map Highlight Feature**: Static golden halo ring overlay for highlighting markers (replaced problematic CSS animation). Conditional "Highlight" (📍) / "Focus" (⌖) button display based on map filtering state.
- **Map Event Info Modal Timeline Structure**: Uses the same timeline design as TimelineModal with vertical line, bullets, date grouping, unified tag styling, clickable tags, and inline edit functionality.
- **Statistics Dashboard Enhancements**: Reorganized layout, added Total Time metric, Hourly and Monthly Visitor Graphs, Full Internationalization, and Peak Concurrent Sessions Tracking.

### Technical Implementations
- **Backend (Go)**: Gorilla Mux HTTP router, RESTful endpoints, Nginx hardened security (XSS, CSP), role-based access control, Nginx performance optimizations (worker processes, epoll, sendfile, gzip, caching), server-side OSM tile caching.
- **Frontend (Vue.js)**: Vue.js 3 with Vite, Vue Router 4 for SPA navigation, session storage for filter conditions and map state, event clustering, reactive data fetching on locale changes, AND logic for multi-tag filtering.
- **Vite Proxy Setup**: Configured Vite dev server to proxy API requests to the backend for seamless communication in the Replit environment.
- **Session Tracking**: Lightweight implementation for both authenticated and anonymous users with a 60-second heartbeat and a statistics dashboard. UUID-based anonymous tracking for privacy. Active window: 5 minutes - sessions become inactive 5 minutes after last heartbeat.
- **Tag Search**: Dynamic search input for tags from visible events with autocomplete, suggestions, and filtering.
- **Service Worker**: Cleanup file unregisters any lingering service workers for tile caching.
- **Anonymous Session Duration Tracking**: Uses `last_seen_at` to calculate session duration.

### Feature Specifications
- **Interactive World Map**: Displays event markers with hover details, click-to-add functionality, and narrative flow visualization between filtered events.
- **Timeline Filtering**: Date range selection (FROM/TO fields) and historical period templates with chronological sorting.
- **Event Management**: CRUD operations for events via an admin panel, including BC date editing.
- **Lens System**: Categorization system for different event types.
- **Tag Management**: CRUD for tags, with cross-page navigation to filter events by tag.
- **User Management**: Admin interface for managing user accounts and access levels (Guest, User, Editor, Admin, Super).
- **Dataset Management**: Import and creation of localized event datasets.
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
- `vue`: Progressive JavaScript framework.
- `vite`: Build tool and dev server.
- `@vitejs/plugin-vue`: Vue plugin for Vite.
- `vue-router`: Official router for Vue.js SPA navigation.
- `leaflet`: Interactive map library for world map visualization.
- `vue-leaflet`: Vue.js wrapper for Leaflet maps.

### Other Integrations
- **OpenStreetMap**: Used for map tiles, with server-side caching via Nginx.