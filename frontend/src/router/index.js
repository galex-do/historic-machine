import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth.js'

// Lazy-loaded components
const MapView = () => import('@/views/MapView.vue')
const AboutPage = () => import('@/views/AboutPage.vue')
const AdminEvents = () => import('@/views/AdminEvents.vue')
const AdminTags = () => import('@/views/AdminTags.vue')
const AdminTemplates = () => import('@/views/AdminTemplates.vue')
const AdminUsers = () => import('@/views/AdminUsers.vue')

const routes = [
  {
    path: '/',
    name: 'Map',
    component: MapView
  },
  {
    path: '/about',
    name: 'About',
    component: AboutPage
  },
  {
    path: '/admin/events',
    name: 'AdminEvents',
    component: AdminEvents,
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/tags',
    name: 'AdminTags',
    component: AdminTags,
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/templates',
    name: 'AdminTemplates',
    component: AdminTemplates,
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: AdminUsers,
    meta: { requiresAdmin: true, requiresSuper: true }
  },
  {
    path: '/admin/datasets',
    name: 'AdminDatasets',
    component: () => import('@/views/DatasetsPanel.vue'),
    meta: { requiresAdmin: true }
  },
  // Legacy routes for backward compatibility
  {
    path: '/events',
    redirect: '/admin/events'
  },
  {
    path: '/datasets',
    redirect: '/admin/datasets'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard for admin routes
router.beforeEach(async (to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAdmin)) {
    const { canAccessAdmin, isSuper, loading, authInitialized, initAuth } = useAuth()
    
    // If authentication hasn't been initialized yet, initialize it
    if (!authInitialized.value && !loading.value) {
      await initAuth()
    }
    
    // Wait for authentication initialization to complete
    while (loading.value) {
      await new Promise(resolve => setTimeout(resolve, 50))
    }
    
    // Check admin access first
    if (!canAccessAdmin.value) {
      // Redirect to map view if user doesn't have admin access
      next({ name: 'Map' })
      return
    }
    
    // Check super access for super-only routes
    if (to.matched.some(record => record.meta.requiresSuper)) {
      if (isSuper.value) {
        next()
      } else {
        // Redirect to admin events if user doesn't have super access
        next({ name: 'AdminEvents' })
      }
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router