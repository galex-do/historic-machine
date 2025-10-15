<template>
  <div class="events-section">
    <!-- Events Header with Filters -->
    <div class="events-header">
      <div class="events-count">
        {{ events.length }} event{{ events.length !== 1 ? 's' : '' }}
      </div>
      <div class="events-filters">
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

    <div class="events-grid">
      <div v-if="events.length === 0" class="no-events">
        <p>No events found for the selected period. Click on the map to add your first historical event!</p>
      </div>
      
      <EventCard
        v-for="event in paginatedEvents"
        :key="event.id"
        :event="event"
        @focus-event="$emit('focus-event', $event)"
      />
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
        Page {{ currentPage }} of {{ totalPages }}
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
</template>

<script>
import EventCard from './EventCard.vue'

export default {
  name: 'EventsGrid',
  components: {
    EventCard
  },
  props: {
    events: {
      type: Array,
      default: () => []
    }
  },
  emits: ['focus-event', 'map-filter-toggle'],
  data() {
    return {
      currentPage: 1,
      eventsPerPage: 3, // Maximum 3 cards to prevent sidebar scrolling
      mapFilterEnabled: false
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
    }
  },
  watch: {
    events() {
      // Reset to first page when events change (e.g., due to filtering)
      this.currentPage = 1
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

.events-count {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
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
  justify-content: center;
  gap: 1rem;
  margin-top: 1rem;
  padding: 1rem 0;
  position: relative;
  z-index: 10;
  pointer-events: auto;
}

.pagination-btn {
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s ease;
  min-width: 2rem;
}

.pagination-btn:hover:not(:disabled) {
  background: #3182ce;
  transform: translateY(-1px);
}

.pagination-btn:disabled {
  background: #a0aec0;
  cursor: not-allowed;
  transform: none;
}

.page-info {
  color: #4a5568;
  font-weight: 500;
  font-size: 0.9rem;
  white-space: nowrap;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .events-section {
    padding: 0.75rem;
  }
  
  .pagination-controls {
    gap: 0.75rem;
  }
  
  .page-info {
    font-size: 0.8rem;
  }
}
</style>