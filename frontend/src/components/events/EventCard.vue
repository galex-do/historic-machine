<template>
  <div class="event-card">
    <div class="event-header">
      <h4 class="event-title">
        <span class="event-emoji">{{ getEventEmoji(event.lens_type) }}</span>
        <a 
          v-if="event.source"
          :href="event.source" 
          target="_blank" 
          rel="noopener noreferrer"
          class="event-title-link" 
          :title="`${event.name} - View source`"
        >
          {{ event.name }}
        </a>
        <span v-else class="event-title-text">{{ event.name }}</span>
      </h4>
      <button 
        class="focus-button" 
        @click="$emit('focus-event', event)" 
        title="Focus on map"
      >
        ⌖
      </button>
    </div>
    <div class="event-description-container" v-if="event.description">
      <button 
        v-if="!description_expanded"
        class="expand-trigger" 
        @click="toggle_description"
        :aria-label="t('showMore')"
        :title="t('showMore')"
        :aria-expanded="false"
      >...</button>
      <p class="event-description expanded" v-if="description_expanded">
        {{ event.description }}
        <button 
          class="collapse-trigger" 
          @click="toggle_description"
          :aria-label="t('showLess')"
          :title="t('showLess')"
          :aria-expanded="true"
        > ▲</button>
      </p>
    </div>
    <div class="event-meta">
      <span class="event-date">{{ formatEventDisplayDate(event.event_date, event.era) }}</span>
      <span class="event-coords">{{ event.latitude.toFixed(2) }}, {{ event.longitude.toFixed(2) }}</span>
    </div>
    
    <!-- Tags Display -->
    <div v-if="event.tags && event.tags.length > 0" class="event-tags">
      <div 
        v-for="tag in visibleTags" 
        :key="tag.id"
        class="tag-chip clickable-tag"
        :style="{ backgroundColor: tag.color }"
        :title="tag.description"
        @click.stop="$emit('tag-clicked', tag)"
      >
        {{ tag.name }}
      </div>
      <div 
        v-if="hasMoreTags && !show_all_tags"
        class="tag-chip more-tags clickable-more-tags"
        :title="`Click to show ${hiddenTagsCount} more tags: ${hiddenTagNames}`"
        @click.stop="toggle_tags_display"
      >
        +{{ hiddenTagsCount }}
      </div>
      <div 
        v-if="show_all_tags && hasMoreTags"
        class="tag-chip show-less-tags clickable-more-tags"
        title="Click to show less tags"
        @click.stop="toggle_tags_display"
      >
        Show less
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { getEventEmoji } from '@/utils/event-utils.js'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'EventCard',
  props: {
    event: {
      type: Object,
      required: true
    }
  },
  emits: ['focus-event', 'tag-clicked'],
  setup() {
    const { formatEventDisplayDate, t } = useLocale()
    const show_all_tags = ref(false)
    const description_expanded = ref(false)
    
    const toggle_tags_display = () => {
      show_all_tags.value = !show_all_tags.value
    }
    
    const toggle_description = () => {
      description_expanded.value = !description_expanded.value
    }
    
    return {
      formatEventDisplayDate,
      t,
      show_all_tags,
      toggle_tags_display,
      description_expanded,
      toggle_description
    }
  },
  computed: {
    visibleTags() {
      if (!this.event.tags || this.event.tags.length === 0) return []
      if (this.show_all_tags) {
        return this.event.tags
      }
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
    getEventEmoji
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

.event-title {
  color: #2d3748;
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  flex: 1;
  line-height: 1.3;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.event-emoji {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.event-title-link {
  color: #3b82f6;
  text-decoration: none;
  transition: all 0.2s;
  font-weight: 600;
}

.event-title-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.event-title-text {
  font-weight: 600;
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

.event-description-container {
  margin: 0 0 0.75rem 0;
}

.event-description {
  color: #4a5568;
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.4;
}

.event-description.expanded {
  margin: 0;
}

.expand-trigger {
  background: none;
  border: none;
  color: #3b82f6;
  cursor: pointer;
  font-weight: 600;
  font-size: 1.2rem;
  padding: 0 0.25rem;
  transition: color 0.2s;
  display: inline;
  line-height: 1;
}

.expand-trigger:hover {
  color: #1d4ed8;
}

.expand-trigger:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
  border-radius: 2px;
}

.collapse-trigger {
  background: none;
  border: none;
  color: #64748b;
  cursor: pointer;
  font-size: 0.875rem;
  padding: 0 0.25rem;
  transition: color 0.2s;
  margin-left: 0.25rem;
  display: inline;
  vertical-align: baseline;
}

.collapse-trigger:hover {
  color: #3b82f6;
}

.collapse-trigger:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
  border-radius: 2px;
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

.clickable-tag {
  cursor: pointer;
  transition: all 0.2s;
}

.clickable-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.25);
  filter: brightness(1.1);
}

.clickable-tag:active {
  transform: translateY(0);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
}

.more-tags {
  background-color: #a0aec0 !important;
  color: white;
  font-style: italic;
}

.clickable-more-tags {
  cursor: pointer;
  transition: all 0.2s;
}

.clickable-more-tags:hover {
  background-color: #718096 !important;
  transform: translateY(-2px);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.25);
}

.clickable-more-tags:active {
  transform: translateY(0);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
}

.show-less-tags {
  background-color: #a0aec0 !important;
  color: white;
  font-style: italic;
}
</style>