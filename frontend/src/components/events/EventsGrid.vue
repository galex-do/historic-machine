<template>
  <div class="events-section">
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
  emits: ['focus-event'],
  data() {
    return {
      currentPage: 1,
      eventsPerPage: 9 // 3 rows × 3 cards per row
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


.events-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 0.75rem;
  flex: 1;
  max-width: 100%;
}

/* Ensure maximum 3 columns */
@supports (grid-template-columns: subgrid) {
  .events-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    max-grid-columns: 3;
  }
}

/* Fallback for browsers without subgrid support */
@media (min-width: 900px) {
  .events-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 600px) and (max-width: 899px) {
  .events-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 599px) {
  .events-grid {
    grid-template-columns: 1fr;
  }
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