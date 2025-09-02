<template>
  <div class="filter-group">
    <div class="date-inputs">
      <div class="date-input-group">
        <label class="filter-label">From:</label>
        <div class="date-input-with-controls">
          <input 
            ref="dateFromInput"
            type="text" 
            :value="dateFromDisplay" 
            @input="handleDateFromInput"
            @blur="$emit('date-from-changed', $event.target.value)"
            placeholder="500 BC, 1066 AD, 1776" 
            class="date-input" 
          />
          <div class="step-controls">
            <button @click="stepDateFrom(-stepSize)" class="step-button" title="Step backward">
              ‹
            </button>
            <button @click="stepDateFrom(stepSize)" class="step-button" title="Step forward">
              ›
            </button>
          </div>
        </div>
      </div>
      <div class="date-input-group">
        <label class="filter-label">To:</label>
        <div class="date-input-with-controls">
          <input 
            ref="dateToInput"
            type="text" 
            :value="dateToDisplay" 
            @input="handleDateToInput"
            @blur="$emit('date-to-changed', $event.target.value)"
            placeholder="500 BC, 1066 AD, 2025" 
            class="date-input" 
          />
          <div class="step-controls">
            <button @click="stepDateTo(-stepSize)" class="step-button" title="Step backward">
              ‹
            </button>
            <button @click="stepDateTo(stepSize)" class="step-button" title="Step forward">
              ›
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Step Size Control -->
    <div class="step-size-control">
      <label class="filter-label">Step Size:</label>
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
</template>

<script>
import { ref } from 'vue'
import { addYearsToHistoricalDate, parseHistoricalDate } from '@/utils/date-utils.js'

export default {
  name: 'CustomDateRange',
  props: {
    dateFromDisplay: {
      type: String,
      default: ''
    },
    dateToDisplay: {
      type: String,
      default: ''
    }
  },
  emits: ['date-from-changed', 'date-to-changed'],
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
    
    // Step date from forward/backward
    const stepDateFrom = (years) => {
      const currentValue = props.dateFromDisplay
      if (!currentValue) return
      
      const newDate = addYearsToHistoricalDate(currentValue, years)
      emit('date-from-changed', newDate)
    }
    
    // Step date to forward/backward
    const stepDateTo = (years) => {
      const currentValue = props.dateToDisplay
      if (!currentValue) return
      
      const newDate = addYearsToHistoricalDate(currentValue, years)
      emit('date-to-changed', newDate)
    }
    
    return {
      stepSize,
      handleDateFromInput,
      handleDateToInput,
      stepDateFrom,
      stepDateTo
    }
  }
}
</script>

<style scoped>
.date-inputs {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.date-input-group {
  flex: 1;
}

.date-input-with-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.date-input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.date-input::placeholder {
  color: #a0aec0;
}

.date-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.step-controls {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.step-button {
  width: 24px;
  height: 20px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  background: #f8f9fa;
  color: #4a5568;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.step-button:hover {
  background: #e2e8f0;
  border-color: #cbd5e0;
}

.step-button:active {
  background: #cbd5e0;
}

.step-size-control {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.step-select {
  padding: 0.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.8rem;
  cursor: pointer;
}

.step-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}
</style>