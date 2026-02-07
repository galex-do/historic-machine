<template>
  <div class="events-section">
    <!-- Events Header with Filters -->
    <div class="events-header">
      <div class="events-info">
        <!-- Merged Pagination with Count -->
        <div class="pagination-controls">
          <button 
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1 || totalPages <= 1"
            class="pagination-btn"
          >
            ‹
          </button>
          
          <span class="page-info">
            {{ paginationDisplay }}
          </span>
          
          <button 
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages || totalPages <= 1"
            class="pagination-btn"
          >
            ›
          </button>
        </div>
      </div>
      <div class="events-filters">
        <button 
          @click="handleShare"
          class="filter-btn share-btn"
          :title="shareCopied ? t('shareCopied') : t('shareUrl')"
        >
          {{ shareCopied ? '✓' : '🔗' }}
        </button>
        <button 
          @click="openTimeline"
          class="filter-btn timeline-btn"
          :disabled="events.length === 0"
          :title="t('timelineView')"
        >
          📅
        </button>
        <button 
          @click="toggleTagFilter"
          class="filter-btn tag-filter"
          :class="{ 'active': tagFilterVisible && selectedTags.length > 0 }"
          :title="tagFilterVisible ? 'Hide tag filters' : 'Show tag filters'"
        >
          🏷️
        </button>
        <button 
          @click="toggleMapFilter"
          class="filter-btn map-filter"
          :class="{ 'active': mapFilterEnabled }"
          title="Filter events by current map view"
        >
          🗺️
        </button>
        <button 
          @click="requestGeolocation"
          class="filter-btn geolocation-btn"
          :class="{ 'loading': geolocationLoading }"
          :disabled="geolocationLoading"
          :title="geolocationLoading ? t('geolocationLoading') : t('geolocationButton')"
        >
          {{ geolocationLoading ? '⏳' : '📍' }}
        </button>
      </div>
    </div>

    <!-- Tag Filter Panel (Collapsible) -->
    <TagFilterPanel
      v-if="tagFilterVisible"
      :selected-tags="selectedTags"
      :available-tags="availableTags"
      :follow-enabled="followEnabled"
      @remove-tag="$emit('remove-tag', $event)"
      @clear-all-tags="$emit('clear-all-tags')"
      @toggle-follow="$emit('toggle-follow')"
      @focus-on-filtered="$emit('focus-on-filtered')"
      @add-tag="$emit('tag-clicked', $event)"
    />

    <div 
      class="events-grid"
      @wheel="handleWheel"
      ref="eventsGrid"
    >
      <div v-if="events.length === 0" class="no-events">
        <p>No events found for the selected period. Click on the map to add your first historical event!</p>
      </div>
      
      <EventCard
        v-for="event in paginatedEvents"
        :key="event.id"
        :event="event"
        :map-filter-enabled="mapFilterEnabled"
        @focus-event="$emit('focus-event', $event)"
        @highlight-event="$emit('highlight-event', $event)"
        @tag-clicked="$emit('tag-clicked', $event)"
        @show-detail="openEventDetail"
      />
    </div>

    <!-- Timeline Modal -->
    <TimelineModal
      :is-open="timelineModalOpen"
      :events="events"
      :preserve-scroll="isBackNavigation"
      @close="timelineModalOpen = false"
      @focus-event="$emit('focus-event', $event)"
      @tag-clicked="$emit('tag-clicked', $event)"
      @show-detail="handleTimelineShowDetail"
    />

    <!-- Event Detail Modal -->
    <EventDetailModal
      :is-open="eventDetailModalOpen"
      :event="selectedDetailEvent"
      :all-events="events"
      :navigation-source="navigationSource"
      @close="closeEventDetail"
      @focus-event="$emit('focus-event', $event)"
      @tag-clicked="$emit('tag-clicked', $event)"
      @select-event="handleSelectRelatedEvent"
      @edit-event="handleEditEvent"
      @back="handleBack"
    />
  </div>
</template>

<script>
import { useLocale } from '@/composables/useLocale.js'
import { useGeolocation } from '@/composables/useGeolocation.js'
import EventCard from './EventCard.vue'
import TagFilterPanel from '../filters/TagFilterPanel.vue'
import TimelineModal from '../timeline/TimelineModal.vue'
import EventDetailModal from './EventDetailModal.vue'

export default {
  name: 'EventsGrid',
  components: {
    EventCard,
    TagFilterPanel,
    TimelineModal,
    EventDetailModal
  },
  props: {
    events: {
      type: Array,
      default: () => []
    },
    selectedTags: {
      type: Array,
      default: () => []
    },
    followEnabled: {
      type: Boolean,
      default: false
    },
    mapFilterEnabled: {
      type: Boolean,
      default: false
    },
    shareCopied: {
      type: Boolean,
      default: false
    }
  },
  emits: ['focus-event', 'highlight-event', 'map-filter-toggle', 'tag-clicked', 'remove-tag', 'clear-all-tags', 'toggle-follow', 'focus-on-filtered', 'geolocate', 'share', 'edit-event', 'back-to-location'],
  setup() {
    const { t } = useLocale()
    const { loading: geolocationLoading, get_current_position } = useGeolocation()
    return { t, geolocationLoading, get_current_position }
  },
  data() {
    const STORAGE_KEY = 'historia_tag_filter_visible'
    const loadTagFilterState = () => {
      try {
        const stored = sessionStorage.getItem(STORAGE_KEY)
        return stored ? JSON.parse(stored) : true // Default to visible
      } catch (error) {
        return true
      }
    }

    return {
      currentPage: 1,
      eventsPerPage: 3, // Maximum 3 cards to prevent sidebar scrolling
      tagFilterVisible: loadTagFilterState(),
      timelineModalOpen: false,
      eventDetailModalOpen: false,
      selectedDetailEvent: {},
      navigationSource: null,
      isBackNavigation: false,
      STORAGE_KEY,
      scrollThrottleTimer: null,
      scrollThrottleDelay: 150 // ms between page changes
    }
  },
  computed: {
    totalPages() {
      return Math.ceil(this.events.length / this.eventsPerPage)
    },
    paginatedEvents() {
      const start = (this.currentPage - 1) * this.eventsPerPage
      const end = start + this.eventsPerPage
      return this.events.slice(start, end)
    },
    paginationDisplay() {
      if (this.events.length === 0) {
        return '0'
      }
      const start = (this.currentPage - 1) * this.eventsPerPage + 1
      const end = Math.min(this.currentPage * this.eventsPerPage, this.events.length)
      return `${start}-${end}/${this.events.length}`
    },
    availableTags() {
      // Extract all unique tags from currently displayed events with usage counts
      const tagCountMap = new Map()
      
      this.events.forEach(event => {
        if (event.tags && Array.isArray(event.tags)) {
          event.tags.forEach(tag => {
            const existing = tagCountMap.get(tag.id)
            if (existing) {
              existing.count++
            } else {
              tagCountMap.set(tag.id, { ...tag, count: 1 })
            }
          })
        }
      })
      
      // Convert to array and sort by frequency (highest first)
      return Array.from(tagCountMap.values()).sort((a, b) => 
        b.count - a.count
      )
    }
  },
  watch: {
    events(newEvents, oldEvents) {
      // Only reset to first page if current page is now invalid
      // This prevents resetting when map view filter changes the events array
      const newTotalPages = Math.ceil(newEvents.length / this.eventsPerPage)
      
      // If current page is beyond the new total pages, reset to the last valid page
      if (this.currentPage > newTotalPages && newTotalPages > 0) {
        this.currentPage = newTotalPages
      } else if (newTotalPages === 0) {
        this.currentPage = 1
      }
      // Otherwise, keep the current page to maintain pagination state
    },
    selectedTags(newTags, oldTags) {
      // Automatically show tag filter panel when a tag is added
      if (newTags.length > oldTags.length) {
        this.tagFilterVisible = true
      }
    },
    tagFilterVisible(newValue) {
      // Save to session storage
      try {
        sessionStorage.setItem(this.STORAGE_KEY, JSON.stringify(newValue))
      } catch (error) {
        console.warn('Error saving tag filter visibility state:', error)
      }
    }
  },
  methods: {
    goToPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page
      }
    },
    handleShare() {
      this.$emit('share')
    },
    toggleMapFilter() {
      this.$emit('map-filter-toggle', !this.mapFilterEnabled)
    },
    toggleTagFilter() {
      this.tagFilterVisible = !this.tagFilterVisible
    },
    openTimeline() {
      this.isBackNavigation = false
      this.timelineModalOpen = true
    },
    openEventDetail(event, source = null) {
      this.selectedDetailEvent = event
      this.navigationSource = source
      this.eventDetailModalOpen = true
    },
    closeEventDetail() {
      this.eventDetailModalOpen = false
      this.navigationSource = null
    },
    handleSelectRelatedEvent(event) {
      this.selectedDetailEvent = event
    },
    handleTimelineShowDetail(event) {
      this.selectedDetailEvent = event
      this.navigationSource = 'timeline'
      this.timelineModalOpen = false
      this.eventDetailModalOpen = true
    },
    handleBack() {
      if (this.navigationSource === 'timeline') {
        this.isBackNavigation = true
        this.timelineModalOpen = true
      } else if (this.navigationSource === 'location') {
        this.$emit('back-to-location')
      }
      this.navigationSource = null
    },
    handleEditEvent(event) {
      this.eventDetailModalOpen = false
      this.$emit('edit-event', event)
    },
    async requestGeolocation() {
      try {
        const coords = await this.get_current_position()
        this.$emit('geolocate', coords)
      } catch (error) {
        alert(error.message)
      }
    },
    handleKeydown(e) {
      // Skip if a modal is open or user is typing in an input
      if (this.timelineModalOpen || this.eventDetailModalOpen) return
      if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return
      
      if (e.key === 'ArrowLeft') {
        e.preventDefault()
        this.goToPage(this.currentPage - 1)
      } else if (e.key === 'ArrowRight') {
        e.preventDefault()
        this.goToPage(this.currentPage + 1)
      }
    },
    handleWheel(e) {
      // Skip if no events or only one page
      if (this.totalPages <= 1) return
      
      // Skip if a modal is open
      if (this.timelineModalOpen || this.eventDetailModalOpen) return
      
      // Throttle scroll events to prevent too rapid page changes
      if (this.scrollThrottleTimer) return
      
      // Determine scroll direction (deltaY > 0 = scroll down = next page)
      const direction = e.deltaY > 0 ? 1 : -1
      const newPage = this.currentPage + direction
      
      // Only process if page change is valid
      if (newPage >= 1 && newPage <= this.totalPages) {
        e.preventDefault()
        this.goToPage(newPage)
        
        // Set throttle timer
        this.scrollThrottleTimer = setTimeout(() => {
          this.scrollThrottleTimer = null
        }, this.scrollThrottleDelay)
      }
    }
  },
  mounted() {
    document.addEventListener('keydown', this.handleKeydown)
  },
  beforeUnmount() {
    document.removeEventListener('keydown', this.handleKeydown)
    if (this.scrollThrottleTimer) {
      clearTimeout(this.scrollThrottleTimer)
    }
  }
}
</script>

<style scoped>
.events-section {
  background: transparent;
  padding: 1rem;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.events-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.events-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.events-count {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
  line-height: 1.25rem;
}

.events-filters {
  display: flex;
  gap: 0.5rem;
}

.filter-btn {
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 0.375rem 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s ease;
  color: #64748b;
}

.filter-btn:hover {
  background: #e2e8f0;
  border-color: #cbd5e1;
}

.filter-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
  box-shadow: 0 2px 4px rgba(59, 130, 246, 0.2);
}

.filter-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.map-filter {
  font-size: 1rem;
  line-height: 1;
}


.events-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  flex: 1;
  overflow-y: auto;
  min-height: 0; /* Allow flex item to shrink and scroll */
}

.no-events {
  text-align: center;
  padding: 2rem;
  color: #718096;
  font-style: italic;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.pagination-btn {
  background: transparent;
  color: #64748b;
  border: none;
  padding: 0;
  cursor: pointer;
  font-size: 1.25rem;
  transition: color 0.2s ease;
  min-width: auto;
  line-height: 1;
}

.pagination-btn:hover:not(:disabled) {
  color: #3b82f6;
}

.pagination-btn:disabled {
  color: #cbd5e1;
  cursor: not-allowed;
}

.page-info {
  color: #64748b;
  font-weight: 400;
  font-size: 0.75rem;
  white-space: nowrap;
  line-height: 1.25rem;
  display: flex;
  align-items: center;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .events-section {
    padding: 0.75rem;
  }
  
  /* Hide event list on mobile - show only counter and controls */
  .events-grid {
    display: none;
  }
  
  /* Hide pagination on mobile since we don't show event cards */
  .pagination-controls {
    display: none;
  }
}
</style>