import { ref, computed } from 'vue'
import { parseDDMMYYYYToISO, formatDateToDDMMYYYY, getTodayISO } from '@/utils/date-utils.js'
import { getAvailableLensTypes } from '@/utils/event-utils.js'

export function useFilters() {
  // Date filtering state (always showing both fields)
  const dateFrom = ref('0001-01-01')
  const dateTo = ref(getTodayISO())
  const dateFromDisplay = ref('01.01.0001')
  const dateToDisplay = ref(formatDateToDDMMYYYY(new Date()))

  // Lens type filtering state
  const selectedLensTypes = ref(['historic', 'political', 'cultural', 'military']) // All selected by default
  const showLensDropdown = ref(false)

  // Available lens types
  const availableLensTypes = computed(() => getAvailableLensTypes())

  // Reset to default date range
  const resetToDefaultDateRange = () => {
    dateFrom.value = '0001-01-01'
    dateTo.value = getTodayISO()
    dateFromDisplay.value = '01.01.0001'
    dateToDisplay.value = formatDateToDDMMYYYY(new Date())
  }

  // Update date from display
  const updateDateFrom = (displayValue) => {
    dateFromDisplay.value = displayValue
    const isoDate = parseDDMMYYYYToISO(displayValue)
    if (isoDate) {
      dateFrom.value = isoDate
    } else {
      // Reset to default if invalid
      dateFromDisplay.value = '01.01.0001'
      dateFrom.value = '0001-01-01'
    }
  }

  // Update date to display
  const updateDateTo = (displayValue) => {
    dateToDisplay.value = displayValue
    const isoDate = parseDDMMYYYYToISO(displayValue)
    if (isoDate) {
      dateTo.value = isoDate
    } else {
      // Reset to today if invalid
      dateToDisplay.value = formatDateToDDMMYYYY(new Date())
      dateTo.value = getTodayISO()
    }
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