<template>
  <div class="events-section">
    <h3 class="section-title">Historical Events ({{ events.length }})</h3>
    
    <div class="events-grid">
      <div v-if="events.length === 0" class="no-events">
        <p>No events found for the selected period. Click on the map to add your first historical event!</p>
      </div>
      
      <EventCard
        v-for="event in events"
        :key="event.id"
        :event="event"
        @focus-event="$emit('focus-event', $event)"
      />
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
  emits: ['focus-event']
}
</script>

<style scoped>
.events-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
  max-height: 40vh;
  overflow-y: auto;
}

.section-title {
  color: #2d3748;
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.events-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1rem;
  align-items: start;
}

.no-events {
  grid-column: 1 / -1;
  text-align: center;
  padding: 2rem;
  color: #718096;
  font-style: italic;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .events-grid {
    grid-template-columns: 1fr;
  }
  
  .events-section {
    max-height: 50vh;
  }
}
</style>