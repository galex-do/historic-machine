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
          :total-items="totalEvents"
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
  name: 'AdminPanel',
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

    const formatDateToDDMMYYYY = (isoDateString) => {
      if (!isoDateString) return ''
      
      try {
        const date = new Date(isoDateString)
        const day = String(date.getUTCDate()).padStart(2, '0')
        const month = String(date.getUTCMonth() + 1).padStart(2, '0')
        const year = String(date.getUTCFullYear()).padStart(4, '0')
        return `${day}/${month}/${year}`
      } catch {
        return ''
      }
    }

    const parseDDMMYYYYToISO = (dateDisplay) => {
      if (!dateDisplay) return null
      
      const parts = dateDisplay.split('/')
      if (parts.length !== 3) return null
      
      const [day, month, year] = parts
      const dayNum = parseInt(day, 10)
      const monthNum = parseInt(month, 10)
      const yearNum = parseInt(year, 10)
      
      // Basic validation
      if (yearNum < 1 || yearNum > 9999) return null
      if (monthNum < 1 || monthNum > 12) return null
      if (dayNum < 1 || dayNum > 31) return null
      
      // For very old years, construct full ISO datetime string manually
      const paddedYear = String(yearNum).padStart(4, '0')
      const paddedMonth = String(monthNum).padStart(2, '0')
      const paddedDay = String(dayNum).padStart(2, '0')
      
      return `${paddedYear}-${paddedMonth}-${paddedDay}T00:00:00Z`
    }

    const updateEventDate = () => {
      const isoDate = parseDDMMYYYYToISO(eventForm.value.date_display)
      if (isoDate) {
        eventForm.value.date = isoDate
      } else {
        eventForm.value.date = ''
      }
    }

    const editEvent = (event) => {
      editingEvent.value = event
      eventForm.value = {
        name: event.name,
        description: event.description,
        date_display: formatDateToDDMMYYYY(event.event_date),
        date: event.event_date,
        era: event.era || 'AD',
        latitude: event.latitude,
        longitude: event.longitude,
        lens_type: event.lens_type,
        selectedTags: event.tags ? [...event.tags] : [],
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
        console.log('Event deleted successfully')
        // Update shared state immediately - this will update both admin and map views
        await handleEventDeleted(event.id)
      } catch (err) {
        console.error('Error deleting event:', err)
        localError.value = 'Failed to delete event'
      } finally {
        localLoading.value = false
      }
    }

    const saveEvent = async () => {
      localLoading.value = true
      localError.value = null
      
      try {
        const eventData = {
          name: eventForm.value.name,
          description: eventForm.value.description,
          event_date: eventForm.value.date, // Use the ISO date directly
          era: eventForm.value.era,
          latitude: parseFloat(eventForm.value.latitude),
          longitude: parseFloat(eventForm.value.longitude),
          lens_type: eventForm.value.lens_type
        }
        
        let eventId
        if (editingEvent.value) {
          // Update existing event
          await apiService.updateEvent(editingEvent.value.id, eventData)
          eventId = editingEvent.value.id
          console.log('Event updated successfully')
        } else {
          // Create new event
          const response = await apiService.createEvent(eventData)
          eventId = response.id
          console.log('Event created successfully')
        }
        
        // Update event tags
        const tagIds = eventForm.value.selectedTags.map(tag => tag.id)
        await setEventTags(eventId, tagIds)
        console.log('Event tags updated successfully')
        
        closeModal()
        await fetchEvents() // Refresh shared events state
        
      } catch (err) {
        console.error('Error saving event:', err)
        localError.value = 'Failed to save event'
      } finally {
        localLoading.value = false
      }
    }

    const closeModal = () => {
      showCreateModal.value = false
      showEditModal.value = false
      editingEvent.value = null
      error.value = null
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
    }

    // Tag functionality
    const filteredTags = computed(() => {
      if (!eventForm.value.tagSearch) return []
      const search = eventForm.value.tagSearch.toLowerCase()
      return allTags.value.filter(tag => 
        tag.name.toLowerCase().includes(search) &&
        !eventForm.value.selectedTags.some(selected => selected.id === tag.id)
      )
    })
    
    const canCreateNewTag = computed(() => {
      if (!eventForm.value.tagSearch.trim()) return false
      const search = eventForm.value.tagSearch.toLowerCase().trim()
      // Check if tag already exists in database or is already selected
      return !allTags.value.some(tag => tag.name.toLowerCase() === search) &&
             !eventForm.value.selectedTags.some(tag => tag.name.toLowerCase() === search)
    })
    
    const addTag = (tag) => {
      if (!eventForm.value.selectedTags.some(selected => selected.id === tag.id)) {
        eventForm.value.selectedTags.push(tag)
        eventForm.value.tagSearch = ''
        error.value = null // Clear any previous errors
      } else {
        error.value = `Tag "${tag.name}" is already selected for this event.`
      }
    }
    
    const removeTag = (tag) => {
      eventForm.value.selectedTags = eventForm.value.selectedTags.filter(t => t.id !== tag.id)
    }
    
    const createAndAddTag = async () => {
      const tagName = eventForm.value.tagSearch.trim()
      if (!tagName) return
      
      // Check for duplicate names (case-insensitive)
      const existingTag = allTags.value.find(tag => 
        tag.name.toLowerCase() === tagName.toLowerCase()
      )
      
      if (existingTag) {
        error.value = `Tag "${tagName}" already exists. Please choose a different name.`
        return
      }
      
      // Check if already selected
      const selectedTag = eventForm.value.selectedTags.find(tag => 
        tag.name.toLowerCase() === tagName.toLowerCase()
      )
      
      if (selectedTag) {
        error.value = `Tag "${tagName}" is already selected for this event.`
        return
      }
      
      try {
        const newTag = await createTag({
          name: tagName,
          color: eventForm.value.newTagColor
        })
        
        if (newTag) {
          eventForm.value.selectedTags.push(newTag)
          eventForm.value.tagSearch = ''
          eventForm.value.newTagColor = '#3B82F6'
          error.value = null // Clear any previous errors
        }
      } catch (err) {
        console.error('Error creating tag:', err)
        // Handle specific error messages
        if (err.message && err.message.includes('duplicate') || err.message.includes('already exists')) {
          error.value = `Tag "${tagName}" already exists. Please choose a different name.`
        } else {
          error.value = 'Failed to create new tag. Please try again.'
        }
      }
    }
    
    const handleTagInput = () => {
      if (filteredTags.value.length > 0) {
        addTag(filteredTags.value[0])
      } else if (canCreateNewTag.value) {
        createAndAddTag()
      } else if (eventForm.value.tagSearch.trim()) {
        // Tag already exists, show helpful message
        const tagName = eventForm.value.tagSearch.trim()
        const existingTag = allTags.value.find(tag => 
          tag.name.toLowerCase() === tagName.toLowerCase()
        )
        if (existingTag) {
          error.value = `Tag "${tagName}" already exists. Use the search suggestions above.`
        }
      }
    }
    
    // Fetch all events for client-side filtering and pagination
    const fetchAllEvents = async () => {
      paginatedLoading.value = true
      try {
        // Use the events composable to get all events
        await fetchEvents()
        allEvents.value = events.value || []
        
        // Events are now sorted in the filteredEvents computed property
      } catch (err) {
        console.error('Error fetching all events:', err)
        localError.value = 'Failed to load events'
      } finally {
        paginatedLoading.value = false
      }
    }
    
    const handlePageChange = (page) => {
      currentPage.value = page
      // No need to refetch - just change page in computed property
    }
    
    const handlePageSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1 // Reset to first page
      // No need to refetch - just change page size in computed property
    }
    
    // Computed sorted events (kept for compatibility but not used in pagination mode)
    const sortedEvents = computed(() => {
      if (!events.value || events.value.length === 0) return []
      
      return [...events.value].sort((a, b) => {
        let aValue, bValue
        
        switch (sortField.value) {
          case 'name':
            aValue = a.name.toLowerCase()
            bValue = b.name.toLowerCase()
            break
          case 'date':
            // Use proper astronomical year for chronological sorting
            const getAstronomicalYear = (dateString, era) => {
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
            
            aValue = getAstronomicalYear(a.event_date, a.era)
            bValue = getAstronomicalYear(b.event_date, b.era)
            break
          case 'type':
            aValue = a.lens_type.toLowerCase()
            bValue = b.lens_type.toLowerCase()
            break
          default:
            return 0
        }
        
        if (aValue < bValue) {
          return sortDirection.value === 'asc' ? -1 : 1
        }
        if (aValue > bValue) {
          return sortDirection.value === 'asc' ? 1 : -1
        }
        return 0
      })
    })
    
    // Toggle sort function
    const toggleSort = (field) => {
      if (sortField.value === field) {
        // Toggle direction if same field
        sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
      } else {
        // Set new field and default to ascending
        sortField.value = field
        sortDirection.value = 'asc'
      }
      // Reset to first page when sorting changes
      currentPage.value = 1
    }
    
    // Import functions
    const triggerFileUpload = () => {
      fileInput.value?.click()
    }
    
    const handleFileUpload = async (event) => {
      const file = event.target.files[0]
      if (!file) return
      
      try {
        loading.value = true
        error.value = null
        
        const text = await file.text()
        const dataset = JSON.parse(text)
        
        if (!dataset.events || !Array.isArray(dataset.events)) {
          throw new Error('Invalid dataset format: expected events array')
        }
        
        // Import events via API with filename
        const response = await apiService.importEvents(dataset.events, file.name)
        
        if (response.success) {
          // Refresh shared events state - this will update both admin and map views
          await fetchEvents()
          alert(`Successfully imported ${response.imported_count} events!`)
        } else {
          throw new Error(response.error || 'Failed to import events')
        }
        
      } catch (err) {
        console.error('Import error:', err)
        localError.value = err.message || 'Failed to import events'
      } finally {
        localLoading.value = false
        // Reset file input
        if (fileInput.value) fileInput.value.value = ''
      }
    }
    
    onMounted(async () => {
      await loadTags()
      await fetchAllEvents()
    })

    return {
      canAccessAdmin,
      events,
      sortedEvents,
      loading: computed(() => loading.value || localLoading.value || paginatedLoading.value),
      error: computed(() => error.value || localError.value),
      showCreateModal,
      showEditModal,
      editingEvent,
      eventForm,
      allTags,
      filteredTags,
      canCreateNewTag,
      sortField,
      sortDirection,
      // Pagination state and functions
      currentPage,
      pageSize,
      totalEvents: filteredTotalEvents,
      displayedEvents: currentPageEvents,
      handlePageChange,
      handlePageSizeChange,
      fetchAllEvents,
      fetchEvents,
      formatDate,
      formatDateWithEra,
      formatDateToDDMMYYYY,
      getContrastColor,
      updateEventDate,
      editEvent,
      deleteEvent,
      saveEvent,
      closeModal,
      addTag,
      removeTag,
      createAndAddTag,
      handleTagInput,
      toggleSort,
      fileInput,
      triggerFileUpload,
      handleFileUpload,
      // Filter state and methods
      selectedLensType,
      showLensDropdown,
      toggleLensDropdown,
      handleLensTypeChange,
      filteredEvents,
      currentPageEvents,
      filteredTotalEvents
    }
  }
}
</script>

<style scoped>
.admin-panel {
  padding: 1.25rem 2rem;
  max-width: 100%;
  margin: 0;
  min-height: calc(100vh - 80px);
  background: #f8f9fa;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #e2e8f0;
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
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(66, 153, 225, 0.2);
}

.create-btn:hover {
  background: #3182ce;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(66, 153, 225, 0.3);
}

.btn-icon {
  font-size: 1rem;
}

.table-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.table-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  padding: 0.75rem 1.25rem;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  border-radius: 12px 12px 0 0;
}

.table-filters {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  color: #718096;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e2e8f0;
  border-top: 3px solid #4299e1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.events-table {
  width: 100%;
  border-collapse: collapse;
}

.events-table th {
  background: #f7fafc;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #4a5568;
  border-bottom: 2px solid #e2e8f0;
}

.event-row {
  transition: background-color 0.2s ease;
}

.event-row:hover {
  background: #f7fafc;
}

.events-table td {
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  vertical-align: top;
}

.event-name {
  font-weight: 600;
  color: #2d3748;
  max-width: 200px;
}

.event-description {
  max-width: 300px;
}

.description-text {
  color: #4a5568;
  line-height: 1.4;
}

.event-date {
  font-family: 'Monaco', 'Consolas', monospace;
  color: #2d3748;
  white-space: nowrap;
}

.event-location {
  font-family: 'Monaco', 'Consolas', monospace;
  color: #4a5568;
  font-size: 0.9rem;
}

.no-location {
  color: #a0aec0;
  font-style: italic;
}

.type-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.type-badge.military { background: #fed7d7; color: #c53030; }
.type-badge.political { background: #bee3f8; color: #2b6cb0; }
.type-badge.historic { background: #faf089; color: #744210; }
.type-badge.cultural { background: #c6f6d5; color: #276749; }

.event-tags {
  max-width: 200px;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.tag-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.no-tags {
  color: #a0aec0;
  font-style: italic;
  font-size: 0.85rem;
}

/* Tags section styling */
.tags-section {
  margin-top: 0.5rem;
}

.selected-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.selected-tag {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 500;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.remove-tag-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.1rem;
  line-height: 1;
  padding: 0;
  margin-left: 0.25rem;
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.remove-tag-btn:hover {
  opacity: 1;
}

.tag-input-section {
  position: relative;
}

.tag-search-input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.9rem;
  transition: border-color 0.2s ease;
}

.tag-search-input:focus {
  outline: none;
  border-color: #4299e1;
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
}

.tag-suggestions {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  max-height: 200px;
  overflow-y: auto;
  z-index: 100000;
  margin-top: 2px;
}

.tag-suggestion {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #f7fafc;
}

.tag-suggestion:hover {
  background-color: #f7fafc;
}

.tag-suggestion:last-child {
  border-bottom: none;
}

.tag-color-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.new-tag-option {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 100000;
  margin-top: 2px;
  padding: 0.75rem;
}

.new-tag-form {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.new-tag-text {
  color: #4a5568;
  font-size: 0.9rem;
  font-weight: 500;
}

.color-picker {
  width: 40px;
  height: 32px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  cursor: pointer;
  padding: 0;
}

.create-tag-btn {
  background: #4299e1;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.create-tag-btn:hover {
  background: #3182ce;
  transform: translateY(-1px);
}

.create-new-suggestion {
  border-top: 1px solid #e2e8f0;
  font-style: italic;
  color: #4a5568;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.create-new-suggestion:hover {
  background-color: #f0f9ff;
}

.inline-color-picker {
  width: 20px;
  height: 20px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  cursor: pointer;
  padding: 0;
  margin-left: 0.5rem;
}

.event-actions {
  white-space: nowrap;
}

.action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 4px;
  margin: 0 0.25rem;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.edit-btn:hover {
  background: #e6fffa;
}

.delete-btn:hover {
  background: #fed7d7;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #718096;
}

/* Modal Styles */
.admin-modal {
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-body {
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}

.modal-header {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 1rem;
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
  color: #a0aec0;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #f7fafc;
  color: #4a5568;
}

/* Sorting Styles */
.sortable-header {
  cursor: pointer;
  user-select: none;
  transition: background-color 0.2s ease;
  position: relative;
}

.sortable-header:hover {
  background: #edf2f7 !important;
}

.sortable-header.active {
  background: #e2e8f0 !important;
  color: #2d3748;
}

.sort-indicator {
  display: inline-block;
  margin-left: 0.5rem;
  font-size: 1rem;
  min-width: 16px;
  font-weight: bold;
}

.sort-arrow {
  color: #4299e1;
  font-weight: bold;
  text-shadow: 0 1px 2px rgba(66, 153, 225, 0.3);
}

.sort-placeholder {
  color: #718096;
  opacity: 0.8;
  transition: all 0.2s ease;
}

.sortable-header:hover .sort-placeholder {
  color: #4299e1;
  opacity: 1;
  transform: scale(1.1);
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.import-btn {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 2px 4px rgba(5, 150, 105, 0.2);
}

.import-btn:hover {
  background: linear-gradient(135deg, #047857 0%, #065f46 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(5, 150, 105, 0.3);
}

.import-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 4px rgba(5, 150, 105, 0.2);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #4a5568;
  font-weight: 600;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.2s ease;
  box-sizing: border-box;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: #4299e1;
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.form-hint {
  color: #718096;
  font-size: 0.8rem;
  margin-top: 0.25rem;
  display: block;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.cancel-btn {
  padding: 0.75rem 1.5rem;
  background: #e2e8f0;
  color: #4a5568;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s ease;
}

.cancel-btn:hover {
  background: #cbd5e0;
}

.submit-btn {
  padding: 0.75rem 1.5rem;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s ease;
}

.submit-btn:hover:not(:disabled) {
  background: #3182ce;
}

.submit-btn:disabled {
  background: #a0aec0;
  cursor: not-allowed;
}

.error-message {
  background: #fed7d7;
  color: #c53030;
  padding: 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  font-weight: 600;
}

@media (max-width: 1024px) {
  .admin-panel {
    padding: 1rem;
  }
  
  .admin-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .events-table {
    font-size: 0.9rem;
  }
  
  .events-table th,
  .events-table td {
    padding: 0.75rem;
  }
}

@media (max-width: 768px) {
  .events-table th:nth-child(2),
  .events-table td:nth-child(2) {
    display: none;
  }
  
  .events-table th:nth-child(4),
  .events-table td:nth-child(4) {
    display: none;
  }
}
</style>