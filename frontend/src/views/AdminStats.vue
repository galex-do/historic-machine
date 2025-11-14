<template>
  <div class="admin-panel stats-page">
    <div class="admin-header">
      <div class="admin-title">
        <h2>{{ t('statsTitle') }}</h2>
        <p class="admin-subtitle">{{ t('statsSubtitle') }}</p>
      </div>
      <div class="stats-controls">
        <button @click="refreshStats" class="refresh-btn" :disabled="loading">
          <span class="btn-icon">🔄</span>
          {{ t('refresh') }}
        </button>
      </div>
    </div>

    <div v-if="loading && !stats" class="loading-state">
      <div class="spinner"></div>
      <p>{{ t('loadingStats') }}</p>
    </div>

    <div v-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="refreshStats" class="retry-btn">{{ t('retryLoad') }}</button>
    </div>

    <div v-if="stats" class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-icon">👥</div>
        <div class="kpi-content">
          <div class="kpi-value">{{ stats.active_users }}</div>
          <div class="kpi-label">{{ t('activeUsers') }}</div>
          <div class="kpi-description">{{ t('activeUsersDesc') }}</div>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon">🔓</div>
        <div class="kpi-content">
          <div class="kpi-value">{{ stats.active_sessions }}</div>
          <div class="kpi-label">{{ t('activeSessions') }}</div>
          <div class="kpi-description">{{ t('activeSessionsDesc') }}</div>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon">📊</div>
        <div class="kpi-content">
          <div class="kpi-value">{{ stats.total_sessions }}</div>
          <div class="kpi-label">{{ t('totalSessions') }}</div>
          <div class="kpi-description">{{ t('totalSessionsDesc') }}</div>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon">⏱️</div>
        <div class="kpi-content">
          <div class="kpi-value">{{ formatDuration(stats.avg_duration) }}</div>
          <div class="kpi-label">{{ t('avgDuration') }}</div>
          <div class="kpi-description">{{ t('avgDurationDesc') }}</div>
        </div>
      </div>
    </div>

    <div v-if="stats" class="last-updated">
      {{ t('lastUpdated') }}: {{ formatTimestamp(lastUpdated) }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import authService from '@/services/authService.js'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'AdminStats',
  setup() {
    const { t } = useLocale()
    const stats = ref(null)
    const loading = ref(false)
    const error = ref(null)
    const lastUpdated = ref(null)

    const refreshStats = async () => {
      loading.value = true
      error.value = null

      try {
        const data = await authService.getSessionStats()
        stats.value = data
        lastUpdated.value = new Date()
      } catch (err) {
        console.error('Failed to fetch session stats:', err)
        error.value = t('statsLoadError')
      } finally {
        loading.value = false
      }
    }

    const formatDuration = (minutes) => {
      if (!minutes || minutes === 0) return '0 min'
      
      const mins = Math.round(minutes)
      if (mins < 60) return `${mins} min`
      
      const hours = Math.floor(mins / 60)
      const remainingMins = mins % 60
      return remainingMins > 0 ? `${hours}h ${remainingMins}m` : `${hours}h`
    }

    const formatTimestamp = (date) => {
      if (!date) return ''
      return date.toLocaleTimeString()
    }

    onMounted(() => {
      refreshStats()
    })

    return {
      t,
      stats,
      loading,
      error,
      lastUpdated,
      refreshStats,
      formatDuration,
      formatTimestamp
    }
  }
}
</script>

<style scoped>
.stats-page {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.admin-title h2 {
  font-size: 2rem;
  font-weight: 700;
  color: #1a202c;
  margin-bottom: 0.5rem;
}

.admin-subtitle {
  font-size: 1rem;
  color: #718096;
}

.stats-controls {
  display: flex;
  gap: 1rem;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: #3182ce;
  transform: translateY(-1px);
}

.refresh-btn:disabled {
  background: #cbd5e0;
  cursor: not-allowed;
}

.btn-icon {
  font-size: 1.2rem;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 3rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #e2e8f0;
  border-top-color: #4299e1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state p {
  color: #e53e3e;
  margin-bottom: 1rem;
}

.retry-btn {
  padding: 0.75rem 1.5rem;
  background: #e53e3e;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.retry-btn:hover {
  background: #c53030;
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.kpi-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 1rem;
  align-items: flex-start;
  transition: transform 0.2s, box-shadow 0.2s;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.kpi-icon {
  font-size: 2.5rem;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ebf8ff;
  border-radius: 12px;
  flex-shrink: 0;
}

.kpi-content {
  flex: 1;
}

.kpi-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1a202c;
  line-height: 1;
  margin-bottom: 0.5rem;
}

.kpi-label {
  font-size: 1rem;
  font-weight: 600;
  color: #4a5568;
  margin-bottom: 0.25rem;
}

.kpi-description {
  font-size: 0.875rem;
  color: #718096;
}

.last-updated {
  text-align: center;
  color: #718096;
  font-size: 0.875rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

@media (max-width: 768px) {
  .stats-page {
    padding: 1rem;
  }

  .admin-header {
    flex-direction: column;
    gap: 1rem;
  }

  .admin-title h2 {
    font-size: 1.5rem;
  }

  .kpi-grid {
    grid-template-columns: 1fr;
  }

  .kpi-card {
    padding: 1rem;
  }

  .kpi-icon {
    font-size: 2rem;
    width: 50px;
    height: 50px;
  }

  .kpi-value {
    font-size: 2rem;
  }
}
</style>
