# Historical Events Mapping Application

## Overview
A comprehensive web application for mapping historical events on an interactive world map with timeline functionality. Users can view, filter, and add historical events with geographical and temporal data. The project aims to provide a rich, interactive experience for exploring history, with capabilities for managing events, tags, and datasets, all within a localized and performant environment. This application has significant market potential for educational platforms, historical research, and general public engagement with historical data.

## Recent Changes (October 23, 2025)

- **Map Marker Tooltips**: Simple tooltips appear when hovering over event markers on map
  - Single events: Shows formatted date + event name (e.g., "01.01.332 BC - Conquest of Alexandria")
  - Clustered events: Shows count (e.g., "5 events found here")
  - Tooltips appear above markers with 90% opacity for clear visibility
  - Uses Leaflet's built-in bindTooltip for smooth hover experience
  - Fixed: Proper tooltip cleanup prevents null reference errors when zooming (unbindTooltip before marker removal)
  - Component: `WorldMap.vue`
- **Clickable Tags in Event Popup**: Tags in map event popups are now clickable and add to tag filter panel
  - Click any tag in the event info modal to instantly filter events by that tag
  - Tags have hover effects (lift, brightness, shadow) to indicate they're interactive
  - Modal automatically closes after clicking a tag for immediate visual feedback
  - Tooltip shows "Click to filter events by '<tag name>'" on hover
  - Emits 'tag-clicked' event from WorldMap component to MapView for filter integration
  - Components: `WorldMap.vue`, `MapView.vue`
- **Chronological Template Sorting**: Date templates now sorted by start date instead of display order
  - Templates within each group now appear in chronological order (oldest to newest)
  - Example: "First Empires" group shows Ancient Kingdom of Egypt (2686 BC), then Akkadian Empire (2334 BC), then Old Assyrian Trade Empire (2025 BC), etc.
  - Makes historical progression more intuitive for users
  - Changed ORDER BY from `display_order` to `start_astronomical_year ASC` in template queries
  - File: `backend/internal/database/repositories/template_repository.go`
- **Assyrian Empire Dataset and Templates**: Comprehensive historical coverage of ancient Assyria added to the application
  - Created dataset with 57 events spanning 2600 BC - 609 BC (foundation of Assur through fall of empire)
  - Covers Old Assyrian period (trade colonies, Kanesh), Middle Assyrian rise (Ashur-uballit I, Tukulti-Ninurta I), and Neo-Assyrian golden age (Tiglath-Pileser III, Sargon II, Sennacherib, Ashurbanipal)
  - Includes major battles (Qarqar, Arpad, Lachish), conquests (Israel, Damascus, Babylon, Egypt), and cultural achievements (Library of Nineveh)
  - Added 5 historical period templates split across groups: "Old Assyrian Trade Empire" in "Bronze Age Powers", and "Middle Assyrian Rise", "Neo-Assyrian Early Expansion", "Neo-Assyrian Golden Age", "Fall of Assyria" in "Imperial Age"
  - Professional English/Russian translations for all events and templates
  - Avoids duplication with existing 7 Assyrian-tagged events in other datasets
  - Dataset: `datasets/assyria_empire.json`, Migration: `backend/migrations/009_add_assyrian_templates.sql`
- **Expandable Tags in Event Modal**: Event info modal now supports expandable tags functionality
  - Click "+N" button to show all tags inline (matching EventCard behavior)
  - "Show less" button appears to collapse back to 3 visible tags
  - Smooth hover effects and visual feedback for clickable tag controls
  - Tooltips show hidden tag names when hovering over "+N" indicator
  - State resets when modal closes for clean UX
  - Component: `WorldMap.vue`
- **Auto-Expand Date Range for Template Groups**: When selecting a template group without choosing a specific template, the date range automatically expands to encompass all templates in that group
  - Example: "Phoenician Mediterranean Empire" group auto-sets range to 1500 BC - 146 BC (full civilization span)
  - Example: "Ancient Greece" group auto-sets range to 800 BC - 31 BC (Archaic Period through Hellenistic Period)
  - Calculates minimum start date and maximum end date across all templates in the group
  - Correctly handles BC/AD era metadata from backend (negates BC years for proper chronological ordering)
  - Uses pre-formatted display dates from backend for accurate localized formatting
  - Enables quick exploration of entire historical periods with one click
  - Components: `useTemplates.js`, `MapView.vue`
- **Narrative Flow Visualization**: Interactive feature to visualize historical event progression
  - "Follow Events" toggle button appears in tag filter panel when tags are selected
  - Draws chronological connections between filtered events with gradient color polylines
  - Color gradient transitions from dark blue (#1e40af) at start to light blue (#93c5fd) at end
  - Gradient makes temporal direction clear when lines converge at same location multiple times
  - Arrow decorators at segment midpoints show travel/progression direction (from older to newer events)
  - Arrows match segment colors for consistent visual flow
  - Automatically skips connections between co-located events (same coordinates)
  - Proper chronological sorting using same BC/AD logic as event list (handles BC dates correctly)
  - Arrow bearing calculation using atan2(Δlat, Δlng) with 180° flip for correct temporal flow
  - Automatically disables when tag filters are removed or cleared
  - Reactive cleanup removes polylines when toggle disabled or filters change
  - Helps visualize historical themes (Hannibal's campaigns, Phoenician expansion, colonization patterns)
  - Components: `WorldMap.vue`, `MapView.vue`, `TagFilterPanel.vue`
- **Two-Pane Hierarchical Date Template Selector**: Replaced dual dropdowns with compact popover selector
  - Single button opens two-pane popover: template groups (left) + specific templates (right)
  - Click group to view its templates, click template to select and close
  - Saves significant filter bar space with improved UX flow
  - Full keyboard support (Escape to close), loading states, responsive design
  - Production-ready viewport positioning: measures actual rendered dimensions, clamps to 16px margins, flips upward when needed
  - Components: `HierarchicalDateTemplateSelector.vue`, updated `DateControlBar.vue`
- **Tag Filtering Logic**: Changed from OR to AND logic for multi-tag filtering
  - Selecting multiple tags now shows only events that have ALL selected tags
  - More intuitive filtering behavior for narrowing down specific events
  - Updated in both map view and admin events page
- **BC Date Editing Fix**: Fixed date input validation to accept 1-4 digit years
  - Admin Events edit form now accepts BC dates with 1-4 digit years (e.g., 146 BC, 44 BC, 3501 BC)
  - Updated regex pattern from `\d{4}` to `\d{1,4}` in `updateEventDate()` function
  - Year values automatically padded to 4 digits for consistency (146 → 0146)
  - File: `frontend/src/views/AdminEvents.vue`
- **Expandable Tags on Event Cards**: Enhanced UX for events with many tags
  - Click "+N" button on event cards to expand and show all tags inline
  - "Show less" button appears to collapse back to 3 visible tags
  - Smooth hover effects and visual feedback for clickable tag controls
  - Maintains all tag filtering and navigation functionality
- **BC Date Chronological Sorting Fix**: Fixed month/day ordering within BC years
  - Updated `astronomical_year` calculation to include month/day precision
  - BC dates now properly sort with earlier months appearing after later months
  - Example: August 332 BC now correctly appears before January 332 BC chronologically
  - Migration: `backend/migrations/008_fix_bc_month_ordering.sql`
- **Phoenician Mediterranean Empire Dataset**: Comprehensive historical dataset covering Phoenician civilization (enriched October 23, 2025)
  - **66 events** spanning 1500 BC to 146 BC (41 original + 25 enrichment events)
  - Covers city-states (Tyre, Sidon, Byblos), alphabet development, colonization, trade networks
  - Added Assyrian domination period (Shalmaneser III tribute 858-841 BC)
  - Added Babylonian siege of Tyre (Nebuchadnezzar's 13-year siege 586-573 BC)
  - Added technological innovations (glass making, bireme warships, transparent glass perfection)
  - Added Persian period events (Phoenician fleet for Persia, Battle of Salamis, Herodotus visit)
  - Expanded Punic Wars coverage (10 additional battles: Cape Ecnomus, Ticinus, Trebia, Trasimene, Metaurus, etc.)
  - Includes Carthaginian ascendancy, Punic Wars, Hannibal's campaigns, economic recovery
  - Professional English/Russian translations for all events
  - Strategic tagging across political, military, cultural, technological, and religious categories
  - File: `datasets/phoenician_mediterranean_empire.json`
- **Phoenician Historical Period Templates**: New date template group for Phoenician civilization
  - Added "Phoenician Mediterranean Empire" template group to Historical Period selector
  - 7 period templates covering full Phoenician timeline (Early City-States, Renaissance, Colonization, Carthaginian Ascendancy, Persian Period, Punic Wars, Hannibalic War)
  - Fully localized in English and Russian
  - Migration: `backend/migrations/007_add_phoenician_templates.sql`

## User Preferences
- Use snake_case naming convention everywhere for elements and functions
- Develop professional, high-end code with proper patterns and templates
- Avoid duplication - refactor similar logic into reusable templates/functions
- Simple, everyday language for communication

## System Architecture

### UI/UX Decisions
- **Interactive World Map**: Utilizes Leaflet for scalable and interactive map visualization.
- **Admin Panel**: Professional table-based interface for managing events, tags, users, and datasets.
- **Localization**: Full internationalization (English/Russian) for all UI elements, historical period templates, and event data. Reactive locale switching for instant updates.
- **Filter Panel**: Kibana-style removable chip interface for tag filtering with session storage persistence.
- **Pagination**: Minimalistic pagination moved to the header for constant visibility.
- **Date Handling**: BC/AD era selectors in forms, localized date formatting, and timezone-safe displays.

### Technical Implementations
- **Backend (Go)**:
    - **Framework**: Gorilla Mux HTTP router.
    - **API Pattern**: RESTful endpoints.
    - **Security**: Editor-level authentication for CRUD operations, Nginx hardening (XSS protection, CSP, buffer limits), and role-based access control.
    - **Performance**: Nginx optimizations (worker processes, epoll, sendfile, gzip, caching), server-side OSM tile caching.
    - **Containerization**: Multi-stage Docker builds for reduced image size, Docker Compose setup.
- **Frontend (Vue.js)**:
    - **Framework**: Vue.js 3 with Vite build system.
    - **Routing**: Vue Router 4 for SPA navigation (Map view, Admin panel).
    - **State Management**: Session storage for filter conditions and map state persistence.
    - **Event Clustering**: Co-located events clustering with count badges.
    - **Reactive Data Fetching**: Event and template data re-fetches automatically on locale changes.

### Feature Specifications
- **Interactive World Map**: Displays event markers with hover details and click-to-add functionality.
- **Timeline Filtering**: Date range selection (FROM/TO fields) and historical period templates.
- **Event Management**: CRUD operations for events via an admin panel.
- **Lens System**: Categorization system for different event types.
- **Tag Management**: CRUD for tags, with cross-page navigation to filter events by tag.
- **User Management**: Admin interface for managing user accounts and access levels (Guest, User, Editor, Admin, Super).
- **Dataset Management**: Import and creation of localized event datasets (e.g., Phoenician Mediterranean Empire, ancient civilizations).
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