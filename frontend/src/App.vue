<template>
  <div id="app">
    <!-- Header -->
    <header class="app-header">
      <h1>🗺️ Historical Events Mapping</h1>
    </header>
    
    <!-- Main Layout: Sidebar + Map -->
    <div class="main-layout">
      <!-- Left Sidebar -->
      <aside class="sidebar">
        <!-- Filters Section -->
        <div class="sidebar-section">
          <h3 class="section-title">🎯 Filters</h3>
          
          <!-- Date Selection Mode -->
          <div class="filter-group">
            <label class="filter-label">Date Selection:</label>
            <select v-model="date_selection_mode" @change="handle_date_mode_change" class="filter-select">
              <option value="historic">Historic Periods</option>
              <option value="custom">Custom Date Range</option>
            </select>
          </div>
          
          <!-- Historic Date Templates -->
          <div v-if="date_selection_mode === 'historic'" class="filter-group">
            <label class="filter-label">Historical Period:</label>
            <select v-model="selected_template_group_id" @change="handle_template_group_change" class="filter-select">
              <option value="">Select a period...</option>
              <option v-for="group in template_groups" :key="group.id" :value="group.id">
                {{ group.name }}
              </option>
            </select>
            
            <div v-if="selected_template_group_id" class="filter-subgroup">
              <label class="filter-label">Specific Period:</label>
              <select v-model="selected_template_id" @change="handle_template_change" class="filter-select">
                <option value="">Select specific period...</option>
                <option v-for="template in available_templates" :key="template.id" :value="template.id">
                  {{ template.name }}
                </option>
              </select>
            </div>
            
            <div v-if="selected_template" class="selected-period">
              <div class="period-info">
                <strong>{{ selected_template.name }}</strong>
                <span class="period-dates">{{ selected_template.start_display_date }} - {{ selected_template.end_display_date }}</span>
                <p v-if="selected_template.description" class="period-desc">{{ selected_template.description }}</p>
              </div>
            </div>
          </div>
          
          <!-- Custom Date Range -->
          <div v-if="date_selection_mode === 'custom'" class="filter-group">
            <div class="date-inputs">
              <div class="date-input-group">
                <label class="filter-label">From:</label>
                <input type="text" v-model="date_from_display" @blur="update_date_from" placeholder="DD.MM.YYYY" class="date-input" />
              </div>
              <div class="date-input-group">
                <label class="filter-label">To:</label>
                <input type="text" v-model="date_to_display" @blur="update_date_to" placeholder="DD.MM.YYYY" class="date-input" />
              </div>
            </div>
          </div>
          
          <!-- Event Types -->
          <div class="filter-group">
            <label class="filter-label">Event Types:</label>
            <div class="multiselect-container">
              <div class="multiselect-dropdown" @click="toggle_lens_dropdown" :class="{ 'active': show_lens_dropdown }">
                <div class="multiselect-display">
                  <span v-if="selected_lens_types.length === 0" class="placeholder">Select event types...</span>
                  <span v-else-if="selected_lens_types.length === available_lens_types.length" class="selection-text">All types selected</span>
                  <div v-else class="selected-badges">
                    <span v-for="type in selected_lens_types" :key="type" :class="['selected-badge', type]">
                      {{ getEventEmoji(type) }} {{ getLensLabel(type) }}
                    </span>
                  </div>
                </div>
                <span class="dropdown-arrow">{{ show_lens_dropdown ? '▲' : '▼' }}</span>
              </div>
              
              <div v-show="show_lens_dropdown" class="multiselect-options">
                <label v-for="lens in available_lens_types" :key="lens.value" class="multiselect-option">
                  <input 
                    type="checkbox" 
                    :value="lens.value" 
                    v-model="selected_lens_types" 
                    @change="filter_events" 
                  />
                  <span class="option-content">
                    <span class="option-emoji">{{ getEventEmoji(lens.value) }}</span>
                    <span class="option-label">{{ lens.label }}</span>
                  </span>
                </label>
              </div>
            </div>
          </div>
          
          <button @click="filter_events" class="filter-button">Apply Filters</button>
        </div>
        
        <!-- Events List Section -->
        <div class="sidebar-section">
          <h3 class="section-title">📍 Events ({{ filtered_events.length }})</h3>
          
          <div class="events-container">
            <div v-if="filtered_events.length === 0" class="no-events">
              <p>No events found. Click on the map to add your first historical event!</p>
            </div>
            
            <div v-for="event in filtered_events" :key="event.id" class="event-card">
              <div class="event-header">
                <span class="event-emoji">{{ getEventEmoji(event.lens_type) }}</span>
                <h4 class="event-title">{{ event.name }}</h4>
              </div>
              <p class="event-description">{{ event.description }}</p>
              <div class="event-meta">
                <span class="event-date">{{ event.display_date || formatDate(event.event_date) }}</span>
                <span class="event-coords">{{ event.latitude.toFixed(2) }}, {{ event.longitude.toFixed(2) }}</span>
              </div>
            </div>
          </div>
        </div>
      </aside>
      
      <!-- Right Map Area -->
      <main class="map-area">
        <WorldMap 
          :events="filtered_events" 
          @event-created="handle_event_created"
        />
      </main>
    </div>
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
      show_lens_dropdown: false,
      
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
    
    // Add click outside handler for dropdown
    document.addEventListener('click', this.handle_click_outside)
  },
  
  beforeUnmount() {
    document.removeEventListener('click', this.handle_click_outside)
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
        // Date filtering using astronomical year comparison for BC dates
        const event_astronomical_year = this.get_astronomical_year(event.event_date, event.era)
        const from_astronomical_year = this.date_from ? this.get_astronomical_year(this.date_from, this.get_era_from_date(this.date_from)) : null
        const to_astronomical_year = this.date_to ? this.get_astronomical_year(this.date_to, this.get_era_from_date(this.date_to)) : null
        
        if (from_astronomical_year !== null && event_astronomical_year < from_astronomical_year) return false
        if (to_astronomical_year !== null && event_astronomical_year > to_astronomical_year) return false
        
        // For same astronomical year, compare the actual dates
        if (from_astronomical_year !== null && event_astronomical_year === from_astronomical_year) {
          if (event.event_date < this.date_from) return false
        }
        if (to_astronomical_year !== null && event_astronomical_year === to_astronomical_year) {
          if (event.event_date > this.date_to) return false
        }
        
        // Lens type filtering
        if (this.selected_lens_types.length > 0 && !this.selected_lens_types.includes(event.lens_type)) {
          return false
        }
        
        return true
      })
      
      const lens_filter_text = this.selected_lens_types.length === 4 ? 'all types' : this.selected_lens_types.join(', ')
      console.log(`Filtering events from ${this.date_from} to ${this.date_to} for lens types: ${lens_filter_text}. Found ${this.filtered_events.length} events.`)
    },
    
    get_astronomical_year(date_string, era) {
      if (!date_string) return null
      
      const year = parseInt(date_string.split('-')[0], 10)
      
      if (era === 'BC') {
        return (year * -1) + 1  // Convert BC year to astronomical year
      } else {
        return year  // AD years are already correct
      }
    },
    
    get_era_from_date(date_string) {
      if (!date_string) return 'AD'
      
      const year = parseInt(date_string.split('-')[0], 10)
      // For our system, years <= 0 are BC, years > 0 are AD
      // But we store BC years as positive numbers with era='BC'
      // So we need to check if this is part of a template selection
      
      // If we're in template mode and have a selected template, use its era
      if (this.date_selection_mode === 'historic' && this.selected_template) {
        if (date_string === this.selected_template.start_date) {
          return this.selected_template.start_era
        }
        if (date_string === this.selected_template.end_date) {
          return this.selected_template.end_era
        }
      }
      
      // Default assumption based on year
      return year < 1000 ? 'BC' : 'AD'  // Heuristic for old dates
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
    },
    
    getEventEmoji(lens_type) {
      const emoji_map = {
        'military': '⚔️',
        'political': '🏛️', 
        'historic': '📜',
        'cultural': '🎭'
      }
      return emoji_map[lens_type] || '📍'
    },
    
    getLensLabel(lens_type) {
      const lens_obj = this.available_lens_types.find(lens => lens.value === lens_type)
      return lens_obj ? lens_obj.label.replace(/[⚔️🏛️📜🎭]\s/, '') : lens_type
    },
    
    toggle_lens_dropdown() {
      this.show_lens_dropdown = !this.show_lens_dropdown
    },
    
    handle_click_outside(event) {
      const dropdown = event.target.closest('.multiselect-container')
      if (!dropdown) {
        this.show_lens_dropdown = false
      }
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
  background: #f5f7fa;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header */
.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 2rem;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.app-header h1 {
  font-size: 1.5rem;
  font-weight: 600;
  text-align: center;
  margin: 0;
}

/* Main Layout */
.main-layout {
  display: flex;
  flex: 1;
  height: calc(100vh - 80px);
  overflow: hidden;
}

/* Sidebar */
.sidebar {
  width: 350px;
  background: #ffffff;
  border-right: 1px solid #e1e8ed;
  overflow-y: auto;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.05);
}

.sidebar-section {
  padding: 1.5rem;
  border-bottom: 1px solid #f1f3f5;
}

.sidebar-section:last-child {
  border-bottom: none;
}

.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Filter Controls */
.filter-group {
  margin-bottom: 1.5rem;
}

.filter-group:last-child {
  margin-bottom: 0;
}

.filter-label {
  display: block;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.filter-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  font-size: 0.9rem;
  color: #2d3748;
  transition: all 0.2s;
}

.filter-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.filter-subgroup {
  margin-top: 1rem;
  padding-left: 1rem;
  border-left: 2px solid #e2e8f0;
}

/* Selected Period */
.selected-period {
  margin-top: 1rem;
  background: #f7fafc;
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid #e2e8f0;
}

.period-info strong {
  display: block;
  color: #2d3748;
  margin-bottom: 0.25rem;
}

.period-dates {
  font-size: 0.8rem;
  color: #667eea;
  font-weight: 500;
}

.period-desc {
  font-size: 0.8rem;
  color: #718096;
  margin-top: 0.5rem;
  line-height: 1.4;
}

/* Date Inputs */
.date-inputs {
  display: flex;
  gap: 1rem;
}

.date-input-group {
  flex: 1;
}

.date-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #2d3748;
  transition: all 0.2s;
}

.date-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* Multi-Select Dropdown */
.multiselect-container {
  position: relative;
}

.multiselect-dropdown {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.2s;
  min-height: 45px;
}

.multiselect-dropdown:hover {
  border-color: #cbd5e0;
}

.multiselect-dropdown.active {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.multiselect-display {
  flex: 1;
}

.placeholder {
  color: #a0aec0;
  font-size: 0.9rem;
}

.selection-text {
  color: #2d3748;
  font-size: 0.9rem;
  font-weight: 500;
}

.selected-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.selected-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 500;
}

.selected-badge.military {
  background: #fed7d7;
  color: #c53030;
}

.selected-badge.political {
  background: #bee3f8;
  color: #2b6cb0;
}

.selected-badge.historic {
  background: #e2e8f0;
  color: #4a5568;
}

.selected-badge.cultural {
  background: #c6f6d5;
  color: #25855a;
}

.dropdown-arrow {
  color: #718096;
  font-size: 0.8rem;
  transition: transform 0.2s;
}

.multiselect-dropdown.active .dropdown-arrow {
  transform: rotate(180deg);
}

.multiselect-options {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 10;
  margin-top: 0.25rem;
  max-height: 200px;
  overflow-y: auto;
}

.multiselect-option {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.multiselect-option:hover {
  background: #f7fafc;
}

.multiselect-option input {
  margin-right: 0.75rem;
  accent-color: #667eea;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.option-emoji {
  font-size: 1.1rem;
}

.option-label {
  font-size: 0.9rem;
  color: #2d3748;
}

/* Filter Button */
.filter-button {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 1rem;
}

.filter-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* Events Container */
.events-container {
  max-height: calc(100vh - 400px);
  overflow-y: auto;
  padding-right: 0.5rem;
}

.events-container::-webkit-scrollbar {
  width: 4px;
}

.events-container::-webkit-scrollbar-track {
  background: #f1f3f5;
  border-radius: 2px;
}

.events-container::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 2px;
}

.events-container::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* Event Cards */
.event-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 0.75rem;
  transition: all 0.2s;
  cursor: pointer;
}

.event-card:hover {
  border-color: #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.event-header {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.event-emoji {
  font-size: 1.25rem;
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.event-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: #2d3748;
  line-height: 1.3;
  margin: 0;
}

.event-description {
  font-size: 0.85rem;
  color: #718096;
  line-height: 1.4;
  margin: 0.5rem 0;
}

.event-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.75rem;
  padding-top: 0.5rem;
  border-top: 1px solid #f1f3f5;
}

.event-date {
  font-size: 0.8rem;
  font-weight: 600;
  color: #667eea;
}

.event-coords {
  font-size: 0.75rem;
  color: #a0aec0;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', monospace;
}

.no-events {
  text-align: center;
  color: #a0aec0;
  font-style: italic;
  padding: 2rem 1rem;
  font-size: 0.9rem;
}

/* Map Area */
.map-area {
  flex: 1;
  background: #ffffff;
  position: relative;
  overflow: hidden;
}

/* Responsive Design */
@media (max-width: 1024px) {
  .main-layout {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    max-height: 50vh;
    border-right: none;
    border-bottom: 1px solid #e1e8ed;
  }
  
  .map-area {
    height: 50vh;
  }
  
  .lens-grid {
    grid-template-columns: 1fr;
  }
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