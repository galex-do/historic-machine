<template>
  <div class="filter-group">
    <div class="date-inputs">
      <div class="date-input-group">
        <label class="filter-label">From:</label>
        <input 
          type="text" 
          :value="dateFromDisplay" 
          @blur="$emit('date-from-changed', $event.target.value)"
          placeholder="DD.MM.YYYY" 
          class="date-input" 
        />
      </div>
      <div class="date-input-group">
        <label class="filter-label">To:</label>
        <input 
          type="text" 
          :value="dateToDisplay" 
          @blur="$emit('date-to-changed', $event.target.value)"
          placeholder="DD.MM.YYYY" 
          class="date-input" 
        />
      </div>
    </div>
  </div>
</template>

<script>
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
  emits: ['date-from-changed', 'date-to-changed']
}
</script>

<style scoped>
.date-inputs {
  display: flex;
  gap: 1rem;
}

.date-input-group {
  flex: 1;
}

.date-input {
  width: 100%;
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
</style>