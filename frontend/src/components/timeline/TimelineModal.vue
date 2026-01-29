<template>
  <div v-if="isOpen" class="timeline_modal_overlay" @click.self="closeModal">
    <div class="timeline_modal">
      <!-- Modal Header -->
      <div class="timeline_modal_header">
        <h2 class="timeline_modal_title">{{ t('historicalEvents') }}</h2>
        <span v-if="totalEventCount > 0" class="timeline_event_count">
          {{ visibleEventCount }} / {{ totalEventCount }}
        </span>
        <button class="close_btn" @click="closeModal" :title="t('close')">×</button>
      </div>

      <!-- Modal Content -->
      <div 
        class="timeline_modal_content" 
        ref="scrollContainer"
        @scroll="handleScroll"
      >
        <div v-if="visibleGroups.length === 0 && !isLoading" class="no_events_message">
          {{ t('noEventsInTimeline') }}
        </div>

        <div v-else class="timeline_container">
          <div v-for="group in visibleGroups" :key="group.date" class="timeline_date_group">
            <!-- Single event on this date: everything on one line -->
            <div v-if="group.events.length === 1" class="timeline_single_event_line">
              <div class="timeline_bullet"></div>
              <span class="timeline_single_text">
                <span class="timeline_date_inline">{{ group.formattedDate }}</span>
                {{ ' ' }}
                <span class="event_icon">{{ getEventEmoji(group.events[0].lens_type) }}</span>
                {{ ' ' }}
                <a v-if="group.events[0].source" 
                   :href="group.events[0].source" 
                   target="_blank" 
                   rel="noopener noreferrer"
                   class="event_name event_name_link">{{ group.events[0].name }}</a>
                <span v-else class="event_name">{{ group.events[0].name }}</span>
                <template v-if="group.events[0].description">
                  {{ ' — ' }}{{ group.events[0].description }}
                </template>
                <template v-if="group.events[0].tags && group.events[0].tags.length > 0">
                  {{ ' ' }}
                  <span
                    v-for="(tag, index) in group.events[0].tags"
                    :key="tag.id"
                    class="event_tag"
                    :style="{ color: tag.color || '#6366f1' }"
                    :title="`Click to filter events by '${tag.name}'`"
                    @click.stop="handleTagClick(tag)"
                  >#{{ tag.name }}{{ index < group.events[0].tags.length - 1 ? ' ' : '' }}</span>
                </template>
                {{ ' ' }}
                <button 
                  class="timeline_focus_btn" 
                  @click="handleFocusEvent(group.events[0])"
                  :title="t('focusOnMap')"
                >
                  ⌖
                </button>
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
                  <!-- Event line: Icon + Name + Description + Tags in single flowing text -->
                  <span class="timeline_single_text">
                    <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
                    {{ ' ' }}
                    <a v-if="event.source" 
                       :href="event.source" 
                       target="_blank" 
                       rel="noopener noreferrer"
                       class="event_name event_name_link">{{ event.name }}</a>
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
                        :title="`Click to filter events by '${tag.name}'`"
                        @click.stop="handleTagClick(tag)"
                      >#{{ tag.name }}{{ index < event.tags.length - 1 ? ' ' : '' }}</span>
                    </template>
                    {{ ' ' }}
                    <button 
                      class="timeline_focus_btn" 
                      @click="handleFocusEvent(event)"
                      :title="t('focusOnMap')"
                    >
                      ⌖
                    </button>
                  </span>
                </div>
              </div>
            </template>
          </div>
          
          <!-- Loading indicator -->
          <div v-if="isLoading" class="timeline_loading">
            {{ t('loading') || 'Loading...' }}
          </div>
          
          <!-- Load more indicator -->
          <div v-if="hasMoreEvents && !isLoading" class="timeline_load_more" ref="loadMoreTrigger">
            <span class="load_more_hint">↓ {{ t('scrollForMore') || 'Scroll for more' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch, onUnmounted, nextTick } from 'vue'
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
  emits: ['close', 'focus-event', 'tag-clicked'],
  setup(props, { emit }) {
    const { t, formatEventDisplayDate } = useLocale()
    const previouslyFocusedElement = ref(null)
    const scrollContainer = ref(null)
    
    const BATCH_SIZE = 50
    const visibleCount = ref(BATCH_SIZE)
    const isLoading = ref(false)
    
    const cachedGroupedEvents = ref([])
    const lastEventsHash = ref('')

    const getEventsHash = (events) => {
      if (!events || events.length === 0) return ''
      return `${events.length}-${events[0]?.id}-${events[events.length - 1]?.id}`
    }

    const getChronologicalValue = (dateString, era) => {
      let year, month, day
      
      if (dateString.startsWith('-')) {
        const parts = dateString.substring(1).split('T')[0].split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10) - 1
        day = parseInt(parts[2], 10)
      } else {
        const parts = dateString.split('T')[0].split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10) - 1
        day = parseInt(parts[2], 10)
      }
      
      if (era === 'BC') {
        return -(year - (month / 12) - (day / 365))
      } else {
        return year + (month / 12) + (day / 365)
      }
    }

    const computeGroupedEvents = (events) => {
      if (!events || events.length === 0) {
        return []
      }

      const sortedEvents = [...events].sort((a, b) => {
        const aValue = getChronologicalValue(a.event_date, a.era)
        const bValue = getChronologicalValue(b.event_date, b.era)
        return aValue - bValue
      })

      const groups = {}
      sortedEvents.forEach(event => {
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
    }

    watch(() => props.events, (newEvents) => {
      const newHash = getEventsHash(newEvents)
      if (newHash !== lastEventsHash.value) {
        lastEventsHash.value = newHash
        cachedGroupedEvents.value = computeGroupedEvents(newEvents)
      }
    }, { immediate: true })

    const allGroupedEvents = computed(() => cachedGroupedEvents.value)

    const totalEventCount = computed(() => props.events?.length || 0)

    const visibleGroups = computed(() => {
      const groups = allGroupedEvents.value
      let eventCount = 0
      const result = []
      
      for (const group of groups) {
        if (eventCount >= visibleCount.value) break
        
        const remainingSlots = visibleCount.value - eventCount
        if (group.events.length <= remainingSlots) {
          result.push(group)
          eventCount += group.events.length
        } else {
          result.push({
            ...group,
            events: group.events.slice(0, remainingSlots)
          })
          eventCount += remainingSlots
          break
        }
      }
      
      return result
    })

    const visibleEventCount = computed(() => {
      return visibleGroups.value.reduce((sum, group) => sum + group.events.length, 0)
    })

    const hasMoreEvents = computed(() => {
      return visibleEventCount.value < totalEventCount.value
    })

    const loadMore = () => {
      if (isLoading.value || !hasMoreEvents.value) return
      
      isLoading.value = true
      
      requestAnimationFrame(() => {
        visibleCount.value += BATCH_SIZE
        isLoading.value = false
      })
    }

    const handleScroll = (e) => {
      const container = e.target
      const scrollBottom = container.scrollHeight - container.scrollTop - container.clientHeight
      
      if (scrollBottom < 100 && hasMoreEvents.value && !isLoading.value) {
        loadMore()
      }
    }

    const closeModal = () => {
      emit('close')
      if (previouslyFocusedElement.value) {
        previouslyFocusedElement.value.focus()
        previouslyFocusedElement.value = null
      }
    }

    const handleFocusEvent = (event) => {
      emit('focus-event', event)
      closeModal()
    }

    const handleTagClick = (tag) => {
      emit('tag-clicked', tag)
      closeModal()
    }

    const handleEscape = (event) => {
      if (event.key === 'Escape' && props.isOpen) {
        closeModal()
      }
    }

    watch(() => props.isOpen, (newValue) => {
      if (newValue) {
        previouslyFocusedElement.value = document.activeElement
        document.addEventListener('keydown', handleEscape)
        visibleCount.value = BATCH_SIZE
        
        nextTick(() => {
          if (scrollContainer.value) {
            scrollContainer.value.scrollTop = 0
          }
        })
      } else {
        document.removeEventListener('keydown', handleEscape)
      }
    })

    onUnmounted(() => {
      document.removeEventListener('keydown', handleEscape)
    })

    return {
      t,
      scrollContainer,
      visibleGroups,
      totalEventCount,
      visibleEventCount,
      hasMoreEvents,
      isLoading,
      closeModal,
      handleFocusEvent,
      handleTagClick,
      handleScroll,
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
  gap: 1rem;
}

.timeline_modal_title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1e293b;
}

.timeline_event_count {
  font-size: 0.875rem;
  color: #64748b;
  background: #f1f5f9;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  margin-left: auto;
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
  margin-left: -6px;
}

.timeline_date {
  font-weight: 600;
  font-size: 0.875rem;
  color: #1e293b;
}

.timeline_single_event_line {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  padding: 0.25rem 0;
  line-height: 1.5;
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

.timeline_single_text .event_name_link {
  font-weight: 700;
  color: #3b82f6;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s;
}

.timeline_single_text .event_name_link:hover {
  color: #2563eb;
  border-bottom-color: #2563eb;
}

.timeline_events_list {
  margin-left: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.timeline_event_line {
  padding: 0.25rem 0;
  padding-left: 0.5rem;
  line-height: 1.5;
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

.timeline_focus_btn {
  background: none;
  border: none;
  color: #64748b;
  font-size: 1rem;
  cursor: pointer;
  padding: 0 0.25rem;
  margin-left: 0.25rem;
  transition: color 0.2s;
  vertical-align: middle;
}

.timeline_focus_btn:hover {
  color: #3b82f6;
}

.timeline_loading {
  text-align: center;
  padding: 1rem;
  color: #64748b;
  font-size: 0.875rem;
}

.timeline_load_more {
  text-align: center;
  padding: 0.75rem;
  margin-top: 0.5rem;
}

.load_more_hint {
  font-size: 0.75rem;
  color: #94a3b8;
}

.timeline_modal_content::-webkit-scrollbar {
  width: 6px;
}

.timeline_modal_content::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.timeline_modal_content::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}
</style>
