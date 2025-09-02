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
  border: none;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 0.9rem;
  backdrop-filter: blur(10px);
}

.date-input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.date-input:focus {
  outline: none;
  background: rgba(255, 255, 255, 0.2);
}
</style>