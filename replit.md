# Historical Events Mapping Application

## Overview
A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data. The project aims to provide a rich, interactive experience for exploring history, with capabilities for managing events, tags, and datasets, all within a localized and performant environment. This application has significant market potential for educational platforms, historical research, and general public engagement with historical data.

## Recent Changes (November 14, 2025)

### Mobile-First Layout Optimization
- **Content-Prioritized Mobile Experience**: Complete redesign of mobile layout for optimal usability on small screens (< 768px)
  - **Events grid above map**: Events section appears first, map below (content-first approach using flexbox order)
  - **Compact control header only**: Shows only event counter (e.g., "52 events") + control buttons (📅 Timeline, 🏷️ Tags, 🗺️ Zone filter)
  - **Event cards hidden on mobile**: Individual event cards removed from mobile view to save screen space
  - **Pagination hidden**: No pagination controls on mobile since event cards are hidden
  - **Content-dependent height**: Events sidebar height set to `auto` - wraps only the controls, no fixed height
  - **Tag filter visibility**: Tag filter panel properly expands above map with `overflow: visible` - no clipping or overlay
  - **Map optimization**: Map gets 60vh+ of screen space for better geographical viewing
  - **Sidebar toggle hidden**: Collapse button (‹) hidden on mobile - events section always visible
  - **Result**: Clean, minimal header with controls + maximum map space for mobile users
  - Components modified: `EventsGrid.vue`, `MapView.vue`

### Visual Polish
- **Timeline date bullets centered**: Applied -6px margin-left to center 14px date bullets (10px + 2px border × 2) on 1px vertical timeline line for better visual alignment
- **Leaflet attribution removed**: Set `attributionControl: false` in WorldMap.vue for cleaner map appearance
- Component modified: `TimelineModal.vue`, `WorldMap.vue`

## Recent Changes (November 12, 2025)

### Compact Event List View
- **Timeline-Style Event Cards**: Complete redesign to use same minimalistic single-line layout as timeline modal
  - **Single-line text flow**: `🏛️ Event Name — Description #tag1 #tag2 ⌖` all flows naturally in one line
  - **Always-visible descriptions**: Full descriptions shown inline for maximum information density
  - **Hashtag tags**: Colored #tag format, clickable for filtering
  - **Date line**: Shown below main text for temporal context
  - **Minimal styling**: 0.5rem vertical padding, 1px bottom border separator
  - **Accessibility**: Semantic `<button>` elements for focus, ARIA labels, keyboard navigation, focus-visible styling
  - **Result**: ~70% height reduction, 3-4x more events visible, stable sidebar height → no map jumping
  - Component redesigned: `EventCard.vue`

### Russian Localization Optimization
- **Shortened Tag Filter Labels**: Reduced label text to fit single-line height in Russian
  - `filteredByTags`: "Фильтр по тегам" → "Теги"
  - `followEvents`: "Следовать" → "Связать"
  - `clearAllTags`, `clearAll`: "Очистить все" → "Очистить"
  - **Result**: Tag filter panel header fits cleanly on one line
  - File modified: `useLocale.js`

## Recent Changes (November 11, 2025)

### Timeline View Modal - Ultra-Compact Design
- **Vertical Timeline Visualization**: Maximizes event density on screen with compact, scannable layout
  - Timeline button (📅) in events-header, disabled when no events visible
  - **Ultra-compact single-line layout**: For single events on a date, everything appears on one continuous line: `○ Date 🏛️ Event Name — Description #tag1 #tag2`
  - **Multi-event layout**: Multiple events on same date show separate date header, then compact event list below
  - Unified typography: 0.875rem (14px) for all text, 1rem icons, 0.75rem tags
  - Reduced spacing: 0.25rem gaps between events (was 1rem), minimal padding throughout
  - Compact tag display: #hashtag format with colored text instead of pill badges
  - Timeline line: 1px (was 2px), 10px date bullets (was 20px)
  - BC/AD aware chronological sorting using timezone-agnostic manual date parsing
  - Expandable descriptions (▲/▼) for text >80 chars with 2-line clamp, improved accessibility contrast
  - ESC key support for closing modal, focus restoration to timeline button
  - Scrollable content area with slim 6px scrollbar
  - Full localization (English/Russian): "Timeline view", "No events to display"
  - Components: `TimelineModal.vue`, `EventsGrid.vue`
  - **Result**: 3-4x more events visible on screen compared to card-based layout
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