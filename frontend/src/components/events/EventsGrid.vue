<template>
  <div class="events-section">
    <!-- Events Header with Filters -->
    <div class="events-header">
      <div class="events-info">
        <div class="events-count">
          {{ events.length }} event{{ events.length !== 1 ? 's' : '' }}
        </div>
        <!-- Pagination Controls -->
        <div v-if="totalPages > 1" class="pagination-controls">
          <button 
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="pagination-btn"
          >
            ‹
          </button>
          
          <span class="page-info">
            {{ currentPage }}/{{ totalPages }}
          </span>
          
          <button 
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="pagination-btn"
          >
            ›
          </button>
        </div>
      </div>
      <div class="events-filters">
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
      @add-tag="$emit('tag-clicked', $event)"
    />

    <div class="events-grid">
      <div v-if="events.length === 0" class="no-events">
        <p>No events found for the selected period. Click on the map to add your first historical event!</p>
      </div>
      
      <EventCard
        v-for="event in paginatedEvents"
        :key="event.id"
        :event="event"
        @focus-event="$emit('focus-event', $event)"
        @tag-clicked="$emit('tag-clicked', $event)"
      />
    </div>

    <!-- Timeline Modal -->
    <TimelineModal
      :is-open="timelineModalOpen"
      :events="events"
      @close="timelineModalOpen = false"
    />
  </div>
</template>

<script>
import { useLocale } from '@/composables/useLocale.js'
import EventCard from './EventCard.vue'
import TagFilterPanel from '../filters/TagFilterPanel.vue'
import TimelineModal from '../timeline/TimelineModal.vue'

export default {
  name: 'EventsGrid',
  components: {
    EventCard,
    TagFilterPanel,
    TimelineModal
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
    }
  },
  emits: ['focus-event', 'map-filter-toggle', 'tag-clicked', 'remove-tag', 'clear-all-tags', 'toggle-follow'],
  setup() {
    // Expose translation function from composable
    const { t } = useLocale()
    return { t }
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
      mapFilterEnabled: false,
      tagFilterVisible: loadTagFilterState(),
      timelineModalOpen: false,
      STORAGE_KEY
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
    availableTags() {
      // Extract all unique tags from currently displayed events
      const tagMap = new Map()
      
      this.events.forEach(event => {
        if (event.tags && Array.isArray(event.tags)) {
          event.tags.forEach(tag => {
            if (!tagMap.has(tag.id)) {
              tagMap.set(tag.id, tag)
            }
          })
        }
      })
      
      // Convert to array and sort alphabetically
      return Array.from(tagMap.values()).sort((a, b) => 
        a.name.localeCompare(b.name)
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
    toggleMapFilter() {
      this.mapFilterEnabled = !this.mapFilterEnabled
      this.$emit('map-filter-toggle', this.mapFilterEnabled)
    },
    toggleTagFilter() {
      this.tagFilterVisible = !this.tagFilterVisible
    },
    openTimeline() {
      this.timelineModalOpen = true
    }
  }
}
</script>

<style scoped>
.events-section {
  background: transparent;
  padding: 1rem;
  height: 100%;
  overflow-y: auto;
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
}
</style>