import { ref, computed } from 'vue'
import apiService from '@/services/api.js'
import { parseHistoricalDate } from '@/utils/date-utils.js'
import { getAstronomicalYear, getEraFromDate } from '@/utils/date-utils.js'

// Shared state - singleton pattern
const events = ref([])
const filteredEvents = ref([])
const loading = ref(false)
const error = ref(null)
const eventsLoaded = ref(false)

export function useEvents() {

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

  // Filter events based on date range and lens types
  const filterEvents = (dateFrom, dateTo, selectedLensTypes, selectedTemplate, dateFromDisplay, dateToDisplay) => {
    // Ensure events is an array before filtering
    if (!Array.isArray(events.value)) {
      console.warn('Events is not an array:', events.value)
      events.value = []
    }
    
    let tempFilteredEvents = events.value.filter(event => {
      // Parse filter dates to understand BC/AD for comparison
      const fromDate = parseHistoricalDate(dateFromDisplay) // Use display value to get proper BC/AD
      const toDate = parseHistoricalDate(dateToDisplay)
      
      // Parse event year correctly - handle negative years for BC dates
      let eventYear
      if (event.event_date.startsWith('-')) {
        // BC date format: "-491-09-12T00:00:00Z" -> year 491
        const yearMatch = event.event_date.match(/^-(\d+)-/)
        eventYear = yearMatch ? parseInt(yearMatch[1], 10) : 0
      } else {
        // AD date format: "1453-05-29T00:00:00Z" -> year 1453
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
      
      return true
    })
    
    // Sort events chronologically (oldest to newest)
    tempFilteredEvents.sort((a, b) => {
      const aDate = new Date(a.event_date)
      const bDate = new Date(b.event_date)
      
      // Handle BC/AD chronological sorting
      if (a.era === 'BC' && b.era === 'BC') {
        // For BC: reverse order (larger year = older)
        return -aDate.getTime() - (-bDate.getTime())
      } else if (a.era === 'AD' && b.era === 'AD') {
        // For AD: normal order (larger year = newer)
        return aDate.getTime() - bDate.getTime()
      } else if (a.era === 'BC' && b.era === 'AD') {
        // BC always comes before AD
        return -1
      } else if (a.era === 'AD' && b.era === 'BC') {
        // AD always comes after BC
        return 1
      }
      return 0
    })
    
    filteredEvents.value = tempFilteredEvents
    
    const lensFilterText = selectedLensTypes.length === 4 ? 'all types' : selectedLensTypes.join(', ')
    console.log(`Filtering events from ${dateFrom} to ${dateTo} for lens types: ${lensFilterText}. Found ${filteredEvents.value.length} events.`)
  }

  // Handle event creation
  const handleEventCreated = async () => {
    await fetchEvents()
    console.log('Events list refreshed after new event creation')
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