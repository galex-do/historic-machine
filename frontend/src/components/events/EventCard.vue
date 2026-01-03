<template>
  <div class="event-card-compact">
    <!-- Title line with icon and name -->
    <div class="event_title_line">
      <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
      <a 
        v-if="event.source"
        :href="event.source" 
        target="_blank" 
        rel="noopener noreferrer"
        class="event_name_link" 
        :title="`${event.name} - View source`"
      >
        {{ event.name }}
      </a>
      <span v-else class="event_name">{{ event.name }}</span>
      <button 
        v-if="mapFilterEnabled"
        class="highlight_btn" 
        @click="$emit('highlight-event', event)" 
        :title="t('highlightOnMap')"
      >
        📍
      </button>
      <button 
        v-else
        class="focus_btn" 
        @click="$emit('focus-event', event)" 
        :title="t('focusOnMap')"
      >
        ⌖
      </button>
    </div>
    
    <!-- Description (truncated, clickable to show full) -->
    <div 
      v-if="event.description" 
      class="event_description"
      @click="$emit('show-detail', event)"
      :title="t('clickToReadMore')"
    >
      {{ event.description }}
    </div>
    
    <!-- Tags row (always visible) -->
    <div v-if="event.tags && event.tags.length > 0" class="event_tags_row">
      <span
        v-for="(tag, index) in event.tags"
        :key="tag.id"
        class="event_tag"
        :style="{ color: tag.color || '#6366f1' }"
        @click.stop="$emit('tag-clicked', tag)"
        :title="tag.description"
      >#{{ tag.name }}</span>
    </div>
    
    <!-- Date below -->
    <div class="event_date_line">
      {{ formatEventDisplayDate(event.event_date, event.era) }}
    </div>
  </div>
</template>

<script>
import { getEventEmoji } from '@/utils/event-utils.js'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'EventCard',
  props: {
    event: {
      type: Object,
      required: true
    },
    mapFilterEnabled: {
      type: Boolean,
      default: false
    }
  },
  emits: ['focus-event', 'highlight-event', 'tag-clicked', 'show-detail'],
  setup() {
    const { formatEventDisplayDate, t } = useLocale()
    
    return {
      formatEventDisplayDate,
      t
    }
  },
  methods: {
    getEventEmoji
  }
}
</script>

<style scoped>
.event-card-compact {
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid #e2e8f0;
  transition: background-color 0.2s;
}

.event-card-compact:hover {
  background-color: #f8f9fa;
}

.event_title_line {
  display: flex;
  align-items: flex-start;
  gap: 0.25rem;
  font-size: 0.875rem;
  line-height: 1.4;
}

.event_icon {
  font-size: 1rem;
  flex-shrink: 0;
}

.event_name {
  font-weight: 600;
  color: #2d3748;
  flex: 1;
  min-width: 0;
}

.event_name_link {
  font-weight: 600;
  color: #3b82f6;
  text-decoration: none;
  transition: color 0.2s;
  flex: 1;
  min-width: 0;
}

.event_name_link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.event_description {
  font-size: 0.8rem;
  line-height: 1.4;
  color: #4a5568;
  margin-top: 0.25rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  cursor: pointer;
  transition: color 0.2s;
}

.event_description:hover {
  color: #3b82f6;
}

.event_tags_row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem 0.5rem;
  margin-top: 0.25rem;
  line-height: 1.2;
}

.event_tag {
  font-size: 0.7rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
  line-height: 1.2;
}

.event_tag:hover {
  opacity: 0.7;
}

.highlight_btn,
.focus_btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 0.9rem;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
  flex-shrink: 0;
  margin-left: auto;
}

.highlight_btn:hover {
  color: #f59e0b;
}

.focus_btn:hover {
  color: #3b82f6;
}

.highlight_btn:focus-visible,
.focus_btn:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
  border-radius: 2px;
}

.event_date_line {
  font-size: 0.7rem;
  color: #718096;
  margin-top: 0.25rem;
  font-weight: 500;
}
</style>
