<template>
  <div class="filter-field">
    <label class="filter-label">Event Type:</label>
    <div class="select-container">
      <div 
        class="select-dropdown" 
        @click="$emit('toggle-dropdown')" 
        :class="{ 'active': showDropdown }"
      >
        <div class="select-display">
          <span class="selection-text">
            {{ getDisplayText() }}
          </span>
        </div>
        <span class="dropdown-arrow">{{ showDropdown ? '▲' : '▼' }}</span>
      </div>
      
      <div v-show="showDropdown" class="select-options">
        <div 
          v-for="option in filterOptions" 
          :key="option.value" 
          class="select-option"
          :class="{ 'selected': selectedLensType === option.value }"
          @click="handleOptionClick(option.value)"
        >
          <span class="option-content">
            <span class="option-label">{{ option.label }}</span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EventTypeFilter',
  props: {
    selectedLensType: {
      type: String,
      default: 'all'
    },
    showDropdown: {
      type: Boolean,
      default: false
    }
  },
  emits: ['toggle-dropdown', 'lens-type-changed'],
  data() {
    return {
      filterOptions: [
        { value: 'all', label: 'All Types' },
        { value: 'historic', label: '📜 Historic' },
        { value: 'political', label: '🏛️ Political' },
        { value: 'cultural', label: '🎭 Cultural' },
        { value: 'scientific', label: '🔬 Scientific' },
        { value: 'military', label: '⚔️ Military' },
        { value: 'religious', label: '⛪ Religious' }
      ]
    }
  },
  methods: {
    getDisplayText() {
      const option = this.filterOptions.find(opt => opt.value === this.selectedLensType)
      return option ? option.label : 'All Types'
    },
    
    handleOptionClick(value) {
      this.$emit('lens-type-changed', value)
      this.$emit('toggle-dropdown') // Close dropdown after selection
    },

    positionDropdown() {
      if (this.showDropdown) {
        this.$nextTick(() => {
          const dropdown = this.$el.querySelector('.select-dropdown')
          const options = this.$el.querySelector('.select-options')
          if (dropdown && options) {
            const rect = dropdown.getBoundingClientRect()
            options.style.top = `${rect.bottom + 2}px`
            options.style.left = `${rect.left}px`
            options.style.width = `${rect.width}px`
          }
        })
      }
    }
  },

  watch: {
    showDropdown: {
      handler() {
        this.positionDropdown()
      },
      immediate: true
    }
  },

  mounted() {
    this.positionDropdown()
    window.addEventListener('scroll', this.positionDropdown)
    window.addEventListener('resize', this.positionDropdown)
  },

  beforeUnmount() {
    window.removeEventListener('scroll', this.positionDropdown)
    window.removeEventListener('resize', this.positionDropdown)
  }
}
</script>

<style scoped>
.select-container {
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

.select-dropdown {
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

.select-dropdown:hover,
.select-dropdown.active {
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.select-display {
  flex: 1;
}

.selection-text {
  color: #2d3748;
  font-size: 0.9rem;
  font-weight: 500;
}

.dropdown-arrow {
  color: #a0aec0;
  font-size: 0.75rem;
  margin-left: 0.5rem;
  transition: transform 0.2s;
}

.select-dropdown.active .dropdown-arrow {
  transform: rotate(180deg);
}

.select-options {
  position: fixed;
  z-index: 100000;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  margin-top: 2px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  max-height: 300px;
  overflow-y: auto;
  min-width: 180px;
}

.select-option {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid #f7fafc;
}

.select-option:last-child {
  border-bottom: none;
}

.select-option:hover {
  background-color: #f7fafc;
}

.select-option.selected {
  background-color: #4f46e5;
  color: white;
}

.select-option.selected:hover {
  background-color: #5a67d8;
}

.option-content {
  display: flex;
  align-items: center;
  width: 100%;
}

.option-label {
  font-size: 0.9rem;
  font-weight: 500;
}
</style>