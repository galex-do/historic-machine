<template>
  <div class="table-pagination">
    <div class="pagination-info">
      <span class="results-info">
        Showing {{ startItem }} to {{ endItem }} of {{ totalItems }} entries
      </span>
      <div class="page-size-selector">
        <label for="page-size">Show:</label>
        <select 
          id="page-size" 
          :value="pageSize" 
          @change="$emit('update:pageSize', parseInt($event.target.value))"
          class="page-size-select"
        >
          <option value="10" selected>10</option>
          <option value="25">25</option>
          <option value="50">50</option>
        </select>
        <span>entries</span>
      </div>
    </div>
    
    <div class="pagination-controls" v-if="totalPages > 1">
      <button 
        class="pagination-btn" 
        :disabled="currentPage === 1"
        @click="$emit('update:currentPage', 1)"
        title="First page"
      >
        ⟪
      </button>
      
      <button 
        class="pagination-btn" 
        :disabled="currentPage === 1"
        @click="$emit('update:currentPage', currentPage - 1)"
        title="Previous page"
      >
        ⟨
      </button>
      
      <div class="page-numbers">
        <button
          v-for="page in visiblePages"
          :key="page"
          class="pagination-btn page-number"
          :class="{ 'active': page === currentPage }"
          @click="$emit('update:currentPage', page)"
        >
          {{ page }}
        </button>
      </div>
      
      <button 
        class="pagination-btn" 
        :disabled="currentPage === totalPages"
        @click="$emit('update:currentPage', currentPage + 1)"
        title="Next page"
      >
        ⟩
      </button>
      
      <button 
        class="pagination-btn" 
        :disabled="currentPage === totalPages"
        @click="$emit('update:currentPage', totalPages)"
        title="Last page"
      >
        ⟫
      </button>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'TablePagination',
  props: {
    currentPage: {
      type: Number,
      required: true
    },
    pageSize: {
      type: Number,
      required: true
    },
    totalItems: {
      type: Number,
      required: true
    }
  },
  emits: ['update:currentPage', 'update:pageSize'],
  setup(props) {
    const totalPages = computed(() => {
      return Math.ceil(props.totalItems / props.pageSize)
    })
    
    const startItem = computed(() => {
      if (props.totalItems === 0) return 0
      return (props.currentPage - 1) * props.pageSize + 1
    })
    
    const endItem = computed(() => {
      const end = props.currentPage * props.pageSize
      return Math.min(end, props.totalItems)
    })
    
    const visiblePages = computed(() => {
      const pages = []
      const total = totalPages.value
      const current = props.currentPage
      
      if (total <= 7) {
        // Show all pages if 7 or fewer
        for (let i = 1; i <= total; i++) {
          pages.push(i)
        }
      } else {
        // Show smart pagination
        if (current <= 4) {
          // Show first 5 pages + ... + last page
          for (let i = 1; i <= 5; i++) {
            pages.push(i)
          }
          if (total > 6) {
            pages.push('...')
            pages.push(total)
          }
        } else if (current >= total - 3) {
          // Show first page + ... + last 5 pages
          pages.push(1)
          if (total > 6) {
            pages.push('...')
          }
          for (let i = total - 4; i <= total; i++) {
            pages.push(i)
          }
        } else {
          // Show first + ... + current-1, current, current+1 + ... + last
          pages.push(1)
          pages.push('...')
          for (let i = current - 1; i <= current + 1; i++) {
            pages.push(i)
          }
          pages.push('...')
          pages.push(total)
        }
      }
      
      return pages.filter(page => page !== '...' || pages.indexOf(page) === pages.lastIndexOf(page))
    })
    
    return {
      totalPages,
      startItem,
      endItem,
      visiblePages
    }
  }
}
</script>

<style scoped>
.table-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  border-top: 1px solid #e5e7eb;
  margin-top: 1rem;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.results-info {
  color: #6b7280;
  font-size: 0.875rem;
}

.page-size-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #6b7280;
}

.page-size-select {
  padding: 0.25rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  background-color: white;
  font-size: 0.875rem;
  cursor: pointer;
}

.page-size-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.pagination-btn {
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  background-color: white;
  color: #374151;
  font-size: 0.875rem;
  cursor: pointer;
  border-radius: 0.375rem;
  transition: all 0.2s ease;
  min-width: 2.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

.pagination-btn:disabled {
  color: #9ca3af;
  cursor: not-allowed;
  background-color: #f9fafb;
}

.pagination-btn.active {
  background-color: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.pagination-btn.active:hover {
  background-color: #2563eb;
}

.page-numbers {
  display: flex;
  gap: 0.25rem;
}

.page-number {
  min-width: 2.5rem;
}

@media (max-width: 768px) {
  .table-pagination {
    flex-direction: column;
    gap: 1rem;
  }
  
  .pagination-info {
    flex-direction: column;
    gap: 0.5rem;
    text-align: center;
  }
  
  .pagination-controls {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>