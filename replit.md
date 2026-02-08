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
- **Interactive World Map**: Utilizes Leaflet for scalable and interactive map visualization with event markers, zoom-dependent clustering, and narrative flow visualization.
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
- **Map Event Info Modal Timeline Structure**: Uses the same timeline design as TimelineModal with vertical line, bullets, date grouping, unified tag styling, clickable tags, inline edit functionality, minified/detailed view toggle (📋/📝), and lazy loading in batches of 50 with scroll-based loading.
- **Pin Mode for Coordinates**: Event form includes a pin button (📍) that closes the modal and enables map-click coordinate selection. Supports "sticky" pin - clicking existing markers/clusters copies their coordinates. Zoom/pan don't interrupt pin mode.

### Technical Implementations
- **Backend (Go)**: Gorilla Mux HTTP router, RESTful endpoints, Nginx hardened security (XSS, CSP), role-based access control, Nginx performance optimizations (worker processes, epoll, sendfile, gzip, caching), server-side OSM tile caching.
- **Frontend (Vue.js)**: Vue.js 3 with Vite, Vue Router 4 for SPA navigation, session storage for filter conditions and map state, zoom-dependent marker clustering (Leaflet MarkerCluster), reactive data fetching on locale changes, AND logic for multi-tag filtering.
- **Shared CSS Architecture**: Centralized reusable styles in `frontend/src/styles/` — `tag-badge.css` (tag pill badges), `timeline.css` (vertical timeline layout), `modal-overlay.css` (modal backdrop). Imported via non-scoped style blocks; component-specific overrides remain in scoped styles. Unified `getTagStyle(tag)` utility in `color-utils.js` handles background, contrast text, and optional inner border for all tag renderings.
- **Vite Proxy Setup**: Configured Vite dev server to proxy API requests to the backend for seamless communication in the Replit environment.
- **Session Tracking**: Lightweight implementation for both authenticated and anonymous users with a 60-second heartbeat. UUID-based anonymous tracking for privacy. Active window: 5 minutes - sessions become inactive 5 minutes after last heartbeat.
- **Prometheus Metrics**: `/metrics` endpoint exposes application metrics in Prometheus format. Includes HTTP request counts/durations, event CRUD counters, login attempts, active sessions, and database entity gauges (events, users, tags, datasets, templates). Metrics collector runs every 30 seconds.
- **Tag Search**: Dynamic search input for tags from visible events with autocomplete, suggestions, and filtering.
- **Progressive Web App (PWA)**: Full PWA support with installability, offline caching for app shell and map tiles, and service worker for network resilience. Includes manifest.json with app icons (72px-512px) and theme colors.
- **Anonymous Session Duration Tracking**: Uses `last_seen_at` to calculate session duration.

### Feature Specifications
- **Interactive World Map**: Displays event markers with hover details, click-to-add functionality, and narrative flow visualization between filtered events.
- **Browser Geolocation**: "Find my location" button (📍) in the events panel header allows users to center and zoom the map on their current position using the browser's Geolocation API. Includes proper error handling for permission denied, unavailable, and timeout scenarios with localized messages.
- **URL Sharing**: Share button (🔗) in events-header copies a URL with current filter state (date range, selected tags, map center/zoom) to clipboard. URLs restore full application state when opened, enabling bookmarking and sharing specific views. URL updates automatically as user navigates.
- **Compact Pagination**: Merged event counter into pagination display (e.g., "1-3/46") for space efficiency in the events-header.
- **Timeline Filtering**: Date range selection (FROM/TO fields) and historical period templates with chronological sorting.
- **Event Management**: CRUD operations for events via an admin panel, including BC date editing.
- **Lens System**: Categorization system for different event types.
- **Tag Management**: CRUD for tags, with cross-page navigation to filter events by tag. Optional `border_color` field for inner border customization using `box-shadow: inset` (no height change). Admin panel supports enable/disable border with color picker. `key_color` boolean flag (false by default) — when true, a small colored dot of the tag's color appears next to event names in timeline/cluster minimalistic views for quick nation/category identification without opening event details.
- **Template Management**: Admin interface for CRUD operations on date templates and template groups. Features dual modal system, BC/AD date handling with chronological sorting, and full localization (English/Russian). Template groups prevent deletion when containing templates.
- **User Management**: Admin interface for managing user accounts and access levels (Guest, User, Editor, Admin, Super).
- **Dataset Management**: Import and creation of localized event datasets. Includes modification tracking - datasets are automatically marked as modified when associated events are created, updated, or deleted. Admin panel shows warning icon for modified datasets with reset button to clear the flag after export.
- **Authentication**: Role-based access control, session persistence.

### System Design Choices
- **Database**: PostgreSQL with Goose migration tool.
- **API**: CORS enabled for frontend integration.
- **Deployment**: Docker Compose for orchestrating backend, frontend, and database services.

## External Dependencies

### Backend Dependencies
- `github.com/gorilla/mux`: HTTP router and URL matcher.
- `github.com/prometheus/client_golang`: Prometheus metrics client for exposing application metrics.
- `PostgreSQL driver`: Database connectivity.
- `Goose`: Database migration tool.

### Frontend Dependencies
- `vue`: Progressive JavaScript framework.
- `vite`: Build tool and dev server.
- `@vitejs/plugin-vue`: Vue plugin for Vite.
- `vue-router`: Official router for Vue.js SPA navigation.
- `leaflet`: Interactive map library for world map visualization.
- `leaflet.markercluster`: Zoom-dependent marker clustering plugin for improved map UX at different zoom levels.
- `vue-leaflet`: Vue.js wrapper for Leaflet maps.

### Other Integrations
- **OpenStreetMap**: Used for map tiles, with server-side caching via Nginx.