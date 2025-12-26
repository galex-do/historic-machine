<template>
  <div class="admin-panel">
    <div class="admin-header">
      <div class="admin-title">
        <h2>{{ t('adminTemplatesTitle') }}</h2>
        <p class="admin-subtitle">{{ t('adminTemplatesSubtitle') }}</p>
      </div>
      <div class="action-buttons" v-if="canAccessAdmin">
        <button @click="showCreateGroupModal = true" class="create-btn secondary">
          <span class="btn-icon">📁</span>
          {{ t('createNewGroup') }}
        </button>
        <button @click="showCreateTemplateModal = true" class="create-btn">
          <span class="btn-icon">➕</span>
          {{ t('createNewTemplate') }}
        </button>
      </div>
    </div>

    <div class="table-container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading templates...</p>
      </div>
      
      <div v-if="!loading && totalTemplates > 0" class="table-controls">
        <div class="table-filters">
          <input 
            v-model="searchQuery"
            type="text"
            :placeholder="t('searchTemplates')"
            class="search-input"
          />
        </div>
        <TablePagination 
          :current-page="currentPage"
          :page-size="pageSize"
          :total-items="filteredTotalTemplates"
          @update:current-page="handlePageChange"
          @update:page-size="handlePageSizeChange"
          :show-full-info="true"
        />
      </div>
      
      <table v-if="!loading" class="events-table">
        <thead>
          <tr>
            <th 
              class="sortable-header" 
              @click="toggleSort('name')"
              :class="{ 'active': sortField === 'name' }"
            >
              {{ t('columnName') }}
              <span class="sort-indicator">
                <span v-if="sortField === 'name'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th 
              class="sortable-header" 
              @click="toggleSort('group_name')"
              :class="{ 'active': sortField === 'group_name' }"
            >
              {{ t('columnGroup') }}
              <span class="sort-indicator">
                <span v-if="sortField === 'group_name'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th 
              class="sortable-header" 
              @click="toggleSort('start_date')"
              :class="{ 'active': sortField === 'start_date' }"
            >
              {{ t('columnDateFrom') }}
              <span class="sort-indicator">
                <span v-if="sortField === 'start_date'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th 
              class="sortable-header" 
              @click="toggleSort('end_date')"
              :class="{ 'active': sortField === 'end_date' }"
            >
              {{ t('columnDateTo') }}
              <span class="sort-indicator">
                <span v-if="sortField === 'end_date'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>{{ t('columnActions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="template in currentPageTemplates" :key="template.id" class="event-row">
            <td class="template-name">{{ template.name }}</td>
            <td class="template-group">
              <span class="group-badge">{{ template.group_name }}</span>
            </td>
            <td class="template-date">{{ template.start_display_date }}</td>
            <td class="template-date">{{ template.end_display_date }}</td>
            <td class="template-actions">
              <div class="actions-wrapper">
                <button @click="editTemplate(template)" class="action-btn edit-btn" title="Edit">
                  ✏️
                </button>
                <button @click="deleteTemplate(template)" class="action-btn delete-btn" title="Delete">
                  🗑️
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="!loading && totalTemplates === 0" class="empty-state">
        <p>No templates found. Create your first template!</p>
      </div>
      
      <div v-if="!loading && totalTemplates > 0 && filteredTemplates.length === 0" class="empty-state">
        <p>No templates match your search criteria.</p>
      </div>
    </div>

    <div class="groups-section" v-if="!loading">
      <h3>{{ t('templateGroupsTitle') }}</h3>
      <div class="groups-list">
        <div v-for="group in templateGroups" :key="group.id" class="group-card">
          <div class="group-info">
            <span class="group-name">{{ group.name }}</span>
            <span class="group-count">{{ getTemplatesInGroup(group.id) }} templates</span>
          </div>
          <div class="group-actions">
            <button @click="editGroup(group)" class="action-btn edit-btn" title="Edit">
              ✏️
            </button>
            <button @click="deleteGroup(group)" class="action-btn delete-btn" title="Delete" :disabled="getTemplatesInGroup(group.id) > 0">
              🗑️
            </button>
          </div>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="showCreateTemplateModal || showEditTemplateModal" class="modal-overlay" @click="closeTemplateModal">
        <div class="modal-content admin-modal wide-modal" @click.stop>
          <div class="modal-header">
            <h3>{{ editingTemplate ? t('editTemplate') : t('createNewTemplate') }}</h3>
            <button @click="closeTemplateModal" class="close-btn">&times;</button>
          </div>
          
          <div class="modal-body">
            <div v-if="error" class="error-message">{{ error }}</div>
            
            <form @submit.prevent="saveTemplate">
              <div class="form-row">
                <div class="form-group">
                  <label for="template-group">{{ t('templateGroup') }} *</label>
                  <select 
                    id="template-group"
                    v-model="templateForm.group_id" 
                    required 
                    class="form-select"
                  >
                    <option value="">{{ t('selectGroup') }}</option>
                    <option v-for="group in templateGroups" :key="group.id" :value="group.id">
                      {{ group.name }}
                    </option>
                  </select>
                </div>
              </div>
              
              <div class="locale-tabs">
                <button type="button" 
                        :class="['locale-tab', { active: activeTemplateLocale === 'en' }]" 
                        @click="activeTemplateLocale = 'en'">
                  English
                </button>
                <button type="button" 
                        :class="['locale-tab', { active: activeTemplateLocale === 'ru' }]" 
                        @click="activeTemplateLocale = 'ru'">
                  Русский
                </button>
              </div>

              <div v-show="activeTemplateLocale === 'en'" class="locale-content">
                <div class="form-group">
                  <label for="template-name-en">{{ t('columnName') }} (English) *</label>
                  <input 
                    id="template-name-en"
                    v-model="templateForm.name_en" 
                    type="text" 
                    required 
                    class="form-input"
                    placeholder="Name in English"
                  />
                </div>
                <div class="form-group">
                  <label for="template-desc-en">{{ t('columnDescription') }} (English)</label>
                  <textarea 
                    id="template-desc-en"
                    v-model="templateForm.description_en" 
                    class="form-textarea"
                    placeholder="Description in English"
                    rows="2"
                  ></textarea>
                </div>
              </div>

              <div v-show="activeTemplateLocale === 'ru'" class="locale-content">
                <div class="form-group">
                  <label for="template-name-ru">{{ t('columnName') }} (Русский)</label>
                  <input 
                    id="template-name-ru"
                    v-model="templateForm.name_ru" 
                    type="text" 
                    class="form-input"
                    placeholder="Название на русском"
                  />
                </div>
                <div class="form-group">
                  <label for="template-desc-ru">{{ t('columnDescription') }} (Русский)</label>
                  <textarea 
                    id="template-desc-ru"
                    v-model="templateForm.description_ru" 
                    class="form-textarea"
                    placeholder="Описание на русском"
                    rows="2"
                  ></textarea>
                </div>
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="start-date">{{ t('startDate') }} *</label>
                  <div class="date-era-input">
                    <input 
                      id="start-date"
                      v-model="templateForm.start_date" 
                      type="date" 
                      required 
                      class="form-input"
                    />
                    <select v-model="templateForm.start_era" class="era-select">
                      <option value="BC">BC</option>
                      <option value="AD">AD</option>
                    </select>
                  </div>
                </div>
                <div class="form-group">
                  <label for="end-date">{{ t('endDate') }} *</label>
                  <div class="date-era-input">
                    <input 
                      id="end-date"
                      v-model="templateForm.end_date" 
                      type="date" 
                      required 
                      class="form-input"
                    />
                    <select v-model="templateForm.end_era" class="era-select">
                      <option value="BC">BC</option>
                      <option value="AD">AD</option>
                    </select>
                  </div>
                </div>
              </div>
              
              <div class="form-actions">
                <button type="button" @click="closeTemplateModal" class="cancel-btn">{{ t('cancel') }}</button>
                <button type="submit" class="submit-btn" :disabled="localLoading">
                  {{ localLoading ? t('saving') : (editingTemplate ? t('updateTemplate') : t('createTemplate')) }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showCreateGroupModal || showEditGroupModal" class="modal-overlay" @click="closeGroupModal">
        <div class="modal-content admin-modal" @click.stop>
          <div class="modal-header">
            <h3>{{ editingGroup ? t('editGroup') : t('createNewGroup') }}</h3>
            <button @click="closeGroupModal" class="close-btn">&times;</button>
          </div>
          
          <div class="modal-body">
            <div v-if="error" class="error-message">{{ error }}</div>
            
            <form @submit.prevent="saveGroup">
              <div class="locale-tabs">
                <button type="button" 
                        :class="['locale-tab', { active: activeGroupLocale === 'en' }]" 
                        @click="activeGroupLocale = 'en'">
                  English
                </button>
                <button type="button" 
                        :class="['locale-tab', { active: activeGroupLocale === 'ru' }]" 
                        @click="activeGroupLocale = 'ru'">
                  Русский
                </button>
              </div>

              <div v-show="activeGroupLocale === 'en'" class="locale-content">
                <div class="form-group">
                  <label for="group-name-en">{{ t('columnName') }} (English) *</label>
                  <input 
                    id="group-name-en"
                    v-model="groupForm.name_en" 
                    type="text" 
                    required 
                    class="form-input"
                    placeholder="Name in English"
                  />
                </div>
                <div class="form-group">
                  <label for="group-desc-en">{{ t('columnDescription') }} (English)</label>
                  <textarea 
                    id="group-desc-en"
                    v-model="groupForm.description_en" 
                    class="form-textarea"
                    placeholder="Description in English"
                    rows="2"
                  ></textarea>
                </div>
              </div>

              <div v-show="activeGroupLocale === 'ru'" class="locale-content">
                <div class="form-group">
                  <label for="group-name-ru">{{ t('columnName') }} (Русский)</label>
                  <input 
                    id="group-name-ru"
                    v-model="groupForm.name_ru" 
                    type="text" 
                    class="form-input"
                    placeholder="Название на русском"
                  />
                </div>
                <div class="form-group">
                  <label for="group-desc-ru">{{ t('columnDescription') }} (Русский)</label>
                  <textarea 
                    id="group-desc-ru"
                    v-model="groupForm.description_ru" 
                    class="form-textarea"
                    placeholder="Описание на русском"
                    rows="2"
                  ></textarea>
                </div>
              </div>
              
              <div class="form-actions">
                <button type="button" @click="closeGroupModal" class="cancel-btn">{{ t('cancel') }}</button>
                <button type="submit" class="submit-btn" :disabled="localLoading">
                  {{ localLoading ? t('saving') : (editingGroup ? t('updateGroup') : t('createGroup')) }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuth } from '@/composables/useAuth.js'
import { useLocale } from '@/composables/useLocale.js'
import apiService from '@/services/api.js'
import TablePagination from '@/components/TablePagination.vue'

export default {
  name: 'AdminTemplates',
  components: {
    TablePagination
  },
  setup() {
    const { canAccessAdmin } = useAuth()
    const { t } = useLocale()
    
    const templates = ref([])
    const templateGroups = ref([])
    const loading = ref(false)
    const localLoading = ref(false)
    const error = ref(null)
    const searchQuery = ref('')
    
    const showCreateTemplateModal = ref(false)
    const showEditTemplateModal = ref(false)
    const editingTemplate = ref(null)
    
    const showCreateGroupModal = ref(false)
    const showEditGroupModal = ref(false)
    const editingGroup = ref(null)
    
    const activeTemplateLocale = ref('en')
    const activeGroupLocale = ref('en')
    
    const sortField = ref('start_date')
    const sortDirection = ref('asc')
    
    const currentPage = ref(1)
    const pageSize = ref(10)
    
    const templateForm = ref({
      group_id: '',
      name_en: '',
      name_ru: '',
      description_en: '',
      description_ru: '',
      start_date: '',
      start_era: 'BC',
      end_date: '',
      end_era: 'BC'
    })
    
    const groupForm = ref({
      name_en: '',
      name_ru: '',
      description_en: '',
      description_ru: ''
    })

    const getChronologicalValue = (dateString, era) => {
      if (!dateString) return 0
      const dateParts = dateString.split('-')
      const year = parseInt(dateParts[0], 10) || 0
      const month = parseInt(dateParts[1], 10) || 1
      const day = parseInt(dateParts[2], 10) || 1
      
      if (era === 'BC') {
        return -(year - (month / 12) - (day / 365))
      } else {
        return year + (month / 12) + (day / 365)
      }
    }

    const filteredTemplates = computed(() => {
      if (!templates.value || !Array.isArray(templates.value)) {
        return []
      }
      
      let filtered = templates.value
      
      if (searchQuery.value.trim()) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(template => 
          template.name?.toLowerCase().includes(query) ||
          template.name_en?.toLowerCase().includes(query) ||
          template.name_ru?.toLowerCase().includes(query) ||
          template.group_name?.toLowerCase().includes(query)
        )
      }
      
      return [...filtered].sort((a, b) => {
        let aValue, bValue
        
        if (sortField.value === 'start_date') {
          aValue = getChronologicalValue(a.start_date, a.start_era)
          bValue = getChronologicalValue(b.start_date, b.start_era)
        } else if (sortField.value === 'end_date') {
          aValue = getChronologicalValue(a.end_date, a.end_era)
          bValue = getChronologicalValue(b.end_date, b.end_era)
        } else if (sortField.value === 'group_name') {
          aValue = a.group_name?.toLowerCase() || ''
          bValue = b.group_name?.toLowerCase() || ''
        } else {
          aValue = a[sortField.value]?.toLowerCase() || ''
          bValue = b[sortField.value]?.toLowerCase() || ''
        }
        
        if (sortDirection.value === 'asc') {
          return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
        } else {
          return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
        }
      })
    })
    
    const filteredTotalTemplates = computed(() => filteredTemplates.value.length)
    const totalTemplates = computed(() => templates.value?.length || 0)
    
    const currentPageTemplates = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filteredTemplates.value.slice(start, end)
    })

    const getTemplatesInGroup = (groupId) => {
      return templates.value.filter(t => t.group_id === groupId).length
    }

    const toggleSort = (field) => {
      if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
      } else {
        sortField.value = field
        sortDirection.value = 'asc'
      }
    }

    const handlePageChange = (page) => {
      currentPage.value = page
    }

    const handlePageSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1
    }

    const loadData = async () => {
      loading.value = true
      error.value = null
      try {
        const [templatesData, groupsData] = await Promise.all([
          apiService.getAllTemplates(),
          apiService.getTemplateGroups()
        ])
        templates.value = templatesData || []
        templateGroups.value = groupsData || []
      } catch (err) {
        console.error('Error loading data:', err)
        error.value = 'Failed to load data'
      } finally {
        loading.value = false
      }
    }

    const editTemplate = (template) => {
      editingTemplate.value = template
      templateForm.value = {
        group_id: template.group_id,
        name_en: template.name_en || template.name || '',
        name_ru: template.name_ru || '',
        description_en: template.description_en || template.description || '',
        description_ru: template.description_ru || '',
        start_date: template.start_date,
        start_era: template.start_era || 'BC',
        end_date: template.end_date,
        end_era: template.end_era || 'BC'
      }
      showEditTemplateModal.value = true
    }

    const deleteTemplate = async (template) => {
      const confirmed = confirm(`Are you sure you want to delete the template "${template.name}"? This action cannot be undone.`)
      if (!confirmed) return
      
      localLoading.value = true
      error.value = null
      try {
        await apiService.deleteTemplate(template.id)
        await loadData()
      } catch (err) {
        console.error('Error deleting template:', err)
        error.value = err.message || 'Failed to delete template'
      } finally {
        localLoading.value = false
      }
    }

    const saveTemplate = async () => {
      localLoading.value = true
      error.value = null
      
      try {
        const data = {
          group_id: parseInt(templateForm.value.group_id),
          name_en: templateForm.value.name_en,
          name_ru: templateForm.value.name_ru || templateForm.value.name_en,
          description_en: templateForm.value.description_en,
          description_ru: templateForm.value.description_ru || templateForm.value.description_en,
          start_date: templateForm.value.start_date,
          start_era: templateForm.value.start_era,
          end_date: templateForm.value.end_date,
          end_era: templateForm.value.end_era
        }
        
        if (editingTemplate.value) {
          await apiService.updateTemplate(editingTemplate.value.id, data)
        } else {
          await apiService.createTemplate(data)
        }
        
        await loadData()
        closeTemplateModal()
      } catch (err) {
        console.error('Error saving template:', err)
        error.value = err.message || 'Failed to save template'
      } finally {
        localLoading.value = false
      }
    }

    const closeTemplateModal = () => {
      showCreateTemplateModal.value = false
      showEditTemplateModal.value = false
      editingTemplate.value = null
      templateForm.value = {
        group_id: '',
        name_en: '',
        name_ru: '',
        description_en: '',
        description_ru: '',
        start_date: '',
        start_era: 'BC',
        end_date: '',
        end_era: 'BC'
      }
      activeTemplateLocale.value = 'en'
      error.value = null
    }

    const editGroup = (group) => {
      editingGroup.value = group
      groupForm.value = {
        name_en: group.name_en || group.name || '',
        name_ru: group.name_ru || '',
        description_en: group.description_en || group.description || '',
        description_ru: group.description_ru || ''
      }
      showEditGroupModal.value = true
    }

    const deleteGroup = async (group) => {
      const templateCount = getTemplatesInGroup(group.id)
      if (templateCount > 0) {
        alert(`Cannot delete group "${group.name}" because it has ${templateCount} template(s). Delete templates first.`)
        return
      }

      const confirmed = confirm(`Are you sure you want to delete the group "${group.name}"? This action cannot be undone.`)
      if (!confirmed) return
      
      localLoading.value = true
      error.value = null
      try {
        await apiService.deleteTemplateGroup(group.id)
        await loadData()
      } catch (err) {
        console.error('Error deleting group:', err)
        error.value = err.message || 'Failed to delete group'
      } finally {
        localLoading.value = false
      }
    }

    const saveGroup = async () => {
      localLoading.value = true
      error.value = null
      
      try {
        const data = {
          name_en: groupForm.value.name_en,
          name_ru: groupForm.value.name_ru || groupForm.value.name_en,
          description_en: groupForm.value.description_en,
          description_ru: groupForm.value.description_ru || groupForm.value.description_en
        }
        
        if (editingGroup.value) {
          await apiService.updateTemplateGroup(editingGroup.value.id, data)
        } else {
          await apiService.createTemplateGroup(data)
        }
        
        await loadData()
        closeGroupModal()
      } catch (err) {
        console.error('Error saving group:', err)
        error.value = err.message || 'Failed to save group'
      } finally {
        localLoading.value = false
      }
    }

    const closeGroupModal = () => {
      showCreateGroupModal.value = false
      showEditGroupModal.value = false
      editingGroup.value = null
      groupForm.value = {
        name_en: '',
        name_ru: '',
        description_en: '',
        description_ru: ''
      }
      activeGroupLocale.value = 'en'
      error.value = null
    }

    onMounted(async () => {
      await loadData()
    })

    return {
      canAccessAdmin,
      loading,
      error,
      templates,
      templateGroups,
      searchQuery,
      sortField,
      sortDirection,
      currentPage,
      pageSize,
      totalTemplates,
      filteredTemplates,
      filteredTotalTemplates,
      currentPageTemplates,
      showCreateTemplateModal,
      showEditTemplateModal,
      editingTemplate,
      templateForm,
      showCreateGroupModal,
      showEditGroupModal,
      editingGroup,
      groupForm,
      activeTemplateLocale,
      activeGroupLocale,
      localLoading,
      toggleSort,
      handlePageChange,
      handlePageSizeChange,
      getTemplatesInGroup,
      editTemplate,
      deleteTemplate,
      saveTemplate,
      closeTemplateModal,
      editGroup,
      deleteGroup,
      saveGroup,
      closeGroupModal,
      t
    }
  }
}
</script>

<style scoped>
.admin-panel {
  padding: 1.25rem 2rem;
  max-width: 100%;
  margin: 0;
  background-color: #f8fafc;
  min-height: calc(100vh - 4rem);
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.admin-title h2 {
  margin: 0;
  color: #2d3748;
  font-size: 2rem;
  font-weight: 700;
}

.admin-subtitle {
  margin: 0.5rem 0 0 0;
  color: #718096;
  font-size: 1rem;
  font-weight: 400;
}

.action-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.create-btn.secondary {
  background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%);
}

.create-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.btn-icon {
  font-size: 1rem;
}

.table-container {
  background: white;
  border-radius: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin-bottom: 2rem;
}

.table-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  flex-wrap: wrap;
  gap: 1rem;
}

.search-input {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  width: 300px;
  max-width: 100%;
}

.search-input:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.events-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.events-table th,
.events-table td {
  padding: 1rem 1.5rem;
  text-align: left;
}

.events-table th {
  background: #f8fafc;
  font-weight: 600;
  color: #374151;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.sortable-header {
  cursor: pointer;
  user-select: none;
  position: relative;
  transition: background-color 0.2s ease;
}

.sortable-header:hover {
  background: #f1f5f9;
}

.sortable-header.active {
  background: #e2e8f0;
}

.sort-indicator {
  margin-left: 0.5rem;
  font-size: 0.8rem;
}

.sort-arrow {
  color: #4f46e5;
  font-weight: bold;
}

.sort-placeholder {
  color: #9ca3af;
}

.event-row {
  border-bottom: 1px solid #f1f5f9;
}

.event-row:hover {
  background: #f8fafc;
}

.template-name {
  font-weight: 600;
  color: #2d3748;
}

.template-group {
  max-width: 200px;
}

.group-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  background: #e0e7ff;
  color: #4338ca;
  border-radius: 0.25rem;
  font-size: 0.8rem;
  font-weight: 500;
}

.template-date {
  color: #4a5568;
  font-size: 0.875rem;
  white-space: nowrap;
}

.template-actions {
  vertical-align: middle;
  text-align: center;
}

.template-actions .actions-wrapper {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  justify-content: center;
}

.action-btn {
  padding: 0.5rem;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.edit-btn {
  background: #fef3c7;
  color: #d97706;
}

.edit-btn:hover:not(:disabled) {
  background: #fde68a;
  transform: scale(1.1);
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
}

.delete-btn:hover:not(:disabled) {
  background: #fecaca;
  transform: scale(1.1);
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 1.5rem;
  color: #6b7280;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid #e5e7eb;
  border-top: 3px solid #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.groups-section {
  background: white;
  border-radius: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

.groups-section h3 {
  margin: 0 0 1rem 0;
  color: #2d3748;
  font-size: 1.25rem;
}

.groups-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.group-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
}

.group-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.group-name {
  font-weight: 600;
  color: #2d3748;
}

.group-count {
  font-size: 0.8rem;
  color: #6b7280;
}

.group-actions {
  display: flex;
  gap: 0.5rem;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100000;
}

.modal-content {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  max-height: 90vh;
  overflow-y: auto;
}

.admin-modal {
  max-width: 500px;
  width: 90%;
}

.admin-modal.wide-modal {
  max-width: 700px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 1.5rem 0;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 1.5rem;
}

.modal-header h3 {
  margin: 0;
  color: #2d3748;
  font-size: 1.25rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #9ca3af;
  padding: 0;
}

.close-btn:hover {
  color: #6b7280;
}

.modal-body {
  padding: 0 1.5rem 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.date-era-input {
  display: flex;
  gap: 0.5rem;
}

.date-era-input .form-input {
  flex: 1;
}

.era-select {
  width: 70px;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  background: white;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.cancel-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  background: white;
  color: #374151;
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.cancel-btn:hover {
  background: #f9fafb;
}

.submit-btn {
  padding: 0.5rem 1rem;
  border: none;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
}

.submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-message {
  padding: 0.75rem;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.375rem;
  color: #dc2626;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .admin-modal.wide-modal {
    max-width: 95%;
  }
}

.locale-tabs {
  display: flex;
  gap: 0.25rem;
  margin-bottom: 1rem;
  border-bottom: 2px solid #e2e8f0;
}

.locale-tab {
  padding: 0.5rem 1rem;
  border: none;
  background: transparent;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  position: relative;
  transition: color 0.2s;
}

.locale-tab:hover {
  color: #4f46e5;
}

.locale-tab.active {
  color: #4f46e5;
}

.locale-tab.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  right: 0;
  height: 2px;
  background: #4f46e5;
}

.locale-content {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
