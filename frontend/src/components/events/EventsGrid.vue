<template>
  <div class="events-section">
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
  background: transparent;
  padding: 1rem;
  height: 100%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
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

/* Responsive adjustments */
@media (max-width: 768px) {
  .events-section {
    padding: 0.75rem;
  }
}
</style>