<template>
  <div v-if="isOpen" class="timeline_modal_overlay" @click.self="closeModal">
    <div class="timeline_modal">
      <!-- Modal Header -->
      <div class="timeline_modal_header">
        <h2 class="timeline_modal_title">{{ t('historicalEvents') }}</h2>
        <button class="close_btn" @click="closeModal" :title="t('close')">×</button>
      </div>

      <!-- Modal Content -->
      <div class="timeline_modal_content">
        <div v-if="groupedEvents.length === 0" class="no_events_message">
          {{ t('noEventsInTimeline') }}
        </div>

        <div v-else class="timeline_container">
          <div v-for="group in groupedEvents" :key="group.date" class="timeline_date_group">
            <!-- Single event on this date: everything on one line -->
            <div v-if="group.events.length === 1" class="timeline_single_event_line">
              <div class="timeline_bullet"></div>
              <span class="timeline_single_text">
                <span class="timeline_date_inline">{{ group.formattedDate }}</span>
                {{ ' ' }}
                <span class="event_icon">{{ getEventEmoji(group.events[0].lens_type) }}</span>
                {{ ' ' }}
                <span class="event_name">{{ group.events[0].name }}</span>
                <template v-if="group.events[0].description">
                  {{ ' — ' }}{{ group.events[0].description }}
                </template>
                <template v-if="group.events[0].tags && group.events[0].tags.length > 0">
                  {{ ' ' }}
                  <span
                    v-for="tag in group.events[0].tags"
                    :key="tag.id"
                    class="event_tag_compact"
                    :style="{ color: tag.color || '#6366f1' }"
                  >{{ tag.name }}</span>
                </template>
              </span>
            </div>

            <!-- Multiple events on this date: separate date header -->
            <template v-else>
              <!-- Date Header -->
              <div class="timeline_date_header">
                <div class="timeline_bullet"></div>
                <div class="timeline_date">{{ group.formattedDate }}</div>
              </div>

              <!-- Events for this date -->
              <div class="timeline_events_list">
                <div 
                  v-for="event in group.events" 
                  :key="event.id"
                  class="timeline_event_line"
                >
                  <!-- Event line: Icon + Name -->
                  <div class="event_main_line">
                    <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
                    <span class="event_name">{{ event.name }}</span>
                  </div>

                  <!-- Description on separate line -->
                  <div 
                    v-if="event.description"
                    class="event_description_block"
                    :class="{ 'expanded': expandedEvents.has(event.id) }"
                  >
                    {{ event.description }}
                    <button 
                      v-if="event.description.length > 80"
                      class="expand_btn"
                      @click="toggleEventExpand(event.id)"
                    >
                      {{ expandedEvents.has(event.id) ? '▲' : '▼' }}
                    </button>
                  </div>

                  <!-- Event Tags - compact inline format -->
                  <div v-if="event.tags && event.tags.length > 0" class="event_tags_inline">
                    <span
                      v-for="tag in event.tags"
                      :key="tag.id"
                      class="event_tag_compact"
                      :style="{ color: tag.color || '#6366f1' }"
                    >
                      {{ tag.name }}
                    </span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useLocale } from '@/composables/useLocale.js'
import { getEventEmoji } from '@/utils/event-utils.js'

export default {
  name: 'TimelineModal',
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    events: {
      type: Array,
      default: () => []
    }
  },
  emits: ['close'],
  setup(props, { emit }) {
    const { t, formatEventDisplayDate } = useLocale()
    const expandedEvents = ref(new Set())
    const previouslyFocusedElement = ref(null)

    // Helper function for chronological sorting with BC/AD support
    // Matches implementation in useEvents.js - parses dates manually to avoid timezone issues
    const getChronologicalValue = (dateString, era) => {
      let year, month, day
      
      if (dateString.startsWith('-')) {
        // Negative year format: "-3501-01-01T00:00:00Z"
        const parts = dateString.substring(1).split('T')[0].split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10) - 1  // Month is 0-indexed
        day = parseInt(parts[2], 10)
      } else {
        // Positive year format: "1453-05-29T00:00:00Z"
        // Parse manually to avoid timezone shifts from Date constructor
        const parts = dateString.split('T')[0].split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10) - 1  // Month is 0-indexed
        day = parseInt(parts[2], 10)
      }
      
      if (era === 'BC') {
        // For BC: larger year number = older (3000 BC < 1000 BC)
        return -(year - (month / 12) - (day / 365))
      } else {
        // For AD: normal positive years
        return year + (month / 12) + (day / 365)
      }
    }

    // Group events by date (day precision)
    const groupedEvents = computed(() => {
      if (!props.events || props.events.length === 0) {
        return []
      }

      // Sort events chronologically using BC/AD aware sorting
      const sortedEvents = [...props.events].sort((a, b) => {
        const aValue = getChronologicalValue(a.event_date, a.era)
        const bValue = getChronologicalValue(b.event_date, b.era)
        return aValue - bValue // Chronological order: older events first
      })

      // Group by date
      const groups = {}
      sortedEvents.forEach(event => {
        // Extract date string (YYYY-MM-DD)
        const eventDate = event.event_date.split('T')[0]
        
        if (!groups[eventDate]) {
          groups[eventDate] = {
            date: eventDate,
            formattedDate: formatEventDisplayDate(event.event_date, event.era),
            events: []
          }
        }
        groups[eventDate].events.push(event)
      })

      return Object.values(groups)
    })

    const toggleEventExpand = (eventId) => {
      if (expandedEvents.value.has(eventId)) {
        expandedEvents.value.delete(eventId)
      } else {
        expandedEvents.value.add(eventId)
      }
      // Force reactivity
      expandedEvents.value = new Set(expandedEvents.value)
    }

    const closeModal = () => {
      emit('close')
      // Reset expanded states when closing
      expandedEvents.value.clear()
      // Restore focus to the element that opened the modal
      if (previouslyFocusedElement.value) {
        previouslyFocusedElement.value.focus()
        previouslyFocusedElement.value = null
      }
    }

    // Handle ESC key to close modal
    const handleEscape = (event) => {
      if (event.key === 'Escape' && props.isOpen) {
        closeModal()
      }
    }

    // Watch for modal open/close to manage focus and keyboard listeners
    watch(() => props.isOpen, (newValue) => {
      if (newValue) {
        // Store currently focused element
        previouslyFocusedElement.value = document.activeElement
        // Add ESC key listener
        document.addEventListener('keydown', handleEscape)
      } else {
        // Remove ESC key listener
        document.removeEventListener('keydown', handleEscape)
      }
    })

    // Cleanup on component unmount
    onUnmounted(() => {
      document.removeEventListener('keydown', handleEscape)
    })

    return {
      t,
      groupedEvents,
      expandedEvents,
      toggleEventExpand,
      closeModal,
      getEventEmoji
    }
  }
}
</script>

<style scoped>
.timeline_modal_overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100000;
  padding: 2rem;
}

.timeline_modal {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  max-width: 900px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.timeline_modal_header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e2e8f0;
}

.timeline_modal_title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1e293b;
}

.close_btn {
  background: none;
  border: none;
  font-size: 2rem;
  color: #64748b;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
  line-height: 1;
  padding: 0;
}

.close_btn:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.timeline_modal_content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 1.5rem;
  font-size: 0.875rem;
}

.no_events_message {
  text-align: center;
  padding: 2rem;
  color: #94a3b8;
  font-size: 0.875rem;
}

.timeline_container {
  position: relative;
  padding-left: 0.5rem;
}

/* Thin vertical line */
.timeline_container::before {
  content: '';
  position: absolute;
  left: 6px;
  top: 0;
  bottom: 0;
  width: 1px;
  background: #cbd5e1;
  z-index: 1;
}

.timeline_date_group {
  position: relative;
  margin-bottom: 0.75rem;
}

.timeline_date_header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
  position: relative;
  z-index: 2;
}

.timeline_bullet {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: white;
  border: 2px solid #3b82f6;
  flex-shrink: 0;
}

.timeline_date {
  font-weight: 600;
  font-size: 0.875rem;
  color: #1e293b;
}

.timeline_single_event_line {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  padding: 0.25rem 0;
  line-height: 1.4;
  position: relative;
  z-index: 2;
}

.timeline_single_text {
  font-size: 0.875rem;
  color: #475569;
  line-height: 1.5;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.timeline_single_text .timeline_date_inline {
  font-weight: 600;
  color: #1e293b;
}

.timeline_single_text .event_name {
  font-weight: 600;
  color: #1e293b;
}

.timeline_events_list {
  margin-left: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.timeline_event_line {
  padding: 0.25rem 0;
  border-left: 1px solid transparent;
  padding-left: 0.5rem;
  transition: border-color 0.2s;
}

.timeline_event_line:hover {
  border-left-color: #e2e8f0;
}

.event_main_line {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  line-height: 1.4;
}

.event_icon {
  font-size: 1rem;
  flex-shrink: 0;
}

.event_name {
  font-weight: 600;
  font-size: 0.875rem;
  color: #1e293b;
  flex-shrink: 0;
}

.event_description_inline {
  color: #475569;
  font-size: 0.875rem;
  font-weight: normal;
}

.event_description_block {
  margin-top: 0.125rem;
  margin-left: 1.5rem;
  color: #475569;
  font-size: 0.875rem;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.event_description_block:not(.expanded) {
  -webkit-line-clamp: 2;
}

.event_description_block.expanded {
  -webkit-line-clamp: unset;
}

.expand_btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 0.75rem;
  cursor: pointer;
  padding: 0 0.25rem;
  margin-left: 0.25rem;
  vertical-align: middle;
  font-weight: bold;
}

.expand_btn:hover {
  color: #3b82f6;
}

.event_tags_inline {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.125rem;
  margin-left: 1.5rem;
  font-size: 0.75rem;
}

.event_tag_compact {
  font-weight: 500;
  opacity: 0.8;
}

.event_tag_compact::before {
  content: '#';
  opacity: 0.6;
}

/* Scrollbar Styling */
.timeline_modal_content::-webkit-scrollbar {
  width: 8px;
}

.timeline_modal_content::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.timeline_modal_content::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

.timeline_modal_content::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
