<template>
  <div class="admin-panel">
    <div class="admin-header">
      <div class="admin-title">
        <h2>Historic Events</h2>
        <p class="admin-subtitle">Manage and create historic events</p>
      </div>
      <button @click="showCreateModal = true" class="create-btn" v-if="canAccessAdmin">
        <span class="btn-icon">➕</span>
        Create New Event
      </button>
    </div>

    <!-- Events Table -->
    <div class="table-container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading events...</p>
      </div>
      
      <table v-else class="events-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Date</th>
            <th>Location</th>
            <th>Type</th>
            <th>Tags</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in events" :key="event.id" class="event-row">
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
      
      <div v-if="!loading && events.length === 0" class="empty-state">
        <p>No events found. Create your first historical event!</p>
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
import apiService from '@/services/api.js'

export default {
  name: 'AdminPanel',
  setup() {
    const { canAccessAdmin } = useAuth()
    const { allTags, loadTags, createTag, setEventTags, getTagsByIds } = useTags()
    
    const events = ref([])
    const loading = ref(false)
    const error = ref(null)
    
    const showCreateModal = ref(false)
    const showEditModal = ref(false)
    const editingEvent = ref(null)
    
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

    const fetchEvents = async () => {
      loading.value = true
      error.value = null
      try {
        const response = await apiService.getEvents()
        events.value = Array.isArray(response) ? response : []
        console.log('Loaded events for admin:', events.value.length)
      } catch (err) {
        console.error('Error fetching events:', err)
        error.value = 'Failed to load events'
      } finally {
        loading.value = false
      }
    }

    const formatDate = (dateString) => {
      try {
        const date = new Date(dateString)
        return date.toLocaleDateString('en-GB', {
          year: 'numeric',
          month: 'short',
          day: 'numeric'
        })
      } catch {
        return 'Invalid Date'
      }
    }

    const formatDateWithEra = (dateString, era) => {
      try {
        const formattedDate = formatDate(dateString)
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
      
      loading.value = true
      error.value = null
      try {
        await apiService.deleteEvent(event.id)
        console.log('Event deleted successfully')
        await fetchEvents() // Refresh the list
      } catch (err) {
        console.error('Error deleting event:', err)
        error.value = 'Failed to delete event'
      } finally {
        loading.value = false
      }
    }

    const saveEvent = async () => {
      loading.value = true
      error.value = null
      
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
        await fetchEvents() // Refresh the list
        
      } catch (err) {
        console.error('Error saving event:', err)
        error.value = 'Failed to save event'
      } finally {
        loading.value = false
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
    
    onMounted(async () => {
      await loadTags()
      await fetchEvents()
    })

    return {
      canAccessAdmin,
      events,
      loading,
      error,
      showCreateModal,
      showEditModal,
      editingEvent,
      eventForm,
      allTags,
      filteredTags,
      canCreateNewTag,
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
      handleTagInput
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
  z-index: 1000;
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
  z-index: 1000;
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