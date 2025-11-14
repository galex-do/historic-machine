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

    <div v-if="stats" class="stats-container">
      <!-- Overall Stats -->
      <div class="stats-section">
        <h3 class="section-title">{{ t('overallStats') }}</h3>
        <div class="kpi-grid">
          <div class="kpi-card highlight">
            <div class="kpi-icon">🌐</div>
            <div class="kpi-content">
              <div class="kpi-value">{{ stats.total_active_users }}</div>
              <div class="kpi-label">{{ t('totalActiveVisitors') }}</div>
              <div class="kpi-description">{{ t('totalActiveVisitorsDesc') }}</div>
            </div>
          </div>

          <div class="kpi-card">
            <div class="kpi-icon">📈</div>
            <div class="kpi-content">
              <div class="kpi-value">{{ stats.anonymous_total_sessions }}</div>
              <div class="kpi-label">{{ t('overallTotalSessions') }}</div>
              <div class="kpi-description">{{ t('overallTotalSessionsDesc') }}</div>
            </div>
          </div>

          <div class="kpi-card">
            <div class="kpi-icon">⌚</div>
            <div class="kpi-content">
              <div class="kpi-value">{{ formatDuration(stats.anonymous_avg_duration) }}</div>
              <div class="kpi-label">{{ t('overallAvgDuration') }}</div>
              <div class="kpi-description">{{ t('overallAvgDurationDesc') }}</div>
            </div>
          </div>

          <div class="kpi-card">
            <div class="kpi-icon">⏱️</div>
            <div class="kpi-content">
              <div class="kpi-value">{{ formatDuration(stats.anonymous_total_time) }}</div>
              <div class="kpi-label">{{ t('overallTotalTime') }}</div>
              <div class="kpi-description">{{ t('overallTotalTimeDesc') }}</div>
            </div>
          </div>

          <div class="kpi-card">
            <div class="kpi-icon">🏆</div>
            <div class="kpi-content">
              <div class="kpi-value">{{ stats.peak_concurrent_sessions }}</div>
              <div class="kpi-label">{{ t('peakConcurrentSessions') }}</div>
              <div class="kpi-description">{{ t('peakConcurrentSessionsDesc') }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Hourly Visitor Graph -->
      <div class="stats-section" v-if="stats.hourly_visitors && stats.hourly_visitors.length > 0">
        <h3 class="section-title">{{ t('hourlyVisitorsTitle') }}</h3>
        <div class="hourly-graph-container">
          <div class="bar-chart">
            <div 
              v-for="(hourStat, index) in stats.hourly_visitors" 
              :key="index"
              class="bar-wrapper"
            >
              <div class="bar-label">{{ formatHourLabel(hourStat.hour) }}</div>
              <div class="bar-column">
                <div 
                  class="bar" 
                  :style="{ height: getBarHeight(hourStat.visitors) + '%' }"
                  :title="formatBarTooltip(hourStat.hour, hourStat.visitors)"
                >
                  <span class="bar-value" v-if="hourStat.visitors > 0">{{ hourStat.visitors }}</span>
                </div>
              </div>
            </div>
          </div>
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
      if (!minutes || minutes === 0) return `0${t('minuteShort')}`
      
      const mins = Math.round(minutes)
      if (mins < 60) return `${mins}${t('minuteShort')}`
      
      const totalHours = Math.floor(mins / 60)
      const remainingMins = mins % 60
      
      // If >= 24 hours (1 day), show days + hours only (more compact)
      if (totalHours >= 24) {
        const days = Math.floor(totalHours / 24)
        const hours = totalHours % 24
        return hours > 0 ? `${days}d ${hours}${t('hourShort')}` : `${days}d`
      }
      
      // Less than 24 hours: show hours + minutes
      return remainingMins > 0 
        ? `${totalHours}${t('hourShort')} ${remainingMins}${t('minuteShort')}` 
        : `${totalHours}${t('hourShort')}`
    }

    const formatTimestamp = (date) => {
      if (!date) return ''
      return date.toLocaleTimeString()
    }

    const formatHourLabel = (hourString) => {
      const date = new Date(hourString)
      return date.getHours().toString().padStart(2, '0') + ':00'
    }

    const formatFullHour = (hourString) => {
      const date = new Date(hourString)
      return date.toLocaleString('en-US', { 
        month: 'short', 
        day: 'numeric', 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    }

    const getBarHeight = (visitors) => {
      if (!stats.value || !stats.value.hourly_visitors) return 0
      const maxVisitors = Math.max(...stats.value.hourly_visitors.map(h => h.visitors), 1)
      return visitors === 0 ? 0 : Math.max((visitors / maxVisitors) * 100, 5)
    }

    const formatBarTooltip = (hourString, visitors) => {
      const visitorText = visitors === 1 ? t('visitor') : t('visitors')
      return `${formatFullHour(hourString)}: ${visitors} ${visitorText}`
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
      formatTimestamp,
      formatHourLabel,
      formatFullHour,
      getBarHeight,
      formatBarTooltip
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

.stats-container {
  margin-bottom: 2rem;
}

.stats-section {
  margin-bottom: 2.5rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e2e8f0;
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
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

.kpi-card.highlight {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.kpi-card.highlight .kpi-icon {
  background: rgba(255, 255, 255, 0.2);
}

.kpi-card.highlight .kpi-value,
.kpi-card.highlight .kpi-label,
.kpi-card.highlight .kpi-description {
  color: white;
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

.hourly-graph-container {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.bar-chart {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  height: 250px;
  gap: 2px;
  padding: 1rem 0;
}

.bar-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.bar-column {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar {
  width: 100%;
  background: linear-gradient(180deg, #4299e1 0%, #3182ce 100%);
  border-radius: 4px 4px 0 0;
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 4px;
  position: relative;
  min-height: 5px;
}

.bar:hover {
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  transform: scaleY(1.05);
  box-shadow: 0 4px 8px rgba(66, 153, 225, 0.3);
}

.bar-value {
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.bar-label {
  font-size: 0.75rem;
  color: #4a5568;
  font-weight: 500;
  white-space: nowrap;
  transform: rotate(-45deg);
  transform-origin: center;
  margin-top: 1.5rem;
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

  .hourly-graph-container {
    padding: 1rem;
  }

  .bar-chart {
    height: 200px;
    gap: 1px;
  }

  .bar-column {
    height: 150px;
  }

  .bar-label {
    font-size: 0.65rem;
  }

  .bar-value {
    font-size: 0.65rem;
  }
}
</style>
