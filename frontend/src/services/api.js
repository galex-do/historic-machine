/**
 * API service for backend communication
 */

import authService from './authService.js'

class ApiService {
  constructor() {
    this.baseURL = this.getBackendUrl()
  }

  getBackendUrl() {
    // Check if running in Docker environment (nginx proxy available)
    if (window.location.host.includes('localhost:3000')) {
      return '/api'
    }
    // Default for Replit development environment
    return 'http://localhost:8080/api'
  }

  async makeRequest(endpoint, options = {}) {
    try {
      const url = `${this.baseURL}${endpoint}`
      console.log(`Making API request to: ${url}`)
      
      // Get authentication headers
      const authHeaders = authService.getHeaders()
      
      const response = await fetch(url, {
        headers: {
          ...authHeaders,
          ...options.headers,
        },
        ...options,
      })

      if (!response.ok) {
        console.error(`HTTP error! status: ${response.status} for ${url}`)
        
        // Try to extract error message from response body
        let errorMessage = `HTTP error! status: ${response.status}`
        try {
          const errorData = await response.text()
          if (errorData) {
            // If response is plain text (like our validation messages), use it directly
            errorMessage = errorData
          }
        } catch (e) {
          // If we can't parse the error response, use the default message
        }
        
        throw new Error(errorMessage)
      }

      const responseData = await response.json()
      console.log(`API response for ${endpoint}:`, responseData)
      
      // Handle both new standardized format and legacy format
      return responseData.data || responseData
    } catch (error) {
      console.error(`API Error for ${endpoint}:`, error.message || error)
      throw error
    }
  }

  // Events API
  async getEvents(page = null, limit = null, sortField = null, sortDirection = null) {
    let endpoint = '/events'
    if (page !== null && limit !== null) {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
      })
      if (sortField) {
        params.append('sort', sortField)
      }
      if (sortDirection) {
        params.append('order', sortDirection)
      }
      endpoint = `/events?${params}`
    }
    return this.makeRequest(endpoint)
  }

  async createEvent(eventData) {
    return this.makeRequest('/events', {
      method: 'POST',
      body: JSON.stringify(eventData),
    })
  }

  async updateEvent(eventId, eventData) {
    return this.makeRequest(`/events/${eventId}`, {
      method: 'PUT',
      body: JSON.stringify(eventData),
    })
  }

  async deleteEvent(eventId) {
    return this.makeRequest(`/events/${eventId}`, {
      method: 'DELETE',
    })
  }

  async getEventById(id) {
    return this.makeRequest(`/events/${id}`)
  }

  async getEventsInBBox(minLat, minLng, maxLat, maxLng) {
    const params = new URLSearchParams({
      min_lat: minLat,
      min_lng: minLng,
      max_lat: maxLat,
      max_lng: maxLng,
    })
    return this.makeRequest(`/events/bbox?${params}`)
  }

  // Template API
  async getTemplateGroups() {
    return this.makeRequest('/date-template-groups')
  }

  async getTemplatesByGroup(groupId) {
    return this.makeRequest(`/date-templates/${groupId}`)
  }

  async getAllTemplates() {
    return this.makeRequest('/date-templates')
  }

  // Tags API
  async getTags() {
    return this.makeRequest('/tags')
  }

  async createTag(tagData) {
    return this.makeRequest('/tags', {
      method: 'POST',
      body: JSON.stringify(tagData),
    })
  }

  async updateTag(id, tagData) {
    return this.makeRequest(`/tags/${id}`, {
      method: 'PUT',
      body: JSON.stringify(tagData),
    })
  }

  async deleteTag(id) {
    return this.makeRequest(`/tags/${id}`, {
      method: 'DELETE',
    })
  }

  async setEventTags(eventId, tagIds) {
    return this.makeRequest(`/events/${eventId}/tags`, {
      method: 'PUT',
      body: JSON.stringify({ tag_ids: tagIds }),
    })
  }

  async addTagToEvent(eventId, tagId) {
    return this.makeRequest(`/events/${eventId}/tags/${tagId}`, {
      method: 'POST',
    })
  }

  async removeTagFromEvent(eventId, tagId) {
    return this.makeRequest(`/events/${eventId}/tags/${tagId}`, {
      method: 'DELETE',
    })
  }

  async importEvents(events, filename = '') {
    return this.makeRequest('/events/import', {
      method: 'POST',
      body: JSON.stringify({ events, filename }),
    })
  }

  // Dataset operations
  async getDatasets() {
    return this.makeRequest('/datasets')
  }

  async deleteDataset(id) {
    return this.makeRequest(`/datasets/${id}`, {
      method: 'DELETE',
    })
  }

  // User management API
  async getUsers() {
    return this.makeRequest('/users')
  }

  async createUser(userData) {
    return this.makeRequest('/users', {
      method: 'POST',
      body: JSON.stringify(userData),
    })
  }

  async updateUser(userId, userData) {
    return this.makeRequest(`/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    })
  }

  async deleteUser(userId) {
    return this.makeRequest(`/users/${userId}`, {
      method: 'DELETE',
    })
  }
}

// Export singleton instance
export default new ApiService()