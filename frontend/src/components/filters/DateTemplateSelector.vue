<template>
  <div class="template-selector-container">
    <!-- Template Group and Specific Period in horizontal layout -->
    <div class="template-controls">
      <!-- Historical Period Dropdown -->
      <div class="filter-field">
        <label class="filter-label">Historical Period:</label>
        <select 
          :value="selectedTemplateGroupId" 
          @change="$emit('template-group-changed', $event.target.value)" 
          class="filter-select"
        >
          <option value="">Default (1 AD - Today)</option>
          <option value="custom">Custom Date Range</option>
          <option v-for="group in templateGroups" :key="group.id" :value="group.id">
            {{ group.name }}
          </option>
        </select>
      </div>
      
      <!-- Specific Period Dropdown (appears to the right) -->
      <div v-if="selectedTemplateGroupId && selectedTemplateGroupId !== 'custom'" class="filter-field specific-period">
        <label class="filter-label">Specific Period:</label>
        <select 
          :value="selectedTemplateId" 
          @change="$emit('template-changed', $event.target.value)" 
          class="filter-select"
        >
          <option value="">Select specific period...</option>
          <option v-for="template in availableTemplates" :key="template.id" :value="template.id">
            {{ template.name }}
          </option>
        </select>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DateTemplateSelector',
  props: {
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
    }
  },
  emits: ['template-group-changed', 'template-changed']
}
</script>

<style scoped>
.template-selector-container {
  display: flex;
  align-items: center;
}

.template-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
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

.filter-select {
  height: 36px;
  padding: 0 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.9rem;
  transition: border-color 0.2s;
  min-width: 200px;
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

.specific-period .filter-select {
  min-width: 220px;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .template-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 0.75rem;
  }
  
  .filter-field {
    width: 100%;
  }
  
  .filter-select {
    min-width: auto;
    width: 100%;
  }
}
</style>