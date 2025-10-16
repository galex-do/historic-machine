<template>
  <div class="tag_filter_panel" v-if="selectedTags.length > 0">
    <div class="tag_filter_header">
      <span class="filter_label">{{ t('filteredByTags') }}:</span>
      <button 
        class="clear_all_btn" 
        @click="$emit('clear-all-tags')"
        :title="t('clearAllTags')"
      >
        {{ t('clearAll') }}
      </button>
    </div>
    <div class="tag_chips_container">
      <div 
        v-for="tag in selectedTags" 
        :key="tag.id"
        class="tag_chip"
        :style="{ 
          backgroundColor: tag.color || '#6366f1',
          borderColor: tag.color || '#6366f1'
        }"
      >
        <span class="tag_name">{{ tag.name }}</span>
        <button 
          class="remove_tag_btn" 
          @click="$emit('remove-tag', tag.id)"
          :aria-label="`${t('remove')} ${tag.name}`"
        >
          ×
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'TagFilterPanel',
  props: {
    selectedTags: {
      type: Array,
      default: () => []
    }
  },
  emits: ['remove-tag', 'clear-all-tags'],
  setup() {
    const { t } = useLocale()
    
    return {
      t
    }
  }
}
</script>

<style scoped>
.tag_filter_panel {
  background: #f8f9fa;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  margin: 0 0 0.5rem 0;
}

.tag_filter_header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.filter_label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #4a5568;
}

.clear_all_btn {
  background: none;
  border: none;
  color: #dc3545;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: all 0.2s;
}

.clear_all_btn:hover {
  background: #fee;
  color: #c82333;
}

.tag_chips_container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag_chip {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-size: 0.75rem;
  font-weight: 500;
  color: white;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.tag_chip:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.tag_name {
  user-select: none;
}

.remove_tag_btn {
  background: rgba(255, 255, 255, 0.3);
  border: none;
  color: white;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 1rem;
  line-height: 1;
  transition: all 0.2s;
  padding: 0;
}

.remove_tag_btn:hover {
  background: rgba(255, 255, 255, 0.5);
  transform: scale(1.1);
}
</style>
