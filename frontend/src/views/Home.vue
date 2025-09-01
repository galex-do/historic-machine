<template>
  <div class="home-page">
    <!-- Main Layout: Sidebar + Map -->
    <div class="main-layout">
      <!-- Left Sidebar (Collapsible) -->
      <aside class="sidebar" :class="{ 'collapsed': sidebar_collapsed }">
        <!-- Sidebar Toggle Button -->
        <button class="sidebar-toggle" @click="toggle_sidebar">
          {{ sidebar_collapsed ? '›' : '‹' }}
        </button>
        
        <!-- Filters Section -->
        <div class="sidebar-section" v-show="!sidebar_collapsed">
          <h3 class="section-title">Filters</h3>
          
          <!-- Date Selection Mode -->
          <div class="filter-group">
            <label class="filter-label">Date Selection:</label>
            <div class="radio-group">
              <label class="radio-option">
                <input type="radio" value="historic" v-model="date_selection_mode" @change="handle_date_selection_change" />
                <span>Historical Periods</span>
              </label>
              <label class="radio-option">
                <input type="radio" value="custom" v-model="date_selection_mode" @change="handle_date_selection_change" />
                <span>Custom Range</span>
              </label>
            </div>
          </div>

          <!-- Historical Templates (when historic mode is selected) -->
          <div v-if="date_selection_mode === 'historic'" class="filter-group">
            <label class="filter-label">Historical Period:</label>
            <select v-model="selected_template_group" @change="apply_template_group" class="form-select">
              <option value="">All History (Year 1 - Present)</option>
              <optgroup v-for="group in template_groups" :key="group.id" :label="group.name">
                <option 
                  v-for="template in group.templates" 
                  :key="template.id" 
                  :value="template.id"
                >
                  {{ template.name }}
                </option>
              </optgroup>
            </select>
          </div>

          <!-- Custom Date Range (when custom mode is selected) -->
          <div v-if="date_selection_mode === 'custom'" class="filter-group">
            <label class="filter-label">Custom Date Range:</label>
            <div class="date-inputs">
              <div class="date-input-group">
                <label class="date-label">From:</label>
                <input 
                  type="date" 
                  v-model="date_from" 
                  @change="filter_events" 
                  class="form-input"
                />
              </div>
              <div class="date-input-group">
                <label class="date-label">To:</label>
                <input 
                  type="date" 
                  v-model="date_to" 
                  @change="filter_events" 
                  class="form-input"
                />
              </div>
            </div>
          </div>

          <!-- Event Type Filter -->
          <div class="filter-group">
            <label class="filter-label">Event Types:</label>
            <div class="multiselect-container">
              <div class="multiselect-header" @click="show_lens_dropdown = !show_lens_dropdown">
                <span class="multiselect-value">
                  <span v-if="selected_lens_types.length === 0" class="placeholder">No types selected</span>
                  <span v-else-if="selected_lens_types.length === available_lens_types.length" class="all-selected">All types selected</span>
                  <span v-else class="selected-count">{{ selected_lens_types.length }} types selected</span>
                </span>
                <span class="multiselect-arrow">{{ show_lens_dropdown ? '▲' : '▼' }}</span>
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
                    <span class="option-label">{{ lens.label }}</span>
                  </span>
                </label>
              </div>
            </div>
          </div>
          
          <button @click="filter_events" class="filter-button">Apply Filters</button>
        </div>
      </aside>
      
      <!-- Right Content Area -->
      <main class="content-area">
        <!-- Map Section -->
        <div class="map-section">
          <WorldMap 
            :events="filtered_events" 
            :focus_event="focus_event"
            @event-created="handle_event_created"
          />
        </div>
        
        <!-- Events List Below Map -->
        <div class="events-section">
          <h3 class="section-title">Historical Events ({{ filtered_events.length }})</h3>
          
          <div class="events-grid">
            <div v-if="filtered_events.length === 0" class="no-events">
              <p>No events found for the selected period. Click on the map to add your first historical event!</p>
            </div>
            
            <div v-for="event in filtered_events" :key="event.id" class="event-card">
              <div class="event-header">
                <span class="event-emoji">{{ getEventEmoji(event.lens_type) }}</span>
                <h4 class="event-title">{{ event.name }}</h4>
                <button class="focus-button" @click="focus_on_event(event)" title="Focus on map">
                  ⌖
                </button>
              </div>
              <p class="event-description">{{ event.description }}</p>
              <div class="event-meta">
                <span class="event-date">{{ event.display_date || formatDate(event.event_date) }}</span>
                <span class="event-coords">{{ event.latitude.toFixed(2) }}, {{ event.longitude.toFixed(2) }}</span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import WorldMap from '../components/WorldMap.vue'

export default {
  name: 'Home',
  components: {
    WorldMap
  },
  data() {
    return {
      events: [],
      filtered_events: [],
      
      // Filter state
      date_from: '0001-01-01', // Start from year 1
      date_to: new Date().toISOString().split('T')[0], // Today
      available_lens_types: [
        { value: 'historic', label: '📜 Historic' },
        { value: 'political', label: '🏛️ Political' },
        { value: 'cultural', label: '🎭 Cultural' },
        { value: 'military', label: '⚔️ Military' }
      ],
      selected_lens_types: ['historic', 'political', 'cultural', 'military'], // All selected by default
      show_lens_dropdown: false,
      sidebar_collapsed: false,
      focus_event: null,
      
      // Date template system
      date_selection_mode: 'historic', // 'historic' or 'custom'
      template_groups: [],
      selected_template_group: '', // ID of selected template, empty for "All History"
    }
  },
  
  async mounted() {
    await this.load_template_groups()
    await this.load_events()
    this.filter_events()
    
    // Close dropdown when clicking outside
    document.addEventListener('click', this.handle_outside_click)
  },
  
  beforeUnmount() {
    document.removeEventListener('click', this.handle_outside_click)
  },
  
  methods: {
    async load_template_groups() {
      try {
        const response = await fetch('http://localhost:8080/api/date-template-groups')
        if (response.ok) {
          this.template_groups = await response.json()
          console.log('Successfully loaded template groups:', this.template_groups.length)
        } else {
          console.error('Failed to load template groups:', response.status)
        }
      } catch (error) {
        console.error('Error loading template groups:', error)
      }
    },

    async load_events() {
      try {
        const response = await fetch('http://localhost:8080/api/events')
        if (response.ok) {
          this.events = await response.json()
          console.log('Successfully loaded events:', this.events.length)
        } else {
          console.error('Failed to load events:', response.status)
        }
      } catch (error) {
        console.error('Error loading events:', error)
      }
    },

    filter_events() {
      let filtered = [...this.events]
      
      // Apply date filtering
      if (this.date_from && this.date_to) {
        filtered = filtered.filter(event => {
          const event_date = new Date(event.event_date)
          const from_date = new Date(this.date_from + 'T00:00:00')
          const to_date = new Date(this.date_to + 'T23:59:59')
          
          return event_date >= from_date && event_date <= to_date
        })
      }
      
      // Apply lens type filtering
      if (this.selected_lens_types.length > 0) {
        filtered = filtered.filter(event => 
          this.selected_lens_types.includes(event.lens_type)
        )
      }
      
      this.filtered_events = filtered
      
      const lens_types_display = this.selected_lens_types.length === this.available_lens_types.length 
        ? 'all types' 
        : this.selected_lens_types.join(', ')
      
      console.log(`Filtering events from ${this.date_from} to ${this.date_to} for lens types: ${lens_types_display}. Found ${filtered.length} events.`)
    },

    apply_template_group() {
      if (!this.selected_template_group) {
        // "All History" selected - show everything from year 1 to present
        this.date_from = '0001-01-01'
        this.date_to = new Date().toISOString().split('T')[0]
      } else {
        // Find the selected template
        let selected_template = null
        for (const group of this.template_groups) {
          const template = group.templates.find(t => t.id == this.selected_template_group)
          if (template) {
            selected_template = template
            break
          }
        }
        
        if (selected_template) {
          this.date_from = selected_template.date_from
          this.date_to = selected_template.date_to
        }
      }
      
      this.filter_events()
    },

    handle_date_selection_change() {
      if (this.date_selection_mode === 'historic') {
        // Reset to current template or "All History"
        this.apply_template_group()
      } else {
        // When switching to custom, keep current dates but don't override with template
        this.filter_events()
      }
    },

    async handle_event_created(new_event) {
      try {
        // Add to local events array
        this.events.push(new_event)
        
        // Re-filter to update display
        this.filter_events()
        
        console.log('Event added successfully:', new_event.name)
      } catch (error) {
        console.error('Error handling event creation:', error)
      }
    },

    getEventEmoji(lens_type) {
      const emoji_map = {
        historic: '📜',
        political: '🏛️', 
        cultural: '🎭',
        military: '⚔️'
      }
      return emoji_map[lens_type] || '📍'
    },

    formatDate(date_string) {
      try {
        const date = new Date(date_string)
        
        // Check if it's a BC date (before year 1)
        if (date.getFullYear() <= 0) {
          // Handle BC dates - convert from astronomical year numbering
          const bc_year = Math.abs(date.getFullYear() - 1)
          return `${bc_year} BC`
        }
        
        // AD dates
        const year = date.getFullYear()
        if (year <= 1000) {
          return `${year} AD`
        } else {
          return year.toString()
        }
      } catch (error) {
        return date_string
      }
    },

    handle_outside_click(event) {
      const multiselect = event.target.closest('.multiselect-container')
      if (!multiselect) {
        this.show_lens_dropdown = false
      }
    },
    
    toggle_sidebar() {
      this.sidebar_collapsed = !this.sidebar_collapsed
      // Trigger map resize after sidebar animation
      this.$nextTick(() => {
        setTimeout(() => {
          this.trigger_map_resize()
        }, 300) // Match transition duration
      })
    },
    
    trigger_map_resize() {
      // Trigger a custom event that the map can listen for
      window.dispatchEvent(new Event('resize'))
    },
    
    focus_on_event(event) {
      this.focus_event = { ...event, timestamp: Date.now() }
    }
  }
}
</script>

<style scoped>
/* Home page styles */
.home-page {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* Main Layout */
.main-layout {
  display: flex;
  flex: 1;
  height: calc(100vh - 80px);
  overflow: hidden;
}

/* Sidebar Styles */
.sidebar {
  width: 350px;
  background: white;
  border-right: 1px solid #e2e8f0;
  transition: all 0.3s ease;
  position: relative;
  display: flex;
  flex-direction: column;
}

.sidebar.collapsed {
  width: 50px;
}

.sidebar-toggle {
  position: absolute;
  top: 50%;
  right: -15px;
  transform: translateY(-50%);
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #667eea;
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: all 0.2s;
}

.sidebar-toggle:hover {
  background: #5a67d8;
  transform: translateY(-50%) scale(1.1);
}

.sidebar-section {
  padding: 1.5rem;
  border-bottom: 1px solid #f7fafc;
}

.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 1rem;
}

.filter-group {
  margin-bottom: 1.5rem;
}

.filter-label {
  display: block;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.radio-option input[type="radio"] {
  margin: 0;
}

.form-select, .form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.form-select:focus, .form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.date-inputs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.date-input-group {
  display: flex;
  flex-direction: column;
}

.date-label {
  font-size: 0.8rem;
  color: #6b7280;
  margin-bottom: 0.25rem;
}

.multiselect-container {
  position: relative;
}

.multiselect-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  cursor: pointer;
  background: white;
  transition: border-color 0.2s;
}

.multiselect-header:hover {
  border-color: #667eea;
}

.multiselect-value {
  font-size: 0.9rem;
}

.placeholder {
  color: #9ca3af;
}

.all-selected, .selected-count {
  color: #4a5568;
}

.multiselect-arrow {
  color: #6b7280;
  font-size: 0.8rem;
}

.multiselect-options {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-top: none;
  border-radius: 0 0 4px 4px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 100;
  max-height: 200px;
  overflow-y: auto;
}

.multiselect-option {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.multiselect-option:hover {
  background: #f7fafc;
}

.multiselect-option input[type="checkbox"] {
  margin-right: 8px;
}

.option-content {
  display: flex;
  align-items: center;
}

.option-label {
  font-size: 0.9rem;
  color: #4a5568;
}

.filter-button {
  width: 100%;
  padding: 10px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: background 0.2s;
}

.filter-button:hover {
  background: #5a67d8;
}

/* Content Area */
.content-area {
  flex: 1;
  background: #ffffff;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
}

/* Map Section */
.map-section {
  flex: 1;
  position: relative;
  z-index: 1;
}

/* Events Section */
.events-section {
  background: #f8f9fa;
  border-top: 1px solid #e2e8f0;
  padding: 1.5rem;
  max-height: 300px;
  overflow-y: auto;
}

.events-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.no-events {
  grid-column: 1 / -1;
  text-align: center;
  padding: 2rem 1rem;
  font-size: 0.9rem;
}

.event-card {
  background: white;
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid #e2e8f0;
  transition: all 0.2s;
  cursor: pointer;
}

.event-card:hover {
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.event-header {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.focus-button {
  background: #f8f9fa;
  color: #667eea;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  width: 28px;
  height: 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s;
  flex-shrink: 0;
  margin-left: auto;
}

.focus-button:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(102, 126, 234, 0.2);
}

.focus-button:active {
  transform: scale(0.95);
}

.event-emoji {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.event-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #2d3748;
  line-height: 1.3;
  flex: 1;
}

.event-description {
  color: #6b7280;
  font-size: 0.9rem;
  line-height: 1.4;
  margin: 0 0 0.75rem 0;
}

.event-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: #9ca3af;
}

.event-date {
  font-weight: 500;
}

.event-coords {
  font-family: monospace;
}

/* Responsive Design */
@media (max-width: 1024px) {
  .main-layout {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    border-right: none;
    border-bottom: 1px solid #e2e8f0;
  }
  
  .sidebar.collapsed {
    height: 60px;
    width: 100%;
  }
  
  .sidebar-toggle {
    right: 1rem;
    top: 1rem;
    transform: none;
  }
  
  .events-grid {
    grid-template-columns: 1fr;
  }
  
  .date-inputs {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .sidebar-section {
    padding: 1rem;
  }
  
  .events-section {
    padding: 1rem;
  }
  
  .event-card {
    padding: 0.75rem;
  }
}
</style>