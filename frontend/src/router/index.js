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
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAdmin)) {
    const { canAccessAdmin } = useAuth()
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