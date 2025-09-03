import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth.js'

// Lazy-loaded components
const MapView = () => import('@/views/MapView.vue')
const AdminPanel = () => import('@/views/AdminPanel.vue')

const routes = [
  {
    path: '/',
    name: 'Map',
    component: MapView
  },
  {
    path: '/events',
    name: 'Events',
    component: AdminPanel,
    meta: { requiresAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard for admin routes
router.beforeEach(async (to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAdmin)) {
    const { canAccessAdmin, loading, authInitialized, initAuth } = useAuth()
    
    // If authentication hasn't been initialized yet, initialize it
    if (!authInitialized.value && !loading.value) {
      await initAuth()
    }
    
    // Wait for authentication initialization to complete
    while (loading.value) {
      await new Promise(resolve => setTimeout(resolve, 50))
    }
    
    if (canAccessAdmin.value) {
      next()
    } else {
      // Redirect to map view if user doesn't have admin access
      next({ name: 'Map' })
    }
  } else {
    next()
  }
})

export default router