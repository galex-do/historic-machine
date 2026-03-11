import { ref, computed, watch } from 'vue'
import apiService from '@/services/api.js'
import { parseHistoricalDate } from '@/utils/date-utils.js'
import { getEraFromDate } from '@/utils/date-utils.js'
import { useLocale } from '@/composables/useLocale.js'

// Shared state - singleton pattern
const events = ref([])
const filteredEvents = ref([])
const loading = ref(false)
const error = ref(null)
const eventsLoaded = ref(false)

export function useEvents() {
  const { locale } = useLocale()

  // Watch locale changes and re-fetch events
  watch(locale, async () => {
    console.log('Locale changed, re-fetching events')
    await fetchEvents()
  })

  // Fetch all events from API
  const fetchEvents = async () => {
    loading.value = true
    error.value = null
    
    try {
      const eventData = await apiService.getEvents()
      events.value = Array.isArray(eventData) ? eventData : []
      eventsLoaded.value = true
      // Don't automatically set filteredEvents - let filtering be applied explicitly
      console.log('Successfully loaded events:', events.value.length)
    } catch (err) {
      console.error('Error fetching events:', err)
      error.value = err.message || 'Failed to fetch events'
      events.value = []
      filteredEvents.value = []
    } finally {
      loading.value = false
    }
  }

  // Filter events based on date range, lens types, and tags
  const filterEvents = (dateFrom, dateTo, selectedLensTypes, selectedTemplate, dateFromDisplay, dateToDisplay, selectedTags = []) => {
    // Ensure events is an array before filtering
    if (!Array.isArray(events.value)) {
      console.warn('Events is not an array:', events.value)
      events.value = []
    }
    
    const to_chronological = (year, era, month, day) => {
      const m = month || 1
      const d = day || 1
      if (era === 'BC') {
        return -(year * 10000 + (13 - m) * 100 + (32 - d))
      }
      return year * 10000 + m * 100 + d
    }

    const parse_event_date_parts = (event) => {
      let year, month, day
      if (event.event_date.startsWith('-')) {
        const parts = event.event_date.substring(1).split('T')[0].split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10)
        day = parseInt(parts[2], 10)
      } else {
        const parts = event.event_date.split('T')[0].split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10)
        day = parseInt(parts[2], 10)
      }
      return { year, month, day, era: event.era }
    }

    const fromDate = parseHistoricalDate(dateFromDisplay)
    const toDate = parseHistoricalDate(dateToDisplay)
    const fromVal = fromDate ? to_chronological(fromDate.year, fromDate.era, fromDate.month, fromDate.day) : null
    const toVal = toDate ? to_chronological(toDate.year, toDate.era, toDate.month, toDate.day) : null

    let tempFilteredEvents = events.value.filter(event => {
      const ep = parse_event_date_parts(event)
      const eventVal = to_chronological(ep.year, ep.era, ep.month, ep.day)

      if (fromVal !== null && eventVal < fromVal) return false
      if (toVal !== null && eventVal > toVal) return false
      
      // Lens type filtering
      if (selectedLensTypes.length > 0 && !selectedLensTypes.includes(event.lens_type)) {
        return false
      }
      
      // Tag filtering (AND logic - event must have all selected tags)
      if (selectedTags.length > 0) {
        const eventTags = event.tags || []
        const hasAllTags = selectedTags.every(selectedTag => 
          eventTags.some(eventTag => eventTag.id === selectedTag.id)
        )
        if (!hasAllTags) {
          return false
        }
      }
      
      return true
    })
    
    // Sort events chronologically (oldest to newest)
    tempFilteredEvents.sort((a, b) => {
      // Use chronological sorting with proper BC/AD handling
      const getChronologicalValue = (dateString, era) => {
        let year, month, day
        
        if (dateString.startsWith('-')) {
          // Negative year format: "-3501-01-01T00:00:00Z" 
          const parts = dateString.substring(1).split('T')[0].split('-')
          year = parseInt(parts[0], 10)
          month = parseInt(parts[1], 10) - 1  // Month is 0-indexed
          day = parseInt(parts[2], 10)
        } else {
          // Positive year format
          const date = new Date(dateString)
          year = date.getFullYear()
          month = date.getMonth()
          day = date.getDate()
        }
        
        if (era === 'BC') {
          // For BC: larger year number = older (3000 BC < 1000 BC)
          // Within a BC year, months run forward: April comes before September
          // So we subtract month/day to make earlier months in the year sort first
          return -(year - (month / 12) - (day / 365))
        } else {
          // For AD: normal positive years
          return year + (month / 12) + (day / 365)
        }
      }
      
      const aYear = getChronologicalValue(a.event_date, a.era)
      const bYear = getChronologicalValue(b.event_date, b.era)
      
      return aYear - bYear // Chronological order: older events first
    })
    
    filteredEvents.value = tempFilteredEvents
    
    const lensFilterText = selectedLensTypes.length === 4 ? 'all types' : selectedLensTypes.join(', ')
    console.log(`Filtering events from ${dateFrom} to ${dateTo} for lens types: ${lensFilterText}. Found ${filteredEvents.value.length} events.`)
  }

  // Handle event creation
  const handleEventCreated = async (newEvent) => {
    // Add the new event to the events array
    if (newEvent && newEvent.id) {
      events.value.push(newEvent)
      // Also add to filtered events so it's immediately visible on the map
      filteredEvents.value.push(newEvent)
      console.log('New event added to arrays:', newEvent.name)
    } else {
      // Fallback: refresh all events if no event data provided
      await fetchEvents()
      console.log('Events list refreshed after new event creation')
    }
  }

  // Handle event update (preserves current filter state)
  const handleEventUpdated = async (updatedEvent) => {
    // Update the event in the current events array
    const eventIndex = events.value.findIndex(e => e.id === updatedEvent.id)
    if (eventIndex !== -1) {
      events.value[eventIndex] = updatedEvent
      console.log('Event updated in place:', updatedEvent.name)
      
      // Update filtered events array if this event is currently visible
      const filteredIndex = filteredEvents.value.findIndex(e => e.id === updatedEvent.id)
      if (filteredIndex !== -1) {
        filteredEvents.value[filteredIndex] = updatedEvent
      }
    } else {
      console.warn('Event not found for update:', updatedEvent.id)
    }
  }

  // Handle event deletion
  const handleEventDeleted = async (deletedEventId) => {
    // Remove from events array
    events.value = events.value.filter(e => e.id !== deletedEventId)
    // Remove from filtered events array  
    filteredEvents.value = filteredEvents.value.filter(e => e.id !== deletedEventId)
    console.log('Event removed from arrays:', deletedEventId)
  }

  return {
    events: computed(() => events.value),
    filteredEvents: computed(() => filteredEvents.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    eventsLoaded: computed(() => eventsLoaded.value),
    fetchEvents,
    filterEvents,
    handleEventCreated,
    handleEventUpdated,
    handleEventDeleted
  }
}