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
      <div class="map-container">
        <div ref="map" class="world-map">
          <!-- World map will be loaded here -->
          <p>World Map Loading... (Map integration coming next)</p>
        </div>
      </div>
      
      <div class="events-list">
        <h3>Historical Events</h3>
        <div v-for="event in filtered_events" :key="event.id" class="event-item">
          <h4>{{ event.name }}</h4>
          <p>{{ event.description }}</p>
          <p>Date: {{ formatDate(event.event_date) }}</p>
          <p>Location: {{ event.latitude }}, {{ event.longitude }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: 'App',
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
      // Simple filtering implementation
      this.filtered_events = this.events
      console.log('Filtering events from', this.date_from, 'to', this.date_to)
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

.map-container {
  height: 400px;
  background: #f0f0f0;
  border: 1px solid #ddd;
  margin-bottom: 20px;
}

.world-map {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
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
</style>