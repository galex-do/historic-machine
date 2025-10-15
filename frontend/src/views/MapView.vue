<template>
  <div class="map-view">
    <!-- Enhanced Filter Bar -->
    <DateControlBar
      :date-from-display="dateFromDisplay"
      :date-to-display="dateToDisplay"
      :template-groups="templateGroups"
      :available-templates="availableTemplates"
      :selected-template-group-id="selectedTemplateGroupId"
      :selected-template-id="selectedTemplateId"
      :selected-template="selectedTemplate"
      @date-from-changed="handleDateFromChange"
      @date-to-changed="handleDateToChange"
      @template-group-changed="handleTemplateGroupChange"
      @template-changed="handleTemplateChange"
      @apply-filters="applyFilters"
    />
    
    <!-- Main Layout: Events Sidebar + Map -->
    <div class="main-layout">
      <!-- Left Events Sidebar (Collapsible) -->
      <aside class="events-sidebar" :class="{ 'collapsed': sidebarCollapsed }">
        <!-- Sidebar Header -->
        <div class="sidebar-header">
          <h3 class="section-title" v-show="!sidebarCollapsed">{{ t('historicalEvents') }}</h3>
          <button class="sidebar-toggle" @click="toggleSidebar">
            {{ sidebarCollapsed ? '›' : '‹' }}
          </button>
        </div>
        
        <!-- Events List in Sidebar -->
        <div class="events-sidebar-content" v-show="!sidebarCollapsed">
          <EventsGrid
            :events="displayedEvents"
            @focus-event="focusOnEvent"
            @map-filter-toggle="handleMapFilterToggle"
          />
        </div>
      </aside>
      
      <!-- Right Map Area (Full Height) -->
      <main class="map-content-area">
        <WorldMap 
          :events="filteredEvents" 
          :focus-event="focusEvent"
          @event-created="handleEventCreated"
          @event-updated="handleEventUpdated"
          @event-deleted="handleEventDeleted"
          @map-bounds-changed="handleMapBoundsChanged"
          @locale-changed="handleLocaleChanged"
          ref="worldMap"
        />
      </main>
    </div>
  </div>
</template>

<script>
import { ref, computed, nextTick, onMounted, watch } from 'vue'
import DateControlBar from '@/components/filters/DateControlBar.vue'
import WorldMap from '@/components/WorldMap.vue'
import EventsGrid from '@/components/events/EventsGrid.vue'
import { useEvents } from '@/composables/useEvents.js'
import { useTemplates } from '@/composables/useTemplates.js'
import { useFilters } from '@/composables/useFilters.js'
import { useTags } from '@/composables/useTags.js'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'MapView',
  components: {
    DateControlBar,
    WorldMap,
    EventsGrid
  },
  setup() {
    // Sidebar state with session storage
    const SIDEBAR_STORAGE_KEY = 'historia_sidebar_collapsed'
    const loadSidebarState = () => {
      try {
        const stored = sessionStorage.getItem(SIDEBAR_STORAGE_KEY)
        return stored ? JSON.parse(stored) : false
      } catch (error) {
        return false
      }
    }
    
    const sidebarCollapsed = ref(loadSidebarState())
    const focusEvent = ref(null)
    const worldMap = ref(null)
    
    // Map filter state
    const mapFilterEnabled = ref(false)
    const mapBounds = ref(null)
    
    // Save sidebar state to session storage
    watch(sidebarCollapsed, (newValue) => {
      try {
        sessionStorage.setItem(SIDEBAR_STORAGE_KEY, JSON.stringify(newValue))
      } catch (error) {
        console.warn('Error saving sidebar state:', error)
      }
    })

    // Composables
    const { 
      filteredEvents, 
      fetchEvents, 
      filterEvents, 
      handleEventCreated,
      handleEventUpdated,
      handleEventDeleted,
      eventsLoaded
    } = useEvents()

    const {
      templateGroups,
      availableTemplates,
      selectedTemplateGroupId,
      selectedTemplateId,
      selectedTemplate,
      fetchTemplateGroups,
      handleTemplateGroupChange: templateGroupChange,
      handleTemplateChange: templateChange,
      applyTemplate
    } = useTemplates()

    const {
      dateFrom,
      dateTo,
      dateFromDisplay,
      dateToDisplay,
      selectedLensTypes,
      showLensDropdown,
      resetToDefaultDateRange,
      updateDateFrom,
      updateDateTo,
      applyTemplateDates,
      toggleLensDropdown,
      handleLensTypesChange,
      closeLensDropdown
    } = useFilters()

    const {
      allTags,
      isLoadingTags,
      loadTags
    } = useTags()

    // Use locale for UI translations
    const { t } = useLocale()

    // Sidebar methods
    const toggleSidebar = () => {
      sidebarCollapsed.value = !sidebarCollapsed.value
      // Trigger map resize after sidebar animation
      nextTick(() => {
        setTimeout(() => {
          triggerMapResize()
        }, 300) // Match transition duration
      })
    }

    const triggerMapResize = () => {
      // Trigger a custom event that the map can listen for
      window.dispatchEvent(new Event('resize'))
    }

    // Template methods
    const handleTemplateGroupChange = async (groupId) => {
      await templateGroupChange(groupId)
      if (!groupId) {
        // Reset to default date range when "Default (1 AD - Today)" is selected
        resetToDefaultDateRange()
      } else if (groupId === 'custom') {
        // Custom mode - keep current dates, don't change anything
        console.log('Switched to Custom date range mode')
      }
      applyFilters()
    }

    const handleTemplateChange = (templateId) => {
      templateChange(templateId)
      // Always apply template dates when a template is selected
      const templateData = applyTemplate()
      if (templateData) {
        applyTemplateDates(templateData)
        console.log(`Selected template: ${selectedTemplate.value.name} (${templateData.displayFrom} - ${templateData.displayTo})`)
      }
      applyFilters()
    }
    
    // Enhanced date update methods that switch to custom mode
    const handleDateFromChange = (newValue, isStepping = false) => {
      // Switch to custom mode when user manually changes dates (but not when stepping)
      if (!isStepping && selectedTemplateGroupId.value && selectedTemplateGroupId.value !== 'custom') {
        templateGroupChange('custom')
        console.log('Switched to Custom mode due to manual date change')
      }
      updateDateFrom(newValue)
      applyFilters(isStepping)
    }
    
    const handleDateToChange = (newValue, isStepping = false) => {
      // Switch to custom mode when user manually changes dates (but not when stepping)
      if (!isStepping && selectedTemplateGroupId.value && selectedTemplateGroupId.value !== 'custom') {
        templateGroupChange('custom')
        console.log('Switched to Custom mode due to manual date change')
      }
      updateDateTo(newValue)
      applyFilters(isStepping)
    }

    // Filter methods
    const applyFilters = (isStepping = false) => {
      filterEvents(
        dateFrom.value,
        dateTo.value,
        selectedLensTypes.value,
        selectedTemplate.value,
        dateFromDisplay.value,
        dateToDisplay.value
      )
      
      // Signal WorldMap component about stepping state
      if (worldMap.value && worldMap.value.setSteppingMode) {
        worldMap.value.setSteppingMode(isStepping)
      }
    }

    // Event focus method
    const focusOnEvent = (event) => {
      focusEvent.value = { ...event, timestamp: Date.now() }
    }

    // Map filter methods
    const handleMapFilterToggle = (enabled) => {
      mapFilterEnabled.value = enabled
      if (enabled && worldMap.value) {
        // Get current map bounds when enabling filter
        worldMap.value.getCurrentBounds().then(bounds => {
          mapBounds.value = bounds
        })
      }
    }

    const handleMapBoundsChanged = (bounds) => {
      if (mapFilterEnabled.value) {
        mapBounds.value = bounds
      }
    }

    const isEventInBounds = (event, bounds) => {
      if (!bounds || !event.latitude || !event.longitude) {
        return true
      }
      
      const lat = parseFloat(event.latitude)
      const lng = parseFloat(event.longitude)
      
      return lat >= bounds.south && 
             lat <= bounds.north && 
             lng >= bounds.west && 
             lng <= bounds.east
    }

    // Computed property for events displayed in sidebar
    const displayedEvents = computed(() => {
      if (!mapFilterEnabled.value || !mapBounds.value) {
        return filteredEvents.value
      }
      
      return filteredEvents.value.filter(event => 
        isEventInBounds(event, mapBounds.value)
      )
    })

    // Click outside to close dropdown
    const handleClickOutside = (event) => {
      // Close lens dropdown if clicking outside
      if (showLensDropdown.value && !event.target.closest('.multiselect-container')) {
        closeLensDropdown()
      }
    }

    // Watch for events loading and apply filters immediately
    watch(eventsLoaded, (loaded) => {
      if (loaded) {
        console.log('Events loaded, applying filters immediately')
        applyFilters()
      }
    })

    // Handle locale changes
    const handleLocaleChanged = async (locale) => {
      console.log('Locale changed in MapView, refetching events for locale:', locale)
      try {
        await fetchEvents() // Refetch events with new locale
        // Reapply current filters to the newly fetched events
        applyFilters()
      } catch (err) {
        console.error('Error refetching events after locale change:', err)
      }
    }

    // Initialize data on mount
    onMounted(async () => {
      // Apply filters if events are already loaded
      if (eventsLoaded.value) {
        applyFilters()
      }
      
      // Add click outside listener
      document.addEventListener('click', handleClickOutside)
    })

    return {
      // Sidebar state
      sidebarCollapsed,
      focusEvent,

      // Events
      filteredEvents,
      handleEventCreated,
      handleEventUpdated,
      handleEventDeleted,
      handleLocaleChanged,

      // Templates
      templateGroups,
      availableTemplates,
      selectedTemplateGroupId,
      selectedTemplateId,
      selectedTemplate,
      handleTemplateGroupChange,
      handleTemplateChange,

      // Filters
      dateFromDisplay,
      dateToDisplay,
      selectedLensTypes,
      showLensDropdown,
      handleDateFromChange,
      handleDateToChange,
      toggleLensDropdown,
      handleLensTypesChange,

      // Tags
      allTags,
      isLoadingTags,

      // Methods
      toggleSidebar,
      applyFilters,
      focusOnEvent,

      // Map filter
      displayedEvents,
      handleMapFilterToggle,
      handleMapBoundsChanged,
      worldMap,

      // Localization
      t
    }
  }
}
</script>

<style scoped>
.map-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 80px); /* Account for header */
}

/* Main Layout */
.main-layout {
  display: flex;
  flex: 1;
  height: calc(100vh - 160px); /* Account for header + enhanced filter bar */
  overflow: hidden;
}

/* Events Sidebar */
.events-sidebar {
  background: #ffffff;
  color: #2d3748;
  position: relative;
  transition: all 0.3s ease;
  width: 350px;
  flex-shrink: 0;
  overflow: hidden;
  border-right: 1px solid #e2e8f0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  z-index: 5;
}

.events-sidebar.collapsed {
  width: 60px;
}

.events-sidebar.collapsed .sidebar-header {
  padding: 1.5rem 0.5rem 1rem;
  justify-content: center;
}

.events-sidebar.collapsed .sidebar-toggle {
  margin: 0;
}

.sidebar-header {
  background: #f8f9fa;
  border-bottom: 1px solid #e2e8f0;
  padding: 1.5rem 1.5rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.section-title {
  color: #2d3748;
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.sidebar-toggle {
  width: 28px;
  height: 28px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #4a5568;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  margin-left: auto;
}

.sidebar-toggle:hover {
  background: #f1f5f9;
  color: #2d3748;
  border-color: #cbd5e0;
}

.events-sidebar-content {
  flex: 1;
  overflow-y: auto;
}

/* Map Content Area (Full Space) */
.map-content-area {
  flex: 1;
  background: #ffffff;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  z-index: 1;
}

/* Responsive Design */
@media (max-width: 1024px) {
  .main-layout {
    flex-direction: column;
    height: calc(100vh - 160px);
  }
  
  .events-sidebar {
    width: 100%;
    height: 40vh;
    order: 2;
  }
  
  .events-sidebar.collapsed {
    height: 60px;
    width: 100%;
  }
  
  .map-content-area {
    flex: 1;
    order: 1;
    min-height: 50vh;
  }
}

@media (max-width: 768px) {
  .main-layout {
    height: calc(100vh - 140px); /* Account for potentially smaller mobile header */
  }
  
  .events-sidebar {
    height: 35vh;
  }
}
</style>