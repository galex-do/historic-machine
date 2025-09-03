<template>
  <div class="datasets-panel">
    <!-- Header Section -->
    <div class="datasets-header">
      <h1 class="page-title">Event Datasets</h1>
      <p class="page-description">
        Manage uploaded event datasets. Deleting a dataset will remove all events imported from it.
      </p>
    </div>

    <!-- Loading State -->
    <div v-if="localLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Loading datasets...</p>
    </div>

    <!-- Error State -->
    <div v-if="localError" class="error-state">
      <div class="error-icon">⚠️</div>
      <p>{{ localError }}</p>
      <button @click="fetchDatasets" class="retry-button">Try Again</button>
    </div>

    <!-- Main Content -->
    <div v-if="!localLoading && !localError" class="datasets-content">
      
      <!-- Datasets Table -->
      <div class="datasets-table-container">
        <table class="datasets-table">
          <thead>
            <tr>
              <th>Filename</th>
              <th>Description</th>
              <th>Event Count</th>
              <th>Uploaded By</th>
              <th>Upload Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="datasets.length === 0">
              <td colspan="6" class="no-datasets">
                <div class="empty-state">
                  <div class="empty-icon">📊</div>
                  <h3>No datasets found</h3>
                  <p>Import events to create your first dataset.</p>
                </div>
              </td>
            </tr>
            <tr v-for="dataset in datasets" :key="dataset.id" class="dataset-row">
              <td class="dataset-filename">
                <div class="filename-cell">
                  <span class="filename">{{ dataset.filename }}</span>
                </div>
              </td>
              <td class="dataset-description">
                {{ dataset.description || 'No description' }}
              </td>
              <td class="dataset-count">
                <span class="count-badge">{{ dataset.event_count }}</span>
              </td>
              <td class="dataset-uploader">
                {{ dataset.uploaded_by_username || `User ${dataset.uploaded_by}` }}
              </td>
              <td class="dataset-date">
                {{ formatDate(dataset.created_at) }}
              </td>
              <td class="dataset-actions">
                <button 
                  @click="confirmDelete(dataset)"
                  class="delete-button"
                  :title="`Delete dataset: ${dataset.filename}`"
                >
                  🗑️
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Delete Dataset</h3>
          <button @click="closeDeleteModal" class="modal-close">×</button>
        </div>
        
        <div class="modal-body">
          <div class="warning-icon">⚠️</div>
          <p class="warning-text">
            Are you sure you want to delete the dataset 
            <strong>"{{ datasetToDelete?.filename }}"</strong>?
          </p>
          <p class="warning-details">
            This will permanently delete:
          </p>
          <ul class="warning-list">
            <li><strong>{{ datasetToDelete?.event_count }} events</strong> imported from this dataset</li>
            <li>The dataset record itself</li>
          </ul>
          <p class="warning-note">
            <strong>Note:</strong> Tags created during import will be preserved.
          </p>
        </div>
        
        <div class="modal-footer">
          <button @click="closeDeleteModal" class="cancel-button">Cancel</button>
          <button 
            @click="deleteDataset" 
            class="confirm-delete-button"
            :disabled="localLoading"
          >
            {{ localLoading ? 'Deleting...' : 'Delete Dataset' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import apiService from '@/services/api.js'

export default {
  name: 'DatasetsPanel',
  setup() {
    const datasets = ref([])
    const localLoading = ref(false)
    const localError = ref(null)
    const showDeleteModal = ref(false)
    const datasetToDelete = ref(null)

    const fetchDatasets = async () => {
      localLoading.value = true
      localError.value = null
      
      try {
        const response = await apiService.getDatasets()
        datasets.value = response || []
        console.log('Loaded datasets:', datasets.value.length)
      } catch (err) {
        console.error('Error fetching datasets:', err)
        localError.value = err.message || 'Failed to load datasets'
      } finally {
        localLoading.value = false
      }
    }

    const confirmDelete = (dataset) => {
      datasetToDelete.value = dataset
      showDeleteModal.value = true
    }

    const closeDeleteModal = () => {
      showDeleteModal.value = false
      datasetToDelete.value = null
    }

    const deleteDataset = async () => {
      if (!datasetToDelete.value) return
      
      localLoading.value = true
      localError.value = null
      
      try {
        await apiService.delete(`/datasets/${datasetToDelete.value.id}`)
        
        // Remove from local list
        datasets.value = datasets.value.filter(d => d.id !== datasetToDelete.value.id)
        
        console.log('Dataset deleted successfully:', datasetToDelete.value.filename)
        closeDeleteModal()
        
      } catch (err) {
        console.error('Error deleting dataset:', err)
        localError.value = err.message || 'Failed to delete dataset'
      } finally {
        localLoading.value = false
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'Unknown'
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    onMounted(() => {
      fetchDatasets()
    })

    return {
      datasets,
      localLoading,
      localError,
      showDeleteModal,
      datasetToDelete,
      fetchDatasets,
      confirmDelete,
      closeDeleteModal,
      deleteDataset,
      formatDate
    }
  }
}
</script>

<style scoped>
.datasets-panel {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.datasets-header {
  margin-bottom: 2rem;
  text-align: center;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 0.5rem;
}

.page-description {
  font-size: 1.1rem;
  color: #64748b;
  max-width: 600px;
  margin: 0 auto;
}

.loading-state, .error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.loading-spinner {
  width: 3rem;
  height: 3rem;
  border: 4px solid #e2e8f0;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.retry-button {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
}

.retry-button:hover {
  background: #2563eb;
}

.datasets-table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.datasets-table {
  width: 100%;
  border-collapse: collapse;
}

.datasets-table th {
  background: #f1f5f9;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

.datasets-table td {
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
}

.dataset-row:hover {
  background: #f8fafc;
}

.filename-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filename {
  font-weight: 600;
  color: #1e293b;
}

.count-badge {
  background: #dbeafe;
  color: #1e40af;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-weight: 600;
  font-size: 0.875rem;
}

.delete-button {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  padding: 0.5rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.2s;
}

.delete-button:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.no-datasets {
  text-align: center;
  padding: 3rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.empty-icon {
  font-size: 4rem;
  opacity: 0.5;
}

.empty-state h3 {
  color: #6b7280;
  margin: 0;
}

.empty-state p {
  color: #9ca3af;
  margin: 0;
}

/* Modal Styles */
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
  border-radius: 12px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
}

.modal-close:hover {
  background: #f3f4f6;
  color: #374151;
}

.modal-body {
  padding: 1.5rem;
  text-align: center;
}

.warning-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.warning-text {
  font-size: 1.1rem;
  color: #374151;
  margin-bottom: 1rem;
}

.warning-details {
  color: #6b7280;
  margin-bottom: 0.5rem;
  text-align: left;
}

.warning-list {
  text-align: left;
  color: #dc2626;
  margin: 1rem 0;
  padding-left: 1.5rem;
}

.warning-note {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 6px;
  padding: 0.75rem;
  color: #0c4a6e;
  font-size: 0.9rem;
  margin-top: 1rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1.5rem;
  border-top: 1px solid #e5e7eb;
}

.cancel-button {
  padding: 0.75rem 1.5rem;
  background: #f9fafb;
  border: 1px solid #d1d5db;
  color: #374151;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.cancel-button:hover {
  background: #f3f4f6;
}

.confirm-delete-button {
  padding: 0.75rem 1.5rem;
  background: #dc2626;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.confirm-delete-button:hover:not(:disabled) {
  background: #b91c1c;
}

.confirm-delete-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Responsive Design */
@media (max-width: 768px) {
  .datasets-panel {
    padding: 1rem;
  }
  
  .page-title {
    font-size: 2rem;
  }
  
  .datasets-table-container {
    overflow-x: auto;
  }
  
  .datasets-table {
    min-width: 600px;
  }
  
  .datasets-table th,
  .datasets-table td {
    padding: 0.75rem 0.5rem;
    font-size: 0.9rem;
  }
  
  .modal-content {
    margin: 1rem;
    width: calc(100% - 2rem);
  }
}
</style>