# Historical Events Mapping Application

## Overview
A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data. The project aims to provide a rich, interactive experience for exploring history, with capabilities for managing events, tags, and datasets, all within a localized and performant environment. This application has significant market potential for educational platforms, historical research, and general public engagement with historical data.

## Recent Changes (November 14, 2025)

### Statistics Dashboard Simplification
- **Removed Authenticated Users Section**: Simplified admin statistics page to show only Overall and Anonymous visitor metrics
  - Focus on anonymous visitor tracking for public historical events site
  - Overall section: Total active visitors (all users combined)
  - Anonymous section: Active visitors, total sessions, average duration
  - **Removed redundant "Active Sessions"**: For anonymous users, active visitors = active sessions (one UUID per visitor), so the duplicate metric was removed
  - Component modified: `AdminStats.vue`

### Session Tracking 5-Minute Active Window Fix
- **Consistent Active Session Logic**: Fixed inconsistency where "active sessions" never expired
  - Both "Active Visitors" and "Active Sessions" now use 5-minute window consistently
  - Sessions automatically become inactive 5 minutes after last heartbeat
  - Prevents inflated "active" counts when browser tabs are closed
  - Backend modified: `user_repository.go` - added 5-minute filter to active session queries

### Service Worker Cleanup
- **Removed Deprecated Tile Caching Service Worker**: Fixed periodic 404 errors for `/sw-tile-cache.js` in nginx logs
  - Previous version used a service worker for client-side tile caching (no longer needed with server-side nginx caching)
  - Created cleanup file at `frontend/public/sw-tile-cache.js` that unregisters any lingering service workers
  - Cleanup worker automatically unregisters itself and clears all caches when browser requests update
  - Eliminates 404 errors without affecting functionality

### Anonymous Session Duration Tracking Fix
- **Removed Unused `ended_at` Field**: Simplified anonymous session tracking by removing the never-set `ended_at` column
  - **Problem**: `ended_at` was never set (anonymous users don't explicitly "log out"), causing average duration to always show 0
  - **Solution**: Use `last_seen_at` instead to calculate session duration: `(last_seen_at - created_at)`
  - **Result**: Average duration now shows actual time from session start to last heartbeat
  - Database migration: `013_remove_anonymous_sessions_ended_at.sql`
  - Backend modified: `user_repository.go`, `auth_service.go`
  - Removed unused `EndAnonymousSession()` function

### Statistics Dashboard Enhancements (November 14, 2025)
- **Reorganized Stats Layout**: Moved session metrics to Overall section for clearer presentation
  - Overall section now shows: Total Active Visitors, Total Sessions, Average Duration, Total Time, Peak Concurrent Sessions
  - Removed redundant Anonymous Visitors section (all relevant data shown in Overall)
  - All metrics clearly labeled as "anonymous-only" except Peak (which tracks all sessions)
  
- **Added Total Time Metric**: Shows cumulative time for all anonymous sessions
  - Backend calculates: SUM of all session durations using `COALESCE(last_seen_at, NOW())`
  - Includes active sessions in calculation (uses NOW() as end time for ongoing sessions)
  - Displayed in localized format (hours/minutes)
  
- **Added Hourly Visitor Graph**: Interactive 24-hour bar chart showing visitor activity
  - Displays last 24 hours with hourly breakdown (oldest→newest, current hour on right)
  - Backend uses PostgreSQL generate_series for complete hour coverage
  - Active sessions included via `COALESCE(last_seen_at, NOW())`
  - Frontend features: CSS-based bars, hover tooltips, responsive scaling, mobile-optimized
  - Fully internationalized: all labels, tooltips, and duration units use i18n
  
- **Full Internationalization**: All new features support English and Russian
  - Added translation keys: overallTotalSessions, overallAvgDuration, overallTotalTime, hourlyVisitorsTitle
  - Duration formatting now uses localized units (min/мин, h/ч)
  - Visitor count pluralization (visitor/visitors, посетитель/посетителей)
  
- **Peak Concurrent Sessions Tracking**: Background service monitors and records maximum concurrent usage
  - Samples every 60 seconds to track peak concurrent sessions (authenticated + anonymous combined)
  - Background PeakStatsService runs as goroutine with graceful shutdown handling
  - Minimal database load: 2 queries/minute regardless of traffic volume
  - Displays lifetime peak in highlighted KPI card: "🏆 Peak Concurrent Sessions"
  
- **Files Modified**:
  - `backend/migrations/014_add_peak_stats.sql`: New peak_stats table for tracking maximum concurrency
  - `backend/internal/services/peak_stats_service.go`: Background goroutine with 60-second ticker
  - `backend/internal/models/user.go`: Added HourlyVisitorStat struct, AnonymousTotalTime, PeakConcurrentSessions fields
  - `backend/internal/database/repositories/user_repository.go`: Added total time, hourly stats, and peak tracking queries
  - `backend/main.go`: Integrated peak stats service with context-based graceful shutdown
  - `frontend/src/views/AdminStats.vue`: Reorganized layout, added bar graph and peak KPI card
  - `frontend/src/composables/useLocale.js`: Added all new translation keys for stats features

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

### Technical Implementations
- **Backend (Go)**: Gorilla Mux HTTP router, RESTful endpoints, Nginx hardened security (XSS, CSP), role-based access control, Nginx performance optimizations (worker processes, epoll, sendfile, gzip, caching), server-side OSM tile caching.
- **Frontend (Vue.js)**: Vue.js 3 with Vite, Vue Router 4 for SPA navigation, session storage for filter conditions and map state, event clustering, reactive data fetching on locale changes, AND logic for multi-tag filtering.
- **Vite Proxy Setup**: Configured Vite dev server to proxy API requests to the backend for seamless communication in the Replit environment.
- **Session Tracking**: Lightweight implementation for both authenticated and anonymous users with a 60-second heartbeat and a statistics dashboard. UUID-based anonymous tracking for privacy. **Active window: 5 minutes** - sessions become inactive 5 minutes after last heartbeat (when browser tab is closed or inactive).
- **Tag Search**: Dynamic search input for tags from visible events with autocomplete, suggestions, and filtering.

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
- `vue@3.5.20`: Progressive JavaScript framework.
- `vite@7.1.3`: Build tool and dev server.
- `@vitejs/plugin-vue@6.0.1`: Vue plugin for Vite.
- `vue-router@4.5.1`: Official router for Vue.js SPA navigation.
- `leaflet@1.9.4`: Interactive map library for world map visualization.
- `vue-leaflet@0.1.0`: Vue.js wrapper for Leaflet maps.

### Other Integrations
- **OpenStreetMap**: Used for map tiles, with server-side caching via Nginx.