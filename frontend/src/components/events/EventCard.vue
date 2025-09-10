<template>
  <div class="event-card">
    <div class="event-header">
      <span class="event-emoji">{{ getEventEmoji(event.lens_type) }}</span>
      <h4 class="event-title">
        {{ event.name }}
        <a 
          v-if="event.source" 
          :href="event.source" 
          target="_blank" 
          rel="noopener noreferrer"
          class="inline-source-link" 
          title="View source"
        >
          🔗
        </a>
      </h4>
      <button 
        class="focus-button" 
        @click="$emit('focus-event', event)" 
        title="Focus on map"
      >
        ⌖
      </button>
    </div>
    <p class="event-description">{{ event.description }}</p>
    <div class="event-meta">
      <span class="event-date">{{ event.display_date || formatDate(event.event_date) }}</span>
      <span class="event-coords">{{ event.latitude.toFixed(2) }}, {{ event.longitude.toFixed(2) }}</span>
    </div>
    
    <!-- Tags Display -->
    <div v-if="event.tags && event.tags.length > 0" class="event-tags">
      <div 
        v-for="tag in visibleTags" 
        :key="tag.id"
        class="tag-chip"
        :style="{ backgroundColor: tag.color }"
        :title="tag.description"
      >
        {{ tag.name }}
      </div>
      <div 
        v-if="hasMoreTags"
        class="tag-chip more-tags"
        :title="`${hiddenTagsCount} more tags: ${hiddenTagNames}`"
      >
        +{{ hiddenTagsCount }}
      </div>
    </div>
  </div>
</template>

<script>
import { getEventEmoji } from '@/utils/event-utils.js'
import { formatDate } from '@/utils/date-utils.js'

export default {
  name: 'EventCard',
  props: {
    event: {
      type: Object,
      required: true
    }
  },
  emits: ['focus-event'],
  computed: {
    visibleTags() {
      if (!this.event.tags || this.event.tags.length === 0) return []
      return this.event.tags.slice(0, 3)
    },
    hasMoreTags() {
      return this.event.tags && this.event.tags.length > 3
    },
    hiddenTagsCount() {
      if (!this.event.tags || this.event.tags.length <= 3) return 0
      return this.event.tags.length - 3
    },
    hiddenTagNames() {
      if (!this.event.tags || this.event.tags.length <= 3) return ''
      return this.event.tags.slice(3).map(tag => tag.name).join(', ')
    }
  },
  methods: {
    getEventEmoji,
    formatDate
  }
}
</script>

<style scoped>
.event-card {
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
  border: 1px solid #e2e8f0;
}

.event-card:hover {
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.event-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.event-emoji {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.event-title {
  color: #2d3748;
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  flex: 1;
  line-height: 1.3;
}

.inline-source-link {
  margin-left: 0.5rem;
  color: #3b82f6;
  text-decoration: none;
  font-size: 0.875rem;
  transition: all 0.2s;
  font-weight: normal;
}

.inline-source-link:hover {
  color: #1d4ed8;
  transform: scale(1.1);
}

.focus-button {
  background: #f8f9fa;
  color: #667eea;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  width: 28px;
  height: 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s;
  flex-shrink: 0;
  margin-left: auto;
}

.focus-button:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(102, 126, 234, 0.2);
}

.focus-button:active {
  transform: scale(0.95);
}

.event-description {
  color: #4a5568;
  margin: 0 0 0.75rem 0;
  font-size: 0.9rem;
  line-height: 1.4;
}

.event-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: #718096;
  padding-top: 0.5rem;
  border-top: 1px solid #e2e8f0;
}

.event-date {
  font-weight: 500;
}

.event-coords {
  font-family: 'Courier New', monospace;
}

.event-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
  margin-top: 0.75rem;
}

.tag-chip {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: white;
  cursor: help;
  transition: all 0.2s;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.tag-chip:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.more-tags {
  background-color: #a0aec0 !important;
  color: white;
  font-style: italic;
}
</style>