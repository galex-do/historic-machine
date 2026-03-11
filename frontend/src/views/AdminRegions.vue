<template>
  <div class="admin-panel">
    <div class="admin-header">
      <div class="admin-title">
        <h2>{{ t('adminRegionsTitle') }}</h2>
        <p class="admin-subtitle">{{ t('adminRegionsSubtitle') }}</p>
      </div>
      <div class="action-buttons" v-if="canAccessAdmin">
        <button @click="openCreateModal" class="create-btn">
          <span class="btn-icon">➕</span>
          {{ t('createRegion') }}
        </button>
      </div>
    </div>

    <div class="table-container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading regions...</p>
      </div>

      <div v-if="!loading && totalRegions > 0" class="table-controls">
        <div class="table-filters">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t('searchRegions')"
            class="search-input"
          />
        </div>
        <TablePagination
          :current-page="currentPage"
          :page-size="pageSize"
          :total-items="filteredTotalRegions"
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
              {{ t('columnName') }}
              <span class="sort-indicator">
                <span v-if="sortField === 'name'" class="sort-arrow">
                  {{ sortDirection === 'asc' ? '▲' : '▼' }}
                </span>
                <span v-else class="sort-placeholder">⇅</span>
              </span>
            </th>
            <th>{{ t('columnColor') }}</th>
            <th>{{ t('linkedTemplates') }}</th>
            <th>{{ t('columnActions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="region in currentPageRegions" :key="region.id" class="event-row">
            <td class="region-name">{{ region.name || region.name_en }}</td>
            <td class="region-color-cell">
              <div class="color-preview">
                <span
                  class="color-swatch"
                  :style="{ backgroundColor: region.color, opacity: region.fill_opacity || 0.2 }"
                ></span>
                <span class="color-code">{{ region.color }}</span>
              </div>
            </td>
            <td class="region-templates">
              <span class="template-count">{{ region.template_ids ? region.template_ids.length : 0 }}</span>
            </td>
            <td class="region-actions">
              <div class="actions-wrapper">
                <button @click="editRegion(region)" class="action-btn edit-btn" title="Edit">
                  ✏️
                </button>
                <button @click="deleteRegion(region)" class="action-btn delete-btn" title="Delete">
                  🗑️
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="!loading && totalRegions === 0" class="empty-state">
        <p>{{ t('noRegions') }}</p>
      </div>

      <div v-if="!loading && totalRegions > 0 && filteredRegions.length === 0" class="empty-state">
        <p>No regions match your search criteria.</p>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click="closeModal">
        <div class="modal-content admin-modal wide-modal" @click.stop>
          <div class="modal-header">
            <h3>{{ editingRegion ? t('editRegion') : t('createRegion') }}</h3>
            <button @click="closeModal" class="close-btn">&times;</button>
          </div>

          <div class="modal-body">
            <div v-if="error" class="error-message">{{ error }}</div>

            <form @submit.prevent="saveRegion">
              <div class="form-row">
                <div class="form-group">
                  <label for="region-name">{{ t('regionName') }} (English) *</label>
                  <input
                    id="region-name"
                    v-model="regionForm.name_en"
                    type="text"
                    required
                    class="form-input"
                    placeholder="Region name in English"
                  />
                </div>
                <div class="form-group">
                  <label for="region-name-ru">{{ t('regionName') }} (Русский)</label>
                  <input
                    id="region-name-ru"
                    v-model="regionForm.name_ru"
                    type="text"
                    class="form-input"
                    placeholder="Название на русском"
                  />
                </div>
              </div>

              <div class="form-group">
                <label for="region-description">{{ t('regionDescription') }}</label>
                <textarea
                  id="region-description"
                  v-model="regionForm.description"
                  class="form-textarea"
                  placeholder="Region description (optional)"
                  rows="2"
                ></textarea>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label for="region-color">{{ t('regionColor') }}</label>
                  <div class="color-input-container">
                    <input
                      id="region-color"
                      v-model="regionForm.color"
                      type="color"
                      class="color-input"
                    />
                    <input
                      v-model="regionForm.color"
                      type="text"
                      class="form-input color-text-input"
                      placeholder="#4f46e5"
                      pattern="^#[0-9A-Fa-f]{6}$"
                    />
                  </div>
                </div>
                <div class="form-group">
                  <label for="region-border-color">{{ t('borderColor') }}</label>
                  <div class="color-input-container">
                    <input
                      id="region-border-color"
                      v-model="regionForm.border_color"
                      type="color"
                      class="color-input"
                    />
                    <input
                      v-model="regionForm.border_color"
                      type="text"
                      class="form-input color-text-input"
                      placeholder="#4f46e5"
                      pattern="^#[0-9A-Fa-f]{6}$"
                    />
                  </div>
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label for="region-opacity">{{ t('fillOpacity') }}: {{ regionForm.fill_opacity }}</label>
                  <input
                    id="region-opacity"
                    v-model.number="regionForm.fill_opacity"
                    type="range"
                    min="0"
                    max="1"
                    step="0.05"
                    class="opacity-slider"
                  />
                </div>
                <div class="form-group">
                  <label for="region-border-width">{{ t('borderWidth') }}</label>
                  <input
                    id="region-border-width"
                    v-model.number="regionForm.border_width"
                    type="number"
                    min="0.5"
                    max="10"
                    step="0.5"
                    class="form-input width-input"
                  />
                </div>
              </div>

              <div class="form-group">
                <label>{{ t('linkedTemplates') }}</label>
                <div class="templates-checklist">
                  <div v-if="allTemplates.length === 0" class="no-templates">
                    {{ t('noTemplatesAvailable') }}
                  </div>
                  <label
                    v-for="tmpl in allTemplates"
                    :key="tmpl.id"
                    class="template-checkbox"
                  >
                    <input
                      type="checkbox"
                      :value="tmpl.id"
                      v-model="regionForm.template_ids"
                    />
                    <span>{{ tmpl.name || tmpl.name_en }}</span>
                  </label>
                </div>
              </div>

              <div class="form-group">
                <label>{{ t('drawRegion') }}</label>
                <div ref="drawMapContainer" class="draw-map-container"></div>
              </div>

              <div class="form-group">
                <label>{{ t('importGeoJSON') }}</label>
                <textarea
                  v-model="geoJSONInput"
                  class="form-textarea geojson-input"
                  :placeholder="t('geoJSONPlaceholder')"
                  rows="4"
                ></textarea>
                <button type="button" @click="applyGeoJSON" class="geojson-btn">
                  {{ t('pasteGeoJSON') }}
                </button>
                <div v-if="geoJSONError" class="geojson-error">{{ geoJSONError }}</div>
              </div>

              <div class="form-actions">
                <button type="button" @click="closeModal" class="cancel-btn">{{ t('cancel') }}</button>
                <button type="submit" class="submit-btn" :disabled="localLoading">
                  {{ localLoading ? t('saving') : (editingRegion ? t('editRegion') : t('createRegion')) }}
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
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useAuth } from '@/composables/useAuth.js'
import { useLocale } from '@/composables/useLocale.js'
import apiService from '@/services/api.js'
import TablePagination from '@/components/TablePagination.vue'
import L from 'leaflet'
import 'leaflet-draw'
import 'leaflet/dist/leaflet.css'
import 'leaflet-draw/dist/leaflet.draw.css'

export default {
  name: 'AdminRegions',
  components: {
    TablePagination
  },
  setup() {
    const { canAccessAdmin } = useAuth()
    const { t } = useLocale()

    const regions = ref([])
    const allTemplates = ref([])
    const loading = ref(false)
    const localLoading = ref(false)
    const error = ref(null)
    const searchQuery = ref('')

    const showModal = ref(false)
    const editingRegion = ref(null)

    const sortField = ref('name')
    const sortDirection = ref('asc')
    const currentPage = ref(1)
    const pageSize = ref(10)

    const drawMapContainer = ref(null)
    let drawMap = null
    let drawnItems = null
    let drawControl = null

    const geoJSONInput = ref('')
    const geoJSONError = ref(null)

    const regionForm = ref({
      name_en: '',
      name_ru: '',
      description: '',
      description_en: '',
      description_ru: '',
      color: '#4f46e5',
      border_color: '#4f46e5',
      fill_opacity: 0.2,
      border_width: 2,
      template_ids: [],
      geojson: null
    })

    const filteredRegions = computed(() => {
      if (!regions.value || !Array.isArray(regions.value)) return []
      let filtered = regions.value
      if (searchQuery.value.trim()) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(r =>
          (r.name && r.name.toLowerCase().includes(query)) ||
          (r.name_en && r.name_en.toLowerCase().includes(query)) ||
          (r.name_ru && r.name_ru.toLowerCase().includes(query))
        )
      }
      return [...filtered].sort((a, b) => {
        let aVal = (a[sortField.value] || '').toString().toLowerCase()
        let bVal = (b[sortField.value] || '').toString().toLowerCase()
        if (sortDirection.value === 'asc') {
          return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
        } else {
          return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
        }
      })
    })

    const filteredTotalRegions = computed(() => filteredRegions.value.length)
    const totalRegions = computed(() => regions.value?.length || 0)

    const currentPageRegions = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filteredRegions.value.slice(start, end)
    })

    const toggleSort = (field) => {
      if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
      } else {
        sortField.value = field
        sortDirection.value = 'asc'
      }
    }

    const handlePageChange = (page) => { currentPage.value = page }
    const handlePageSizeChange = (size) => { pageSize.value = size; currentPage.value = 1 }

    const loadData = async () => {
      loading.value = true
      error.value = null
      try {
        const [regionsData, templatesData] = await Promise.all([
          apiService.getAllRegions(),
          apiService.getAllTemplates()
        ])
        regions.value = regionsData || []
        allTemplates.value = templatesData || []
      } catch (err) {
        console.error('Error loading regions:', err)
        error.value = t('failedToLoadRegions')
      } finally {
        loading.value = false
      }
    }

    const initDrawMap = () => {
      if (drawMap) {
        drawMap.remove()
        drawMap = null
      }

      nextTick(() => {
        if (!drawMapContainer.value) return

        drawMap = L.map(drawMapContainer.value).setView([30, 35], 3)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          attribution: '&copy; OpenStreetMap contributors',
          maxZoom: 18
        }).addTo(drawMap)

        drawnItems = new L.FeatureGroup()
        drawMap.addLayer(drawnItems)

        drawControl = new L.Control.Draw({
          position: 'topright',
          draw: {
            polygon: {
              allowIntersection: false,
              showArea: true,
              shapeOptions: {
                color: regionForm.value.border_color,
                fillColor: regionForm.value.color,
                fillOpacity: regionForm.value.fill_opacity,
                weight: regionForm.value.border_width
              }
            },
            polyline: false,
            circle: false,
            rectangle: false,
            marker: false,
            circlemarker: false
          },
          edit: {
            featureGroup: drawnItems,
            remove: true
          }
        })
        drawMap.addControl(drawControl)

        drawMap.on(L.Draw.Event.CREATED, (e) => {
          drawnItems.addLayer(e.layer)
          updateGeoJSONFromDrawn()
        })

        drawMap.on(L.Draw.Event.EDITED, () => {
          updateGeoJSONFromDrawn()
        })

        drawMap.on(L.Draw.Event.DELETED, () => {
          updateGeoJSONFromDrawn()
        })

        if (regionForm.value.geojson) {
          try {
            const geo = typeof regionForm.value.geojson === 'string'
              ? JSON.parse(regionForm.value.geojson)
              : regionForm.value.geojson
            const geoLayer = L.geoJSON(geo, {
              style: {
                color: regionForm.value.border_color,
                fillColor: regionForm.value.color,
                fillOpacity: regionForm.value.fill_opacity,
                weight: regionForm.value.border_width
              }
            })
            geoLayer.eachLayer((layer) => {
              drawnItems.addLayer(layer)
            })
            if (drawnItems.getLayers().length > 0) {
              drawMap.fitBounds(drawnItems.getBounds(), { padding: [20, 20] })
            }
          } catch (e) {
            console.error('Error loading GeoJSON onto map:', e)
          }
        }

        setTimeout(() => {
          drawMap.invalidateSize()
        }, 200)
      })
    }

    const updateGeoJSONFromDrawn = () => {
      if (!drawnItems) return
      const layers = drawnItems.getLayers()
      if (layers.length === 0) {
        regionForm.value.geojson = null
        return
      }
      const geojson = drawnItems.toGeoJSON()
      regionForm.value.geojson = geojson
    }

    const applyGeoJSON = () => {
      geoJSONError.value = null
      if (!geoJSONInput.value.trim()) return
      try {
        const parsed = JSON.parse(geoJSONInput.value)
        if (!parsed.type || !['Feature', 'FeatureCollection', 'Polygon', 'MultiPolygon', 'GeometryCollection'].includes(parsed.type)) {
          geoJSONError.value = t('invalidGeoJSON')
          return
        }
        regionForm.value.geojson = parsed

        if (drawnItems) {
          drawnItems.clearLayers()
          const geoLayer = L.geoJSON(parsed, {
            style: {
              color: regionForm.value.border_color,
              fillColor: regionForm.value.color,
              fillOpacity: regionForm.value.fill_opacity,
              weight: regionForm.value.border_width
            }
          })
          geoLayer.eachLayer((layer) => {
            drawnItems.addLayer(layer)
          })
          if (drawnItems.getLayers().length > 0 && drawMap) {
            drawMap.fitBounds(drawnItems.getBounds(), { padding: [20, 20] })
          }
        }
        geoJSONInput.value = ''
      } catch (e) {
        geoJSONError.value = t('invalidGeoJSON')
      }
    }

    const openCreateModal = () => {
      editingRegion.value = null
      regionForm.value = {
        name_en: '',
        name_ru: '',
        description: '',
        description_en: '',
        description_ru: '',
        color: '#4f46e5',
        border_color: '#4f46e5',
        fill_opacity: 0.2,
        border_width: 2,
        template_ids: [],
        geojson: null
      }
      geoJSONInput.value = ''
      geoJSONError.value = null
      error.value = null
      showModal.value = true
      nextTick(() => initDrawMap())
    }

    const editRegion = async (region) => {
      editingRegion.value = region
      regionForm.value = {
        name_en: region.name_en || region.name || '',
        name_ru: region.name_ru || '',
        description: region.description || '',
        description_en: region.description_en || '',
        description_ru: region.description_ru || '',
        color: region.color || '#4f46e5',
        border_color: region.border_color || '#4f46e5',
        fill_opacity: region.fill_opacity != null ? region.fill_opacity : 0.2,
        border_width: region.border_width != null ? region.border_width : 2,
        template_ids: region.template_ids || [],
        geojson: region.geojson || null
      }
      geoJSONInput.value = ''
      geoJSONError.value = null
      error.value = null
      showModal.value = true
      nextTick(() => initDrawMap())
    }

    const deleteRegion = async (region) => {
      const confirmed = confirm(t('confirmDeleteRegion'))
      if (!confirmed) return

      localLoading.value = true
      error.value = null
      try {
        await apiService.deleteRegion(region.id)
        await loadData()
      } catch (err) {
        console.error('Error deleting region:', err)
        error.value = t('failedToDeleteRegion')
      } finally {
        localLoading.value = false
      }
    }

    const saveRegion = async () => {
      localLoading.value = true
      error.value = null

      try {
        updateGeoJSONFromDrawn()

        const payload = {
          name: regionForm.value.name_en,
          name_en: regionForm.value.name_en,
          name_ru: regionForm.value.name_ru,
          description: regionForm.value.description,
          description_en: regionForm.value.description_en,
          description_ru: regionForm.value.description_ru,
          color: regionForm.value.color,
          border_color: regionForm.value.border_color,
          fill_opacity: regionForm.value.fill_opacity,
          border_width: regionForm.value.border_width,
          geojson: regionForm.value.geojson,
          template_ids: regionForm.value.template_ids
        }

        if (editingRegion.value) {
          await apiService.updateRegion(editingRegion.value.id, payload)
        } else {
          await apiService.createRegion(payload)
        }

        await loadData()
        closeModal()
      } catch (err) {
        console.error('Error saving region:', err)
        error.value = editingRegion.value ? t('failedToUpdateRegion') : t('failedToCreateRegion')
      } finally {
        localLoading.value = false
      }
    }

    const closeModal = () => {
      showModal.value = false
      editingRegion.value = null
      error.value = null
      geoJSONInput.value = ''
      geoJSONError.value = null
      if (drawMap) {
        drawMap.remove()
        drawMap = null
        drawnItems = null
        drawControl = null
      }
    }

    onMounted(async () => {
      await loadData()
    })

    return {
      canAccessAdmin,
      loading,
      localLoading,
      error,
      showModal,
      editingRegion,
      regionForm,
      searchQuery,
      sortField,
      sortDirection,
      currentPage,
      pageSize,
      totalRegions,
      filteredRegions,
      filteredTotalRegions,
      currentPageRegions,
      allTemplates,
      drawMapContainer,
      geoJSONInput,
      geoJSONError,
      toggleSort,
      handlePageChange,
      handlePageSizeChange,
      openCreateModal,
      editRegion,
      deleteRegion,
      saveRegion,
      closeModal,
      applyGeoJSON,
      t
    }
  }
}
</script>

<style scoped>
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

.search-input {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  width: 300px;
  max-width: 100%;
}

.search-input:focus {
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
}

.sort-placeholder {
  color: #9ca3af;
}

.event-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background-color 0.15s ease;
}

.event-row:hover {
  background-color: #f8fafc;
}

.color-preview {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.color-swatch {
  display: inline-block;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #d1d5db;
}

.color-code {
  font-family: monospace;
  font-size: 0.85rem;
  color: #4a5568;
}

.template-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 0.5rem;
  background: #eef2ff;
  color: #4f46e5;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
}

.actions-wrapper {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  padding: 0.4rem 0.6rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
}

.action-btn:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem;
  color: #718096;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid #e2e8f0;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #718096;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #374151;
  font-size: 0.9rem;
}

.form-input {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.9rem;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.form-textarea {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.9rem;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
}

.form-textarea:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.color-input-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.color-input {
  width: 3rem;
  height: 2.5rem;
  padding: 2px;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  cursor: pointer;
  background: white;
}

.color-text-input {
  flex: 1;
  max-width: 200px;
}

.opacity-slider {
  width: 100%;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  background: #e2e8f0;
  border-radius: 3px;
  outline: none;
  cursor: pointer;
}

.opacity-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  background: #4f46e5;
  border-radius: 50%;
  cursor: pointer;
}

.opacity-slider::-moz-range-thumb {
  width: 18px;
  height: 18px;
  background: #4f46e5;
  border-radius: 50%;
  cursor: pointer;
  border: none;
}

.width-input {
  max-width: 120px;
}

.templates-checklist {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  padding: 0.5rem;
}

.template-checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.5rem;
  border-radius: 0.25rem;
  cursor: pointer;
  transition: background 0.15s ease;
  font-weight: 400;
  font-size: 0.875rem;
  color: #4a5568;
}

.template-checkbox:hover {
  background: #f7fafc;
}

.template-checkbox input[type="checkbox"] {
  accent-color: #4f46e5;
}

.no-templates {
  padding: 1rem;
  text-align: center;
  color: #a0aec0;
  font-size: 0.875rem;
}

.draw-map-container {
  width: 100%;
  height: 350px;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  z-index: 0;
}

.geojson-input {
  font-family: 'Courier New', monospace;
  font-size: 0.8rem;
}

.geojson-btn {
  margin-top: 0.5rem;
  padding: 0.5rem 1rem;
  background: #eef2ff;
  color: #4f46e5;
  border: 1px solid #c7d2fe;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.geojson-btn:hover {
  background: #e0e7ff;
}

.geojson-error {
  margin-top: 0.5rem;
  color: #dc2626;
  font-size: 0.85rem;
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

  .search-input {
    width: 100%;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
