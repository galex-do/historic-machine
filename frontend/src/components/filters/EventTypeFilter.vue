<template>
  <div class="filter-group">
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

.multiselect-dropdown {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  backdrop-filter: blur(10px);
}

.multiselect-dropdown:hover,
.multiselect-dropdown.active {
  background: rgba(255, 255, 255, 0.2);
}

.multiselect-display {
  flex: 1;
}

.placeholder {
  color: rgba(255, 255, 255, 0.6);
  font-style: italic;
}

.selection-text {
  color: white;
}

.selected-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.selected-badge {
  background: rgba(255, 255, 255, 0.2);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
  color: white;
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
  background: rgba(45, 55, 72, 0.95);
  border-radius: 6px;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 100;
  max-height: 200px;
  overflow-y: auto;
}

.multiselect-option {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s;
  color: white;
}

.multiselect-option:hover {
  background: rgba(255, 255, 255, 0.1);
}

.multiselect-option input[type="checkbox"] {
  margin-right: 0.75rem;
  transform: scale(1.2);
}

.option-content {
  flex: 1;
}

.option-label {
  font-weight: 500;
}
</style>