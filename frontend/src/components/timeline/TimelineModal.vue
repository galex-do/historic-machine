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
            <!-- Date Bullet -->
            <div class="timeline_bullet_wrapper">
              <div class="timeline_bullet"></div>
              <div class="timeline_date">{{ group.formattedDate }}</div>
            </div>

            <!-- Events for this date -->
            <div class="timeline_events_wrapper">
              <div 
                v-for="event in group.events" 
                :key="event.id"
                class="timeline_event_card"
              >
                <!-- Event Header: Icon, Name -->
                <div class="timeline_event_header">
                  <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
                  <h3 class="event_name">{{ event.name }}</h3>
                </div>

                <!-- Event Description -->
                <div class="event_description_wrapper">
                  <p 
                    class="event_description"
                    :class="{ 'expanded': expandedEvents.has(event.id) }"
                  >
                    {{ event.description }}
                  </p>
                  <button 
                    v-if="event.description && event.description.length > 150"
                    class="expand_btn"
                    @click="toggleEventExpand(event.id)"
                  >
                    {{ expandedEvents.has(event.id) ? t('showLess') : t('showMore') }}
                  </button>
                </div>

                <!-- Event Tags -->
                <div v-if="event.tags && event.tags.length > 0" class="event_tags">
                  <div
                    v-for="tag in event.tags"
                    :key="tag.id"
                    class="event_tag"
                    :style="{ 
                      backgroundColor: tag.color || '#6366f1',
                      borderColor: tag.color || '#6366f1'
                    }"
                  >
                    {{ tag.name }}
                  </div>
                </div>
              </div>
            </div>
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
  z-index: 1000;
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
  padding: 2rem;
}

.no_events_message {
  text-align: center;
  padding: 3rem 1rem;
  color: #94a3b8;
  font-size: 1rem;
}

.timeline_container {
  position: relative;
  padding-left: 2rem;
}

/* Vertical line */
.timeline_container::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #cbd5e0;
}

.timeline_date_group {
  position: relative;
  margin-bottom: 2rem;
}

.timeline_bullet_wrapper {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.timeline_bullet {
  position: absolute;
  left: 6px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  border: 3px solid #3b82f6;
  z-index: 1;
}

.timeline_date {
  margin-left: 2.5rem;
  font-size: 0.95rem;
  font-weight: 600;
  color: #1e293b;
  background: #f8fafc;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
}

.timeline_events_wrapper {
  margin-left: 2.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.timeline_event_card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem 1.25rem;
  transition: all 0.2s;
}

.timeline_event_card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: #cbd5e0;
}

.timeline_event_header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.event_icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.event_name {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 600;
  color: #1e293b;
  line-height: 1.4;
}

.event_description_wrapper {
  margin-bottom: 0.75rem;
}

.event_description {
  margin: 0;
  color: #475569;
  font-size: 0.925rem;
  line-height: 1.6;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.event_description.expanded {
  display: block;
  -webkit-line-clamp: unset;
}

.expand_btn {
  background: none;
  border: none;
  color: #3b82f6;
  font-size: 0.875rem;
  cursor: pointer;
  padding: 0.25rem 0;
  margin-top: 0.5rem;
  font-weight: 500;
  transition: color 0.2s;
}

.expand_btn:hover {
  color: #2563eb;
  text-decoration: underline;
}

.event_tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.event_tag {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.65rem;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-size: 0.75rem;
  font-weight: 500;
  color: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
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
