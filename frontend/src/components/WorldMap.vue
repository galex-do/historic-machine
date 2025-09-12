<template>
  <div class="map-container">
    <div ref="map" class="leaflet-map"></div>
    
    <!-- Event Creation Modal (only show if user can create events) -->
    <div v-if="show_event_modal && canCreateEvents" class="modal-overlay" @click="close_modal">
      <div class="modal-content" @click.stop>
        <h3>{{ editing_event ? 'Edit Historical Event' : 'Add Historical Event' }}</h3>
        <form @submit.prevent="create_event">
          <!-- Locale Selection Tabs -->
          <div class="locale-tabs">
            <button type="button" 
                    :class="['locale-tab', { active: activeLocaleTab === 'en' }]" 
                    @click="activeLocaleTab = 'en'">
              English
            </button>
            <button type="button" 
                    :class="['locale-tab', { active: activeLocaleTab === 'ru' }]" 
                    @click="activeLocaleTab = 'ru'">
              Русский
            </button>
          </div>

          <!-- English Fields -->
          <div v-show="activeLocaleTab === 'en'" class="locale-content">
            <div class="form-group">
              <label>Event Name (English):</label>
              <input type="text" v-model="new_event.name_en" required />
            </div>
            <div class="form-group">
              <label>Description (English):</label>
              <textarea v-model="new_event.description_en" rows="3" required></textarea>
            </div>
          </div>

          <!-- Russian Fields -->
          <div v-show="activeLocaleTab === 'ru'" class="locale-content">
            <div class="form-group">
              <label>Event Name (Russian):</label>
              <input type="text" v-model="new_event.name_ru" />
            </div>
            <div class="form-group">
              <label>Description (Russian):</label>
              <textarea v-model="new_event.description_ru" rows="3"></textarea>
            </div>
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
              <option 
                v-for="lensType in getAvailableLensTypes()" 
                :key="lensType.value" 
                :value="lensType.value"
              >
                {{ lensType.label }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Source (optional):</label>
            <input type="url" v-model="new_event.source" placeholder="https://example.com/source" />
            <small class="field-hint">HTTP/HTTPS link to the source where information about this event was found</small>
          </div>
          <div class="form-group">
            <label>Dataset (optional):</label>
            <select v-model="new_event.dataset_id">
              <option :value="null">No dataset assigned</option>
              <option 
                v-for="dataset in datasets" 
                :key="dataset.id" 
                :value="dataset.id"
              >
                {{ dataset.filename }} ({{ dataset.event_count }} events)
              </option>
            </select>
            <small class="field-hint">Optionally assign this event to an existing dataset for organization</small>
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

    <!-- Event Info Modal (shows event details without coordinate corruption) -->
    <div v-if="show_event_info_modal" class="modal-overlay" @click="close_event_info_modal">
      <div class="modal-content event-info-modal" @click.stop>
        <div class="modal-header">
          <h3 v-if="selected_events.length === 1">
            {{ get_event_emoji(selected_events[0].lens_type) }}
            <a 
              v-if="selected_events[0].source"
              :href="selected_events[0].source" 
              target="_blank" 
              rel="noopener noreferrer"
              class="modal-event-title-link" 
              :title="`${selected_events[0].name} - View source`"
            >
              {{ selected_events[0].name }}
            </a>
            <span v-else class="modal-event-title-text">{{ selected_events[0].name }}</span>
          </h3>
          <h3 v-else>{{ selected_events.length }} Events at this Location</h3>
          <div class="header-actions">
            <span v-if="selected_events.length === 1 && canEditEvents" @click="edit_event_from_info(selected_events[0].id)" class="edit-icon header-edit" title="Edit event">✏️</span>
            <button class="close-button" @click="close_event_info_modal">×</button>
          </div>
        </div>
        
        <div class="events-list">
          <div v-for="(event, index) in selected_events" :key="event.id" class="event-info-item">
            <!-- Show event header only for multiple events -->
            <div v-if="selected_events.length > 1" class="event-header">
              <h4>
                {{ get_event_emoji(event.lens_type) }}
                <a 
                  v-if="event.source"
                  :href="event.source" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  class="modal-event-header-link" 
                  :title="`${event.name} - View source`"
                >
                  {{ event.name }}
                </a>
                <span v-else>{{ event.name }}</span>
              </h4>
              <span v-if="canEditEvents" @click="edit_event_from_info(event.id)" class="edit-icon" title="Edit event">✏️</span>
            </div>
            <div class="event-date">{{ event.display_date || format_date(event.event_date) }}</div>
            <p class="event-description">{{ event.description }}</p>
            <!-- Event Tags (max 3) -->
            <div v-if="event.tags && event.tags.length > 0" class="event-tags">
              <span 
                v-for="tag in event.tags.slice(0, 3)" 
                :key="tag.id"
                class="event-tag"
                :style="{ backgroundColor: tag.color, color: getContrastColor(tag.color) }"
              >
                {{ tag.name }}
              </span>
              <span v-if="event.tags.length > 3" class="more-tags">+{{ event.tags.length - 3 }} more</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { useAuth } from '@/composables/useAuth.js'
import { useTags } from '@/composables/useTags.js'
import apiService from '@/services/api.js'
import { getEventEmoji, getAvailableLensTypes } from '@/utils/event-utils.js'

export default {
  name: 'WorldMap',
  setup() {
    const { canCreateEvents, canEditEvents, isGuest } = useAuth()
    const { allTags, loadTags } = useTags()
    return {
      canCreateEvents,
      canEditEvents,
      isGuest,
      allTags,
      loadTags
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
      show_event_info_modal: false, // New modal for event info
      selected_events: [], // Events to show in info modal
      editing_event: null, // Store the event being edited
      is_stepping: false, // Track if current update is from date stepping
      activeLocaleTab: 'en', // Track active locale tab in form
      new_event: {
        name: '',
        description: '',
        date: '', // Internal ISO format
        date_display: '', // Display format DD.MM.YYYY
        era: 'AD',
        latitude: 0,
        longitude: 0,
        lens_type: 'historic',
        source: '',
        dataset_id: null
      },
      datasets: [] // Available datasets for selection
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
  
  async mounted() {
    this.initialize_map()
    this.add_event_markers()
    await this.fetchDatasets()
    
    // Add resize observer to handle layout changes
    this.resize_observer = new ResizeObserver(() => {
      this.handle_resize()
    })
    this.resize_observer.observe(this.$refs.map)
    
    // Listen for window resize events (triggered by sidebar toggle)
    window.addEventListener('resize', this.handle_resize)
    
    // Listen for locale changes to refetch data
    this.handleLocaleChange = (event) => {
      console.log('Locale changed in WorldMap, refetching data for locale:', event.detail)
      // Emit event to parent (MapView) to refetch events with new locale
      this.$emit('locale-changed', event.detail)
    }
    window.addEventListener('localeChanged', this.handleLocaleChange)
  },
  
  beforeUnmount() {
    if (this.resize_observer) {
      this.resize_observer.disconnect()
    }
    window.removeEventListener('resize', this.handle_resize)
    
    // Remove locale change listener
    if (this.handleLocaleChange) {
      window.removeEventListener('localeChanged', this.handleLocaleChange)
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
      this.new_event.name_en = ''
      this.new_event.name_ru = ''
      this.new_event.description_en = ''
      this.new_event.description_ru = ''
      this.new_event.date = ''
      this.new_event.date_display = ''
      this.new_event.era = 'AD'
      this.new_event.lens_type = 'historic'
      this.new_event.source = ''
      this.new_event.dataset_id = null
      
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
      
      // Populate form with event data including locale-specific fields
      this.new_event.name = event.name
      this.new_event.description = event.description
      this.new_event.name_en = event.name_en || event.name || ''
      this.new_event.name_ru = event.name_ru || event.name || ''
      this.new_event.description_en = event.description_en || event.description || ''
      this.new_event.description_ru = event.description_ru || event.description || ''
      this.new_event.latitude = event.latitude
      this.new_event.longitude = event.longitude
      this.new_event.era = event.era || 'AD'
      this.new_event.lens_type = event.lens_type
      this.new_event.source = event.source || ''
      this.new_event.dataset_id = event.dataset_id || null
      
      // Handle date formatting
      this.new_event.date = event.event_date.split('T')[0]
      this.new_event.date_display = this.format_date(event.event_date)
      
      // Show modal
      this.show_event_modal = true
    },
    
    add_event_markers() {      
      // Clear existing markers first
      this.markers.forEach(marker => {
        if (this.map && this.map.hasLayer(marker)) {
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
            
            // Use click event instead of popup to avoid coordinate corruption
            marker.on('click', () => {
              this.show_events_info(eventGroup.events)
            })
            
            this.markers.push(marker)
          })
        } catch (error) {
          console.error('Error adding markers:', error)
        }
        
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
        // Use locale-specific fields as primary, fallback to legacy fields
        const name_en = this.new_event.name_en || this.new_event.name || ''
        const name_ru = this.new_event.name_ru || this.new_event.name || ''
        const description_en = this.new_event.description_en || this.new_event.description || ''
        const description_ru = this.new_event.description_ru || this.new_event.description || ''
        
        const event_data = {
          name: name_en, // Use English as legacy default
          description: description_en, // Use English as legacy default
          name_en: name_en,
          name_ru: name_ru,
          description_en: description_en,
          description_ru: description_ru,
          latitude: this.new_event.latitude,
          longitude: this.new_event.longitude,
          event_date: new Date(this.new_event.date).toISOString(),
          era: this.new_event.era,
          lens_type: this.new_event.lens_type,
          source: this.new_event.source || null,
          dataset_id: this.new_event.dataset_id || null
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
        name_en: '',
        name_ru: '',
        description_en: '',
        description_ru: '',
        date: '',
        date_display: '',
        era: 'AD',
        latitude: 0,
        longitude: 0,
        source: '',
        lens_type: 'historic',
        dataset_id: null
      }
    },

    async fetchDatasets() {
      try {
        const response = await apiService.getDatasets()
        this.datasets = Array.isArray(response) ? response.filter(d => d && d.id) : []
        console.log('Loaded datasets for event form:', this.datasets.length)
      } catch (err) {
        console.error('Error fetching datasets:', err)
        this.datasets = []
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

    get_event_emoji(lensType) {
      return getEventEmoji(lensType)
    },

    getAvailableLensTypes() {
      return getAvailableLensTypes()
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
      // Use the utility function for consistency
      const emoji = getEventEmoji(lens_type)
      
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


    // Method to set stepping mode from parent component
    setSteppingMode(isStepping) {
      this.is_stepping = isStepping
    },

    // Show events info in custom modal (no coordinate corruption)
    show_events_info(events) {
      // Ensure tags are loaded first
      if (this.allTags.length === 0) {
        this.loadTags()
      }
      
      // Enrich events with their tag information
      this.selected_events = events.map(event => ({
        ...event,
        tags: event.tags || [] // Ensure tags array exists
      }))
      this.show_event_info_modal = true
    },

    // Close event info modal
    close_event_info_modal() {
      this.show_event_info_modal = false
      this.selected_events = []
    },

    // Edit event from info modal
    edit_event_from_info(eventId) {
      this.close_event_info_modal()
      this.edit_event(eventId)
    },

    // Get contrast color for text on colored backgrounds
    getContrastColor(hexColor) {
      if (!hexColor) return '#000000'
      
      // Remove # if present
      const color = hexColor.replace('#', '')
      
      // Convert to RGB
      const r = parseInt(color.substr(0, 2), 16)
      const g = parseInt(color.substr(2, 2), 16)
      const b = parseInt(color.substr(4, 2), 16)
      
      // Calculate luminance
      const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
      
      // Return black for light colors, white for dark colors
      return luminance > 0.5 ? '#000000' : '#ffffff'
    }
  }
}
</script>

<style scoped>
.event-tags {
  margin-top: 0.75rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  align-items: center;
}

.event-tag {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.more-tags {
  font-size: 0.75rem;
  color: #6b7280;
  font-style: italic;
}
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

/* Locale tabs styling */
.locale-tabs {
  display: flex;
  margin-bottom: 15px;
  border-bottom: 1px solid #ddd;
}

.locale-tab {
  flex: 1;
  padding: 8px 16px;
  border: none;
  background: #f5f5f5;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.locale-tab:hover {
  background: #e9e9e9;
}

.locale-tab.active {
  background: white;
  border-bottom-color: #007cba;
  color: #007cba;
}

.locale-content {
  margin-bottom: 10px;
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
  cursor: pointer;
  transition: transform 0.2s ease;
}

:deep(.emoji-marker:hover) {
  transform: scale(1.1);
}

/* Event Info Modal Styles */
.event-info-modal {
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  width: 90vw;
  padding: 1rem 1.5rem;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding: 0.75rem 0;
  border-bottom: 2px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  color: #2d3748;
  font-size: 1.25rem;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  color: #a0aec0;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
}

.close-button:hover {
  background: #f1f5f9;
  color: #4a5568;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-edit {
  font-size: 18px;
  cursor: pointer;
  opacity: 0.7;
  transition: all 0.2s;
  padding: 0.375rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
}

.header-edit:hover {
  opacity: 1;
  background: #f1f5f9;
  transform: scale(1.1);
}

.events-list {
  max-height: 60vh;
  overflow-y: auto;
}

.event-info-item {
  padding: 0.75rem 0;
}

.event-info-item:not(:last-child) {
  border-bottom: 1px solid #e2e8f0;
}

.event-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.75rem;
}

.event-header h4 {
  margin: 0;
  color: #2d3748;
  font-size: 1.1rem;
  font-weight: 600;
  flex: 1;
  margin-right: 0.5rem;
}

.edit-icon {
  font-size: 16px;
  cursor: pointer;
  opacity: 0.7;
  transition: all 0.2s;
  padding: 0.25rem;
  border-radius: 4px;
  flex-shrink: 0;
}

.edit-icon:hover {
  opacity: 1;
  background: #f1f5f9;
  transform: scale(1.1);
}

.modal-event-title-link {
  color: #3b82f6;
  text-decoration: none;
  transition: all 0.2s;
  font-weight: normal;
  margin-left: 0.5rem;
}

.modal-event-title-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.modal-event-title-text {
  margin-left: 0.5rem;
}

.modal-event-header-link {
  color: #3b82f6;
  text-decoration: none;
  transition: all 0.2s;
  font-weight: 600;
  margin-left: 0.5rem;
}

.modal-event-header-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.event-description {
  color: #4a5568;
  margin: 0.5rem 0 0 0;
  line-height: 1.5;
}

.event-date {
  color: #6b7280;
  font-size: 0.9rem;
  margin: 0.25rem 0 0.5rem 0;
  font-weight: 500;
}

.single-event-edit {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 0.5rem;
}

.event-details p {
  margin: 0.25rem 0;
  font-size: 0.9rem;
  color: #6b7280;
}

.event-details strong {
  color: #374151;
}

.event-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 1rem 0;
}
</style>