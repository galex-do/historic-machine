<template>
  <div class="tag_filter_panel">
    <!-- Tag Search Input -->
    <div class="tag_search_container" ref="searchContainerRef">
      <input
        v-model="searchQuery"
        type="text"
        class="tag_search_input"
        :placeholder="t('searchTags')"
        @focus="handleFocus"
        @blur="handleBlur"
      />
    </div>
    
    <!-- Tag Suggestions Dropdown (Teleported to body) -->
    <Teleport to="body">
      <div v-if="showSuggestions && filteredAvailableTags.length > 0" class="tag_suggestions" :style="dropdownPosition">
        <div
          v-for="tag in filteredAvailableTags"
          :key="tag.id"
          class="tag_suggestion_item"
          @mousedown.prevent="addTag(tag)"
          :style="{ 
            borderLeftColor: tag.color || '#6366f1'
          }"
        >
          <span class="tag_suggestion_name">{{ tag.name }}</span>
          <span class="tag_suggestion_count">({{ getTagEventCount(tag.id) }})</span>
        </div>
      </div>
      
      <div v-if="showSuggestions && searchQuery && filteredAvailableTags.length === 0" class="no_suggestions" :style="dropdownPosition">
        {{ t('noTagsFound') }}
      </div>
    </Teleport>

    <!-- Selected Tags Section -->
    <div v-if="selectedTags.length > 0" class="selected_tags_section">
      <div class="tag_filter_header">
        <span class="filter_label">{{ t('filteredByTags') }}:</span>
        <div class="header_actions">
          <button 
            class="follow_toggle_btn"
            :class="{ 'active': followEnabled }"
            @click="$emit('toggle-follow')"
            :title="followEnabled ? t('disableNarrativeFlow') : t('enableNarrativeFlow')"
          >
            {{ followEnabled ? '🔗' : '○' }} {{ t('followEvents') }}
          </button>
          <button 
            class="clear_all_btn" 
            @click="$emit('clear-all-tags')"
            :title="t('clearAllTags')"
          >
            {{ t('clearAll') }}
          </button>
        </div>
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
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'TagFilterPanel',
  props: {
    selectedTags: {
      type: Array,
      default: () => []
    },
    availableTags: {
      type: Array,
      default: () => []
    },
    followEnabled: {
      type: Boolean,
      default: false
    }
  },
  emits: ['remove-tag', 'clear-all-tags', 'toggle-follow', 'add-tag'],
  setup(props, { emit }) {
    const { t } = useLocale()
    
    const searchQuery = ref('')
    const showSuggestions = ref(false)
    const searchContainerRef = ref(null)
    const dropdownPosition = ref({})
    
    // Filter available tags based on search query and exclude already selected tags
    const filteredAvailableTags = computed(() => {
      const query = searchQuery.value.toLowerCase().trim()
      const selectedTagIds = new Set(props.selectedTags.map(tag => tag.id))
      
      return props.availableTags
        .filter(tag => {
          // Exclude already selected tags
          if (selectedTagIds.has(tag.id)) {
            return false
          }
          
          // If no search query, show all available tags
          if (!query) {
            return true
          }
          
          // Filter by search query
          return tag.name.toLowerCase().includes(query)
        })
        .slice(0, 10) // Limit to 10 suggestions for performance
    })
    
    const calculateDropdownPosition = () => {
      if (searchContainerRef.value) {
        const rect = searchContainerRef.value.getBoundingClientRect()
        dropdownPosition.value = {
          position: 'fixed',
          top: `${rect.bottom + 4}px`,
          left: `${rect.left}px`,
          width: `${rect.width}px`
        }
      }
    }
    
    const handleFocus = () => {
      showSuggestions.value = true
      calculateDropdownPosition()
    }
    
    const addTag = (tag) => {
      emit('add-tag', tag)
      searchQuery.value = '' // Clear search after adding
      showSuggestions.value = false
    }
    
    const handleBlur = () => {
      // Delay hiding to allow click events to fire
      setTimeout(() => {
        showSuggestions.value = false
      }, 200)
    }
    
    // Count how many events have this tag (would need events prop for accuracy)
    // For now, just show that it's available
    const getTagEventCount = (tagId) => {
      // This could be enhanced to show actual count if we pass events
      return '✓'
    }
    
    return {
      t,
      searchQuery,
      showSuggestions,
      searchContainerRef,
      dropdownPosition,
      filteredAvailableTags,
      addTag,
      handleFocus,
      handleBlur,
      getTagEventCount
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

.tag_search_container {
  position: relative;
  margin-bottom: 0.75rem;
}

.tag_search_input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  font-size: 0.875rem;
  outline: none;
  transition: all 0.2s;
  background: white;
}

.tag_search_input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.tag_search_input::placeholder {
  color: #94a3b8;
}

.tag_suggestions {
  background: white;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-height: 300px;
  overflow-y: auto;
  z-index: 10000;
}

.tag_suggestion_item {
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-left: 3px solid;
  transition: all 0.2s;
}

.tag_suggestion_item:hover {
  background: #f1f5f9;
}

.tag_suggestion_name {
  font-size: 0.875rem;
  color: #334155;
  font-weight: 500;
}

.tag_suggestion_count {
  font-size: 0.75rem;
  color: #64748b;
  margin-left: 0.5rem;
}

.no_suggestions {
  padding: 0.75rem;
  text-align: center;
  color: #94a3b8;
  font-size: 0.875rem;
  background: white;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  z-index: 10000;
}

.selected_tags_section {
  border-top: 1px solid #e2e8f0;
  padding-top: 0.75rem;
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

.header_actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.follow_toggle_btn {
  background: #e2e8f0;
  border: 1px solid #cbd5e0;
  color: #4a5568;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  padding: 0.25rem 0.65rem;
  border-radius: 4px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.follow_toggle_btn:hover {
  background: #cbd5e0;
  border-color: #a0aec0;
}

.follow_toggle_btn.active {
  background: #3b82f6;
  border-color: #2563eb;
  color: white;
}

.follow_toggle_btn.active:hover {
  background: #2563eb;
  border-color: #1d4ed8;
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
