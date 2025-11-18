<template>
  <div class="event-card-compact">
    <!-- Single line with all event info -->
    <span class="event_line_text">
      <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
      {{ ' ' }}
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
      <template v-if="event.description">
        {{ ' — ' }}{{ event.description }}
      </template>
      <template v-if="event.tags && event.tags.length > 0">
        {{ ' ' }}
        <span
          v-for="(tag, index) in event.tags"
          :key="tag.id"
          class="event_tag"
          :style="{ color: tag.color || '#6366f1' }"
          @click.stop="$emit('tag-clicked', tag)"
          :title="tag.description"
        >#{{ tag.name }}{{ index < event.tags.length - 1 ? ' ' : '' }}</span>
      </template>
      {{ ' ' }}
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
    </span>
    
    <!-- Date below the line -->
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
  emits: ['focus-event', 'highlight-event', 'tag-clicked'],
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

.event_line_text {
  font-size: 0.875rem;
  line-height: 1.5;
  color: #2d3748;
  display: block;
  word-wrap: break-word;
}

.event_icon {
  font-size: 1rem;
}

.event_name {
  font-weight: 600;
  color: #2d3748;
}

.event_name_link {
  font-weight: 600;
  color: #3b82f6;
  text-decoration: none;
  transition: color 0.2s;
}

.event_name_link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.event_tag {
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}

.event_tag:hover {
  opacity: 0.7;
}

.highlight_btn,
.focus_btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 1rem;
  cursor: pointer;
  padding: 0 0.25rem;
  transition: color 0.2s;
  vertical-align: middle;
  display: inline;
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
  font-size: 0.75rem;
  color: #718096;
  margin-top: 0.25rem;
  font-weight: 500;
}
</style>
