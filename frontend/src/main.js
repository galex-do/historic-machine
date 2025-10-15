import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index.js'

// Register Service Worker for tile caching
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw-tile-cache.js')
      .then((registration) => {
        console.log('Service Worker registered for tile caching:', registration.scope)
      })
      .catch((error) => {
        console.error('Service Worker registration failed:', error)
      })
  })
}

createApp(App).use(router).mount('#app')