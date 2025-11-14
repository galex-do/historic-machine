// Service worker cleanup file
// This file exists to unregister any previously registered service workers
// and prevent 404 errors in nginx logs

self.addEventListener('install', function(event) {
  // Skip waiting to activate immediately
  self.skipWaiting();
});

self.addEventListener('activate', function(event) {
  // Unregister this service worker
  event.waitUntil(
    self.registration.unregister().then(function() {
      // Clear all caches
      return caches.keys().then(function(cacheNames) {
        return Promise.all(
          cacheNames.map(function(cacheName) {
            return caches.delete(cacheName);
          })
        );
      });
    })
  );
});
