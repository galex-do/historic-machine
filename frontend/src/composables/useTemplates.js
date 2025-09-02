import { ref, computed } from 'vue'
import apiService from '@/services/api.js'

export function useTemplates() {
  const templateGroups = ref([])
  const availableTemplates = ref([])
  const selectedTemplateGroupId = ref('')
  const selectedTemplateId = ref('')
  const loading = ref(false)
  const error = ref(null)

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
      console.log('Using fallback sample template groups for development')
      
      // Fallback sample data for development
      templateGroups.value = [
        { id: 1, name: "Roman Empire", description: "Major periods of the Roman Empire", display_order: 1 },
        { id: 2, name: "Medieval Period", description: "Key events of the Middle Ages", display_order: 2 },
        { id: 3, name: "Ancient Greece", description: "Classical Greek period", display_order: 3 }
      ]
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
      error.value = err.message || 'Failed to fetch templates'
      console.log('Using fallback sample templates for development')
      
      // Fallback sample data for development
      if (groupId == 1) { // Roman Empire
        availableTemplates.value = [
          {
            id: 1,
            name: "Roman Republic",
            description: "Period of the Roman Republic",
            start_date: "0509-01-01T00:00:00Z",
            end_date: "0049-01-01T00:00:00Z",
            start_display_date: "01.01.509 BC",
            end_display_date: "01.01.49 BC"
          },
          {
            id: 2,
            name: "Roman Empire",
            description: "Imperial period of Rome",
            start_date: "0027-01-01T00:00:00Z",
            end_date: "0476-01-01T00:00:00Z",
            start_display_date: "01.01.27 BC",
            end_display_date: "01.01.476 AD"
          }
        ]
      } else if (groupId == 2) { // Medieval Period
        availableTemplates.value = [
          {
            id: 3,
            name: "Early Middle Ages",
            description: "Dark Ages period",
            start_date: "0476-01-01T00:00:00Z",
            end_date: "1000-01-01T00:00:00Z",
            start_display_date: "01.01.476 AD",
            end_display_date: "01.01.1000 AD"
          }
        ]
      } else {
        availableTemplates.value = []
      }
    } finally {
      loading.value = false
    }
  }

  // Handle template group change
  const handleTemplateGroupChange = async (groupId) => {
    selectedTemplateGroupId.value = groupId
    selectedTemplateId.value = ''
    await fetchTemplatesForGroup(groupId)
  }

  // Handle template change
  const handleTemplateChange = (templateId) => {
    selectedTemplateId.value = templateId
  }

  // Apply template to date range
  const applyTemplate = () => {
    const template = selectedTemplate.value
    if (!template) return null

    return {
      dateFrom: template.start_date,
      dateTo: template.end_date,
      displayFrom: template.start_display_date,
      displayTo: template.end_display_date
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