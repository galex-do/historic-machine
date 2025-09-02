/**
 * API service for backend communication
 */

class ApiService {
  constructor() {
    this.baseURL = this.getBackendUrl()
  }

  getBackendUrl() {
    const host = window.location.hostname
    const protocol = window.location.protocol
    
    // For local development
    if (host === 'localhost' || host === '127.0.0.1') {
      return 'http://localhost:8080/api'
    }
    
    // For Docker environment (nginx proxy available)
    if (window.location.host.includes('localhost:3000')) {
      return '/api'
    }
    
    // For Replit environment - use internal network address
    if (host.includes('replit.dev')) {
      // In Replit, use the internal network address for backend
      return 'http://0.0.0.0:8080/api'
    }
    
    // Fallback for other environments
    return `${protocol}//${host}:8080/api`
  }

  async makeRequest(endpoint, options = {}) {
    try {
      const url = `${this.baseURL}${endpoint}`
      console.log(`Making API request to: ${url}`)
      
      const response = await fetch(url, {
        method: options.method || 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          ...options.headers,
        },
        body: options.body,
        mode: 'cors',
        credentials: 'omit',
        cache: 'no-cache'
      })

      console.log(`Response status: ${response.status}`)
      console.log(`Response headers:`, [...response.headers.entries()])

      if (!response.ok) {
        const errorText = await response.text()
        console.error(`HTTP error! status: ${response.status} for ${url}. Response: ${errorText}`)
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const responseData = await response.json()
      console.log(`API response for ${endpoint}:`, responseData)
      
      // Handle both new standardized format and legacy format
      return responseData.data || responseData
    } catch (error) {
      console.error(`API Error for ${endpoint}:`, error.message || error)
      console.error(`Full error:`, error)
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