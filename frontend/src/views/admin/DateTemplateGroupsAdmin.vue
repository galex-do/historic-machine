<template>
  <div class="date-template-groups-admin">
    <div class="page-header">
      <h2>Date Template Groups Management</h2>
      <button @click="show_create_modal = true" class="btn-primary">Add New Group</button>
    </div>

    <!-- Groups Table -->
    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Description</th>
            <th>Templates Count</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="group in groups" :key="group.id">
            <td>{{ group.id }}</td>
            <td>{{ group.name }}</td>
            <td class="description-cell">{{ group.description }}</td>
            <td>{{ group.templates ? group.templates.length : 0 }}</td>
            <td class="actions-cell">
              <button @click="edit_group(group)" class="btn-edit">Edit</button>
              <button @click="delete_group(group.id)" class="btn-delete">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="show_create_modal || show_edit_modal" class="modal-overlay" @click="close_modal">
      <div class="modal-content" @click.stop>
        <h3>{{ show_create_modal ? 'Add New Group' : 'Edit Group' }}</h3>
        
        <form @submit.prevent="submit_group">
          <div class="form-group">
            <label>Name:</label>
            <input v-model="group_form.name" type="text" required class="form-input" />
          </div>
          
          <div class="form-group">
            <label>Description:</label>
            <textarea v-model="group_form.description" required class="form-textarea"></textarea>
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
  name: 'DateTemplateGroupsAdmin',
  data() {
    return {
      groups: [],
      show_create_modal: false,
      show_edit_modal: false,
      group_form: {
        id: null,
        name: '',
        description: ''
      }
    }
  },
  
  async mounted() {
    await this.load_groups()
  },
  
  methods: {
    async load_groups() {
      try {
        const response = await fetch('http://localhost:8080/api/date-template-groups')
        if (response.ok) {
          this.groups = await response.json()
        } else {
          console.error('Failed to load date template groups')
        }
      } catch (error) {
        console.error('Error loading date template groups:', error)
      }
    },
    
    edit_group(group) {
      this.group_form = { ...group }
      this.show_edit_modal = true
    },
    
    async submit_group() {
      try {
        const url = this.show_create_modal 
          ? 'http://localhost:8080/api/date-template-groups'
          : `http://localhost:8080/api/date-template-groups/${this.group_form.id}`
        
        const method = this.show_create_modal ? 'POST' : 'PUT'
        
        const response = await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(this.group_form)
        })
        
        if (response.ok) {
          await this.load_groups()
          this.close_modal()
        } else {
          console.error('Failed to save group')
        }
      } catch (error) {
        console.error('Error saving group:', error)
      }
    },
    
    async delete_group(id) {
      if (!confirm('Are you sure you want to delete this date template group?')) return
      
      try {
        const response = await fetch(`http://localhost:8080/api/date-template-groups/${id}`, {
          method: 'DELETE'
        })
        
        if (response.ok) {
          await this.load_groups()
        } else {
          console.error('Failed to delete group')
        }
      } catch (error) {
        console.error('Error deleting group:', error)
      }
    },
    
    close_modal() {
      this.show_create_modal = false
      this.show_edit_modal = false
      this.group_form = {
        id: null,
        name: '',
        description: ''
      }
    }
  }
}
</script>

<style scoped>
.date-template-groups-admin {
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

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.form-input, .form-textarea {
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