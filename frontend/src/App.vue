<template>
  <div id="app">
    <header>
      <h1>Historical Events Mapping</h1>
      <div class="timeline-controls">
        <label>From: <input type="text" v-model="date_from_display" @blur="update_date_from" placeholder="DD.MM.YYYY" /></label>
        <label>To: <input type="text" v-model="date_to_display" @blur="update_date_to" placeholder="DD.MM.YYYY" /></label>
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
      date_from: '2000-01-01', // Internal ISO format for filtering
      date_to: '', // Will be set in created() hook
      date_from_display: '01.01.2000', // Display format DD.MM.YYYY
      date_to_display: '' // Will be set in created() hook
    }
  },
  created() {
    // Set today's date in DD.MM.YYYY format
    this.date_to_display = this.format_date_to_ddmmyyyy(new Date())
    this.date_to = this.get_today_iso()
  },
  async mounted() {
    await this.fetch_events()
    // Apply initial filter with default date range
    this.filter_events()
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
      return this.format_date_to_ddmmyyyy(new Date(dateString))
    },
    
    format_date_to_ddmmyyyy(date) {
      const day = String(date.getDate()).padStart(2, '0')
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const year = date.getFullYear()
      return `${day}.${month}.${year}`
    },
    
    parse_ddmmyyyy_to_iso(dateStr) {
      if (!dateStr || !dateStr.match(/^\d{1,2}\.\d{1,2}\.\d{4}$/)) {
        return null
      }
      const [day, month, year] = dateStr.split('.')
      const date = new Date(year, month - 1, day)
      return date.toISOString().split('T')[0]
    },
    
    get_today_iso() {
      return new Date().toISOString().split('T')[0]
    },
    
    update_date_from() {
      const iso_date = this.parse_ddmmyyyy_to_iso(this.date_from_display)
      if (iso_date) {
        this.date_from = iso_date
      } else {
        // Reset to default if invalid
        this.date_from_display = '01.01.2000'
        this.date_from = '2000-01-01'
      }
    },
    
    update_date_to() {
      const iso_date = this.parse_ddmmyyyy_to_iso(this.date_to_display)
      if (iso_date) {
        this.date_to = iso_date
      } else {
        // Reset to today if invalid
        this.date_to_display = this.format_date_to_ddmmyyyy(new Date())
        this.date_to = this.get_today_iso()
      }
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