import { ref, computed, onMounted } from 'vue'
import authService from '@/services/authService.js'

const user = ref(null)
const isAuthenticated = ref(false)
const loading = ref(false)
const error = ref(null)

export function useAuth() {
  
  // Initialize authentication state
  const initAuth = async () => {
    loading.value = true
    error.value = null
    
    try {
      if (authService.isAuthenticated()) {
        const currentUser = await authService.getCurrentUser()
        if (currentUser) {
          user.value = currentUser
          isAuthenticated.value = true
        } else {
          user.value = null
          isAuthenticated.value = false
        }
      } else {
        user.value = null
        isAuthenticated.value = false
      }
    } catch (err) {
      console.error('Auth initialization error:', err)
      error.value = err.message || 'Failed to initialize authentication'
      user.value = null
      isAuthenticated.value = false
    } finally {
      loading.value = false
    }
  }

  // Login function
  const login = async (username, password) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await authService.login(username, password)
      user.value = response.user
      isAuthenticated.value = true
      console.log('Login successful:', response.user.username)
      return response
    } catch (err) {
      console.error('Login error:', err)
      error.value = err.message || 'Login failed'
      user.value = null
      isAuthenticated.value = false
      throw err
    } finally {
      loading.value = false
    }
  }

  // Register function
  const register = async (username, password, email = '') => {
    loading.value = true
    error.value = null
    
    try {
      const response = await authService.register(username, password, email)
      console.log('Registration successful:', response.username)
      return response
    } catch (err) {
      console.error('Registration error:', err)
      error.value = err.message || 'Registration failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Logout function
  const logout = async () => {
    loading.value = true
    error.value = null
    
    try {
      await authService.logout()
      user.value = null
      isAuthenticated.value = false
      console.log('Logout successful')
    } catch (err) {
      console.error('Logout error:', err)
      error.value = err.message || 'Logout failed'
      // Still clear local state even if API call fails
      user.value = null
      isAuthenticated.value = false
    } finally {
      loading.value = false
    }
  }

  // Change password function
  const changePassword = async (currentPassword, newPassword) => {
    loading.value = true
    error.value = null
    
    try {
      await authService.changePassword(currentPassword, newPassword)
      console.log('Password changed successfully')
      return true
    } catch (err) {
      console.error('Change password error:', err)
      error.value = err.message || 'Password change failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Clear error
  const clearError = () => {
    error.value = null
  }

  // Computed properties for permissions
  const isGuest = computed(() => !isAuthenticated.value)
  const canCreateEvents = computed(() => authService.canCreateEvents())
  const isAdmin = computed(() => authService.isAdmin())
  const isSuper = computed(() => authService.isSuper())

  // Initialize on mount
  onMounted(() => {
    initAuth()
  })

  return {
    // State
    user: computed(() => user.value),
    isAuthenticated: computed(() => isAuthenticated.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    
    // Computed permissions
    isGuest,
    canCreateEvents,
    isAdmin,
    isSuper,
    
    // Methods
    login,
    register,
    logout,
    changePassword,
    clearError,
    initAuth
  }
}