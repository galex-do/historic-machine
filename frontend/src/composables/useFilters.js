import { ref, computed } from 'vue'
import { parseDDMMYYYYToISO, formatDateToDDMMYYYY, getTodayISO, parseHistoricalDate, formatHistoricalDate, historicalDateToISO } from '@/utils/date-utils.js'
import { getAvailableLensTypes } from '@/utils/event-utils.js'

// Shared state - singleton pattern  
// Date filtering state (always showing both fields)
const dateFrom = ref('0001-01-01')
const dateTo = ref(getTodayISO())
const dateFromDisplay = ref('1 AD')
const dateToDisplay = ref('2025 AD')

// Lens type filtering state
const selectedLensTypes = ref(['historic', 'political', 'cultural', 'military']) // All selected by default
const showLensDropdown = ref(false)

export function useFilters() {

  // Available lens types
  const availableLensTypes = computed(() => getAvailableLensTypes())

  // Reset to default date range
  const resetToDefaultDateRange = () => {
    dateFrom.value = '0001-01-01'
    dateTo.value = getTodayISO()
    dateFromDisplay.value = '1 AD'
    dateToDisplay.value = '2025 AD'
  }

  // Update date from display
  const updateDateFrom = (displayValue) => {
    // Always update the display value to maintain user input
    dateFromDisplay.value = displayValue
    
    // Try to parse as historical date first
    const historicalDate = parseHistoricalDate(displayValue)
    if (historicalDate) {
      dateFrom.value = historicalDate.isoDate
      return
    }
    
    // Fallback to DD.MM.YYYY format
    const isoDate = parseDDMMYYYYToISO(displayValue)
    if (isoDate) {
      dateFrom.value = isoDate
    }
    // Note: Don't reset display value on invalid input to preserve user's typing
  }

  // Update date to display
  const updateDateTo = (displayValue) => {
    // Always update the display value to maintain user input
    dateToDisplay.value = displayValue
    
    // Try to parse as historical date first
    const historicalDate = parseHistoricalDate(displayValue)
    if (historicalDate) {
      dateTo.value = historicalDate.isoDate
      return
    }
    
    // Fallback to DD.MM.YYYY format
    const isoDate = parseDDMMYYYYToISO(displayValue)
    if (isoDate) {
      dateTo.value = isoDate
    }
    // Note: Don't reset display value on invalid input to preserve user's typing
  }

  // Apply template dates
  const applyTemplateDates = (templateData) => {
    if (!templateData) return

    dateFrom.value = templateData.dateFrom
    dateTo.value = templateData.dateTo
    dateFromDisplay.value = templateData.displayFrom
    dateToDisplay.value = templateData.displayTo
  }

  // Toggle lens dropdown
  const toggleLensDropdown = () => {
    showLensDropdown.value = !showLensDropdown.value
  }

  // Handle lens types change
  const handleLensTypesChange = (newLensTypes) => {
    selectedLensTypes.value = newLensTypes
  }

  // Close dropdown when clicking outside
  const closeLensDropdown = () => {
    showLensDropdown.value = false
  }

  return {
    // State
    dateFrom: computed(() => dateFrom.value),
    dateTo: computed(() => dateTo.value),
    dateFromDisplay: computed(() => dateFromDisplay.value),
    dateToDisplay: computed(() => dateToDisplay.value),
    selectedLensTypes: computed(() => selectedLensTypes.value),
    showLensDropdown: computed(() => showLensDropdown.value),
    availableLensTypes,

    // Methods
    resetToDefaultDateRange,
    updateDateFrom,
    updateDateTo,
    applyTemplateDates,
    toggleLensDropdown,
    handleLensTypesChange,
    closeLensDropdown
  }
}