<template>
  <div>
    <!-- Template Group Selection -->
    <div class="filter-group">
      <label class="filter-label">Historical Period:</label>
      <select 
        :value="selectedTemplateGroupId" 
        @change="$emit('template-group-changed', $event.target.value)" 
        class="filter-select"
      >
        <option value="">Default (1 AD - Today)</option>
        <option v-for="group in templateGroups" :key="group.id" :value="group.id">
          {{ group.name }}
        </option>
      </select>
    </div>
    
    <!-- Specific Template Selection -->
    <div v-if="selectedTemplateGroupId" class="filter-subgroup">
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
    
    <!-- Selected Period Info -->
    <div v-if="selectedTemplate" class="selected-period">
      <div class="period-info">
        <strong>{{ selectedTemplate.name }}</strong>
        <span class="period-dates">
          {{ selectedTemplate.start_display_date }} - {{ selectedTemplate.end_display_date }}
        </span>
        <p v-if="selectedTemplate.description" class="period-desc">
          {{ selectedTemplate.description }}
        </p>
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

.filter-subgroup {
  margin-top: 1rem;
  padding-left: 1rem;
  border-left: 2px solid #e2e8f0;
}

.selected-period {
  margin-top: 1rem;
  padding: 1rem;
  background: #f8f9fa;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
}

.period-info strong {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
  color: #2d3748;
}

.period-dates {
  display: block;
  font-size: 0.9rem;
  color: #667eea;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.period-desc {
  font-size: 0.85rem;
  color: #718096;
  margin: 0;
  line-height: 1.4;
}
</style>