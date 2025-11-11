# Historical Events Mapping Application

## Overview
A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data. The project aims to provide a rich, interactive experience for exploring history, with capabilities for managing events, tags, and datasets, all within a localized and performant environment. This application has significant market potential for educational platforms, historical research, and general public engagement with historical data.

## Recent Changes (November 11, 2025)

### Timeline View Modal
- **Vertical Timeline Visualization**: Large modal displaying filtered events in chronological order
  - Timeline button (📅) in events-header, disabled when no events visible
  - Vertical timeline with connecting line and date bullets
  - Event grouping by date (day precision) with formatted date headers
  - BC/AD aware chronological sorting using timezone-agnostic manual date parsing
  - Each event shows: lens type emoji, name, expandable description, colored tag badges
  - "Show more"/"Show less" for long descriptions (>150 chars, 3-line clamp)
  - ESC key support for closing modal
  - Focus restoration to timeline button after closing
  - Scrollable content area with custom scrollbar styling
  - Full localization (English/Russian): "Timeline view", "Show more", "Show less", "No events to display"
  - Components: `TimelineModal.vue`, `EventsGrid.vue`
  - Critical implementation: Uses manual ISO date parsing (`split('-')`) for both BC and AD dates to prevent UTC timezone shifts that break chronological ordering and date grouping

### Tag Search Functionality
- **Dynamic Tag Search**: Search input field that adapts to currently displayed events
  - Autocomplete dropdown shows tags from visible events only (recalculates when event set changes)
  - Filters out already selected tags from suggestions
  - Shows up to 10 tag suggestions with colored left border indicators matching tag colors
  - Click-to-add functionality - clicking a tag adds it to filters and closes dropdown
  - Search query filtering - type to narrow down tag suggestions
  - Empty state message when no matching tags found ("No tags found")
  - Full localization support (English: "Search tags...", Russian: "Поиск тегов...")
  - Components modified: `TagFilterPanel.vue` (search UI + logic), `EventsGrid.vue` (availableTags computed)
  - Reactive behavior: tag list automatically updates when event set changes via date filter, map filter, or tag filters

### Tag Filter Panel UI Reorganization
- **Improved Events Sidebar Layout**: Collapsible tag filtering with better visual hierarchy
  - Moved TagFilterPanel from above EventsGrid to below events-header (after event count/pagination)
  - Added tag filter toggle button (🏷️) in events-header next to map filter button (🗺️)
  - Toggle button shows active state when filters are visible and tags are selected
  - Tag filter panel now collapsible for cleaner UI when not needed
  - Visibility state persisted in session storage (key: `historia_tag_filter_visible`, default: visible)
  - Better visual grouping: all header actions (tag filter, map filter) together
  - All tag filtering functionality preserved: click, remove, clear all, narrative flow
  - Components modified: `EventsGrid.vue`, `MapView.vue`

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
- **Timeline View Modal**: Vertical chronological visualization of filtered events with date grouping, expandable descriptions, BC/AD aware sorting, and accessibility features (ESC key, focus management).
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