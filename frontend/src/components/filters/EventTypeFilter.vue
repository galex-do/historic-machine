<template>
  <div class="filter-field">
    <label class="filter-label">Event Types:</label>
    <div class="multiselect-container">
      <div 
        class="multiselect-dropdown" 
        @click="$emit('toggle-dropdown')" 
        :class="{ 'active': showDropdown }"
      >
        <div class="multiselect-display">
          <span v-if="selectedLensTypes.length === 0" class="placeholder">
            Select event types...
          </span>
          <span v-else-if="selectedLensTypes.length === availableLensTypes.length" class="selection-text">
            All types selected
          </span>
          <div v-else class="selected-badges">
            <span v-for="type in selectedLensTypes" :key="type" :class="['selected-badge', type]">
              {{ getEventEmoji(type) }} {{ getLensLabel(type) }}
            </span>
          </div>
        </div>
        <span class="dropdown-arrow">{{ showDropdown ? '▲' : '▼' }}</span>
      </div>
      
      <div v-show="showDropdown" class="multiselect-options">
        <label v-for="lens in availableLensTypes" :key="lens.value" class="multiselect-option">
          <input 
            type="checkbox" 
            :value="lens.value" 
            :checked="selectedLensTypes.includes(lens.value)"
            @change="handleLensTypeChange" 
          />
          <span class="option-content">
            <span class="option-label">{{ lens.label }}</span>
          </span>
        </label>
      </div>
    </div>
  </div>
</template>

<script>
import { getEventEmoji, getLensLabel, getAvailableLensTypes } from '@/utils/event-utils.js'

export default {
  name: 'EventTypeFilter',
  props: {
    selectedLensTypes: {
      type: Array,
      default: () => []
    },
    showDropdown: {
      type: Boolean,
      default: false
    }
  },
  emits: ['toggle-dropdown', 'lens-types-changed'],
  computed: {
    availableLensTypes() {
      return getAvailableLensTypes()
    }
  },
  methods: {
    getEventEmoji,
    getLensLabel,
    handleLensTypeChange(event) {
      const value = event.target.value
      const isChecked = event.target.checked
      
      let newSelectedTypes = [...this.selectedLensTypes]
      
      if (isChecked) {
        if (!newSelectedTypes.includes(value)) {
          newSelectedTypes.push(value)
        }
      } else {
        newSelectedTypes = newSelectedTypes.filter(type => type !== value)
      }
      
      this.$emit('lens-types-changed', newSelectedTypes)
    }
  }
}
</script>

<style scoped>
.multiselect-container {
  position: relative;
}

.filter-field {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}

.filter-label {
  color: #4a5568;
  font-weight: 500;
  font-size: 0.9rem;
  white-space: nowrap;
  margin: 0;
}

.multiselect-dropdown {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 36px;
  padding: 0 0.75rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 180px;
}

.multiselect-dropdown:hover,
.multiselect-dropdown.active {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.multiselect-display {
  flex: 1;
}

.placeholder {
  color: #a0aec0;
  font-style: italic;
}

.selection-text {
  color: #2d3748;
  font-weight: 500;
}

.selected-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.selected-badge {
  background: #f0f4ff;
  color: #667eea;
  border: 1px solid #ddd6fe;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 500;
}

.dropdown-arrow {
  margin-left: 0.5rem;
  transition: transform 0.2s;
}

.multiselect-options {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  max-height: 200px;
  overflow-y: auto;
  margin-top: 0.25rem;
}

.multiselect-option {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s;
  color: #2d3748;
}

.multiselect-option:hover {
  background: #f8f9fa;
}

.multiselect-option input[type="checkbox"] {
  margin-right: 0.75rem;
  transform: scale(1.2);
  accent-color: #667eea;
}

.option-content {
  flex: 1;
}

.option-label {
  font-weight: 500;
  color: #2d3748;
}
</style>