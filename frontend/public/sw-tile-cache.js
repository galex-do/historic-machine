// Service Worker for caching OpenStreetMap tiles
const CACHE_NAME = 'osm-tiles-v1'
const TIMESTAMP_DB_NAME = 'tile-timestamps'
const TIMESTAMP_STORE_NAME = 'timestamps'
const TILE_CACHE_DURATION = 7 * 24 * 60 * 60 * 1000 // 7 days in milliseconds

// Helper to send messages to all clients
function notifyClients(message) {
  self.clients.matchAll().then(clients => {
    clients.forEach(client => {
      client.postMessage({type: 'SW_LOG', message})
    })
  })
}

// Open IndexedDB for persistent timestamp storage
function openTimestampDB() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(TIMESTAMP_DB_NAME, 1)
    
    request.onerror = () => reject(request.error)
    request.onsuccess = () => resolve(request.result)
    
    request.onupgradeneeded = (event) => {
      const db = event.target.result
      if (!db.objectStoreNames.contains(TIMESTAMP_STORE_NAME)) {
        db.createObjectStore(TIMESTAMP_STORE_NAME)
      }
    }
  })
}

// Get timestamp for a tile
async function getTimestamp(url) {
  const db = await openTimestampDB()
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([TIMESTAMP_STORE_NAME], 'readonly')
    const store = transaction.objectStore(TIMESTAMP_STORE_NAME)
    const request = store.get(url)
    
    request.onsuccess = () => resolve(request.result)
    request.onerror = () => reject(request.error)
  })
}

// Set timestamp for a tile
async function setTimestamp(url, timestamp) {
  const db = await openTimestampDB()
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([TIMESTAMP_STORE_NAME], 'readwrite')
    const store = transaction.objectStore(TIMESTAMP_STORE_NAME)
    const request = store.put(timestamp, url)
    
    request.onsuccess = () => resolve()
    request.onerror = () => reject(request.error)
  })
}

self.addEventListener('install', (event) => {
  notifyClients('Installing Service Worker...')
  self.skipWaiting() // Activate immediately
})

self.addEventListener('activate', (event) => {
  notifyClients('Service Worker activated!')
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME && cacheName.startsWith('osm-tiles-')) {
            notifyClients(`Deleting old cache: ${cacheName}`)
            return caches.delete(cacheName)
          }
        })
      )
    }).then(() => {
      notifyClients('Taking control of all pages')
      return self.clients.claim()
    })
  )
})

self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)
  
  // Match OSM tile pattern: a.tile.openstreetmap.org, b.tile.openstreetmap.org, etc.
  const isTile = (url.hostname.endsWith('.tile.openstreetmap.org') || url.hostname === 'tile.openstreetmap.org') 
                  && url.pathname.match(/\/\d+\/\d+\/\d+\.png$/)
  
  if (isTile) {
    event.respondWith(
      (async () => {
        const cache = await caches.open(CACHE_NAME)
        const cachedResponse = await cache.match(event.request)
        const cacheKey = event.request.url
        
        // Check if we have a cached tile and if it's still fresh
        if (cachedResponse) {
          const cachedTime = await getTimestamp(cacheKey)
          const now = Date.now()
          
          if (cachedTime && (now - cachedTime < TILE_CACHE_DURATION)) {
            notifyClients(`✓ Cache HIT: ${url.pathname}`)
            return cachedResponse
          } else {
            notifyClients(`⏰ Expired: ${url.pathname}`)
          }
        }
        
        // Fetch from network
        notifyClients(`↓ Fetching: ${url.pathname}`)
        try {
          const networkResponse = await fetch(event.request)
          
          // Only cache successful responses
          if (networkResponse && networkResponse.status === 200) {
            // Store the cache timestamp
            await setTimestamp(cacheKey, Date.now())
            
            // Clone and cache the response
            await cache.put(event.request, networkResponse.clone())
            notifyClients(`✓ Cached: ${url.pathname}`)
          }
          
          return networkResponse
        } catch (error) {
          notifyClients(`✗ Network error: ${url.pathname}`)
          
          // Return stale cache if available (offline fallback)
          if (cachedResponse) {
            notifyClients(`⚠️ Serving stale: ${url.pathname}`)
            return cachedResponse
          }
          
          throw error
        }
      })()
    )
  }
})

notifyClients('Service Worker script loaded')
