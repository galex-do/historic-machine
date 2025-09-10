<template>
  <div class="admin-panel">
    <div class="admin-header">
      <div class="admin-title">
        <h2>Users Management</h2>
        <p class="admin-subtitle">Manage user accounts and access levels</p>
      </div>
      <div class="action-buttons" v-if="isSuper">
        <button @click="showCreateModal = true" class="create-btn">
          <span class="btn-icon">➕</span>
          Create New User
        </button>
      </div>
    </div>

    <!-- Users Table -->
    <div class="table-container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading users...</p>
      </div>
      
      <!-- Table Controls & Pagination -->
      <div v-if="!loading && totalUsers > 0" class="table-controls">
        <div class="table-filters">
          <input 
            v-model="searchQuery"
            type="text"
            placeholder="Search users..."
            class="search-input"
          />
          <select v-model="accessLevelFilter" class="filter-select">
            <option value="all">All Access Levels</option>
            <option value="guest">Guest</option>
            <option value="user">User</option>
            <option value="editor">Editor</option>
            <option value="admin">Admin</option>
            <option value="super">Super</option>
          </select>
        </div>
        <TablePagination 
          :current-page="currentPage"
          :page-size="pageSize"
          :total-items="filteredTotalUsers"
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
              @click="toggleSort('username')"
              :class="{ 'active': sortField === 'username' }"
            >
              Username
              <span class="sort-indicator">
                <span v-if="sortField === 'username'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>Email</th>
            <th 
              class="sortable-header" 
              @click="toggleSort('access_level')"
              :class="{ 'active': sortField === 'access_level' }"
            >
              Access Level
              <span class="sort-indicator">
                <span v-if="sortField === 'access_level'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>Status</th>
            <th 
              class="sortable-header" 
              @click="toggleSort('created_at')"
              :class="{ 'active': sortField === 'created_at' }"
            >
              Created
              <span class="sort-indicator">
                <span v-if="sortField === 'created_at'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>Last Active</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="userItem in currentPageUsers" :key="userItem.id" class="event-row">
            <td class="user-username">
              <div class="username-cell">
                {{ userItem.username }}
                <span v-if="userItem.id === user?.id" class="current-user-badge">You</span>
              </div>
            </td>
            <td class="user-email">{{ userItem.email || 'No email' }}</td>
            <td class="user-access-level">
              <span class="access-badge" :class="userItem.access_level">{{ userItem.access_level.toUpperCase() }}</span>
            </td>
            <td class="user-status">
              <span class="status-badge" :class="{ 'active': userItem.is_active, 'inactive': !userItem.is_active }">
                {{ userItem.is_active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="user-created">{{ formatDate(userItem.created_at) }}</td>
            <td class="user-last-active">{{ formatLastActive(userItem.last_active_at) }}</td>
            <td class="user-actions">
              <button @click="editUser(userItem)" class="action-btn edit-btn" title="Edit" :disabled="!canEditUser(userItem)">
                ✏️
              </button>
              <button @click="deleteUser(userItem)" class="action-btn delete-btn" title="Delete" :disabled="!canDeleteUser(userItem)">
                🗑️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="!loading && totalUsers === 0" class="empty-state">
        <p>No users found.</p>
      </div>
      
      <div v-if="!loading && totalUsers > 0 && filteredUsers.length === 0" class="empty-state">
        <p>No users match your search criteria.</p>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content admin-modal" @click.stop>
          <div class="modal-header">
            <h3>{{ editingUser ? 'Edit User' : 'Create New User' }}</h3>
            <button @click="closeModal" class="close-btn">&times;</button>
          </div>
          
          <div class="modal-body">
            <div v-if="error" class="validation-warning">{{ error }}</div>
            
            <form @submit.prevent="saveUser">
              <div class="form-row">
                <div class="form-group">
                  <label for="user-username">Username *</label>
                  <input 
                    id="user-username"
                    v-model="userForm.username" 
                    type="text" 
                    required 
                    class="form-input"
                    placeholder="Enter username"
                    :disabled="!!editingUser"
                  />
                  <small v-if="editingUser" class="form-hint">Username cannot be changed</small>
                </div>
                <div class="form-group">
                  <label for="user-email">Email</label>
                  <input 
                    id="user-email"
                    v-model="userForm.email" 
                    type="email" 
                    class="form-input"
                    placeholder="Enter email address"
                  />
                </div>
              </div>
              
              <div class="form-row" v-if="!editingUser">
                <div class="form-group">
                  <label for="user-password">Password *</label>
                  <input 
                    id="user-password"
                    v-model="userForm.password" 
                    type="password" 
                    required 
                    class="form-input"
                    placeholder="Enter password"
                  />
                  <small class="form-hint">Minimum 8 characters, include uppercase, lowercase, numbers, and special characters</small>
                </div>
                <div class="form-group">
                  <label for="user-confirm-password">Confirm Password *</label>
                  <input 
                    id="user-confirm-password"
                    v-model="userForm.confirmPassword" 
                    type="password" 
                    required 
                    class="form-input"
                    placeholder="Confirm password"
                  />
                </div>
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="user-access-level">Access Level *</label>
                  <select 
                    id="user-access-level"
                    v-model="userForm.access_level" 
                    required 
                    class="form-input"
                    :disabled="!canChangeAccessLevel(editingUser)"
                  >
                    <option value="user">User</option>
                    <option value="editor">Editor</option>
                    <option value="admin">Admin</option>
                  </select>
                  <small v-if="editingUser?.id === user?.id" class="form-hint">You cannot change your own access level</small>
                  <small v-else class="form-hint">Choose the appropriate access level for this user</small>
                </div>
                <div class="form-group">
                  <label for="user-status">Status</label>
                  <select 
                    id="user-status"
                    v-model="userForm.is_active" 
                    class="form-input"
                    :disabled="editingUser?.id === user?.id"
                  >
                    <option :value="true">Active</option>
                    <option :value="false">Inactive</option>
                  </select>
                  <small v-if="editingUser?.id === user?.id" class="form-hint">You cannot deactivate yourself</small>
                </div>
              </div>
              
              <div class="form-actions">
                <button type="button" @click="closeModal" class="cancel-btn">Cancel</button>
                <button type="submit" class="submit-btn" :disabled="localLoading || !isFormValid">
                  {{ localLoading ? 'Saving...' : (editingUser ? 'Update User' : 'Create User') }}
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
import apiService from '@/services/api.js'
import TablePagination from '@/components/TablePagination.vue'

export default {
  name: 'AdminUsers',
  components: {
    TablePagination
  },
  setup() {
    const { user, isSuper } = useAuth()
    
    const localLoading = ref(false)
    const localError = ref(null)
    const searchQuery = ref('')
    const accessLevelFilter = ref('all')
    
    const showCreateModal = ref(false)
    const showEditModal = ref(false)
    const editingUser = ref(null)
    
    // Sorting state
    const sortField = ref('username')
    const sortDirection = ref('asc')
    
    // Pagination state
    const currentPage = ref(1)
    const pageSize = ref(10)
    
    const allUsers = ref([])
    
    const userForm = ref({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      access_level: 'user',
      is_active: true
    })

    // Filtered and sorted users
    const filteredUsers = computed(() => {
      if (!allUsers.value || !Array.isArray(allUsers.value)) {
        return []
      }
      
      let filtered = allUsers.value
      
      // Apply search filter
      if (searchQuery.value.trim()) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(user => 
          user.username.toLowerCase().includes(query) ||
          (user.email && user.email.toLowerCase().includes(query))
        )
      }
      
      // Apply access level filter
      if (accessLevelFilter.value !== 'all') {
        filtered = filtered.filter(user => user.access_level === accessLevelFilter.value)
      }
      
      // Apply sorting
      return [...filtered].sort((a, b) => {
        let aValue, bValue
        
        if (sortField.value === 'created_at') {
          aValue = new Date(a.created_at).getTime()
          bValue = new Date(b.created_at).getTime()
        } else {
          aValue = a[sortField.value] || ''
          bValue = b[sortField.value] || ''
          if (typeof aValue === 'string') aValue = aValue.toLowerCase()
          if (typeof bValue === 'string') bValue = bValue.toLowerCase()
        }
        
        if (sortDirection.value === 'asc') {
          return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
        } else {
          return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
        }
      })
    })
    
    const filteredTotalUsers = computed(() => filteredUsers.value.length)
    const totalUsers = computed(() => allUsers.value?.length || 0)
    
    // Current page of filtered users
    const currentPageUsers = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filteredUsers.value.slice(start, end)
    })

    const isFormValid = computed(() => {
      if (!userForm.value.username) return false
      if (!editingUser.value) {
        if (!userForm.value.password || !userForm.value.confirmPassword) return false
        if (userForm.value.password !== userForm.value.confirmPassword) return false
      }
      return true
    })

    const formatDate = (dateString) => {
      if (!dateString) return 'Unknown'
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

    const formatLastActive = (dateString) => {
      if (!dateString) return 'Never'
      try {
        const date = new Date(dateString)
        const now = new Date()
        const diffMs = now - date
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
        
        if (diffDays === 0) return 'Today'
        if (diffDays === 1) return 'Yesterday'
        if (diffDays < 7) return `${diffDays} days ago`
        if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`
        if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`
        return `${Math.floor(diffDays / 365)} years ago`
      } catch {
        return 'Unknown'
      }
    }

    const canEditUser = (userItem) => {
      if (!isSuper.value) return false
      return true // Super can edit anyone
    }

    const canDeleteUser = (userItem) => {
      if (!isSuper.value) return false
      if (userItem.id === user.value?.id) return false // Cannot delete yourself
      return true
    }

    const canChangeAccessLevel = (userItem) => {
      if (!isSuper.value) return false
      if (userItem?.id === user.value?.id) return false // Cannot change your own access level
      return true
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

    const loadUsers = async () => {
      localLoading.value = true
      localError.value = null
      try {
        const users = await apiService.getUsers()
        allUsers.value = Array.isArray(users) ? users : []
      } catch (err) {
        console.error('Error loading users:', err)
        localError.value = err.message || 'Failed to load users'
      } finally {
        localLoading.value = false
      }
    }

    const editUser = (userItem) => {
      editingUser.value = userItem
      userForm.value = {
        username: userItem.username,
        email: userItem.email || '',
        password: '',
        confirmPassword: '',
        access_level: userItem.access_level,
        is_active: userItem.is_active
      }
      showEditModal.value = true
    }

    const deleteUser = async (userItem) => {
      const confirmed = confirm(`Are you sure you want to delete user "${userItem.username}"? This action cannot be undone.`)
      if (!confirmed) return
      
      localLoading.value = true
      localError.value = null
      try {
        await apiService.deleteUser(userItem.id)
        await loadUsers() // Refresh users list
        console.log('User deleted successfully')
      } catch (err) {
        console.error('Error deleting user:', err)
        localError.value = err.message || 'Failed to delete user'
      } finally {
        localLoading.value = false
      }
    }

    const saveUser = async () => {
      localLoading.value = true
      localError.value = null
      
      try {
        if (editingUser.value) {
          // Update existing user
          const updateData = {
            email: userForm.value.email,
            access_level: userForm.value.access_level,
            is_active: userForm.value.is_active
          }
          await apiService.updateUser(editingUser.value.id, updateData)
        } else {
          // Create new user
          await apiService.createUser(userForm.value)
        }
        
        await loadUsers() // Refresh users list
        closeModal()
        console.log('User saved successfully')
      } catch (err) {
        console.error('Error saving user:', err)
        localError.value = err.message || 'Failed to save user'
      } finally {
        localLoading.value = false
      }
    }

    const closeModal = () => {
      showCreateModal.value = false
      showEditModal.value = false
      editingUser.value = null
      userForm.value = {
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
        access_level: 'user',
        is_active: true
      }
      localError.value = null
    }

    // Load initial data
    onMounted(async () => {
      if (isSuper.value) {
        await loadUsers()
      }
    })

    return {
      user,
      isSuper,
      loading: computed(() => localLoading.value),
      error: computed(() => localError.value),
      showCreateModal,
      showEditModal,
      editingUser,
      userForm,
      searchQuery,
      accessLevelFilter,
      sortField,
      sortDirection,
      currentPage,
      pageSize,
      totalUsers,
      filteredUsers,
      filteredTotalUsers,
      currentPageUsers,
      isFormValid,
      formatDate,
      formatLastActive,
      canEditUser,
      canDeleteUser,
      canChangeAccessLevel,
      toggleSort,
      handlePageChange,
      handlePageSizeChange,
      editUser,
      deleteUser,
      saveUser,
      closeModal,
      localLoading
    }
  }
}
</script>

<style scoped>
/* Reuse most styles from AdminTags.vue with additional user-specific styles */
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

.create-btn:hover {
  background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
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

.table-filters {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
}

.search-input,
.filter-select {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.search-input {
  width: 250px;
  max-width: 100%;
}

.filter-select {
  min-width: 180px;
}

.search-input:focus,
.filter-select:focus {
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

.event-row:hover {
  background: #f8fafc;
}

.username-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-username {
  font-weight: 600;
  color: #2d3748;
}

.current-user-badge {
  background: #dbeafe;
  color: #1e40af;
  font-size: 0.6rem;
  padding: 0.125rem 0.375rem;
  border-radius: 9999px;
  font-weight: 600;
  text-transform: uppercase;
}

.user-email {
  color: #4a5568;
  font-size: 0.875rem;
}

.access-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.access-badge.guest { background: #f3f4f6; color: #374151; }
.access-badge.user { background: #f0f9ff; color: #0284c7; }
.access-badge.editor { background: #f0fdf4; color: #16a34a; }
.access-badge.admin { background: #fef3c7; color: #d97706; }
.access-badge.super { background: #fdf2f8; color: #db2777; }

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-badge.active { background: #f0fdf4; color: #16a34a; }
.status-badge.inactive { background: #fef2f2; color: #dc2626; }

.user-created,
.user-last-active {
  color: #4a5568;
  font-size: 0.875rem;
}

.user-actions {
  display: flex;
  gap: 0.5rem;
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

/* Modal styles - same as AdminTags.vue */
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
  max-width: 600px;
  width: 90%;
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
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #6b7280;
}

.modal-body {
  padding: 0 1.5rem 1.5rem;
}

.form-row {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.form-row .form-group {
  flex: 1;
  min-width: 200px;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #374151;
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.form-input:disabled {
  background-color: #f9fafb;
  color: #6b7280;
}

.form-hint {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #6b7280;
}

.form-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e2e8f0;
}

.cancel-btn,
.submit-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 0.375rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.cancel-btn {
  background: #f3f4f6;
  color: #374151;
}

.cancel-btn:hover {
  background: #e5e7eb;
}

.submit-btn {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.validation-warning {
  background: #fffbeb;
  color: #d97706;
  padding: 0.75rem 1rem;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  border: 1px solid #fed7aa;
  border-left: 4px solid #f59e0b;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.validation-warning::before {
  content: "⚠️";
  font-size: 1rem;
}

.error-message {
  background: #fef2f2;
  color: #dc2626;
  padding: 0.75rem 1rem;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  border: 1px solid #fecaca;
}

@media (max-width: 1024px) {
  .admin-panel {
    padding: 1rem;
  }
  
  .admin-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .table-controls {
    flex-direction: column;
    align-items: stretch;
  }
  
  .table-filters {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .search-input,
  .filter-select {
    width: 100%;
  }
}
</style>