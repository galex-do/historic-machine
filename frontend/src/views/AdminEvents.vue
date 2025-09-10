<template>
  <div class="admin-panel">
    <div class="admin-header">
      <div class="admin-title">
        <h2>Historic Events</h2>
        <p class="admin-subtitle">Manage and create historic events</p>
      </div>
      <div class="action-buttons" v-if="canAccessAdmin">
        <button @click="showCreateModal = true" class="create-btn">
          <span class="btn-icon">➕</span>
          Create New Event
        </button>
        <button @click="triggerFileUpload" class="import-btn">
          <span class="btn-icon">📁</span>
          Import Events
        </button>
        <input 
          ref="fileInput" 
          type="file" 
          accept=".json" 
          @change="handleFileUpload" 
          style="display: none"
        />
      </div>
    </div>

    <!-- Events Table -->
    <div class="table-container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading events...</p>
      </div>
      
      <!-- Table Controls & Pagination -->
      <div v-if="!loading && totalEvents > 0" class="table-controls">
        <div class="table-filters">
          <EventTypeFilter
            :selected-lens-type="selectedLensType"
            :show-dropdown="showLensDropdown"
            @toggle-dropdown="toggleLensDropdown"
            @lens-type-changed="handleLensTypeChange"
          />
        </div>
        <TablePagination 
          :current-page="currentPage"
          :page-size="pageSize"
          :total-items="filteredTotalEvents"
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
              Name
              <span class="sort-indicator">
                <span v-if="sortField === 'name'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>Description</th>
            <th 
              class="sortable-header" 
              @click="toggleSort('date')"
              :class="{ 'active': sortField === 'date' }"
            >
              Date
              <span class="sort-indicator">
                <span v-if="sortField === 'date'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>Location</th>
            <th>Type</th>
            <th>Tags</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in currentPageEvents" :key="event.id" class="event-row">
            <td class="event-name">{{ event.name }}</td>
            <td class="event-description">
              <span class="description-text" :title="event.description">
                {{ event.description.length > 100 ? event.description.substring(0, 100) + '...' : event.description }}
              </span>
            </td>
            <td class="event-date">{{ formatDateWithEra(event.event_date, event.era) }}</td>
            <td class="event-location">
              <span v-if="event.latitude && event.longitude">
                {{ event.latitude.toFixed(2) }}, {{ event.longitude.toFixed(2) }}
              </span>
              <span v-else class="no-location">No location</span>
            </td>
            <td class="event-type">
              <span class="type-badge" :class="event.lens_type">{{ event.lens_type }}</span>
            </td>
            <td class="event-tags">
              <div class="tags-container">
                <span 
                  v-for="tag in event.tags" 
                  :key="tag.id" 
                  class="tag-badge"
                  :style="{ backgroundColor: tag.color, color: getContrastColor(tag.color) }"
                >
                  {{ tag.name }}
                </span>
                <span v-if="!event.tags || event.tags.length === 0" class="no-tags">No tags</span>
              </div>
            </td>
            <td class="event-actions">
              <button @click="editEvent(event)" class="action-btn edit-btn" title="Edit">
                ✏️
              </button>
              <button @click="deleteEvent(event)" class="action-btn delete-btn" title="Delete">
                🗑️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="!loading && totalEvents === 0" class="empty-state">
        <p>No events found. Create your first historical event!</p>
      </div>
      
      <div v-if="!loading && totalEvents > 0 && filteredEvents.length === 0" class="empty-state">
        <p>No events match the selected filters. Try adjusting your event type selection.</p>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content admin-modal" @click.stop>
          <div class="modal-header">
            <h3>{{ editingEvent ? 'Edit Event' : 'Create New Event' }}</h3>
            <button @click="closeModal" class="close-btn">&times;</button>
          </div>
          
          <div class="modal-body">
            <div v-if="error" class="error-message">{{ error }}</div>
            
            <form @submit.prevent="saveEvent">
              <div class="form-row">
                <div class="form-group">
                  <label for="event-name">Event Name *</label>
                  <input 
                    id="event-name"
                    v-model="eventForm.name" 
                    type="text" 
                    required 
                    class="form-input"
                    placeholder="Enter event name"
                  />
                </div>
                <div class="form-group">
                  <label for="event-type">Type *</label>
                  <select id="event-type" v-model="eventForm.lens_type" required class="form-input">
                    <option value="">Select type</option>
                    <option value="military">Military</option>
                    <option value="political">Political</option>
                    <option value="historic">Historic</option>
                    <option value="cultural">Cultural</option>
                    <option value="scientific">Scientific</option>
                  </select>
                </div>
              </div>
              
              <div class="form-group">
                <label for="event-description">Description *</label>
                <textarea 
                  id="event-description"
                  v-model="eventForm.description" 
                  required 
                  class="form-textarea"
                  placeholder="Enter event description"
                  rows="4"
                ></textarea>
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="event-date">Date *</label>
                  <input 
                    id="event-date"
                    v-model="eventForm.date_display" 
                    type="text" 
                    required 
                    class="form-input"
                    placeholder="DD/MM/YYYY"
                    @input="updateEventDate"
                  />
                  <small class="form-hint">Format: DD/MM/YYYY (e.g., 15/03/1066)</small>
                </div>
                <div class="form-group">
                  <label for="event-era">Era *</label>
                  <select 
                    id="event-era"
                    v-model="eventForm.era" 
                    required
                    class="form-input"
                  >
                    <option value="BC">BC (Before Christ)</option>
                    <option value="AD">AD (Anno Domini)</option>
                  </select>
                  <small class="form-hint">Select BC for dates before year 1, AD for dates after</small>
                </div>
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="event-latitude">Latitude</label>
                  <input 
                    id="event-latitude"
                    v-model.number="eventForm.latitude" 
                    type="number" 
                    step="any"
                    class="form-input"
                    placeholder="e.g., 51.5074"
                  />
                </div>
                <div class="form-group">
                  <label for="event-longitude">Longitude</label>
                  <input 
                    id="event-longitude"
                    v-model.number="eventForm.longitude" 
                    type="number" 
                    step="any"
                    class="form-input"
                    placeholder="e.g., -0.1278"
                  />
                </div>
              </div>
              
              <!-- Tags Section -->
              <div class="form-group">
                <label for="event-tags">Tags</label>
                <div class="tags-section">
                  <!-- Selected Tags Display -->
                  <div v-if="eventForm.selectedTags.length > 0" class="selected-tags">
                    <div 
                      v-for="tag in eventForm.selectedTags" 
                      :key="tag.id"
                      class="selected-tag"
                      :style="{ backgroundColor: tag.color, color: getContrastColor(tag.color) }"
                    >
                      {{ tag.name }}
                      <button 
                        type="button" 
                        @click="removeTag(tag)"
                        class="remove-tag-btn"
                        :style="{ color: getContrastColor(tag.color) }"
                      >
                        ×
                      </button>
                    </div>
                  </div>
                  
                  <!-- Tag Search and Add -->
                  <div class="tag-input-section">
                    <input 
                      v-model="eventForm.tagSearch"
                      type="text"
                      class="tag-search-input"
                      placeholder="Search tags or type new tag name..."
                      @keydown.enter.prevent="handleTagInput"
                    />
                    
                    <!-- Tag Suggestions -->
                    <div v-if="filteredTags.length > 0 && eventForm.tagSearch" class="tag-suggestions">
                      <div 
                        v-for="tag in filteredTags.slice(0, 5)" 
                        :key="tag.id"
                        class="tag-suggestion"
                        @click="addTag(tag)"
                      >
                        <span 
                          class="tag-color-indicator"
                          :style="{ backgroundColor: tag.color }"
                        ></span>
                        {{ tag.name }}
                      </div>
                      <!-- Create new tag option at bottom of suggestions -->
                      <div v-if="canCreateNewTag" class="tag-suggestion create-new-suggestion" @click="createAndAddTag">
                        <span class="tag-color-indicator" :style="{ backgroundColor: eventForm.newTagColor }"></span>
                        Create "{{ eventForm.tagSearch }}"
                        <input 
                          v-model="eventForm.newTagColor"
                          type="color"
                          class="inline-color-picker"
                          @click.stop
                        />
                      </div>
                    </div>
                    
                    <!-- Create New Tag Option -->
                    <div v-if="canCreateNewTag && !filteredTags.length" class="new-tag-option">
                      <div class="new-tag-form">
                        <span class="new-tag-text">Create new tag:</span>
                        <input 
                          v-model="eventForm.newTagColor"
                          type="color"
                          class="color-picker"
                        />
                        <button 
                          type="button"
                          @click="createAndAddTag"
                          class="create-tag-btn"
                        >
                          Add "{{ eventForm.tagSearch }}"
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="form-actions">
                <button type="button" @click="closeModal" class="cancel-btn">Cancel</button>
                <button type="submit" class="submit-btn" :disabled="loading">
                  {{ loading ? 'Saving...' : (editingEvent ? 'Update Event' : 'Create Event') }}
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
import { useTags } from '@/composables/useTags.js'
import { useEvents } from '@/composables/useEvents.js'
import { getAvailableLensTypes } from '@/utils/event-utils.js'
import apiService from '@/services/api.js'
import TablePagination from '@/components/TablePagination.vue'
import EventTypeFilter from '@/components/filters/EventTypeFilter.vue'

export default {
  name: 'AdminEvents',
  components: {
    TablePagination,
    EventTypeFilter
  },
  setup() {
    const { canAccessAdmin } = useAuth()
    const { allTags, loadTags, createTag, setEventTags, getTagsByIds } = useTags()
    const { events, fetchEvents, loading, error, handleEventDeleted } = useEvents()
    
    // Local filter state for admin panel (single-select)
    const selectedLensType = ref('all')
    const showLensDropdown = ref(false)
    
    // Filter methods
    const toggleLensDropdown = () => {
      showLensDropdown.value = !showLensDropdown.value
    }
    
    const handleLensTypeChange = (newLensType) => {
      selectedLensType.value = newLensType
      currentPage.value = 1 // Reset to first page when filter changes
    }
    
    const localLoading = ref(false)
    const localError = ref(null)
    
    const showCreateModal = ref(false)
    const showEditModal = ref(false)
    const editingEvent = ref(null)
    const fileInput = ref(null)
    
    // Sorting state
    const sortField = ref('date')
    const sortDirection = ref('asc')
    
    // Pagination state
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalEvents = ref(0)
    const displayedEvents = ref([])
    const paginatedLoading = ref(false)
    
    // All events (for filtering before pagination)
    const allEvents = ref([])
    
    // Filtered and sorted events based on lens type selection and current sort settings
    const filteredEvents = computed(() => {
      if (!allEvents.value || !Array.isArray(allEvents.value) || allEvents.value.length === 0) {
        return []
      }
      
      let filtered = allEvents.value
      
      // Apply lens type filter
      if (selectedLensType.value !== 'all') {
        filtered = filtered.filter(event => 
          event && event.lens_type && event.lens_type === selectedLensType.value
        )
      }
      
      // Apply sorting
      return [...filtered].sort((a, b) => {
        let aValue, bValue
        
        if (sortField.value === 'date') {
          // Use proper astronomical year for chronological sorting
          const getAstronomicalYear = (dateString) => {
            if (dateString.startsWith('-')) {
              // BC date format: "-754-04-21T00:00:00Z" 
              const parts = dateString.substring(1).split('T')[0].split('-')
              const year = parseInt(parts[0], 10)
              const month = parseInt(parts[1], 10)
              const day = parseInt(parts[2], 10)
              
              // Convert to astronomical year: 754 BC = -753 astronomical year
              // Add month/day for sub-year precision
              return -(year - 1) - (month / 12) - (day / 365)
            } else {
              // AD date: use normal date parsing
              const date = new Date(dateString)
              return date.getFullYear() + (date.getMonth() / 12) + (date.getDate() / 365)
            }
          }
          
          aValue = getAstronomicalYear(a.event_date || a.date || '0')
          bValue = getAstronomicalYear(b.event_date || b.date || '0')
        } else {
          aValue = a[sortField.value] || ''
          bValue = b[sortField.value] || ''
        }
        
        if (sortDirection.value === 'asc') {
          return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
        } else {
          return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
        }
      })
    })
    
    // Update total events to reflect filtered count
    const filteredTotalEvents = computed(() => filteredEvents.value.length)
    
    // Current page of filtered events
    const currentPageEvents = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filteredEvents.value.slice(start, end)
    })
    
    const eventForm = ref({
      name: '',
      description: '',
      date_display: '',
      date: '',
      era: 'AD',
      latitude: null,
      longitude: null,
      lens_type: '',
      selectedTags: [],
      newTagName: '',
      newTagColor: '#3B82F6',
      tagSearch: ''
    })

    // All the helper functions from AdminPanel.vue would go here
    // (formatDate, formatDateWithEra, getContrastColor, etc.)
    
    const formatDate = (dateString) => {
      if (!dateString) return 'Invalid Date'
      
      try {
        let year, month, day
        
        if (dateString.startsWith('-')) {
          // BC date format: "-491-09-12T00:00:00Z"
          const parts = dateString.substring(1).split('T')[0].split('-')
          year = parseInt(parts[0], 10)
          month = parseInt(parts[1], 10)
          day = parseInt(parts[2], 10)
        } else {
          // AD date format: "1453-05-29T00:00:00Z"
          const date = new Date(dateString)
          if (isNaN(date.getTime())) return 'Invalid Date'
          
          return date.toLocaleDateString('en-GB', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
          })
        }
        
        // For BC dates, format manually
        const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 
                           'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
        return `${day} ${monthNames[month - 1]} ${year}`
      } catch {
        return 'Invalid Date'
      }
    }

    const formatDateWithEra = (dateString, era) => {
      try {
        const formattedDate = formatDate(dateString)
        if (formattedDate === 'Invalid Date') return 'Invalid Date'
        return `${formattedDate} ${era || 'AD'}`
      } catch {
        return 'Invalid Date'
      }
    }

    const getContrastColor = (hexColor) => {
      if (!hexColor) return '#000000'
      
      // Remove # if present
      const color = hexColor.replace('#', '')
      
      // Convert to RGB
      const r = parseInt(color.substr(0, 2), 16)
      const g = parseInt(color.substr(2, 2), 16)
      const b = parseInt(color.substr(4, 2), 16)
      
      // Calculate luminance
      const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
      
      // Return black for light colors, white for dark colors
      return luminance > 0.5 ? '#000000' : '#ffffff'
    }

    // Event management functions
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

    const editEvent = (event) => {
      editingEvent.value = event
      eventForm.value = {
        name: event.name,
        description: event.description,
        date_display: event.display_date || '',
        date: event.event_date,
        era: event.era || 'AD',
        latitude: event.latitude,
        longitude: event.longitude,
        lens_type: event.lens_type,
        selectedTags: event.tags || [],
        newTagName: '',
        newTagColor: '#3B82F6',
        tagSearch: ''
      }
      showEditModal.value = true
    }

    const deleteEvent = async (event) => {
      const confirmed = confirm(`Are you sure you want to delete "${event.name}"? This action cannot be undone.`)
      if (!confirmed) return
      
      localLoading.value = true
      localError.value = null
      try {
        await apiService.deleteEvent(event.id)
        await fetchEvents() // Refresh events list
        allEvents.value = events.value || []
        console.log('Event deleted successfully')
      } catch (err) {
        console.error('Error deleting event:', err)
        localError.value = err.message || 'Failed to delete event'
      } finally {
        localLoading.value = false
      }
    }

    const closeModal = () => {
      showCreateModal.value = false
      showEditModal.value = false
      editingEvent.value = null
      eventForm.value = {
        name: '',
        description: '',
        date_display: '',
        date: '',
        era: 'AD',
        latitude: null,
        longitude: null,
        lens_type: '',
        selectedTags: [],
        newTagName: '',
        newTagColor: '#3B82F6',
        tagSearch: ''
      }
      localError.value = null
    }

    const triggerFileUpload = () => {
      fileInput.value?.click()
    }

    // Tag management functions and computed properties
    const filteredTags = computed(() => {
      if (!allTags.value || !eventForm.value.tagSearch) return []
      const searchTerm = eventForm.value.tagSearch.toLowerCase()
      return allTags.value.filter(tag => 
        tag.name.toLowerCase().includes(searchTerm) &&
        !eventForm.value.selectedTags.some(selected => selected.id === tag.id)
      )
    })

    const canCreateNewTag = computed(() => {
      if (!eventForm.value.tagSearch) return false
      const searchTerm = eventForm.value.tagSearch.toLowerCase()
      return !allTags.value.some(tag => tag.name.toLowerCase() === searchTerm)
    })

    const addTag = (tag) => {
      if (!eventForm.value.selectedTags.some(selected => selected.id === tag.id)) {
        eventForm.value.selectedTags.push(tag)
      }
      eventForm.value.tagSearch = ''
    }

    const removeTag = (tagToRemove) => {
      eventForm.value.selectedTags = eventForm.value.selectedTags.filter(
        tag => tag.id !== tagToRemove.id
      )
    }

    const handleTagInput = () => {
      if (eventForm.value.tagSearch.trim()) {
        if (filteredTags.value.length > 0) {
          addTag(filteredTags.value[0])
        } else if (canCreateNewTag.value) {
          createAndAddTag()
        }
      }
    }

    const createAndAddTag = async () => {
      if (!eventForm.value.tagSearch.trim()) return
      
      try {
        const newTag = await apiService.createTag({
          name: eventForm.value.tagSearch.trim(),
          description: `Auto-generated tag for ${eventForm.value.tagSearch.trim()}`,
          color: eventForm.value.newTagColor
        })
        
        // Add to allTags and selected tags
        allTags.value.push(newTag)
        eventForm.value.selectedTags.push(newTag)
        eventForm.value.tagSearch = ''
        eventForm.value.newTagColor = '#3B82F6' // Reset color
      } catch (err) {
        console.error('Error creating tag:', err)
        localError.value = err.message || 'Failed to create tag'
      }
    }

    // Load initial data
    onMounted(async () => {
      try {
        await fetchEvents()
        allEvents.value = events.value || []
        totalEvents.value = allEvents.value.length
        await loadTags()
      } catch (err) {
        console.error('Error loading admin events data:', err)
        localError.value = 'Failed to load data'
      }
    })

    return {
      canAccessAdmin,
      loading: computed(() => loading.value || localLoading.value),
      error: computed(() => error.value || localError.value),
      showCreateModal,
      showEditModal,
      editingEvent,
      eventForm,
      fileInput,
      formatDateWithEra,
      getContrastColor,
      filteredEvents,
      currentPageEvents,
      filteredTotalEvents,
      handleLensTypeChange,
      toggleLensDropdown,
      selectedLensType,
      showLensDropdown,
      sortField,
      sortDirection,
      currentPage,
      pageSize,
      totalEvents,
      // Tag management
      filteredTags,
      canCreateNewTag,
      // Functions
      toggleSort,
      handlePageChange,
      handlePageSizeChange,
      editEvent,
      deleteEvent,
      closeModal,
      triggerFileUpload,
      // Tag functions
      addTag,
      removeTag,
      handleTagInput,
      createAndAddTag
    }
  }
}
</script>

<style scoped>
/* Use the same styles as AdminPanel.vue */
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

.create-btn,
.import-btn {
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
}

.create-btn {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.create-btn:hover {
  background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.import-btn {
  background: linear-gradient(135deg, #059669 0%, #0d9488 100%);
  color: white;
}

.import-btn:hover {
  background: linear-gradient(135deg, #047857 0%, #0f766e 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(5, 150, 105, 0.3);
}

.btn-icon {
  font-size: 1rem;
}

.table-container {
  background: white;
  border-radius: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
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

.events-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.events-table th,
.events-table td {
  padding: 1rem 1.5rem;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
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

.event-row:hover {
  background: #f8fafc;
}

.event-row td {
  vertical-align: middle;
}

.event-name {
  font-weight: 600;
  color: #2d3748;
}

.event-description {
  max-width: 300px;
}

.description-text {
  color: #4a5568;
  line-height: 1.4;
}

.event-date {
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  color: #2d3748;
  font-weight: 500;
}

.event-location {
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  color: #4a5568;
  font-size: 0.85rem;
}

.no-location {
  color: #a0aec0;
  font-style: italic;
}

.type-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.type-badge.military { background: #fef2f2; color: #dc2626; }
.type-badge.political { background: #f0f9ff; color: #2563eb; }
.type-badge.historic { background: #f9fafb; color: #374151; }
.type-badge.cultural { background: #f0fdf4; color: #16a34a; }
.type-badge.scientific { background: #fefce8; color: #ca8a04; }

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  max-width: 200px;
}

.tag-badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
}

.no-tags {
  color: #a0aec0;
  font-style: italic;
  font-size: 0.8rem;
}

.event-actions {
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

.edit-btn {
  background: #fef3c7;
  color: #d97706;
}

.edit-btn:hover {
  background: #fde68a;
  transform: scale(1.1);
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
}

.delete-btn:hover {
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

/* Modal styles would be included here - same as AdminPanel.vue */
</style>