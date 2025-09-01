<template>
  <div class="events-admin">
    <div class="page-header">
      <h2>Events Management</h2>
      <button @click="show_create_modal = true" class="btn-primary">Add New Event</button>
    </div>

    <!-- Events Table -->
    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Description</th>
            <th>Date</th>
            <th>Location</th>
            <th>Type</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in events" :key="event.id">
            <td>{{ event.id }}</td>
            <td>{{ event.name }}</td>
            <td class="description-cell">{{ event.description }}</td>
            <td>{{ formatDate(event.event_date) }}</td>
            <td>{{ event.latitude.toFixed(2) }}, {{ event.longitude.toFixed(2) }}</td>
            <td>
              <span class="lens-badge" :class="`lens-${event.lens_type}`">
                {{ getLensTypeDisplay(event.lens_type) }}
              </span>
            </td>
            <td class="actions-cell">
              <button @click="edit_event(event)" class="btn-edit">Edit</button>
              <button @click="delete_event(event.id)" class="btn-delete">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="show_create_modal || show_edit_modal" class="modal-overlay" @click="close_modal">
      <div class="modal-content" @click.stop>
        <h3>{{ show_create_modal ? 'Add New Event' : 'Edit Event' }}</h3>
        
        <form @submit.prevent="submit_event">
          <div class="form-group">
            <label>Name:</label>
            <input v-model="event_form.name" type="text" required class="form-input" />
          </div>
          
          <div class="form-group">
            <label>Description:</label>
            <textarea v-model="event_form.description" required class="form-textarea"></textarea>
          </div>
          
          <div class="form-group">
            <label>Date:</label>
            <input v-model="event_form.event_date" type="date" required class="form-input" />
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>Latitude:</label>
              <input v-model.number="event_form.latitude" type="number" step="any" required class="form-input" />
            </div>
            <div class="form-group">
              <label>Longitude:</label>
              <input v-model.number="event_form.longitude" type="number" step="any" required class="form-input" />
            </div>
          </div>
          
          <div class="form-group">
            <label>Event Type:</label>
            <select v-model="event_form.lens_type" required class="form-select">
              <option value="historic">📜 Historic</option>
              <option value="political">🏛️ Political</option>
              <option value="cultural">🎭 Cultural</option>
              <option value="military">⚔️ Military</option>
            </select>
          </div>
          
          <div class="modal-actions">
            <button type="button" @click="close_modal" class="btn-secondary">Cancel</button>
            <button type="submit" class="btn-primary">{{ show_create_modal ? 'Create' : 'Update' }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EventsAdmin',
  data() {
    return {
      events: [],
      show_create_modal: false,
      show_edit_modal: false,
      event_form: {
        id: null,
        name: '',
        description: '',
        event_date: '',
        latitude: 0,
        longitude: 0,
        lens_type: 'historic'
      }
    }
  },
  
  async mounted() {
    await this.load_events()
  },
  
  methods: {
    async load_events() {
      try {
        const response = await fetch('http://localhost:8080/api/events')
        if (response.ok) {
          this.events = await response.json()
        } else {
          console.error('Failed to load events')
        }
      } catch (error) {
        console.error('Error loading events:', error)
      }
    },
    
    edit_event(event) {
      this.event_form = { ...event }
      this.show_edit_modal = true
    },
    
    async submit_event() {
      try {
        const url = this.show_create_modal 
          ? 'http://localhost:8080/api/events'
          : `http://localhost:8080/api/events/${this.event_form.id}`
        
        const method = this.show_create_modal ? 'POST' : 'PUT'
        
        const response = await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(this.event_form)
        })
        
        if (response.ok) {
          await this.load_events()
          this.close_modal()
        } else {
          console.error('Failed to save event')
        }
      } catch (error) {
        console.error('Error saving event:', error)
      }
    },
    
    async delete_event(id) {
      if (!confirm('Are you sure you want to delete this event?')) return
      
      try {
        const response = await fetch(`http://localhost:8080/api/events/${id}`, {
          method: 'DELETE'
        })
        
        if (response.ok) {
          await this.load_events()
        } else {
          console.error('Failed to delete event')
        }
      } catch (error) {
        console.error('Error deleting event:', error)
      }
    },
    
    close_modal() {
      this.show_create_modal = false
      this.show_edit_modal = false
      this.event_form = {
        id: null,
        name: '',
        description: '',
        event_date: '',
        latitude: 0,
        longitude: 0,
        lens_type: 'historic'
      }
    },
    
    formatDate(date_string) {
      try {
        const date = new Date(date_string)
        return date.toLocaleDateString()
      } catch (error) {
        return date_string
      }
    },
    
    getLensTypeDisplay(lens_type) {
      const types = {
        historic: '📜 Historic',
        political: '🏛️ Political',
        cultural: '🎭 Cultural',
        military: '⚔️ Military'
      }
      return types[lens_type] || lens_type
    }
  }
}
</script>

<style scoped>
.events-admin {
  max-width: 1200px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h2 {
  margin: 0;
  color: #2d3748;
}

.table-container {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
}

.admin-table th,
.admin-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.admin-table th {
  background: #f7fafc;
  font-weight: 600;
  color: #4a5568;
}

.description-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lens-badge {
  padding: 4px 8px;
  border-radius: 16px;
  font-size: 0.875rem;
  font-weight: 500;
}

.lens-historic { background: #e6fffa; color: #234e52; }
.lens-political { background: #fef5e7; color: #744210; }
.lens-cultural { background: #f0fff4; color: #22543d; }
.lens-military { background: #fed7d7; color: #742a2a; }

.actions-cell {
  white-space: nowrap;
}

.btn-primary, .btn-edit, .btn-delete, .btn-secondary {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  margin-right: 8px;
  transition: background 0.2s;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5a67d8;
}

.btn-edit {
  background: #4299e1;
  color: white;
}

.btn-edit:hover {
  background: #3182ce;
}

.btn-delete {
  background: #f56565;
  color: white;
}

.btn-delete:hover {
  background: #e53e3e;
}

.btn-secondary {
  background: #e2e8f0;
  color: #4a5568;
}

.btn-secondary:hover {
  background: #cbd5e0;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content h3 {
  margin: 0 0 1.5rem 0;
  color: #2d3748;
}

.form-group {
  margin-bottom: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.form-input, .form-textarea, .form-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 1rem;
}

.form-textarea {
  min-height: 80px;
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}
</style>