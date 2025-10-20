<template>
  <div class="filter-field">
    <label class="filter-label">{{ t('tags') }}:</label>
    <div class="select-container">
      <div 
        ref="dropdownButton"
        class="select-dropdown" 
        @click="$emit('toggle-dropdown')" 
        :class="{ 'active': showDropdown, 'highlight': isHighlighted }"
      >
        <div class="select-display">
          <span class="selection-text">
            {{ getDisplayText() }}
          </span>
          <span v-if="selectedTags.length > 0" class="tag-count-badge">
            {{ selectedTags.length }}
          </span>
        </div>
        <span class="dropdown-arrow">{{ showDropdown ? '▲' : '▼' }}</span>
      </div>
      
      <div v-show="showDropdown" class="select-options">
        <!-- Search bar -->
        <div class="tag-search">
          <input 
            v-model="searchQuery"
            type="text"
            class="tag-search-input"
            :placeholder="t('searchTags')"
            @click.stop
          />
        </div>
        
        <!-- Clear all button -->
        <div v-if="selectedTags.length > 0" class="clear-all-section">
          <button 
            class="clear-all-btn"
            @click.stop="handleClearAll"
          >
            {{ t('clearAllTags') }}
          </button>
        </div>
        
        <!-- Tag options with checkboxes -->
        <div class="tags-list">
          <div 
            v-for="tag in filteredTags" 
            :key="tag.id" 
            class="tag-option"
            @click.stop="handleTagToggle(tag)"
          >
            <input 
              type="checkbox"
              class="tag-checkbox"
              :checked="isTagSelected(tag)"
              @click.stop="handleTagToggle(tag)"
            />
            <span 
              class="tag-color-swatch"
              :style="{ backgroundColor: tag.color }"
            ></span>
            <span class="tag-name">{{ tag.name }}</span>
          </div>
          <div v-if="filteredTags.length === 0" class="no-tags-found">
            {{ t('noTagsFound') }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useLocale } from '../../composables/useLocale'

export default {
  name: 'TagFilterDropdown',
  props: {
    selectedTags: {
      type: Array,
      default: () => []
    },
    showDropdown: {
      type: Boolean,
      default: false
    },
    allTags: {
      type: Array,
      default: () => []
    },
    isHighlighted: {
      type: Boolean,
      default: false
    }
  },
  emits: ['toggle-dropdown', 'tag-toggled', 'clear-all'],
  setup() {
    const { t } = useLocale()
    return { t }
  },
  data() {
    return {
      searchQuery: ''
    }
  },
  computed: {
    filteredTags() {
      if (!this.searchQuery.trim()) {
        return this.allTags
      }
      const query = this.searchQuery.toLowerCase()
      return this.allTags.filter(tag => 
        tag.name.toLowerCase().includes(query)
      )
    }
  },
  methods: {
    getDisplayText() {
      if (this.selectedTags.length === 0) {
        return this.t('allTags')
      }
      if (this.selectedTags.length === 1) {
        return this.selectedTags[0].name
      }
      // Show first tag name + count for multiple tags
      return `${this.selectedTags[0].name} +${this.selectedTags.length - 1}`
    },
    
    isTagSelected(tag) {
      return this.selectedTags.some(t => t.id === tag.id)
    },
    
    handleTagToggle(tag) {
      this.$emit('tag-toggled', tag)
    },
    
    handleClearAll() {
      this.$emit('clear-all')
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
            options.style.minWidth = `${Math.max(rect.width, 280)}px`
          }
        })
      }
    }
  },

  watch: {
    showDropdown: {
      handler(newVal) {
        this.positionDropdown()
        // Clear search when closing dropdown
        if (!newVal) {
          this.searchQuery = ''
        }
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
  min-width: 200px;
}

.select-dropdown:hover,
.select-dropdown.active {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.select-dropdown.highlight {
  border-color: #48bb78;
  box-shadow: 0 0 0 3px rgba(72, 187, 120, 0.2);
  animation: pulse-green 0.6s ease-in-out;
}

@keyframes pulse-green {
  0%, 100% {
    box-shadow: 0 0 0 3px rgba(72, 187, 120, 0.2);
  }
  50% {
    box-shadow: 0 0 0 6px rgba(72, 187, 120, 0.3);
  }
}

.select-display {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.selection-text {
  color: #2d3748;
  font-size: 0.9rem;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  background: #667eea;
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: 10px;
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
  min-width: 280px;
}

.tag-search {
  padding: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
}

.tag-search-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 0.875rem;
  outline: none;
  transition: border-color 0.2s;
}

.tag-search-input:focus {
  border-color: #667eea;
}

.clear-all-section {
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid #e2e8f0;
}

.clear-all-btn {
  width: 100%;
  padding: 0.5rem;
  background: #f7fafc;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-all-btn:hover {
  background: #edf2f7;
  border-color: #cbd5e0;
}

.tags-list {
  max-height: 300px;
  overflow-y: auto;
}

.tag-option {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid #f7fafc;
}

.tag-option:last-child {
  border-bottom: none;
}

.tag-option:hover {
  background-color: #f7fafc;
}

.tag-checkbox {
  width: 16px;
  height: 16px;
  cursor: pointer;
  flex-shrink: 0;
}

.tag-color-swatch {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  flex-shrink: 0;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.tag-name {
  font-size: 0.9rem;
  color: #2d3748;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-tags-found {
  padding: 1rem;
  text-align: center;
  color: #a0aec0;
  font-size: 0.875rem;
}
</style>
