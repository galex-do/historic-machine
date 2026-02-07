<template>
  <div v-if="isOpen" class="timeline_modal_overlay" @click.self="closeModal">
    <div class="timeline_modal">
      <!-- Modal Header -->
      <div class="timeline_modal_header">
        <h2 class="timeline_modal_title">{{ t('historicalEvents') }}</h2>
        <span v-if="totalEventCount > 0" class="timeline_event_count">
          {{ visibleEventCount }} / {{ totalEventCount }}
        </span>
        <button 
          class="toggle_details_btn" 
          @click="toggleDetails"
          :title="showDetails ? t('hideDetails') : t('showDetails')"
        >
          {{ showDetails ? '📝' : '📋' }}
        </button>
        <button class="close_btn" @click="closeModal" :title="t('close')">×</button>
      </div>

      <!-- Modal Content -->
      <div 
        class="timeline_modal_content" 
        ref="scrollContainer"
        @scroll="handleScroll"
      >
        <div v-if="visibleYearGroups.length === 0 && !isLoading" class="no_events_message">
          {{ t('noEventsInTimeline') }}
        </div>

        <div v-else class="timeline_container">
          <div v-for="yearGroup in visibleYearGroups" :key="yearGroup.yearKey" class="timeline_year_group" :data-year-key="yearGroup.yearKey">
            <!-- Single event in this year: everything on one line -->
            <div v-if="yearGroup.totalEvents === 1" class="timeline_single_event_line">
              <div class="timeline_bullet"></div>
              <span class="timeline_single_text">
                <span class="timeline_date_inline">{{ yearGroup.dateGroups[0].events[0]._formattedDate }}</span>
                {{ ' ' }}
                <span class="event_icon">{{ getEventEmoji(yearGroup.dateGroups[0].events[0].lens_type) }}</span>
                {{ ' ' }}
                <span 
                   class="event_name event_name_link"
                   @click="handleShowDetail(yearGroup.dateGroups[0].events[0])"
                >{{ yearGroup.dateGroups[0].events[0].name }}</span>
                <template v-if="showDetails && yearGroup.dateGroups[0].events[0].description">
                  {{ ' — ' }}{{ yearGroup.dateGroups[0].events[0].description }}
                </template>
                <template v-if="showDetails && yearGroup.dateGroups[0].events[0].tags && yearGroup.dateGroups[0].events[0].tags.length > 0">
                  {{ ' ' }}
                  <span
                    v-for="tag in yearGroup.dateGroups[0].events[0].tags"
                    :key="tag.id"
                    class="event_tag_badge"
                    :style="{ backgroundColor: tag.color || '#6366f1', color: getContrastColor(tag.color || '#6366f1') }"
                    :title="`Click to filter events by '${tag.name}'`"
                    @click.stop="handleTagClick(tag)"
                  >{{ tag.name }}</span>
                </template>
                <template v-if="showDetails">
                  {{ ' ' }}
                  <button 
                    class="timeline_focus_btn" 
                    @click="handleFocusEvent(yearGroup.dateGroups[0].events[0])"
                    :title="t('focusOnMap')"
                  >
                    ⌖
                  </button>
                </template>
              </span>
            </div>

            <!-- Multiple events in this year -->
            <template v-else>
              <!-- Year Header -->
              <div class="timeline_date_header">
                <div class="timeline_bullet"></div>
                <div class="timeline_date">{{ yearGroup.formattedYear }}</div>
              </div>

              <!-- Date subgroups within year -->
              <div class="timeline_events_list">
                <template v-for="dateGroup in yearGroup.dateGroups" :key="dateGroup.date">
                  <!-- Date subheader (only for non-Jan-1 dates with multiple events, or mixed dates) -->
                  <div v-if="dateGroup.showDateHeader" class="timeline_date_subheader">
                    {{ dateGroup.formattedDate }}
                  </div>

                  <div 
                    v-for="event in dateGroup.events" 
                    :key="event.id"
                    class="timeline_event_line"
                    :class="{ 'timeline_event_line_indented': dateGroup.showDateHeader }"
                  >
                    <span class="timeline_single_text">
                      <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
                      {{ ' ' }}
                      <span 
                         class="event_name event_name_link"
                         @click="handleShowDetail(event)"
                      >{{ event.name }}</span>
                      <template v-if="showDetails && event.description">
                        {{ ' — ' }}{{ event.description }}
                      </template>
                      <template v-if="showDetails && event.tags && event.tags.length > 0">
                        {{ ' ' }}
                        <span
                          v-for="tag in event.tags"
                          :key="tag.id"
                          class="event_tag_badge"
                          :style="{ backgroundColor: tag.color || '#6366f1', color: getContrastColor(tag.color || '#6366f1') }"
                          :title="`Click to filter events by '${tag.name}'`"
                          @click.stop="handleTagClick(tag)"
                        >{{ tag.name }}</span>
                      </template>
                      <template v-if="showDetails">
                        {{ ' ' }}
                        <button 
                          class="timeline_focus_btn" 
                          @click="handleFocusEvent(event)"
                          :title="t('focusOnMap')"
                        >
                          ⌖
                        </button>
                      </template>
                    </span>
                  </div>
                </template>
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
import { getContrastColor } from '@/utils/color-utils.js'

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
    },
    preserveScroll: {
      type: Boolean,
      default: false
    }
  },
  emits: ['close', 'focus-event', 'tag-clicked', 'show-detail'],
  setup(props, { emit }) {
    const { t, formatEventDisplayDate, formatDayMonth, currentLocale } = useLocale()
    const previouslyFocusedElement = ref(null)
    const scrollContainer = ref(null)
    
    const BATCH_SIZE = 50
    const visibleCount = ref(BATCH_SIZE)
    const showDetails = ref(false)
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

    const getYearKey = (isoDateString, era) => {
      let year
      if (isoDateString.startsWith('-')) {
        year = parseInt(isoDateString.substring(1).split('T')[0].split('-')[0], 10)
      } else {
        year = parseInt(isoDateString.split('T')[0].split('-')[0], 10)
      }
      return `${year}_${era || 'AD'}`
    }

    const getFormattedYear = (isoDateString, era) => {
      let year
      if (isoDateString.startsWith('-')) {
        year = parseInt(isoDateString.substring(1).split('T')[0].split('-')[0], 10)
      } else {
        year = parseInt(isoDateString.split('T')[0].split('-')[0], 10)
      }
      const eraLabel = era === 'BC' ? t('eraBC') : t('eraAD')
      return `${year} ${eraLabel}`
    }

    const isJanFirst = (isoDateString) => {
      const datePart = isoDateString.startsWith('-')
        ? isoDateString.substring(1).split('T')[0]
        : isoDateString.split('T')[0]
      const parts = datePart.split('-')
      return parseInt(parts[1], 10) === 1 && parseInt(parts[2], 10) === 1
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

      const yearGroups = new Map()
      sortedEvents.forEach(event => {
        const yearKey = getYearKey(event.event_date, event.era)
        if (!yearGroups.has(yearKey)) {
          yearGroups.set(yearKey, {
            yearKey,
            formattedYear: getFormattedYear(event.event_date, event.era),
            events: []
          })
        }
        yearGroups.get(yearKey).events.push(event)
      })

      return Array.from(yearGroups.values()).map(yg => {
        const dateMap = new Map()
        yg.events.forEach(event => {
          const eventDate = event.event_date.split('T')[0] + '_' + (event.era || 'AD')
          if (!dateMap.has(eventDate)) {
            dateMap.set(eventDate, {
              date: eventDate,
              formattedDate: formatDayMonth(event.event_date),
              isYearOnly: isJanFirst(event.event_date),
              events: []
            })
          }
          dateMap.get(eventDate).events.push({
            ...event,
            _formattedDate: formatEventDisplayDate(event.event_date, event.era)
          })
        })

        const dateGroups = Array.from(dateMap.values())
        const hasSpecificDates = dateGroups.some(dg => !dg.isYearOnly)
        dateGroups.forEach(dg => {
          dg.showDateHeader = !dg.isYearOnly && (hasSpecificDates || dateGroups.length > 1)
        })

        return {
          yearKey: yg.yearKey,
          formattedYear: yg.formattedYear,
          totalEvents: yg.events.length,
          dateGroups
        }
      })
    }

    watch([() => props.events, currentLocale], ([newEvents]) => {
      const newHash = getEventsHash(newEvents) + '_' + currentLocale.value?.code
      if (newHash !== lastEventsHash.value) {
        lastEventsHash.value = newHash
        cachedGroupedEvents.value = computeGroupedEvents(newEvents)
      }
    }, { immediate: true })

    const allGroupedEvents = computed(() => cachedGroupedEvents.value)

    const totalEventCount = computed(() => props.events?.length || 0)

    const visibleYearGroups = computed(() => {
      const groups = allGroupedEvents.value
      let eventCount = 0
      const result = []
      
      for (const yearGroup of groups) {
        if (eventCount >= visibleCount.value) break
        
        const remainingSlots = visibleCount.value - eventCount
        if (yearGroup.totalEvents <= remainingSlots) {
          result.push(yearGroup)
          eventCount += yearGroup.totalEvents
        } else {
          const trimmedDateGroups = []
          let used = 0
          for (const dg of yearGroup.dateGroups) {
            if (used >= remainingSlots) break
            const slotsLeft = remainingSlots - used
            if (dg.events.length <= slotsLeft) {
              trimmedDateGroups.push(dg)
              used += dg.events.length
            } else {
              trimmedDateGroups.push({ ...dg, events: dg.events.slice(0, slotsLeft) })
              used += slotsLeft
              break
            }
          }
          result.push({ ...yearGroup, totalEvents: used, dateGroups: trimmedDateGroups })
          eventCount += used
          break
        }
      }
      
      return result
    })

    const visibleEventCount = computed(() => {
      return visibleYearGroups.value.reduce((sum, yg) => sum + yg.totalEvents, 0)
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

    const toggleDetails = () => {
      const container = scrollContainer.value
      if (!container) {
        showDetails.value = !showDetails.value
        return
      }
      const containerRect = container.getBoundingClientRect()
      const groups = container.querySelectorAll('.timeline_year_group[data-year-key]')
      let anchor = null
      for (const el of groups) {
        const elRect = el.getBoundingClientRect()
        const relativeTop = elRect.top - containerRect.top
        if (relativeTop <= 20) {
          anchor = { key: el.dataset.yearKey, relativeTop }
        } else {
          if (!anchor) anchor = { key: el.dataset.yearKey, relativeTop }
          break
        }
      }
      showDetails.value = !showDetails.value
      nextTick(() => {
        if (anchor) {
          const el = container.querySelector(`.timeline_year_group[data-year-key="${anchor.key}"]`)
          if (el) {
            const newRect = el.getBoundingClientRect()
            const newRelativeTop = newRect.top - containerRect.top
            container.scrollTop += newRelativeTop - anchor.relativeTop
          }
        }
      })
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

    const savedScrollPosition = ref(0)
    const savedVisibleCount = ref(0)

    const handleShowDetail = (event) => {
      // Save scroll position and loaded count before closing
      if (scrollContainer.value) {
        savedScrollPosition.value = scrollContainer.value.scrollTop
      }
      savedVisibleCount.value = visibleCount.value
      emit('show-detail', event)
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
        
        if (props.preserveScroll && savedVisibleCount.value > 0) {
          // Back navigation - restore loaded count first
          visibleCount.value = savedVisibleCount.value
        } else {
          // Fresh open
          visibleCount.value = BATCH_SIZE
        }
        
        nextTick(() => {
          if (scrollContainer.value) {
            if (props.preserveScroll && savedScrollPosition.value > 0) {
              // Restore saved position (back navigation)
              scrollContainer.value.scrollTop = savedScrollPosition.value
            } else {
              // Fresh open - reset to top
              scrollContainer.value.scrollTop = 0
              savedScrollPosition.value = 0
              savedVisibleCount.value = 0
            }
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
      visibleYearGroups,
      totalEventCount,
      visibleEventCount,
      hasMoreEvents,
      isLoading,
      showDetails,
      toggleDetails,
      closeModal,
      handleFocusEvent,
      handleTagClick,
      handleShowDetail,
      handleScroll,
      getEventEmoji,
      getContrastColor
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

.toggle_details_btn {
  background: none;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 1.1rem;
  color: #64748b;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  transition: all 0.2s;
}

.toggle_details_btn:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
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

.timeline_year_group {
  position: relative;
  margin-bottom: 0.25rem;
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
  cursor: pointer;
  transition: all 0.2s;
}

.timeline_single_text .event_name_link:hover {
  color: #2563eb;
}

.timeline_events_list {
  margin-left: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.timeline_date_subheader {
  font-weight: 600;
  font-size: 0.8rem;
  color: #475569;
  padding: 0.15rem 0 0.05rem 0.35rem;
  margin-top: 0.15rem;
}

.timeline_event_line {
  padding: 0.15rem 0;
  padding-left: 0.35rem;
  line-height: 1.5;
}

.timeline_event_line_indented {
  padding-left: 0.75rem;
}

.event_tag_badge {
  display: inline-block;
  padding: 0.05rem 0.4rem;
  border-radius: 9999px;
  font-size: 0.65rem;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
  white-space: nowrap;
  vertical-align: middle;
  margin: 0 0.1rem;
}

.event_tag_badge:hover {
  opacity: 0.75;
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
