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
    
    let tempFilteredEvents = events.value.filter(event => {
      // Parse filter dates to understand BC/AD for comparison
      const fromDate = parseHistoricalDate(dateFromDisplay) // Use display value to get proper BC/AD
      const toDate = parseHistoricalDate(dateToDisplay)
      
      
      // Parse event year correctly from negative and positive years
      let eventYear
      if (event.event_date.startsWith('-')) {
        // Negative year format: "-3501-01-01T00:00:00Z" -> year 3501
        const yearMatch = event.event_date.match(/^-(\d+)-/)
        eventYear = yearMatch ? parseInt(yearMatch[1], 10) : 0
      } else {
        // Positive year format: "1453-05-29T00:00:00Z" -> year 1453
        eventYear = parseInt(event.event_date.split('-')[0], 10)
      }
      const eventEra = event.era
      
      // Date filtering with proper BC/AD logic
      // For BC: smaller number = more recent (25 BC is after 500 BC)
      // For AD: larger number = more recent (500 AD is after 25 AD)
      
      if (fromDate) {
        if (fromDate.era === 'BC' && eventEra === 'BC') {
          // Both BC: event must be same year or more recent (smaller year number)
          if (eventYear > fromDate.year) return false
        } else if (fromDate.era === 'BC' && eventEra === 'AD') {
          // Event is AD, filter is BC: AD events are always after BC, so include
        } else if (fromDate.era === 'AD' && eventEra === 'BC') {
          // Event is BC, filter is AD: BC events are always before AD, so exclude
          return false
        } else {
          // Both AD: normal comparison
          if (eventYear < fromDate.year) return false
        }
      }
      
      if (toDate) {
        if (toDate.era === 'BC' && eventEra === 'BC') {
          // Both BC: event must be same year or older (larger year number)
          if (eventYear < toDate.year) return false
        } else if (toDate.era === 'BC' && eventEra === 'AD') {
          // Event is AD, filter is BC: AD events are always after BC, so exclude
          return false
        } else if (toDate.era === 'AD' && eventEra === 'BC') {
          // Event is BC, filter is AD: BC events are always before AD, so include
        } else {
          // Both AD: normal comparison
          if (eventYear > toDate.year) return false
        }
      }
      
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