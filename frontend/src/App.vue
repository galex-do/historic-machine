<template>
  <div id="app">
    <header>
      <h1>Historical Events Mapping</h1>
      <div class="timeline-controls">
        <!-- Date Selection Mode -->
        <div class="date-mode-selector">
          <label>Date Selection:</label>
          <select v-model="date_selection_mode" @change="handle_date_mode_change">
            <option value="historic">Historic Periods</option>
            <option value="custom">Custom Date Range</option>
          </select>
        </div>
        
        <!-- Historic Date Templates -->
        <div v-if="date_selection_mode === 'historic'" class="historic-date-controls">
          <div class="template-group-selector">
            <label>Historical Period:</label>
            <select v-model="selected_template_group_id" @change="handle_template_group_change">
              <option value="">Select a historical period...</option>
              <option v-for="group in template_groups" :key="group.id" :value="group.id">
                {{ group.name }}
              </option>
            </select>
          </div>
          
          <div v-if="selected_template_group_id" class="template-selector">
            <label>Specific Period:</label>
            <select v-model="selected_template_id" @change="handle_template_change">
              <option value="">Select specific period...</option>
              <option v-for="template in available_templates" :key="template.id" :value="template.id">
                {{ template.name }} ({{ template.start_display_date }} - {{ template.end_display_date }})
              </option>
            </select>
          </div>
          
          <div v-if="selected_template" class="selected-period-display">
            <strong>Selected Period:</strong> {{ selected_template.name }}<br>
            <strong>Duration:</strong> {{ selected_template.start_display_date }} to {{ selected_template.end_display_date }}
            <p v-if="selected_template.description" class="period-description">{{ selected_template.description }}</p>
          </div>
        </div>
        
        <!-- Custom Date Range (existing functionality) -->
        <div v-if="date_selection_mode === 'custom'" class="custom-date-controls">
          <label>From: <input type="text" v-model="date_from_display" @blur="update_date_from" placeholder="DD.MM.YYYY" /></label>
          <label>To: <input type="text" v-model="date_to_display" @blur="update_date_to" placeholder="DD.MM.YYYY" /></label>
        </div>
        
        <!-- Event Type Lens Selector -->
        <div class="lens-selector">
          <label>Event Types:</label>
          <div class="lens-options">
            <label v-for="lens in available_lens_types" :key="lens.value" class="lens-checkbox">
              <input type="checkbox" :value="lens.value" v-model="selected_lens_types" @change="filter_events" />
              <span :class="['lens-label', lens.value]">{{ lens.label }}</span>
            </label>
          </div>
        </div>
        
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
          <p><strong>Date:</strong> {{ event.display_date || formatDate(event.event_date) }}</p>
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
      date_from: '0001-01-01', // Internal ISO format for filtering - start from year 1
      date_to: '', // Will be set in created() hook
      date_from_display: '01.01.0001', // Display format DD.MM.YYYY - start from year 1
      date_to_display: '', // Will be set in created() hook
      available_lens_types: [
        { value: 'historic', label: '📜 Historic' },
        { value: 'political', label: '🏛️ Political' },
        { value: 'cultural', label: '🎭 Cultural' },
        { value: 'military', label: '⚔️ Military' }
      ],
      selected_lens_types: ['historic', 'political', 'cultural', 'military'], // All selected by default
      
      // Date template system
      date_selection_mode: 'historic', // 'historic' or 'custom'
      template_groups: [],
      available_templates: [],
      selected_template_group_id: '',
      selected_template_id: '',
      selected_template: null
    }
  },
  created() {
    // Set today's date in DD.MM.YYYY format
    this.date_to_display = this.format_date_to_ddmmyyyy(new Date())
    this.date_to = this.get_today_iso()
  },
  async mounted() {
    await this.fetch_template_groups()
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
      this.filtered_events = this.events.filter(event => {
        // Date filtering
        const event_date = new Date(event.event_date)
        const from_date = this.date_from ? new Date(this.date_from) : null
        const to_date = this.date_to ? new Date(this.date_to) : null
        
        if (from_date && event_date < from_date) return false
        if (to_date && event_date > to_date) return false
        
        // Lens type filtering
        if (this.selected_lens_types.length > 0 && !this.selected_lens_types.includes(event.lens_type)) {
          return false
        }
        
        return true
      })
      
      const lens_filter_text = this.selected_lens_types.length === 4 ? 'all types' : this.selected_lens_types.join(', ')
      console.log(`Filtering events from ${this.date_from} to ${this.date_to} for lens types: ${lens_filter_text}. Found ${this.filtered_events.length} events.`)
    },
    
    async handle_event_created(new_event) {
      // Refresh events list after new event is created
      await this.fetch_events()
      console.log('Events list refreshed after new event creation')
    },
    
    formatDate(dateString) {
      // For very old dates, parse manually to avoid Date object issues
      if (dateString.startsWith('00') || dateString.startsWith('01')) {
        return this.format_date_to_ddmmyyyy(dateString)
      }
      return this.format_date_to_ddmmyyyy(new Date(dateString))
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
    
    get_today_iso() {
      return new Date().toISOString().split('T')[0]
    },
    
    update_date_from() {
      const iso_date = this.parse_ddmmyyyy_to_iso(this.date_from_display)
      if (iso_date) {
        this.date_from = iso_date
      } else {
        // Reset to default if invalid
        this.date_from_display = '01.01.0001'
        this.date_from = '0001-01-01'
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
    },
    
    // Date template methods
    async fetch_template_groups() {
      try {
        const backend_url = this.get_backend_url()
        const response = await fetch(`${backend_url}/date-template-groups`)
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        this.template_groups = await response.json()
        console.log('Successfully loaded template groups:', this.template_groups.length)
      } catch (error) {
        console.error('Error fetching template groups:', error)
      }
    },
    
    async fetch_templates_for_group(group_id) {
      try {
        const backend_url = this.get_backend_url()
        const response = await fetch(`${backend_url}/date-templates/${group_id}`)
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        this.available_templates = await response.json()
        console.log('Successfully loaded templates for group:', this.available_templates.length)
      } catch (error) {
        console.error('Error fetching templates for group:', error)
        this.available_templates = []
      }
    },
    
    handle_date_mode_change() {
      // Clear template selections when switching modes
      if (this.date_selection_mode === 'historic') {
        this.selected_template_group_id = ''
        this.selected_template_id = ''
        this.selected_template = null
        this.available_templates = []
        // Set default date range to encompass all history
        this.date_from = '0001-01-01'
        this.date_to = this.get_today_iso()
      } else {
        // Reset to custom date defaults
        this.date_from = '0001-01-01'
        this.date_to = this.get_today_iso()
        this.date_from_display = '01.01.0001'
        this.date_to_display = this.format_date_to_ddmmyyyy(new Date())
      }
      this.filter_events()
    },
    
    async handle_template_group_change() {
      this.selected_template_id = ''
      this.selected_template = null
      
      if (this.selected_template_group_id) {
        await this.fetch_templates_for_group(this.selected_template_group_id)
      } else {
        this.available_templates = []
      }
      
      // Reset to full range when no specific template selected
      this.date_from = '0001-01-01'
      this.date_to = this.get_today_iso()
      this.filter_events()
    },
    
    handle_template_change() {
      if (this.selected_template_id) {
        this.selected_template = this.available_templates.find(t => t.id == this.selected_template_id)
        
        if (this.selected_template) {
          // Set date range based on selected template
          this.date_from = this.selected_template.start_date
          this.date_to = this.selected_template.end_date
          console.log(`Selected template: ${this.selected_template.name} (${this.selected_template.start_display_date} - ${this.selected_template.end_display_date})`)
        }
      } else {
        this.selected_template = null
        // Reset to full range
        this.date_from = '0001-01-01'
        this.date_to = this.get_today_iso()
      }
      this.filter_events()
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

.timeline-controls input[type="text"] {
  margin: 0 5px;
  padding: 5px;
}

.timeline-controls button {
  padding: 5px 15px;
  background: #3498db;
  color: white;
  border: none;
  cursor: pointer;
  margin-left: 10px;
}

.lens-selector {
  margin: 10px 0;
  padding: 10px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
}

.lens-selector > label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
}

.lens-options {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.lens-checkbox {
  display: flex;
  align-items: center;
  margin-right: 0 !important;
  cursor: pointer;
}

.lens-checkbox input[type="checkbox"] {
  margin-right: 6px;
  margin-left: 0;
}

.lens-label {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: bold;
  transition: all 0.3s ease;
}

.lens-label.historic {
  background: #8e44ad;
  color: white;
}

.lens-label.political {
  background: #e74c3c;
  color: white;
}

.lens-label.cultural {
  background: #f39c12;
  color: white;
}

.lens-label.military {
  background: #2c3e50;
  color: white;
}

.lens-checkbox input:not(:checked) + .lens-label {
  opacity: 0.4;
  background: #bdc3c7 !important;
  color: #7f8c8d !important;
}

/* Date selection mode styling */
.date-mode-selector {
  margin-bottom: 15px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 6px;
}

.date-mode-selector label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
}

.date-mode-selector select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  color: #2c3e50;
  font-size: 14px;
}

/* Historic date controls */
.historic-date-controls {
  background: rgba(255, 255, 255, 0.1);
  padding: 15px;
  border-radius: 6px;
  margin-bottom: 15px;
}

.template-group-selector,
.template-selector {
  margin-bottom: 15px;
}

.template-group-selector label,
.template-selector label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: white;
}

.template-group-selector select,
.template-selector select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  color: #2c3e50;
  font-size: 14px;
}

.selected-period-display {
  background: rgba(255, 255, 255, 0.2);
  padding: 15px;
  border-radius: 6px;
  border-left: 4px solid #3498db;
  margin-top: 15px;
}

.selected-period-display strong {
  color: #ecf0f1;
}

.period-description {
  margin-top: 8px;
  font-style: italic;
  color: #bdc3c7;
  font-size: 13px;
}

/* Custom date controls */
.custom-date-controls {
  background: rgba(255, 255, 255, 0.1);
  padding: 15px;
  border-radius: 6px;
  margin-bottom: 15px;
}

.custom-date-controls label {
  margin-right: 15px;
  color: white;
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