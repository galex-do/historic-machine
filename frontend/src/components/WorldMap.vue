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
            <input type="text" v-model="new_event.date_display" @blur="update_event_date" placeholder="DD.MM.YYYY" required />
          </div>
          <div class="form-group">
            <label>Event Type:</label>
            <select v-model="new_event.lens_type">
              <option value="historic">📜 Historic</option>
              <option value="political">🏛️ Political</option>
              <option value="cultural">🎭 Cultural</option>
              <option value="military">⚔️ Military</option>
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
    },
    focus_event: {
      type: Object,
      default: null
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
        date: '', // Internal ISO format
        date_display: '', // Display format DD.MM.YYYY
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
    events: {
      handler() {
        if (this.map) {
          this.add_event_markers()
          this.fit_map_to_events()
        }
      },
      deep: true,
      immediate: false
    },
    focus_event: {
      handler(new_focus_event) {
        if (new_focus_event && this.map) {
          this.center_on_event(new_focus_event)
        }
      },
      deep: true,
      immediate: false
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
      this.new_event.date_display = ''
      this.new_event.lens_type = 'historic'
      
      // Show modal
      this.show_event_modal = true
    },
    
    add_event_markers() {
      // Clear existing markers
      this.markers.forEach(marker => {
        if (this.map && this.map.hasLayer(marker)) {
          this.map.removeLayer(marker)
        }
      })
      this.markers = []
      
      // Wait for next tick to ensure map is ready
      this.$nextTick(() => {
        // Add markers for each event with emoji icons
        this.events.forEach(event => {
          if (!event.latitude || !event.longitude || !this.map) return
          
          const emoji_icon = this.create_emoji_marker_icon(event.lens_type)
          const marker = L.marker([event.latitude, event.longitude], { 
            icon: emoji_icon,
            riseOnHover: true
          }).addTo(this.map)
          
          // Add popup with event information
          marker.bindPopup(`
            <div class="event-popup">
              <h4>${event.name}</h4>
              <p>${event.description}</p>
              <p><strong>Date:</strong> ${event.display_date || this.format_date(event.event_date)}</p>
              <p><strong>Type:</strong> ${event.lens_type}</p>
            </div>
          `)
          
          this.markers.push(marker)
        })
        
        // Invalidate map size to ensure proper rendering
        if (this.map) {
          setTimeout(() => {
            this.map.invalidateSize(true)
          }, 100)
        }
      })
    },
    
    fit_map_to_events() {
      // Add comprehensive checks to prevent Leaflet errors
      if (!this.map || !this.map._loaded || !this.events || this.events.length === 0) {
        return
      }
      
      // Wait for next tick to ensure map is fully rendered
      this.$nextTick(() => {
        try {
          // If only one event, center on it with reasonable zoom
          if (this.events.length === 1) {
            const event = this.events[0]
            if (event.latitude != null && event.longitude != null) {
              this.map.setView([event.latitude, event.longitude], 6)
            }
            return
          }
          
          // For multiple events, calculate bounds
          const lats = this.events.map(event => event.latitude).filter(lat => lat != null && !isNaN(lat))
          const lngs = this.events.map(event => event.longitude).filter(lng => lng != null && !isNaN(lng))
          
          if (lats.length === 0 || lngs.length === 0) {
            return
          }
          
          const minLat = Math.min(...lats)
          const maxLat = Math.max(...lats)
          const minLng = Math.min(...lngs)
          const maxLng = Math.max(...lngs)
          
          // Ensure bounds are valid
          if (minLat === maxLat && minLng === maxLng) {
            // Single point, just center on it
            this.map.setView([minLat, minLng], 6)
            return
          }
          
          // Create bounds with some padding
          const bounds = [
            [minLat, minLng],
            [maxLat, maxLng]
          ]
          
          // Fit the map to show all events with padding
          this.map.fitBounds(bounds, {
            padding: [20, 20],
            maxZoom: 8 // Don't zoom in too much for close events
          })
        } catch (error) {
          console.warn('Error fitting map to events:', error)
          // Fallback to world view if there's an error
          this.map.setView([20, 0], 2)
        }
      })
    },
    
    center_on_event(event) {
      if (!this.map || !event.latitude || !event.longitude) {
        return
      }
      
      try {
        // Center the map on the specific event with a good zoom level
        this.map.setView([event.latitude, event.longitude], 8, {
          animate: true,
          duration: 1.0
        })
      } catch (error) {
        console.warn('Error centering on event:', error)
      }
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
      // For very old dates, parse manually to avoid Date object issues
      if (date_string.startsWith('00') || date_string.startsWith('01')) {
        return this.format_date_to_ddmmyyyy(date_string)
      }
      return this.format_date_to_ddmmyyyy(new Date(date_string))
    },
    
    format_date_to_ddmmyyyy(date) {
      // Handle both Date objects and ISO strings
      let day, month, year
      
      if (typeof date === 'string') {
        // Parse ISO string manually for very old dates
        const [year_str, month_str, day_str] = date.split('-')
        day = parseInt(day_str, 10)
        month = parseInt(month_str, 10) 
        year = parseInt(year_str, 10)
      } else {
        // Regular Date object
        day = date.getDate()
        month = date.getMonth() + 1
        year = date.getFullYear()
      }
      
      const padded_day = String(day).padStart(2, '0')
      const padded_month = String(month).padStart(2, '0')
      
      return `${padded_day}.${padded_month}.${year}`
    },
    
    parse_ddmmyyyy_to_iso(dateStr) {
      // Allow years from 1-4 digits (year 1 to 9999)
      if (!dateStr || !dateStr.match(/^\d{1,2}\.\d{1,2}\.\d{1,4}$/)) {
        return null
      }
      const [day, month, year] = dateStr.split('.')
      const year_num = parseInt(year, 10)
      const month_num = parseInt(month, 10)
      const day_num = parseInt(day, 10)
      
      // Basic validation
      if (year_num < 1 || year_num > 9999) return null
      if (month_num < 1 || month_num > 12) return null
      if (day_num < 1 || day_num > 31) return null
      
      // For very old years, construct ISO string manually to avoid Date object issues
      const padded_year = String(year_num).padStart(4, '0')
      const padded_month = String(month_num).padStart(2, '0')
      const padded_day = String(day_num).padStart(2, '0')
      
      return `${padded_year}-${padded_month}-${padded_day}`
    },
    
    update_event_date() {
      const iso_date = this.parse_ddmmyyyy_to_iso(this.new_event.date_display)
      if (iso_date) {
        this.new_event.date = iso_date
      } else {
        // Clear if invalid
        this.new_event.date_display = ''
        this.new_event.date = ''
      }
    },
    
    create_emoji_marker_icon(lens_type) {
      // Define emojis for different event types
      const emoji_map = {
        'military': '⚔️',     // Crossed swords
        'political': '🏛️',   // Classical building/government
        'historic': '📜',     // Ancient scroll/manuscript
        'cultural': '🎭'      // Theater masks
      }
      
      const emoji = emoji_map[lens_type] || '📍' // Default location pin
      
      return L.divIcon({
        html: `<div class="emoji-marker" data-lens="${lens_type}">${emoji}</div>`,
        className: 'emoji-marker-container',
        iconSize: [30, 30],
        iconAnchor: [15, 30],
        popupAnchor: [0, -30],
        // Ensure proper positioning
        bgPos: [15, 30]
      })
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

/* Emoji marker styling */
:deep(.emoji-marker-container) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
}

:deep(.emoji-marker) {
  font-size: 24px;
  text-align: center;
  line-height: 30px;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.5);
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
}
</style>