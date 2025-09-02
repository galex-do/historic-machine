<template>
  <div class="date-control-bar">
    <div class="date-controls-container">
      <h3 class="controls-title">Filters</h3>
      
      <!-- Historical Period Template Selector -->
      <div class="template-section">
        <DateTemplateSelector
          :template-groups="templateGroups"
          :available-templates="availableTemplates"
          :selected-template-group-id="selectedTemplateGroupId"
          :selected-template-id="selectedTemplateId"
          :selected-template="selectedTemplate"
          @template-group-changed="$emit('template-group-changed', $event)"
          @template-changed="$emit('template-changed', $event)"
        />
      </div>
      
      <!-- Date Range Controls -->
      <div class="date-range-section">
        <div class="date-inputs">
          <div class="date-input-group">
            <label class="filter-label">From:</label>
            <input 
              ref="dateFromInput"
              type="text" 
              :value="dateFromDisplay" 
              @input="handleDateFromInput"
              @blur="$emit('date-from-changed', $event.target.value)"
              placeholder="500 BC, 1066 AD, 1776" 
              class="date-input" 
            />
          </div>
          
          <div class="date-input-group">
            <label class="filter-label">To:</label>
            <input 
              ref="dateToInput"
              type="text" 
              :value="dateToDisplay" 
              @input="handleDateToInput"
              @blur="$emit('date-to-changed', $event.target.value)"
              placeholder="500 BC, 1066 AD, 2025" 
              class="date-input" 
            />
          </div>
        </div>
        
        <!-- Centralized Step Controls -->
        <div class="step-controls-section">
          <div class="step-controls">
            <button @click="stepBothDates(-stepSize)" class="step-button" title="Step both dates backward">
              ‹‹
            </button>
            <button @click="stepBothDates(stepSize)" class="step-button" title="Step both dates forward">
              ››
            </button>
          </div>
          
          <div class="step-size-control">
            <label class="filter-label">Step:</label>
            <select v-model="stepSize" class="step-select">
              <option :value="1">1 year</option>
              <option :value="5">5 years</option>
              <option :value="10">10 years</option>
              <option :value="25">25 years</option>
              <option :value="50">50 years</option>
              <option :value="100">100 years</option>
            </select>
          </div>
        </div>
      </div>
      
      <!-- Event Type Filter -->
      <div class="event-type-section">
        <EventTypeFilter
          :selected-lens-types="selectedLensTypes"
          :show-dropdown="showLensDropdown"
          @toggle-dropdown="$emit('toggle-lens-dropdown')"
          @lens-types-changed="$emit('lens-types-changed', $event)"
        />
      </div>
      
      <!-- Apply Filters Button -->
      <div class="apply-section">
        <button @click="$emit('apply-filters')" class="apply-button">Apply</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { addYearsToHistoricalDate, parseHistoricalDate } from '@/utils/date-utils.js'
import DateTemplateSelector from './DateTemplateSelector.vue'
import EventTypeFilter from './EventTypeFilter.vue'

export default {
  name: 'DateControlBar',
  components: {
    DateTemplateSelector,
    EventTypeFilter
  },
  props: {
    dateFromDisplay: {
      type: String,
      default: ''
    },
    dateToDisplay: {
      type: String,
      default: ''
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
    'date-from-changed', 
    'date-to-changed', 
    'template-group-changed', 
    'template-changed', 
    'toggle-lens-dropdown', 
    'lens-types-changed',
    'apply-filters'
  ],
  setup(props, { emit }) {
    const stepSize = ref(10) // Default step size: 10 years
    
    // Handle real-time input validation
    const handleDateFromInput = (event) => {
      const value = event.target.value
      // Only update if it's a valid historical date or empty
      if (!value || parseHistoricalDate(value)) {
        // Let the parent handle the change through blur event
      }
    }
    
    const handleDateToInput = (event) => {
      const value = event.target.value
      // Only update if it's a valid historical date or empty
      if (!value || parseHistoricalDate(value)) {
        // Let the parent handle the change through blur event
      }
    }
    
    // Step both dates forward/backward together
    const stepBothDates = (years) => {
      const fromValue = props.dateFromDisplay
      const toValue = props.dateToDisplay
      
      if (fromValue) {
        const newFromDate = addYearsToHistoricalDate(fromValue, years)
        emit('date-from-changed', newFromDate)
      }
      
      if (toValue) {
        const newToDate = addYearsToHistoricalDate(toValue, years)
        emit('date-to-changed', newToDate)
      }
    }
    
    return {
      stepSize,
      handleDateFromInput,
      handleDateToInput,
      stepBothDates
    }
  }
}
</script>

<style scoped>
.date-control-bar {
  background: #ffffff;
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  padding: 1rem 1.5rem;
  margin-bottom: 0;
}

.date-controls-container {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.controls-title {
  color: #2d3748;
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  white-space: nowrap;
}

.date-range-section {
  display: flex;
  align-items: center;
  gap: 2rem;
  flex: 1;
}

.date-inputs {
  display: flex;
  gap: 1.5rem;
  flex: 1;
}

.date-input-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
}

.filter-label {
  color: #4a5568;
  font-weight: 500;
  font-size: 0.9rem;
  white-space: nowrap;
  margin: 0;
}

.date-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.9rem;
  transition: border-color 0.2s;
  min-width: 150px;
}

.date-input::placeholder {
  color: #a0aec0;
}

.date-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.template-section {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.event-type-section {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.apply-section {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.apply-button {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.apply-button:hover {
  background: #5a67d8;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(102, 126, 234, 0.3);
}

.step-controls-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.step-controls {
  display: flex;
  gap: 4px;
}

.step-button {
  width: 36px;
  height: 32px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #f8f9fa;
  color: #4a5568;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.step-button:hover {
  background: #e2e8f0;
  border-color: #cbd5e0;
  color: #2d3748;
}

.step-button:active {
  background: #cbd5e0;
}

.step-size-control {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.step-select {
  padding: 0.5rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.9rem;
  cursor: pointer;
  min-width: 90px;
}

.step-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

/* Responsive design */
@media (max-width: 1200px) {
  .date-controls-container {
    gap: 1rem;
  }
  
  .date-range-section {
    flex-direction: column;
    gap: 0.75rem;
  }
}

@media (max-width: 768px) {
  .date-control-bar {
    padding: 1rem;
  }
  
  .date-controls-container {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .template-section,
  .event-type-section,
  .apply-section {
    width: 100%;
  }
  
  .date-range-section {
    flex-direction: column;
    gap: 1rem;
  }
  
  .date-inputs {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>