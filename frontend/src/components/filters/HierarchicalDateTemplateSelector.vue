<template>
  <div class="hierarchical-selector" ref="selectorRoot">
    <!-- Trigger Button -->
    <button 
      @click="togglePopover" 
      class="selector-button"
      :class="{ 'active': isOpen }"
      type="button"
    >
      <span class="button-label">{{ t('historicalPeriod') }}:</span>
      <span class="button-value">{{ displayText }}</span>
      <span class="button-icon">{{ isOpen ? '▲' : '▼' }}</span>
    </button>

    <!-- Two-Pane Popover -->
    <Teleport to="body">
      <div 
        v-if="isOpen" 
        class="popover-overlay"
        @click="closePopover"
      >
        <div 
          class="popover-container" 
          @click.stop
          :style="popoverPosition"
        >
          <!-- Left Pane: Template Groups -->
          <div class="popover-pane groups-pane">
            <div class="pane-header">{{ t('selectPeriodGroup') }}</div>
            <div class="pane-content">
              <!-- Default option -->
              <div 
                class="group-item"
                :class="{ 'selected': selectedTemplateGroupId === '' }"
                @click="selectGroup('')"
              >
                <span class="item-icon">🌍</span>
                <span class="item-name">{{ t('defaultPeriod') }}</span>
              </div>
              
              <!-- Custom option -->
              <div 
                class="group-item"
                :class="{ 'selected': selectedTemplateGroupId === 'custom' }"
                @click="selectGroup('custom')"
              >
                <span class="item-icon">✏️</span>
                <span class="item-name">{{ t('customDateRange') }}</span>
              </div>
              
              <div class="group-divider"></div>
              
              <!-- Template Groups -->
              <div 
                v-for="group in templateGroups" 
                :key="group.id"
                class="group-item"
                :class="{ 
                  'selected': selectedTemplateGroupId === group.id,
                  'active': hoveredGroupId === group.id
                }"
                @click="selectGroup(group.id)"
                @mouseenter="hoveredGroupId = group.id"
              >
                <span class="item-icon">📅</span>
                <span class="item-name">{{ group.name }}</span>
                <span class="item-arrow">›</span>
              </div>
            </div>
          </div>

          <!-- Right Pane: Templates -->
          <div 
            v-if="showTemplatesPane" 
            class="popover-pane templates-pane"
          >
            <div class="pane-header">{{ selectedGroupName }}</div>
            <div class="pane-content">
              <div v-if="loadingTemplates" class="loading-state">
                <div class="spinner-small"></div>
                <span>{{ t('loading') }}...</span>
              </div>
              
              <div v-else-if="availableTemplates.length === 0" class="empty-state">
                {{ t('noTemplatesAvailable') }}
              </div>
              
              <div 
                v-else
                v-for="template in availableTemplates" 
                :key="template.id"
                class="template-item"
                :class="{ 'selected': selectedTemplateId === template.id }"
                @click="selectTemplate(template.id)"
              >
                <span class="item-icon">⏳</span>
                <div class="template-info">
                  <span class="template-name">{{ template.name }}</span>
                  <span class="template-dates">{{ formatTemplateDates(template) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'HierarchicalDateTemplateSelector',
  props: {
    templateGroups: {
      type: Array,
      default: () => []
    },
    availableTemplates: {
      type: Array,
      default: () => []
    },
    selectedTemplateGroupId: {
      type: [String, Number],
      default: ''
    },
    selectedTemplateId: {
      type: [String, Number],
      default: ''
    },
    selectedTemplate: {
      type: Object,
      default: null
    },
    loadingTemplates: {
      type: Boolean,
      default: false
    }
  },
  emits: ['template-group-changed', 'template-changed'],
  setup(props, { emit }) {
    const { t } = useLocale()
    const isOpen = ref(false)
    const hoveredGroupId = ref(null)
    const selectorRoot = ref(null)
    const popoverPosition = ref({})

    // Display text for the button
    const displayText = computed(() => {
      if (props.selectedTemplate) {
        const groupName = selectedGroupName.value
        const templateName = props.selectedTemplate.name
        return `${groupName} › ${templateName}`
      }
      
      if (props.selectedTemplateGroupId === 'custom') {
        return t('customDateRange')
      }
      
      if (props.selectedTemplateGroupId) {
        const group = props.templateGroups.find(g => g.id == props.selectedTemplateGroupId)
        return group ? group.name : t('defaultPeriod')
      }
      
      return t('defaultPeriod')
    })

    // Get selected group name
    const selectedGroupName = computed(() => {
      if (props.selectedTemplateGroupId === 'custom') {
        return t('customDateRange')
      }
      if (props.selectedTemplateGroupId === '') {
        return t('defaultPeriod')
      }
      const group = props.templateGroups.find(g => g.id == props.selectedTemplateGroupId)
      return group ? group.name : t('selectPeriodGroup')
    })

    // Show templates pane only when a group with templates is selected
    const showTemplatesPane = computed(() => {
      return props.selectedTemplateGroupId && 
             props.selectedTemplateGroupId !== '' && 
             props.selectedTemplateGroupId !== 'custom'
    })

    // Format template date range
    const formatTemplateDates = (template) => {
      if (template.start_display_date && template.end_display_date) {
        return `${template.start_display_date} - ${template.end_display_date}`
      }
      return ''
    }

    // Toggle popover
    const togglePopover = () => {
      if (isOpen.value) {
        closePopover()
      } else {
        openPopover()
      }
    }

    // Open popover and calculate position
    const openPopover = () => {
      if (selectorRoot.value) {
        const rect = selectorRoot.value.getBoundingClientRect()
        popoverPosition.value = {
          top: `${rect.bottom + 8}px`,
          left: `${rect.left}px`
        }
      }
      isOpen.value = true
    }

    // Close popover
    const closePopover = () => {
      isOpen.value = false
      hoveredGroupId.value = null
    }

    // Select a group
    const selectGroup = (groupId) => {
      emit('template-group-changed', groupId)
      
      // If selecting default or custom, close immediately
      if (groupId === '' || groupId === 'custom') {
        closePopover()
      }
      // Otherwise, keep open to select template
    }

    // Select a template
    const selectTemplate = (templateId) => {
      emit('template-changed', templateId)
      closePopover()
    }

    // Handle Escape key
    const handleEscape = (event) => {
      if (event.key === 'Escape' && isOpen.value) {
        closePopover()
      }
    }

    // Recalculate position on window resize
    const handleResize = () => {
      if (isOpen.value && selectorRoot.value) {
        const rect = selectorRoot.value.getBoundingClientRect()
        popoverPosition.value = {
          top: `${rect.bottom + 8}px`,
          left: `${rect.left}px`
        }
      }
    }

    onMounted(() => {
      document.addEventListener('keydown', handleEscape)
      window.addEventListener('resize', handleResize)
    })

    onUnmounted(() => {
      document.removeEventListener('keydown', handleEscape)
      window.removeEventListener('resize', handleResize)
    })

    // Auto-hover first group when templates pane should show
    watch(() => props.selectedTemplateGroupId, (newGroupId) => {
      if (newGroupId && newGroupId !== '' && newGroupId !== 'custom') {
        hoveredGroupId.value = newGroupId
      }
    })

    return {
      isOpen,
      hoveredGroupId,
      selectorRoot,
      popoverPosition,
      displayText,
      selectedGroupName,
      showTemplatesPane,
      formatTemplateDates,
      togglePopover,
      closePopover,
      selectGroup,
      selectTemplate,
      t
    }
  }
}
</script>

<style scoped>
.hierarchical-selector {
  position: relative;
  display: inline-block;
}

.selector-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  height: 36px;
  padding: 0 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #2d3748;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  min-width: 200px;
}

.selector-button:hover {
  border-color: #cbd5e0;
  background: #f7fafc;
}

.selector-button.active {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.button-label {
  color: #718096;
  font-weight: 500;
}

.button-value {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  text-align: left;
  color: #2d3748;
}

.button-icon {
  color: #a0aec0;
  font-size: 0.75rem;
  transition: transform 0.2s;
}

.popover-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.1);
}

.popover-container {
  position: absolute;
  display: flex;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  max-height: 500px;
  z-index: 1001;
}

.popover-pane {
  display: flex;
  flex-direction: column;
  background: #ffffff;
}

.groups-pane {
  width: 280px;
  border-right: 1px solid #e2e8f0;
}

.templates-pane {
  width: 320px;
}

.pane-header {
  padding: 0.75rem 1rem;
  background: #f7fafc;
  border-bottom: 1px solid #e2e8f0;
  font-weight: 600;
  font-size: 0.85rem;
  color: #4a5568;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.pane-content {
  flex: 1;
  overflow-y: auto;
  max-height: 440px;
}

.group-item,
.template-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  cursor: pointer;
  transition: all 0.15s;
  border-bottom: 1px solid #f7fafc;
}

.group-item:hover,
.template-item:hover {
  background: #f7fafc;
}

.group-item.selected,
.template-item.selected {
  background: #edf2f7;
  border-left: 3px solid #667eea;
  padding-left: calc(1rem - 3px);
}

.group-item.active {
  background: #f7fafc;
}

.item-icon {
  font-size: 1.1rem;
  flex-shrink: 0;
}

.item-name {
  flex: 1;
  color: #2d3748;
  font-size: 0.9rem;
}

.item-arrow {
  color: #a0aec0;
  font-size: 1rem;
  margin-left: auto;
}

.group-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 0.5rem 0;
}

.template-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.template-name {
  color: #2d3748;
  font-size: 0.9rem;
  font-weight: 500;
}

.template-dates {
  color: #718096;
  font-size: 0.8rem;
}

.loading-state,
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 2rem;
  color: #718096;
  font-size: 0.9rem;
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid #e2e8f0;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Scrollbar styling */
.pane-content::-webkit-scrollbar {
  width: 6px;
}

.pane-content::-webkit-scrollbar-track {
  background: #f7fafc;
}

.pane-content::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 3px;
}

.pane-content::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* Responsive - mobile fallback */
@media (max-width: 768px) {
  .popover-container {
    left: 1rem !important;
    right: 1rem;
    max-width: calc(100vw - 2rem);
  }
  
  .groups-pane,
  .templates-pane {
    width: 100%;
    max-width: none;
  }
  
  .templates-pane {
    display: none;
  }
  
  .popover-container {
    flex-direction: column;
  }
}
</style>
