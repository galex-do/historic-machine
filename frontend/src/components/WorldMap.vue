<template>
  <div class="map-container">
    <div ref="map" class="leaflet-map"></div>
    
    <!-- Event Creation Modal -->
    <div v-if="show_event_modal" class="modal-overlay" @click="close_modal">
      <div class="modal-content" @click.stop>
        <h3>Add Historical Event</h3>
        <form @submit.prevent="create_event">
          <div class="form-group">
            <label>Event Name:</label>
            <input type="text" v-model="new_event.name" required />
          </div>
          <div class="form-group">
            <label>Description:</label>
            <textarea v-model="new_event.description" rows="3" required></textarea>
          </div>
          <div class="form-group">
            <label>Date:</label>
            <input type="date" v-model="new_event.date" required />
          </div>
          <div class="form-group">
            <label>Event Type:</label>
            <select v-model="new_event.lens_type">
              <option value="historic">Historic</option>
              <option value="political">Political</option>
              <option value="cultural">Cultural</option>
              <option value="military">Military</option>
            </select>
          </div>
          <div class="coordinates-info">
            <p>Location: {{ new_event.latitude.toFixed(4) }}, {{ new_event.longitude.toFixed(4) }}</p>
          </div>
          <div class="modal-actions">
            <button type="button" @click="close_modal" class="btn-cancel">Cancel</button>
            <button type="submit" class="btn-create">Create Event</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

export default {
  name: 'WorldMap',
  props: {
    events: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      map: null,
      markers: [],
      show_event_modal: false,
      new_event: {
        name: '',
        description: '',
        date: '',
        latitude: 0,
        longitude: 0,
        lens_type: 'historic'
      }
    }
  },
  mounted() {
    this.initialize_map()
    this.add_event_markers()
  },
  watch: {
    events() {
      this.add_event_markers()
    }
  },
  methods: {
    initialize_map() {
      // Create map centered on world view
      this.map = L.map(this.$refs.map).setView([20, 0], 2)
      
      // Add OpenStreetMap tile layer (free, no API key needed)
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 18,
        minZoom: 2
      }).addTo(this.map)
      
      // Add click event for creating new events
      this.map.on('click', this.handle_map_click)
      
      // Fix for default marker icon in Leaflet with bundlers - use local assets
      delete L.Icon.Default.prototype._getIconUrl
      L.Icon.Default.mergeOptions({
        iconRetinaUrl: '/images/marker-icon-2x.png',
        iconUrl: '/images/marker-icon.png',
        shadowUrl: '/images/marker-shadow.png',
      })
    },
    
    handle_map_click(event) {
      const { lat, lng } = event.latlng
      
      // Set coordinates for new event
      this.new_event.latitude = lat
      this.new_event.longitude = lng
      
      // Reset form fields
      this.new_event.name = ''
      this.new_event.description = ''
      this.new_event.date = ''
      this.new_event.lens_type = 'historic'
      
      // Show modal
      this.show_event_modal = true
    },
    
    add_event_markers() {
      // Clear existing markers
      this.markers.forEach(marker => this.map.removeLayer(marker))
      this.markers = []
      
      // Add markers for each event
      this.events.forEach(event => {
        const marker = L.marker([event.latitude, event.longitude]).addTo(this.map)
        
        // Add popup with event information
        marker.bindPopup(`
          <div class="event-popup">
            <h4>${event.name}</h4>
            <p>${event.description}</p>
            <p><strong>Date:</strong> ${this.format_date(event.event_date)}</p>
            <p><strong>Type:</strong> ${event.lens_type}</p>
          </div>
        `)
        
        this.markers.push(marker)
      })
    },
    
    async create_event() {
      try {
        const event_data = {
          name: this.new_event.name,
          description: this.new_event.description,
          latitude: this.new_event.latitude,
          longitude: this.new_event.longitude,
          event_date: new Date(this.new_event.date).toISOString(),
          lens_type: this.new_event.lens_type
        }
        
        const backend_url = this.get_backend_url()
        const response = await fetch(`${backend_url}/events`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(event_data)
        })
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        const created_event = await response.json()
        console.log('Event created successfully:', created_event)
        
        // Emit event to parent component to refresh events list
        this.$emit('event-created', created_event)
        
        // Close modal
        this.close_modal()
        
      } catch (error) {
        console.error('Error creating event:', error)
        alert('Failed to create event. Please try again.')
      }
    },
    
    close_modal() {
      this.show_event_modal = false
    },
    
    get_backend_url() {
      // Check if running in Docker environment (nginx proxy available)
      if (window.location.host.includes('localhost:3000')) {
        return '/api'
      }
      // Default for Replit development environment
      return 'http://localhost:8080/api'
    },
    
    format_date(date_string) {
      return new Date(date_string).toLocaleDateString()
    }
  }
}
</script>

<style scoped>
.map-container {
  position: relative;
  height: 100%;
  width: 100%;
}

.leaflet-map {
  height: 500px;
  width: 100%;
  border-radius: 8px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

.modal-content {
  background: white;
  padding: 30px;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #2c3e50;
  text-align: center;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #34495e;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 2px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #3498db;
}

.coordinates-info {
  background: #f8f9fa;
  padding: 10px;
  border-radius: 6px;
  margin: 15px 0;
}

.coordinates-info p {
  margin: 0;
  font-family: monospace;
  color: #6c757d;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.btn-cancel,
.btn-create {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.3s;
}

.btn-cancel {
  background: #6c757d;
  color: white;
}

.btn-cancel:hover {
  background: #5a6268;
}

.btn-create {
  background: #28a745;
  color: white;
}

.btn-create:hover {
  background: #218838;
}

/* Event popup styling */
:deep(.leaflet-popup-content) {
  margin: 8px 12px;
}

:deep(.event-popup h4) {
  margin: 0 0 8px 0;
  color: #2c3e50;
}

:deep(.event-popup p) {
  margin: 4px 0;
  font-size: 13px;
}
</style>