<template>
  <div id="app">
    <header>
      <h1>Historical Events Mapping</h1>
      <div class="timeline-controls">
        <label>From: <input type="date" v-model="date_from" /></label>
        <label>To: <input type="date" v-model="date_to" /></label>
        <button @click="filter_events">Filter Events</button>
      </div>
    </header>
    
    <main>
      <div class="map-section">
        <WorldMap 
          :events="filtered_events" 
          @event-created="handle_event_created"
        />
      </div>
      
      <div class="events-list">
        <h3>Historical Events ({{ filtered_events.length }})</h3>
        <div v-if="filtered_events.length === 0" class="no-events">
          <p>No events found. Click on the map to add your first historical event!</p>
        </div>
        <div v-for="event in filtered_events" :key="event.id" class="event-item">
          <h4>{{ event.name }}</h4>
          <p>{{ event.description }}</p>
          <p><strong>Date:</strong> {{ formatDate(event.event_date) }}</p>
          <p><strong>Location:</strong> {{ event.latitude.toFixed(4) }}, {{ event.longitude.toFixed(4) }}</p>
          <p><strong>Type:</strong> {{ event.lens_type }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import WorldMap from './components/WorldMap.vue'

export default {
  name: 'App',
  components: {
    WorldMap
  },
  data() {
    return {
      events: [],
      filtered_events: [],
      date_from: '',
      date_to: ''
    }
  },
  async mounted() {
    await this.fetch_events()
  },
  methods: {
    get_backend_url() {
      // Check if running in Docker environment (nginx proxy available)
      if (window.location.host.includes('localhost:3000')) {
        return '/api'
      }
      // Default for Replit development environment
      return 'http://localhost:8080/api'
    },
    async fetch_events() {
      try {
        const backend_url = this.get_backend_url()
        const response = await fetch(`${backend_url}/events`)
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        this.events = await response.json()
        this.filtered_events = this.events
        console.log('Successfully loaded events:', this.events.length)
      } catch (error) {
        console.error('Error fetching events:', error)
        // Fallback sample data for development
        this.events = [
          {
            id: 1,
            name: 'Fall of Constantinople',
            description: 'Ottoman Empire conquered Byzantine Empire',
            latitude: 41.0082,
            longitude: 28.9784,
            event_date: '1453-05-29T00:00:00Z',
            lens_type: 'historic'
          }
        ]
        this.filtered_events = this.events
      }
    },
    filter_events() {
      if (!this.date_from && !this.date_to) {
        this.filtered_events = this.events
        return
      }
      
      this.filtered_events = this.events.filter(event => {
        const event_date = new Date(event.event_date)
        const from_date = this.date_from ? new Date(this.date_from) : null
        const to_date = this.date_to ? new Date(this.date_to) : null
        
        if (from_date && event_date < from_date) return false
        if (to_date && event_date > to_date) return false
        
        return true
      })
      
      console.log(`Filtering events from ${this.date_from} to ${this.date_to}. Found ${this.filtered_events.length} events.`)
    },
    
    async handle_event_created(new_event) {
      // Refresh events list after new event is created
      await this.fetch_events()
      console.log('Events list refreshed after new event creation')
    },
    
    formatDate(dateString) {
      return new Date(dateString).toLocaleDateString()
    }
  }
}
</script>

<style>
#app {
  font-family: Arial, sans-serif;
  padding: 20px;
}

header {
  background: #2c3e50;
  color: white;
  padding: 20px;
  margin-bottom: 20px;
}

.timeline-controls {
  margin-top: 10px;
}

.timeline-controls label {
  margin-right: 15px;
}

.timeline-controls input {
  margin: 0 5px;
  padding: 5px;
}

.timeline-controls button {
  padding: 5px 15px;
  background: #3498db;
  color: white;
  border: none;
  cursor: pointer;
}

.map-section {
  margin-bottom: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.events-list {
  background: #f9f9f9;
  padding: 20px;
}

.event-item {
  background: white;
  margin-bottom: 10px;
  padding: 15px;
  border-left: 4px solid #3498db;
}

.event-item h4 {
  margin-top: 0;
  color: #2c3e50;
}

.no-events {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}
</style>