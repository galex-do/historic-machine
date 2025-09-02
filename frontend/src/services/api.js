/**
 * API service for backend communication
 */

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
      
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      })

      if (!response.ok) {
        console.error(`HTTP error! status: ${response.status} for ${url}`)
        throw new Error(`HTTP error! status: ${response.status}`)
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
  async getEvents() {
    return this.makeRequest('/events')
  }

  async createEvent(eventData) {
    return this.makeRequest('/events', {
      method: 'POST',
      body: JSON.stringify(eventData),
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
}

// Export singleton instance
export default new ApiService()