<template>
  <div id="app">
    <!-- Header -->
    <AppHeader />
    
    <!-- Main Layout: Sidebar + Map -->
    <div class="main-layout">
      <!-- Left Sidebar (Collapsible) -->
      <SidebarFilters
        :collapsed="sidebarCollapsed"
        :template-groups="templateGroups"
        :available-templates="availableTemplates"
        :selected-template-group-id="selectedTemplateGroupId"
        :selected-template-id="selectedTemplateId"
        :selected-template="selectedTemplate"
        :selected-lens-types="selectedLensTypes"
        :show-lens-dropdown="showLensDropdown"
        @toggle="toggleSidebar"
        @template-group-changed="handleTemplateGroupChange"
        @template-changed="handleTemplateChange"
        @toggle-lens-dropdown="toggleLensDropdown"
        @lens-types-changed="handleLensTypesChange"
        @apply-filters="applyFilters"
      />
      
      <!-- Right Content Area -->
      <main class="content-area">
        <!-- Date Control Bar Above Map -->
        <DateControlBar
          :date-from-display="dateFromDisplay"
          :date-to-display="dateToDisplay"
          @date-from-changed="updateDateFrom"
          @date-to-changed="updateDateTo"
        />
        
        <!-- Map Section -->
        <div class="map-section">
          <WorldMap 
            :events="filteredEvents" 
            :focus-event="focusEvent"
            @event-created="handleEventCreated"
          />
        </div>
        
        <!-- Events List Below Map -->
        <EventsGrid
          :events="filteredEvents"
          @focus-event="focusOnEvent"
        />
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, nextTick } from 'vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import SidebarFilters from '@/components/filters/SidebarFilters.vue'
import DateControlBar from '@/components/filters/DateControlBar.vue'
import WorldMap from '@/components/WorldMap.vue'
import EventsGrid from '@/components/events/EventsGrid.vue'
import { useEvents } from '@/composables/useEvents.js'
import { useTemplates } from '@/composables/useTemplates.js'
import { useFilters } from '@/composables/useFilters.js'

export default {
  name: 'App',
  components: {
    AppHeader,
    SidebarFilters,
    DateControlBar,
    WorldMap,
    EventsGrid
  },
  setup() {
    // Sidebar state
    const sidebarCollapsed = ref(false)
    const focusEvent = ref(null)

    // Composables
    const { 
      filteredEvents, 
      fetchEvents, 
      filterEvents, 
      handleEventCreated 
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

    // Filter methods
    const applyFilters = () => {
      filterEvents(
        dateFrom.value,
        dateTo.value,
        selectedLensTypes.value,
        selectedTemplate.value
      )
    }

    // Event focus method
    const focusOnEvent = (event) => {
      focusEvent.value = { ...event, timestamp: Date.now() }
    }

    // Click outside to close dropdown
    const handleClickOutside = (event) => {
      // Close lens dropdown if clicking outside
      if (showLensDropdown.value && !event.target.closest('.multiselect-container')) {
        closeLensDropdown()
      }
    }

    // Initialize data on mount
    onMounted(async () => {
      await Promise.all([
        fetchEvents(),
        fetchTemplateGroups()
      ])
      
      // Apply initial filters
      applyFilters()
      
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
      updateDateFrom,
      updateDateTo,
      toggleLensDropdown,
      handleLensTypesChange,

      // Methods
      toggleSidebar,
      applyFilters,
      focusOnEvent
    }
  }
}
</script>

<style>
/* Global app styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f7fafc;
  color: #2d3748;
  line-height: 1.6;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Main Layout */
.main-layout {
  display: flex;
  flex: 1;
  height: calc(100vh - 80px);
  overflow: hidden;
}

/* Content Area */
.content-area {
  flex: 1;
  background: #ffffff;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
}

/* Map Section */
.map-section {
  flex: 1;
  position: relative;
  z-index: 1;
}

/* Responsive Design */
@media (max-width: 1024px) {
  .main-layout {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    max-height: 60vh;
    overflow-y: auto;
  }
  
  .sidebar.collapsed {
    height: 60px;
    max-height: 60px;
  }
  
  .content-area {
    flex: 1;
    min-height: 40vh;
  }
}

@media (max-width: 768px) {
  .main-layout {
    height: calc(100vh - 60px);
  }
}
</style>