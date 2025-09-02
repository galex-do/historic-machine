import { ref, computed } from 'vue'
import apiService from '@/services/api.js'
import { getAstronomicalYear, getEraFromDate } from '@/utils/date-utils.js'

export function useEvents() {
  const events = ref([])
  const filteredEvents = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Fetch all events from API
  const fetchEvents = async () => {
    loading.value = true
    error.value = null
    
    try {
      const eventData = await apiService.getEvents()
      events.value = Array.isArray(eventData) ? eventData : []
      filteredEvents.value = events.value
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
  const filterEvents = (dateFrom, dateTo, selectedLensTypes, selectedTemplate) => {
    // Ensure events is an array before filtering
    if (!Array.isArray(events.value)) {
      console.warn('Events is not an array:', events.value)
      events.value = []
    }
    
    filteredEvents.value = events.value.filter(event => {
      // Convert both event date and filter dates to astronomical years for proper BC/AD comparison
      const eventDate = event.event_date.split('T')[0] // Get just the date part: "0753-04-21"
      const eventYear = parseInt(eventDate.split('-')[0], 10)
      
      // Convert event to astronomical year (BC dates become negative)
      const eventAstronomicalDate = event.era === 'BC' 
        ? `-${String(eventYear - 1).padStart(4, '0')}-01-01`  // 753 BC -> -0752-01-01
        : eventDate // AD dates stay as-is
      
      // Date filtering with proper BC/AD comparison using astronomical dates
      if (dateFrom && eventAstronomicalDate < dateFrom) {
        return false
      }
      if (dateTo && eventAstronomicalDate > dateTo) {
        return false
      }
      
      // Lens type filtering
      if (selectedLensTypes.length > 0 && !selectedLensTypes.includes(event.lens_type)) {
        return false
      }
      
      return true
    })
    
    const lensFilterText = selectedLensTypes.length === 4 ? 'all types' : selectedLensTypes.join(', ')
    console.log(`Filtering events from ${dateFrom} to ${dateTo} for lens types: ${lensFilterText}. Found ${filteredEvents.value.length} events.`)
  }

  // Handle event creation
  const handleEventCreated = async () => {
    await fetchEvents()
    console.log('Events list refreshed after new event creation')
  }

  return {
    events: computed(() => events.value),
    filteredEvents: computed(() => filteredEvents.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    fetchEvents,
    filterEvents,
    handleEventCreated
  }
}