// Service Worker for caching OpenStreetMap tiles
const CACHE_NAME = 'osm-tiles-v1'
const TILE_CACHE_DURATION = 7 * 24 * 60 * 60 * 1000 // 7 days in milliseconds

// Cache tile requests
self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)
  
  // Only cache OSM tile requests
  if (url.hostname.includes('tile.openstreetmap.org') && url.pathname.match(/\/\d+\/\d+\/\d+\.png$/)) {
    event.respondWith(
      caches.open(CACHE_NAME).then((cache) =>
        cache.match(event.request).then((cachedResponse) => {
          // Check if we have a cached tile
          if (cachedResponse) {
            // Check if cached tile is still fresh (within 7 days)
            const cachedTime = new Date(cachedResponse.headers.get('sw-cached-time'))
            const now = new Date()
            
            if (now - cachedTime < TILE_CACHE_DURATION) {
              console.log('[SW] Serving cached tile:', url.pathname)
              return cachedResponse
            } else {
              console.log('[SW] Cached tile expired, fetching new:', url.pathname)
            }
          }
          
          // Fetch from network if not cached or expired
          return fetch(event.request)
            .then((networkResponse) => {
              // Clone the response before caching
              const responseToCache = networkResponse.clone()
              
              // Add custom header with cache time
              const headers = new Headers(responseToCache.headers)
              headers.append('sw-cached-time', new Date().toISOString())
              
              const modifiedResponse = new Response(responseToCache.body, {
                status: responseToCache.status,
                statusText: responseToCache.statusText,
                headers: headers
              })
              
              cache.put(event.request, modifiedResponse)
              console.log('[SW] Cached new tile:', url.pathname)
              
              return networkResponse
            })
            .catch((error) => {
              console.error('[SW] Fetch failed for tile:', url.pathname, error)
              
              // Try to return stale cache if network fails
              if (cachedResponse) {
                console.log('[SW] Network failed, serving stale cached tile:', url.pathname)
                return cachedResponse
              }
              
              throw error
            })
        })
      )
    )
  }
})

// Clean up old caches on activation
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME && cacheName.startsWith('osm-tiles-')) {
            console.log('[SW] Deleting old cache:', cacheName)
            return caches.delete(cacheName)
          }
        })
      )
    })
  )
})

self.addEventListener('install', () => {
  console.log('[SW] Service Worker installed for tile caching')
  self.skipWaiting()
})

console.log('[SW] Service Worker loaded')
