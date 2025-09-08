<template>
  <div class="map-container">
    <div ref="map" class="leaflet-map"></div>
    
    <!-- Event Creation Modal (only show if user can create events) -->
    <div v-if="show_event_modal && canCreateEvents" class="modal-overlay" @click="close_modal">
      <div class="modal-content" @click.stop>
        <h3>{{ editing_event ? 'Edit Historical Event' : 'Add Historical Event' }}</h3>
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
            <label>Era:</label>
            <select v-model="new_event.era">
              <option value="BC">BC (Before Christ)</option>
              <option value="AD">AD (Anno Domini)</option>
            </select>
          </div>
          <div class="form-group">
            <label>Event Type:</label>
            <select v-model="new_event.lens_type">
              <option value="historic">📜 Historic</option>
              <option value="political">🏛️ Political</option>
              <option value="cultural">🎭 Cultural</option>
              <option value="military">⚔️ Military</option>
              <option value="scientific">🔬 Scientific</option>
            </select>
          </div>
          <div class="coordinates-info">
            <p>Location: {{ new_event.latitude.toFixed(4) }}, {{ new_event.longitude.toFixed(4) }}</p>
          </div>
          <div class="modal-actions">
            <div class="left-actions">
              <button v-if="editing_event" type="button" @click="delete_event" class="btn-delete">Delete Event</button>
            </div>
            <div class="right-actions">
              <button type="button" @click="close_modal" class="btn-cancel">Cancel</button>
              <button type="submit" class="btn-create">{{ editing_event ? 'Update Event' : 'Create Event' }}</button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { useAuth } from '@/composables/useAuth.js'
import apiService from '@/services/api.js'

export default {
  name: 'WorldMap',
  setup() {
    const { canCreateEvents, canEditEvents, isGuest } = useAuth()
    return {
      canCreateEvents,
      canEditEvents,
      isGuest
    }
  },
  props: {
    events: {
      type: Array,
      default: () => []
    },
    focusEvent: {
      type: Object,
      default: null
    }
  },
  emits: ['event-created', 'event-updated', 'event-deleted', 'map-bounds-changed'],
  data() {
    return {
      map: null,
      markers: [],
      resize_observer: null,
      show_event_modal: false,
      editing_event: null, // Store the event being edited
      is_stepping: false, // Track if current update is from date stepping
      new_event: {
        name: '',
        description: '',
        date: '', // Internal ISO format
        date_display: '', // Display format DD.MM.YYYY
        era: 'AD',
        latitude: 0,
        longitude: 0,
        lens_type: 'historic'
      }
    }
  },
  watch: {
    events: {
      handler() {
        if (this.map) {
          this.add_event_markers()
          // Only recenter map if not stepping through dates
          if (!this.is_stepping) {
            this.fit_map_to_events()
          }
          // Reset stepping flag after processing
          this.is_stepping = false
        }
      },
      deep: true,
      immediate: false
    },
    focusEvent: {
      handler(new_focus_event) {
        if (new_focus_event && this.map) {
          this.center_on_event(new_focus_event)
        }
      },
      deep: true,
      immediate: false
    },
    // Update marker popups when authentication state changes
    canEditEvents: {
      handler() {
        if (this.map && this.markers.length > 0) {
          console.log('Auth state changed, updating popup content')
          this.update_marker_popups()
        }
      },
      immediate: false
    }
  },
  
  mounted() {
    this.initialize_map()
    this.add_event_markers()
    this.setup_coordinate_system_protection()
    
    // Add resize observer to handle layout changes
    this.resize_observer = new ResizeObserver(() => {
      this.handle_resize()
    })
    this.resize_observer.observe(this.$refs.map)
    
    // Listen for window resize events (triggered by sidebar toggle)
    window.addEventListener('resize', this.handle_resize)
    
    // Set up global function for edit buttons in popups
    window.editEvent = (eventId) => {
      this.edit_event(eventId)
    }
  },
  
  beforeUnmount() {
    if (this.resize_observer) {
      this.resize_observer.disconnect()
    }
    window.removeEventListener('resize', this.handle_resize)
    
    // Clean up coordinate system health check
    if (this.coordinate_health_check) {
      clearInterval(this.coordinate_health_check)
    }
    
    // Clean up global function
    if (window.editEvent) {
      delete window.editEvent
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

      // Add map bounds change listeners for the filter
      this.map.on('moveend', this.handle_bounds_change)
      this.map.on('zoomend', this.handle_bounds_change)
      
      // Fix for default marker icon in Leaflet with bundlers - use local assets
      delete L.Icon.Default.prototype._getIconUrl
      L.Icon.Default.mergeOptions({
        iconRetinaUrl: '/images/marker-icon-2x.png',
        iconUrl: '/images/marker-icon.png',
        shadowUrl: '/images/marker-shadow.png',
      })
    },
    
    handle_map_click(event) {
      // Check if user can create events
      if (!this.canCreateEvents) {
        return
      }

      // Check if any popup is currently open
      if (this.map._popup && this.map._popup._isOpen) {
        // If popup is open, close it instead of opening event creation
        this.map.closePopup()
        return
      }

      const { lat, lng } = event.latlng
      
      // Set coordinates for new event
      this.new_event.latitude = lat
      this.new_event.longitude = lng
      
      // Reset form fields
      this.new_event.name = ''
      this.new_event.description = ''
      this.new_event.date = ''
      this.new_event.date_display = ''
      this.new_event.era = 'AD'
      this.new_event.lens_type = 'historic'
      
      // Show modal
      this.show_event_modal = true
    },

    // Get current map bounds
    getCurrentBounds() {
      if (!this.map) {
        return Promise.resolve(null)
      }
      
      const bounds = this.map.getBounds()
      return Promise.resolve({
        north: bounds.getNorth(),
        south: bounds.getSouth(),
        east: bounds.getEast(),
        west: bounds.getWest()
      })
    },

    // Handle map bounds changes
    handle_bounds_change() {
      if (!this.map) return
      
      const bounds = this.map.getBounds()
      const boundsData = {
        north: bounds.getNorth(),
        south: bounds.getSouth(),
        east: bounds.getEast(),
        west: bounds.getWest()
      }
      
      this.$emit('map-bounds-changed', boundsData)
    },

    edit_event(eventId) {
      // Find the event to edit
      const event = this.events.find(e => e.id === eventId)
      if (!event) {
        console.error('Event not found:', eventId)
        return
      }

      // Set editing mode
      this.editing_event = event
      
      // Populate form with event data
      this.new_event.name = event.name
      this.new_event.description = event.description
      this.new_event.latitude = event.latitude
      this.new_event.longitude = event.longitude
      this.new_event.era = event.era || 'AD'
      this.new_event.lens_type = event.lens_type
      
      // Handle date formatting
      this.new_event.date = event.event_date.split('T')[0]
      this.new_event.date_display = this.format_date(event.event_date)
      
      // Show modal
      this.show_event_modal = true
    },
    
    add_event_markers() {
      this.comprehensive_popup_cleanup()
      
      // Clear existing markers first
      this.markers.forEach(marker => {
        if (this.map && this.map.hasLayer(marker)) {
          // Properly unbind popup before removing marker
          if (marker.getPopup()) {
            marker.unbindPopup()
          }
          this.map.removeLayer(marker)
        }
      })
      this.markers = []
      
      // Wait for next tick to ensure map is ready
      this.$nextTick(() => {
        // Double-check map state before proceeding
        if (!this.map || !this.map._loaded) {
          console.warn('Map not ready for marker addition')
          return
        }
        
        try {
          // Group events by location (same coordinates)
          const locationGroups = this.group_events_by_location(this.events)
          
          // Add markers for each location group
          Object.values(locationGroups).forEach(eventGroup => {
            if (!eventGroup.latitude || !eventGroup.longitude || !this.map) return
            
            // Validate coordinates before creating marker
            const lat = parseFloat(eventGroup.latitude)
            const lng = parseFloat(eventGroup.longitude)
            
            if (isNaN(lat) || isNaN(lng) || lat < -90 || lat > 90 || lng < -180 || lng > 180) {
              console.warn('Invalid coordinates:', lat, lng)
              return
            }
            
            const emoji_icon = this.create_emoji_marker_icon(
              eventGroup.events.length > 1 ? 'multiple' : eventGroup.events[0].lens_type
            )
            
            const marker = L.marker([lat, lng], { 
              icon: emoji_icon,
              riseOnHover: true
            }).addTo(this.map)
            
            // Create popup content for single or multiple events
            const popupContent = this.create_popup_content(eventGroup.events)
            marker.bindPopup(popupContent)
            
            this.markers.push(marker)
          })
        } catch (error) {
          console.error('Error adding markers:', error)
          // Force complete map refresh on error
          this.force_map_refresh()
        }
        
        // Invalidate map size to ensure proper rendering
        if (this.map) {
          setTimeout(() => {
            this.map.invalidateSize(true)
          }, 100)
        }
      })
    },

    // Comprehensive popup cleanup to prevent coordinate system corruption
    comprehensive_popup_cleanup() {
      if (!this.map) return
      
      try {
        // Close any open popups
        this.map.closePopup()
        
        // Clear popup reference from map
        if (this.map._popup) {
          this.map._popup = null
        }
        
        // Remove all popup-related events
        this.map.off('popupopen')
        this.map.off('popupclose')
        
        // Force coordinate system recalculation
        this.map.invalidateSize(false)
        
        // Additional cleanup for any lingering popup DOM elements
        const popupElements = document.querySelectorAll('.leaflet-popup')
        popupElements.forEach(popup => {
          if (popup.parentNode) {
            popup.parentNode.removeChild(popup)
          }
        })
        
      } catch (error) {
        console.warn('Error during popup cleanup:', error)
        // If cleanup fails, force complete map refresh
        this.force_map_refresh()
      }
    },

    // Force complete map refresh when coordinate system gets corrupted
    force_map_refresh() {
      if (!this.map) return
      
      try {
        // Get current map state
        const center = this.map.getCenter()
        const zoom = this.map.getZoom()
        
        // Force complete map invalidation
        this.map.invalidateSize(true)
        
        // Reset view to ensure coordinate system is recalculated
        setTimeout(() => {
          if (this.map) {
            this.map.setView(center, zoom, { reset: true })
          }
        }, 50)
      } catch (error) {
        console.error('Error during force map refresh:', error)
      }
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
    
    handle_resize() {
      if (this.map) {
        // Small delay to ensure DOM has updated
        this.$nextTick(() => {
          try {
            this.map.invalidateSize()
            // Force a redraw to prevent rendering issues
            setTimeout(() => {
              if (this.map) {
                this.map.invalidateSize()
              }
            }, 100)
          } catch (error) {
            console.warn('Error resizing map:', error)
          }
        })
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
          era: this.new_event.era,
          lens_type: this.new_event.lens_type
        }
        
        let response
        if (this.editing_event) {
          // Update existing event
          response = await apiService.updateEvent(this.editing_event.id, event_data)
          console.log('Event updated successfully:', response)
          this.$emit('event-updated', response)
        } else {
          // Create new event
          response = await apiService.createEvent(event_data)
          console.log('Event created successfully:', response)
          this.$emit('event-created', response)
        }
        
        // Close modal
        this.close_modal()
        
      } catch (error) {
        console.error('Error saving event:', error)
        if (error.message.includes('401')) {
          alert('Authentication failed. Please log in again.')
        } else {
          const action = this.editing_event ? 'update' : 'create'
          alert(`Failed to ${action} event. Please try again.`)
        }
      }
    },

    async delete_event() {
      if (!this.editing_event) return
      
      const confirmed = confirm(`Are you sure you want to delete "${this.editing_event.name}"? This action cannot be undone.`)
      if (!confirmed) return
      
      try {
        await apiService.deleteEvent(this.editing_event.id)
        console.log('Event deleted successfully')
        
        // Emit event to parent component to refresh events list
        this.$emit('event-deleted', this.editing_event.id)
        
        // Close modal
        this.close_modal()
        
      } catch (error) {
        console.error('Error deleting event:', error)
        if (error.message.includes('401')) {
          alert('Authentication failed. Please log in again.')
        } else {
          alert('Failed to delete event. Please try again.')
        }
      }
    },
    
    close_modal() {
      this.show_event_modal = false
      this.editing_event = null
      
      // Reset form
      this.new_event = {
        name: '',
        description: '',
        date: '',
        date_display: '',
        era: 'AD',
        latitude: 0,
        longitude: 0,
        lens_type: 'historic'
      }
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
        'scientific': '🔬',   // Microscope
        'cultural': '🎭'      // Theater masks
      }
      
      const emoji = emoji_map[lens_type] || '📍' // Default location pin (also for multiple events)
      
      return L.divIcon({
        html: `<div class="emoji-marker" data-lens="${lens_type}">${emoji}</div>`,
        className: 'emoji-marker-container',
        iconSize: [30, 30],
        iconAnchor: [15, 30],
        popupAnchor: [0, -30],
        // Ensure proper positioning
        bgPos: [15, 30]
      })
    },

    group_events_by_location(events) {
      const groups = {}
      const tolerance = 0.0001 // Small tolerance for coordinate matching
      
      events.forEach(event => {
        if (!event.latitude || !event.longitude) return
        
        // Create a location key with rounded coordinates
        const lat = Math.round(event.latitude / tolerance) * tolerance
        const lng = Math.round(event.longitude / tolerance) * tolerance
        const locationKey = `${lat},${lng}`
        
        if (!groups[locationKey]) {
          groups[locationKey] = {
            latitude: event.latitude,
            longitude: event.longitude,
            events: []
          }
        }
        
        groups[locationKey].events.push(event)
      })
      
      return groups
    },

    update_marker_popups() {
      // Update popup content for existing markers without recreating them
      const locationGroups = this.group_events_by_location(this.events)
      const groupsArray = Object.values(locationGroups)
      
      this.markers.forEach((marker, index) => {
        if (index < groupsArray.length) {
          const eventGroup = groupsArray[index]
          const popupContent = this.create_popup_content(eventGroup.events)
          marker.setPopupContent(popupContent)
        }
      })
    },

    create_popup_content(events) {
      if (events.length === 1) {
        const event = events[0]
        const editIcon = this.canEditEvents ? 
          `<span onclick="window.editEvent(${event.id})" class="edit-icon" title="Edit event">✏️</span>` : ''
        
        return `
          <div class="event-popup">
            <div class="event-header">
              <h4>${event.name} ${editIcon}</h4>
            </div>
            <p>${event.description}</p>
            <p><strong>Date:</strong> ${event.display_date || this.format_date(event.event_date)}</p>
            <p><strong>Type:</strong> ${event.lens_type}</p>
          </div>
        `
      } else {
        // Multiple events at same location
        const eventsHtml = events.map((event, index) => {
          const editIcon = this.canEditEvents ? 
            `<span onclick="window.editEvent(${event.id})" class="edit-icon" title="Edit event">✏️</span>` : ''
          
          return `
            <div class="event-item">
              <div class="event-header">
                <h4>${event.name} ${editIcon}</h4>
              </div>
              <p>${event.description}</p>
              <p><strong>Date:</strong> ${event.display_date || this.format_date(event.event_date)}</p>
              <p><strong>Type:</strong> ${event.lens_type}</p>
              ${index < events.length - 1 ? '<div class="event-divider"></div>' : ''}
            </div>
          `
        }).join('')
        
        return `
          <div class="event-popup multi-event">
            <div class="popup-header">
              <h3>${events.length} Events at this location</h3>
            </div>
            ${eventsHtml}
          </div>
        `
      }
    },

    // Method to set stepping mode from parent component
    setSteppingMode(isStepping) {
      this.is_stepping = isStepping
    },

    // Setup coordinate system protection to prevent _latLngToNewLayerPoint errors
    setup_coordinate_system_protection() {
      if (!this.map) return

      // Store original map state for recovery
      this.original_map_state = {
        center: this.map.getCenter(),
        zoom: this.map.getZoom()
      }

      // Override popup handling to prevent coordinate corruption
      this.map.on('popupopen', (e) => {
        try {
          // Force immediate coordinate system refresh when popup opens
          setTimeout(() => {
            if (this.map && this.map._loaded) {
              this.map.invalidateSize(false)
            }
          }, 0)
        } catch (error) {
          console.warn('Error during popup open handling:', error)
        }
      })

      this.map.on('popupclose', (e) => {
        try {
          // Clean up coordinate system after popup closes
          setTimeout(() => {
            if (this.map && this.map._loaded) {
              this.repair_coordinate_system()
            }
          }, 0)
        } catch (error) {
          console.warn('Error during popup close handling:', error)
        }
      })

      // Monitor zoom events for coordinate system corruption
      this.map.on('zoom', (e) => {
        this.detect_and_fix_coordinate_corruption()
      })

      this.map.on('zoomstart', (e) => {
        this.detect_and_fix_coordinate_corruption()
      })

      this.map.on('zoomend', (e) => {
        this.detect_and_fix_coordinate_corruption()
      })

      // Periodic coordinate system health check
      this.coordinate_health_check = setInterval(() => {
        this.detect_and_fix_coordinate_corruption()
      }, 5000) // Check every 5 seconds
    },

    // Detect and fix coordinate system corruption
    detect_and_fix_coordinate_corruption() {
      if (!this.map || !this.map._loaded) return

      try {
        // Test coordinate system health by trying a simple transformation
        const testLatLng = L.latLng(0, 0)
        this.map.latLngToContainerPoint(testLatLng)
      } catch (error) {
        if (error.message && error.message.includes('_latLngToNewLayerPoint')) {
          console.log('Coordinate system corruption detected, repairing...')
          this.repair_coordinate_system()
        }
      }
    },

    // Repair corrupted coordinate system
    repair_coordinate_system() {
      if (!this.map || !this.map._loaded) return

      try {
        // Get current state
        const currentCenter = this.map.getCenter()
        const currentZoom = this.map.getZoom()

        // Force complete coordinate system reset
        this.map.invalidateSize(true)
        
        // Reset coordinate transformation matrices
        if (this.map._renderer && this.map._renderer._transformMatrix) {
          delete this.map._renderer._transformMatrix
        }

        // Clear any cached coordinate transformations
        if (this.map._pixelTransform) {
          this.map._pixelTransform = null
        }

        // Force view reset to recalculate coordinate system
        setTimeout(() => {
          if (this.map && this.map._loaded) {
            this.map.setView(currentCenter, currentZoom, {
              reset: true,
              animate: false
            })
          }
        }, 10)

        console.log('Coordinate system repaired')
      } catch (error) {
        console.error('Failed to repair coordinate system:', error)
        // Last resort: force complete map refresh
        this.emergency_map_reset()
      }
    },

    // Emergency map reset as last resort
    emergency_map_reset() {
      if (!this.map) return

      try {
        console.log('Performing emergency map reset...')
        
        // Save current state
        const currentCenter = this.map.getCenter()
        const currentZoom = this.map.getZoom()
        const currentEvents = [...this.events]

        // Remove all layers
        this.map.eachLayer((layer) => {
          if (layer !== this.map._layers[this.map._leaflet_id]) {
            this.map.removeLayer(layer)
          }
        })

        // Clear all markers
        this.markers = []

        // Force complete map reset
        this.map.remove()
        
        // Reinitialize map
        setTimeout(() => {
          this.initialize_map()
          // Restore view
          if (this.map) {
            this.map.setView(currentCenter, currentZoom)
            // Restore markers
            this.add_event_markers()
          }
        }, 100)

        console.log('Emergency map reset completed')
      } catch (error) {
        console.error('Emergency map reset failed:', error)
      }
    }
  }
}
</script>

<style scoped>
.map-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
}

.leaflet-map {
  height: 100%;
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
  z-index: 9999;
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
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
}

.left-actions {
  display: flex;
  gap: 10px;
}

.right-actions {
  display: flex;
  gap: 10px;
}

.btn-cancel,
.btn-create,
.btn-delete {
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

.btn-delete {
  background: #dc3545;
  color: white;
}

.btn-delete:hover {
  background: #c82333;
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

/* Event popup header */
:deep(.event-header) {
  margin-bottom: 8px;
}

:deep(.event-header h4) {
  margin: 0;
  display: inline;
}

/* Edit icon styling in popup */
:deep(.edit-icon) {
  font-size: 14px;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.2s;
  margin-left: 6px;
  display: inline;
}

:deep(.edit-icon:hover) {
  opacity: 1;
}

/* Multi-event popup styling */
:deep(.multi-event) {
  max-width: 300px;
}

:deep(.popup-header h3) {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: #2c3e50;
  text-align: center;
  border-bottom: 1px solid #eee;
  padding-bottom: 8px;
}

:deep(.event-item) {
  margin-bottom: 12px;
}

:deep(.event-item:last-child) {
  margin-bottom: 0;
}

:deep(.event-divider) {
  height: 1px;
  background: #eee;
  margin: 12px 0;
}
</style>