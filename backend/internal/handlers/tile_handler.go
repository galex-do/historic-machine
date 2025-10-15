package handlers

import (
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "strconv"
)

// TileHandler handles map tile caching and proxying
type TileHandler struct {
        cacheDir string
}

// NewTileHandler creates a new tile handler with specified cache directory
func NewTileHandler(cacheDir string) *TileHandler {
        // Create cache directory if it doesn't exist
        if err := os.MkdirAll(cacheDir, 0755); err != nil {
                log.Printf("Warning: Failed to create tile cache directory: %v", err)
        }
        
        return &TileHandler{
                cacheDir: cacheDir,
        }
}

// GetTile handles tile requests with caching
func (h *TileHandler) GetTile(w http.ResponseWriter, r *http.Request) {
        // Extract tile coordinates from URL path
        // Format: /tiles/{z}/{x}/{y}.png
        vars := r.URL.Query()
        z := vars.Get("z")
        x := vars.Get("x")
        y := vars.Get("y")
        
        if z == "" || x == "" || y == "" {
                http.Error(w, "Missing tile coordinates", http.StatusBadRequest)
                return
        }
        
        // Validate coordinates
        if _, err := strconv.Atoi(z); err != nil {
                http.Error(w, "Invalid zoom level", http.StatusBadRequest)
                return
        }
        if _, err := strconv.Atoi(x); err != nil {
                http.Error(w, "Invalid x coordinate", http.StatusBadRequest)
                return
        }
        if _, err := strconv.Atoi(y); err != nil {
                http.Error(w, "Invalid y coordinate", http.StatusBadRequest)
                return
        }
        
        // Build cache file path
        cacheFilePath := filepath.Join(h.cacheDir, z, x, fmt.Sprintf("%s.png", y))
        
        // Check if tile exists in cache
        if _, err := os.Stat(cacheFilePath); err == nil {
                // Serve from cache
                http.ServeFile(w, r, cacheFilePath)
                return
        }
        
        // Tile not in cache, fetch from OpenStreetMap
        // Use round-robin subdomains (a, b, c)
        subdomains := []string{"a", "b", "c"}
        subdomain := subdomains[(hashCoordinates(z, x, y)) % len(subdomains)]
        
        osmURL := fmt.Sprintf("https://%s.tile.openstreetmap.org/%s/%s/%s.png", subdomain, z, x, y)
        
        // Create HTTP request with proper headers (required by OSM tile usage policy)
        req, err := http.NewRequest("GET", osmURL, nil)
        if err != nil {
                log.Printf("Error creating request: %v", err)
                http.Error(w, "Failed to create request", http.StatusInternalServerError)
                return
        }
        
        // Set User-Agent header (required by OSM) - use generic browser UA
        req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HistoriaExMachina/1.0; +https://github.com/historia-ex-machina)")
        if referer := r.Header.Get("Referer"); referer != "" {
                req.Header.Set("Referer", referer)
        }
        
        // Fetch tile from OSM
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                log.Printf("Error fetching tile from OSM: %v", err)
                http.Error(w, "Failed to fetch tile", http.StatusInternalServerError)
                return
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
                http.Error(w, fmt.Sprintf("OSM returned status %d", resp.StatusCode), resp.StatusCode)
                return
        }
        
        // Create cache directory structure
        cacheDir := filepath.Join(h.cacheDir, z, x)
        if err := os.MkdirAll(cacheDir, 0755); err != nil {
                log.Printf("Warning: Failed to create cache directory: %v", err)
        }
        
        // Create cache file
        cacheFile, err := os.Create(cacheFilePath)
        if err != nil {
                log.Printf("Warning: Failed to create cache file: %v", err)
                // Continue serving even if caching fails
                w.Header().Set("Content-Type", "image/png")
                io.Copy(w, resp.Body)
                return
        }
        defer cacheFile.Close()
        
        // Write to both cache and response
        w.Header().Set("Content-Type", "image/png")
        w.Header().Set("Cache-Control", "public, max-age=2592000") // 30 days
        
        // Use TeeReader to write to both cache and response simultaneously
        teeReader := io.TeeReader(resp.Body, cacheFile)
        io.Copy(w, teeReader)
        
        log.Printf("Cached new tile: %s/%s/%s.png", z, x, y)
}

// hashCoordinates creates a simple hash from coordinates for subdomain selection
func hashCoordinates(z, x, y string) int {
        zInt, _ := strconv.Atoi(z)
        xInt, _ := strconv.Atoi(x)
        yInt, _ := strconv.Atoi(y)
        return (zInt + xInt + yInt)
}
