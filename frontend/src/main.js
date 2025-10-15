import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index.js'

// Register Service Worker for tile caching
if ('serviceWorker' in navigator) {
  // Listen for messages from Service Worker
  navigator.serviceWorker.addEventListener('message', (event) => {
    if (event.data.type === 'SW_LOG') {
      console.log('[SW]', event.data.message)
    }
  })
  
  window.addEventListener('load', () => {
    // Unregister old service workers and register fresh one
    navigator.serviceWorker.getRegistrations().then((registrations) => {
      registrations.forEach((registration) => {
        registration.unregister()
        console.log('🔄 Unregistered old Service Worker')
      })
    }).then(() => {
      // Register new Service Worker
      navigator.serviceWorker.register('/sw-tile-cache.js')
        .then((registration) => {
          console.log('✅ Service Worker registered for tile caching:', registration.scope)
          
          // Force update check
          registration.update()
          
          // Listen for updates
          registration.addEventListener('updatefound', () => {
            console.log('🔄 Service Worker update found, installing...')
          })
        })
        .catch((error) => {
          console.error('❌ Service Worker registration failed:', error)
        })
    })
  })
}

createApp(App).use(router).mount('#app')