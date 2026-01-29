const CACHE_NAME = 'historia-v1';
const STATIC_CACHE = 'historia-static-v1';
const TILE_CACHE = 'historia-tiles-v1';

const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/manifest.json'
];

const TILE_URL_PATTERN = /^https:\/\/[abc]\.tile\.openstreetmap\.org\//;
const MAX_TILE_CACHE_SIZE = 1000;

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(STATIC_CACHE)
      .then((cache) => {
        console.log('[SW] Caching static assets');
        return cache.addAll(STATIC_ASSETS);
      })
      .then(() => self.skipWaiting())
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames
            .filter((name) => {
              return name.startsWith('historia-') && 
                     name !== STATIC_CACHE && 
                     name !== TILE_CACHE;
            })
            .map((name) => {
              console.log('[SW] Deleting old cache:', name);
              return caches.delete(name);
            })
        );
      })
      .then(() => self.clients.claim())
  );
});

self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url);
  
  if (TILE_URL_PATTERN.test(event.request.url)) {
    event.respondWith(handleTileRequest(event.request));
    return;
  }
  
  if (event.request.method !== 'GET') {
    return;
  }
  
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(handleApiRequest(event.request));
    return;
  }
  
  event.respondWith(handleStaticRequest(event.request));
});

async function handleStaticRequest(request) {
  try {
    const networkResponse = await fetch(request);
    
    if (networkResponse.ok) {
      const cache = await caches.open(STATIC_CACHE);
      cache.put(request, networkResponse.clone());
    }
    
    return networkResponse;
  } catch (error) {
    const cachedResponse = await caches.match(request);
    
    if (cachedResponse) {
      return cachedResponse;
    }
    
    if (request.destination === 'document') {
      return caches.match('/index.html');
    }
    
    throw error;
  }
}

async function handleApiRequest(request) {
  try {
    const networkResponse = await fetch(request);
    return networkResponse;
  } catch (error) {
    return new Response(
      JSON.stringify({ error: 'Offline', message: 'Network unavailable' }),
      { 
        status: 503,
        headers: { 'Content-Type': 'application/json' }
      }
    );
  }
}

async function handleTileRequest(request) {
  const cache = await caches.open(TILE_CACHE);
  const cachedResponse = await cache.match(request);
  
  if (cachedResponse) {
    fetchAndCacheTile(request, cache);
    return cachedResponse;
  }
  
  try {
    const networkResponse = await fetch(request);
    
    if (networkResponse.ok) {
      cache.put(request, networkResponse.clone());
      trimTileCache(cache);
    }
    
    return networkResponse;
  } catch (error) {
    return new Response('', { status: 408, statusText: 'Tile unavailable offline' });
  }
}

async function fetchAndCacheTile(request, cache) {
  try {
    const networkResponse = await fetch(request);
    if (networkResponse.ok) {
      cache.put(request, networkResponse.clone());
    }
  } catch (error) {
  }
}

async function trimTileCache(cache) {
  const keys = await cache.keys();
  if (keys.length > MAX_TILE_CACHE_SIZE) {
    const deleteCount = keys.length - MAX_TILE_CACHE_SIZE;
    for (let i = 0; i < deleteCount; i++) {
      await cache.delete(keys[i]);
    }
  }
}

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});
