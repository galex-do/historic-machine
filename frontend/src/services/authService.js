/**
 * Authentication service for handling login, logout, and user management
 */

const API_BASE = 'http://localhost:8080/api'

class AuthService {
  constructor() {
    this.token = localStorage.getItem('auth_token')
    this.user = this.token ? JSON.parse(localStorage.getItem('auth_user') || 'null') : null
  }

  // Get authorization headers
  getHeaders() {
    const headers = {
      'Content-Type': 'application/json'
    }
    
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }
    
    return headers
  }

  // Login user
  async login(username, password) {
    try {
      console.log('Attempting login with:', { username, password: '***' })
      console.log('Making request to:', `${API_BASE}/auth/login`)
      
      const requestBody = JSON.stringify({ username, password })
      console.log('Request body:', requestBody)
      
      const response = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: requestBody
      })

      console.log('Response status:', response.status)
      console.log('Response headers:', response.headers)

      if (!response.ok) {
        const errorText = await response.text()
        console.error('Login failed with status:', response.status, 'Error:', errorText)
        throw new Error(errorText || 'Login failed')
      }

      const data = await response.json()
      console.log('Login successful, response data:', data)
      
      // Store token and user data
      this.token = data.token
      this.user = data.user
      localStorage.setItem('auth_token', this.token)
      localStorage.setItem('auth_user', JSON.stringify(this.user))
      
      return data
    } catch (error) {
      console.error('Login error details:', error.message, error)
      throw error
    }
  }

  // Register user
  async register(username, password, email = '') {
    try {
      const response = await fetch(`${API_BASE}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password, email })
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(errorText || 'Registration failed')
      }

      const data = await response.json()
      return data
    } catch (error) {
      console.error('Registration error:', error)
      throw error
    }
  }

  // Logout user
  async logout() {
    try {
      if (this.token) {
        await fetch(`${API_BASE}/auth/logout`, {
          method: 'POST',
          headers: this.getHeaders()
        })
      }
    } catch (error) {
      console.error('Logout error:', error)
      // Continue with local logout even if API call fails
    } finally {
      // Always clear local storage
      this.token = null
      this.user = null
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
    }
  }

  // Get current user info
  async getCurrentUser() {
    if (!this.token) {
      return null
    }

    try {
      const response = await fetch(`${API_BASE}/auth/me`, {
        method: 'GET',
        headers: this.getHeaders()
      })

      if (!response.ok) {
        // Token might be expired or invalid
        this.logout()
        return null
      }

      const user = await response.json()
      this.user = user
      localStorage.setItem('auth_user', JSON.stringify(user))
      return user
    } catch (error) {
      console.error('Get current user error:', error)
      this.logout()
      return null
    }
  }

  // Change password
  async changePassword(currentPassword, newPassword) {
    try {
      const response = await fetch(`${API_BASE}/auth/change-password`, {
        method: 'POST',
        headers: this.getHeaders(),
        body: JSON.stringify({
          current_password: currentPassword,
          new_password: newPassword
        })
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(errorText || 'Password change failed')
      }

      return true
    } catch (error) {
      console.error('Change password error:', error)
      throw error
    }
  }

  // Check if user is authenticated
  isAuthenticated() {
    return !!this.token && !!this.user
  }

  // Check if user is guest
  isGuest() {
    return !this.isAuthenticated()
  }

  // Check if user can create events
  canCreateEvents() {
    return this.isAuthenticated() && this.user && this.user.access_level !== 'guest'
  }

  // Check if user has admin access
  isAdmin() {
    return this.isAuthenticated() && this.user && (this.user.access_level === 'admin' || this.user.access_level === 'super')
  }

  // Check if user has super access
  isSuper() {
    return this.isAuthenticated() && this.user && this.user.access_level === 'super'
  }

  // Get user data
  getUser() {
    return this.user
  }

  // Get token
  getToken() {
    return this.token
  }
}

// Export singleton instance
export default new AuthService()