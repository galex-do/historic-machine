import { ref, computed } from 'vue'
import apiService from '@/services/api.js'
import { formatHistoricalDate } from '@/utils/date-utils.js'

// Shared state - singleton pattern
const templateGroups = ref([])
const availableTemplates = ref([])
const selectedTemplateGroupId = ref('')
const selectedTemplateId = ref('')
const loading = ref(false)
const error = ref(null)

export function useTemplates() {

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
    selectedTemplateGroupId.value = groupId
    selectedTemplateId.value = ''
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
    applyTemplate
  }
}