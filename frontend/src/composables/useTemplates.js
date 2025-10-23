import { ref, computed, watch } from 'vue'
import apiService from '@/services/api.js'
import { formatHistoricalDate } from '@/utils/date-utils.js'
import { useLocale } from '@/composables/useLocale.js'

// Session storage keys for templates
const TEMPLATE_STORAGE_KEYS = {
  SELECTED_TEMPLATE_GROUP_ID: 'historia_selected_template_group_id',
  SELECTED_TEMPLATE_ID: 'historia_selected_template_id'
}

// Load template state from session storage
const loadTemplateFromStorage = (key, defaultValue) => {
  try {
    const stored = sessionStorage.getItem(key)
    return stored ? JSON.parse(stored) : defaultValue
  } catch (error) {
    console.warn('Error loading template from session storage:', error)
    return defaultValue
  }
}

// Save template state to session storage
const saveTemplateToStorage = (key, value) => {
  try {
    sessionStorage.setItem(key, JSON.stringify(value))
  } catch (error) {
    console.warn('Error saving template to session storage:', error)
  }
}

// Shared state - singleton pattern
const templateGroups = ref([])
const availableTemplates = ref([])
const selectedTemplateGroupId = ref(loadTemplateFromStorage(TEMPLATE_STORAGE_KEYS.SELECTED_TEMPLATE_GROUP_ID, ''))
const selectedTemplateId = ref(loadTemplateFromStorage(TEMPLATE_STORAGE_KEYS.SELECTED_TEMPLATE_ID, ''))
const loading = ref(false)
const error = ref(null)

export function useTemplates() {
  const { locale } = useLocale()

  // Setup watchers to save template state to session storage
  watch(selectedTemplateGroupId, (newValue) => {
    saveTemplateToStorage(TEMPLATE_STORAGE_KEYS.SELECTED_TEMPLATE_GROUP_ID, newValue)
  })

  watch(selectedTemplateId, (newValue) => {
    saveTemplateToStorage(TEMPLATE_STORAGE_KEYS.SELECTED_TEMPLATE_ID, newValue)
  })

  // Watch locale changes and re-fetch template groups
  watch(locale, async () => {
    console.log('Locale changed, re-fetching template groups')
    await fetchTemplateGroups()
  })

  // Computed to get selected template object
  const selectedTemplate = computed(() => {
    if (!selectedTemplateId.value) return null
    return availableTemplates.value.find(template => template.id == selectedTemplateId.value) || null
  })

  // Fetch template groups from API
  const fetchTemplateGroups = async () => {
    loading.value = true
    error.value = null
    
    try {
      const groupData = await apiService.getTemplateGroups()
      templateGroups.value = Array.isArray(groupData) ? groupData : []
      console.log('Successfully loaded template groups:', templateGroups.value.length)
      
      // If there's a stored template group selection, restore it
      if (selectedTemplateGroupId.value && selectedTemplateGroupId.value !== 'custom') {
        console.log('Restoring template group selection:', selectedTemplateGroupId.value)
        await fetchTemplatesForGroup(selectedTemplateGroupId.value)
      }
    } catch (err) {
      console.error('Error fetching template groups:', err)
      error.value = err.message || 'Failed to fetch template groups'
      templateGroups.value = []
    } finally {
      loading.value = false
    }
  }

  // Fetch templates for a specific group
  const fetchTemplatesForGroup = async (groupId) => {
    if (!groupId) {
      availableTemplates.value = []
      return
    }

    loading.value = true
    error.value = null
    
    try {
      const templateData = await apiService.getTemplatesByGroup(groupId)
      availableTemplates.value = Array.isArray(templateData) ? templateData : []
      console.log('Successfully loaded templates for group:', availableTemplates.value.length)
      
      // If there's a stored template selection and it exists in the loaded templates, verify it's still valid
      if (selectedTemplateId.value) {
        const templateExists = availableTemplates.value.find(t => t.id == selectedTemplateId.value)
        if (templateExists) {
          console.log('Restored template selection:', templateExists.name)
        } else {
          console.log('Stored template ID not found, clearing selection')
          selectedTemplateId.value = ''
        }
      }
    } catch (err) {
      console.error('Error fetching templates for group:', err)
      error.value = err.message
      availableTemplates.value = []
    } finally {
      loading.value = false
    }
  }

  // Handle template group change
  const handleTemplateGroupChange = async (groupId) => {
    const previousGroupId = selectedTemplateGroupId.value
    selectedTemplateGroupId.value = groupId
    
    // Only reset template selection if this is a change to a different group
    if (groupId !== previousGroupId) {
      selectedTemplateId.value = ''
    }
    
    if (groupId && groupId !== 'custom') {
      await fetchTemplatesForGroup(groupId)
    } else {
      // Reset to default when "Default (1 AD - Today)" is selected
      // Or when "Custom" is selected (no templates to fetch)
      availableTemplates.value = []
    }
  }

  // Handle template change
  const handleTemplateChange = (templateId) => {
    selectedTemplateId.value = templateId
  }

  // Apply template to date range
  const applyTemplate = () => {
    const template = selectedTemplate.value
    if (!template) return null

    // Convert backend dates to BC/AD format if needed
    const displayFrom = template.start_display_date || formatHistoricalDate(template.start_date)
    const displayTo = template.end_display_date || formatHistoricalDate(template.end_date)

    return {
      dateFrom: template.start_date,
      dateTo: template.end_date,
      displayFrom,
      displayTo
    }
  }

  // Calculate date range that encompasses all templates in a group
  const getGroupDateRange = () => {
    if (availableTemplates.value.length === 0) {
      return null
    }

    // Helper to convert ISO date to comparable numerical value
    // BC dates: "-YYYY-MM-DD" -> negative value
    // AD dates: "YYYY-MM-DD" -> positive value
    const toComparableValue = (isoDate) => {
      if (!isoDate) return 0
      
      let year, month, day
      
      if (isoDate.startsWith('-')) {
        // BC date: "-1499-01-01"
        const parts = isoDate.substring(1).split('-')
        year = -parseInt(parts[0], 10)
        month = parseInt(parts[1], 10) || 1
        day = parseInt(parts[2], 10) || 1
      } else {
        // AD date: "0146-03-15"
        const parts = isoDate.split('-')
        year = parseInt(parts[0], 10)
        month = parseInt(parts[1], 10) || 1
        day = parseInt(parts[2], 10) || 1
      }
      
      // Convert to numerical value for comparison
      // For BC: more negative = earlier (1500 BC < 146 BC)
      // For AD: more positive = later (146 AD > 1 AD)
      return year + (month / 12) + (day / 365)
    }

    // Find min start_date and max end_date across all templates
    let minStartDate = null
    let maxEndDate = null
    let minStartValue = null
    let maxEndValue = null

    availableTemplates.value.forEach(template => {
      const startValue = toComparableValue(template.start_date)
      const endValue = toComparableValue(template.end_date)

      if (minStartValue === null || startValue < minStartValue) {
        minStartValue = startValue
        minStartDate = template.start_date
      }
      if (maxEndValue === null || endValue > maxEndValue) {
        maxEndValue = endValue
        maxEndDate = template.end_date
      }
    })

    if (!minStartDate || !maxEndDate) {
      return null
    }

    // Format for display using existing utility
    const displayFrom = formatHistoricalDate(minStartDate)
    const displayTo = formatHistoricalDate(maxEndDate)

    return {
      dateFrom: minStartDate,
      dateTo: maxEndDate,
      displayFrom,
      displayTo
    }
  }

  return {
    templateGroups: computed(() => templateGroups.value),
    availableTemplates: computed(() => availableTemplates.value),
    selectedTemplateGroupId: computed(() => selectedTemplateGroupId.value),
    selectedTemplateId: computed(() => selectedTemplateId.value),
    selectedTemplate,
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    fetchTemplateGroups,
    fetchTemplatesForGroup,
    handleTemplateGroupChange,
    handleTemplateChange,
    applyTemplate,
    getGroupDateRange
  }
}