<template>
  <aside class="sidebar" :class="{ 'collapsed': collapsed }">
    <!-- Sidebar Toggle Button -->
    <button class="sidebar-toggle" @click="$emit('toggle')">
      {{ collapsed ? '›' : '‹' }}
    </button>
    
    <!-- Filters Section -->
    <div class="sidebar-section" v-show="!collapsed">
      <h3 class="section-title">Filters</h3>
      
      <!-- Date Selection Mode -->
      <div class="filter-group">
        <label class="filter-label">Date Selection:</label>
        <select 
          :value="dateSelectionMode" 
          @change="$emit('date-mode-changed', $event.target.value)" 
          class="filter-select"
        >
          <option value="historic">Historic Periods</option>
          <option value="custom">Custom Date Range</option>
        </select>
      </div>
      
      <!-- Historic Date Templates -->
      <DateTemplateSelector
        v-if="dateSelectionMode === 'historic'"
        :template-groups="templateGroups"
        :available-templates="availableTemplates"
        :selected-template-group-id="selectedTemplateGroupId"
        :selected-template-id="selectedTemplateId"
        :selected-template="selectedTemplate"
        @template-group-changed="$emit('template-group-changed', $event)"
        @template-changed="$emit('template-changed', $event)"
      />
      
      <!-- Custom Date Range -->
      <CustomDateRange
        v-if="dateSelectionMode === 'custom'"
        :date-from-display="dateFromDisplay"
        :date-to-display="dateToDisplay"
        @date-from-changed="$emit('date-from-changed', $event)"
        @date-to-changed="$emit('date-to-changed', $event)"
      />
      
      <!-- Event Types -->
      <EventTypeFilter
        :selected-lens-types="selectedLensTypes"
        :show-dropdown="showLensDropdown"
        @toggle-dropdown="$emit('toggle-lens-dropdown')"
        @lens-types-changed="$emit('lens-types-changed', $event)"
      />
      
      <button @click="$emit('apply-filters')" class="filter-button">Apply Filters</button>
    </div>
  </aside>
</template>

<script>
import DateTemplateSelector from './DateTemplateSelector.vue'
import CustomDateRange from './CustomDateRange.vue'
import EventTypeFilter from './EventTypeFilter.vue'

export default {
  name: 'SidebarFilters',
  components: {
    DateTemplateSelector,
    CustomDateRange,
    EventTypeFilter
  },
  props: {
    collapsed: {
      type: Boolean,
      required: true
    },
    dateSelectionMode: {
      type: String,
      required: true
    },
    templateGroups: {
      type: Array,
      default: () => []
    },
    availableTemplates: {
      type: Array,
      default: () => []
    },
    selectedTemplateGroupId: {
      type: [String, Number],
      default: ''
    },
    selectedTemplateId: {
      type: [String, Number],
      default: ''
    },
    selectedTemplate: {
      type: Object,
      default: null
    },
    dateFromDisplay: {
      type: String,
      default: ''
    },
    dateToDisplay: {
      type: String,
      default: ''
    },
    selectedLensTypes: {
      type: Array,
      default: () => []
    },
    showLensDropdown: {
      type: Boolean,
      default: false
    }
  },
  emits: [
    'toggle',
    'date-mode-changed',
    'template-group-changed',
    'template-changed',
    'date-from-changed',
    'date-to-changed',
    'toggle-lens-dropdown',
    'lens-types-changed',
    'apply-filters'
  ]
}
</script>

<style scoped>
/* White sidebar styling */
.sidebar {
  background: #ffffff;
  color: #2d3748;
  padding: 0;
  position: relative;
  transition: all 0.3s ease;
  width: 350px;
  flex-shrink: 0;
  overflow: hidden;
  border-right: 1px solid #e2e8f0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.sidebar.collapsed {
  width: 60px;
}

.sidebar-toggle {
  position: absolute;
  top: 20px;
  right: 20px;
  background: #f8f9fa;
  color: #667eea;
  border: 1px solid #e2e8f0;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s;
  z-index: 10;
}

.sidebar-toggle:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
  transform: scale(1.1);
}

.sidebar-section {
  padding: 4rem 2rem 2rem;
}

.section-title {
  color: #2d3748;
  margin: 0 0 1.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.filter-group {
  margin-bottom: 1.5rem;
}

.filter-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.filter-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.filter-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.filter-select option {
  background: #ffffff;
  color: #2d3748;
}

.filter-button {
  width: 100%;
  padding: 0.75rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 1rem;
}

.filter-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}
</style>