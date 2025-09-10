<template>
  <div class="filter-field">
    <label class="filter-label">Dataset:</label>
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
          v-for="option in datasetOptions" 
          :key="option.value" 
          class="select-option"
          :class="{ 'selected': selectedDataset === option.value }"
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
  name: 'DatasetFilter',
  props: {
    selectedDataset: {
      type: String,
      default: 'all'
    },
    showDropdown: {
      type: Boolean,
      default: false
    },
    datasets: {
      type: Array,
      default: () => []
    }
  },
  emits: ['toggle-dropdown', 'dataset-changed'],
  computed: {
    datasetOptions() {
      const options = [
        { value: 'all', label: '📊 All Datasets' },
        { value: 'none', label: '🚫 No Dataset' }
      ]
      
      // Add dataset options
      if (this.datasets && this.datasets.length > 0) {
        this.datasets.forEach(dataset => {
          options.push({
            value: dataset.id.toString(),
            label: `📁 ${dataset.filename}`
          })
        })
      }
      
      return options
    }
  },
  methods: {
    getDisplayText() {
      const option = this.datasetOptions.find(opt => opt.value === this.selectedDataset)
      return option ? option.label : '📊 All Datasets'
    },
    
    handleOptionClick(value) {
      this.$emit('dataset-changed', value)
      this.$emit('toggle-dropdown') // Close dropdown after selection
    }
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
  min-width: 200px;
  width: auto;
}

.select-dropdown:hover,
.select-dropdown.active {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
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
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 100000;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  margin-top: 2px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  max-height: 300px;
  overflow-y: auto;
  min-width: 280px;
  width: max-content;
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
  background-color: #667eea;
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
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 220px;
}
</style>