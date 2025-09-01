package main

import (
        "database/sql"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "os"
        "strconv"
        "time"

        "github.com/gorilla/mux"
        "github.com/twpayne/go-geom"
        "github.com/twpayne/go-geom/encoding/ewkb"
        _ "github.com/lib/pq"
)

type HistoricalEvent struct {
        ID          int       `json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        Latitude    float64   `json:"latitude"`
        Longitude   float64   `json:"longitude"`
        EventDate   time.Time `json:"event_date"`
        Era         string    `json:"era"`
        LensType    string    `json:"lens_type"`
        DisplayDate string    `json:"display_date,omitempty"`
}

type DatabaseConfig struct {
        Host     string
        Port     string
        User     string
        Password string
        DBName   string
        SSLMode  string
}

var db *sql.DB

func main() {
        // Initialize database connection
        var err error
        db, err = connect_database()
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }
        defer db.Close()

        router := mux.NewRouter()
        
        // Enable CORS
        router.Use(func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        w.Header().Set("Access-Control-Allow-Origin", "*")
                        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
                        
                        if r.Method == "OPTIONS" {
                                return
                        }
                        
                        next.ServeHTTP(w, r)
                })
        })

        // API endpoints
        router.HandleFunc("/api/events", get_events).Methods("GET")
        router.HandleFunc("/api/events", create_event).Methods("POST")
        router.HandleFunc("/api/events/{id}", get_event_by_id).Methods("GET")
        
        // Spatial query endpoints
        router.HandleFunc("/api/events/bbox", get_events_in_bbox).Methods("GET")
        router.HandleFunc("/api/events/radius", get_events_in_radius).Methods("GET")

        log.Println("Server starting on :8080")
        log.Fatal(http.ListenAndServe(":8080", router))
}

// Database functions
func get_database_config() DatabaseConfig {
        // Check if DATABASE_URL is provided (for Replit/production)
        if database_url := getEnv("DATABASE_URL", ""); database_url != "" {
                // Use DATABASE_URL directly for connection
                return DatabaseConfig{}
        }
        
        return DatabaseConfig{
                Host:     getEnv("DB_HOST", "localhost"),
                Port:     getEnv("DB_PORT", "5432"),
                User:     getEnv("DB_USER", "postgres"),
                Password: getEnv("DB_PASSWORD", "password"),
                DBName:   getEnv("DB_NAME", "historical_events"),
                SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        }
}

func getEnv(key, fallback string) string {
        if value, exists := os.LookupEnv(key); exists {
                return value
        }
        return fallback
}

func connect_database() (*sql.DB, error) {
        // Check if DATABASE_URL is provided (for Replit/production)
        var connection_string string
        if database_url := getEnv("DATABASE_URL", ""); database_url != "" {
                connection_string = database_url
        } else {
                config := get_database_config()
                connection_string = fmt.Sprintf(
                        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
                        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
                )
        }

        // Try to connect with retries for Docker environment
        var database *sql.DB
        var err error
        
        for i := 0; i < 10; i++ {
                database, err = sql.Open("postgres", connection_string)
                if err != nil {
                        log.Printf("Failed to open database connection: %v", err)
                        time.Sleep(2 * time.Second)
                        continue
                }
                
                err = database.Ping()
                if err != nil {
                        log.Printf("Failed to ping database (attempt %d): %v", i+1, err)
                        time.Sleep(2 * time.Second)
                        continue
                }
                
                break
        }
        
        if err != nil {
                return nil, fmt.Errorf("failed to connect to database after retries: %v", err)
        }

        log.Println("Successfully connected to PostgreSQL database")
        return database, nil
}

func get_events_from_db(database *sql.DB) ([]HistoricalEvent, error) {
        // Use the view that handles BC/AD dates properly
        query := `
                SELECT id, name, description, latitude, longitude, event_date, era, lens_type, display_date
                FROM events_with_display_dates 
                ORDER BY astronomical_year ASC`
        
        rows, err := database.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        
        var events []HistoricalEvent
        
        for rows.Next() {
                var event HistoricalEvent
                
                err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Latitude, 
                                                &event.Longitude, &event.EventDate, &event.Era, &event.LensType, &event.DisplayDate)
                if err != nil {
                        return nil, err
                }
                
                
                events = append(events, event)
        }
        
        return events, nil
}

// Spatial query functions
func get_events_in_bounding_box(database *sql.DB, minLat, minLng, maxLat, maxLng float64) ([]HistoricalEvent, error) {
        query := `
                SELECT id, name, description, latitude, longitude, event_date, lens_type,
                       ST_X(location::geometry) as lng, ST_Y(location::geometry) as lat
                FROM events 
                WHERE location && ST_MakeEnvelope($1, $2, $3, $4, 4326)
                ORDER BY event_date DESC`
        
        rows, err := database.Query(query, minLng, minLat, maxLng, maxLat)
        if err != nil {
                return nil, fmt.Errorf("bounding box query failed: %v", err)
        }
        defer rows.Close()
        
        var events []HistoricalEvent
        for rows.Next() {
                var event HistoricalEvent
                var lng, lat float64
                
                err := rows.Scan(&event.ID, &event.Name, &event.Description, 
                                                &event.Latitude, &event.Longitude, &event.EventDate, &event.LensType,
                                                &lng, &lat)
                if err != nil {
                        log.Printf("Error scanning bounding box event: %v", err)
                        continue
                }
                
                // Update with PostGIS coordinates for accuracy
                event.Longitude = lng
                event.Latitude = lat
                events = append(events, event)
        }
        
        return events, nil
}

func create_event_with_location(database *sql.DB, event HistoricalEvent) (*HistoricalEvent, error) {
        query := `
                INSERT INTO events (name, description, latitude, longitude, event_date, era, lens_type) 
                VALUES ($1, $2, $3, $4, $5, $6, $7) 
                RETURNING id`
        
        var created_event HistoricalEvent = event
        
        err := database.QueryRow(query, event.Name, event.Description, event.Latitude, 
                                           event.Longitude, event.EventDate, event.Era, event.LensType).
                Scan(&created_event.ID)
        
        if err != nil {
                return nil, fmt.Errorf("failed to create event with location: %v", err)
        }
        
        return &created_event, nil
}

func validate_coordinates(lat, lng float64) error {
        if lat < -90 || lat > 90 {
                return fmt.Errorf("latitude must be between -90 and 90 degrees, got: %f", lat)
        }
        if lng < -180 || lng > 180 {
                return fmt.Errorf("longitude must be between -180 and 180 degrees, got: %f", lng)
        }
        return nil
}

func parse_coordinate_string(coord string) (float64, error) {
        value, err := strconv.ParseFloat(coord, 64)
        if err != nil {
                return 0, fmt.Errorf("invalid coordinate format: %s", coord)
        }
        return value, nil
}

// HTTP handlers
func get_events(w http.ResponseWriter, r *http.Request) {
        events, err := get_events_from_db(db)
        if err != nil {
                http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
                log.Printf("Error fetching events: %v", err)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events)
}

func create_event(w http.ResponseWriter, r *http.Request) {
        var new_event HistoricalEvent
        err := json.NewDecoder(r.Body).Decode(&new_event)
        if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        // Set default era if not provided
        if new_event.Era == "" {
                new_event.Era = "AD"
        }
        
        // Validate coordinates
        err = validate_coordinates(new_event.Latitude, new_event.Longitude)
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }
        
        // Create event with era support
        created_event, err := create_event_with_location(db, new_event)
        if err != nil {
                http.Error(w, "Failed to create event", http.StatusInternalServerError)
                log.Printf("Error creating event: %v", err)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(created_event)
}

func get_event_by_id(w http.ResponseWriter, r *http.Request) {
        events, err := get_events_from_db(db)
        if err != nil {
                http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
                return
        }
        
        if len(events) == 0 {
                http.Error(w, "Event not found", http.StatusNotFound)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events[0])
}

// Spatial query handlers
func get_events_in_bbox(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        
        minLat, err := parse_coordinate_string(query.Get("min_lat"))
        if err != nil {
                http.Error(w, "Invalid min_lat parameter", http.StatusBadRequest)
                return
        }
        
        minLng, err := parse_coordinate_string(query.Get("min_lng"))
        if err != nil {
                http.Error(w, "Invalid min_lng parameter", http.StatusBadRequest)
                return
        }
        
        maxLat, err := parse_coordinate_string(query.Get("max_lat"))
        if err != nil {
                http.Error(w, "Invalid max_lat parameter", http.StatusBadRequest)
                return
        }
        
        maxLng, err := parse_coordinate_string(query.Get("max_lng"))
        if err != nil {
                http.Error(w, "Invalid max_lng parameter", http.StatusBadRequest)
                return
        }
        
        events, err := get_events_in_bounding_box(db, minLat, minLng, maxLat, maxLng)
        if err != nil {
                http.Error(w, "Failed to fetch events in bounding box", http.StatusInternalServerError)
                log.Printf("Bounding box query error: %v", err)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events)
}

func get_events_in_radius(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        
        centerLat, err := parse_coordinate_string(query.Get("lat"))
        if err != nil {
                http.Error(w, "Invalid lat parameter", http.StatusBadRequest)
                return
        }
        
        centerLng, err := parse_coordinate_string(query.Get("lng"))
        if err != nil {
                http.Error(w, "Invalid lng parameter", http.StatusBadRequest)
                return
        }
        
        radius, err := parse_coordinate_string(query.Get("radius"))
        if err != nil || radius <= 0 {
                http.Error(w, "Invalid radius parameter (must be > 0)", http.StatusBadRequest)
                return
        }
        
        // Simple fallback to bounding box for demo (radius function needs more implementation)
        // Convert radius to rough bounding box 
        degreeRadius := radius / 111000 // Rough conversion from meters to degrees
        events, err := get_events_in_bounding_box(db, centerLat-degreeRadius, centerLng-degreeRadius, centerLat+degreeRadius, centerLng+degreeRadius)
        if err != nil {
                http.Error(w, "Failed to fetch events within radius", http.StatusInternalServerError)
                log.Printf("Radius query error: %v", err)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events)
}

// Helper functions for geometry handling (unused but available for future use)
func point_from_coordinates(lat, lng float64) *geom.Point {
        return geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{lng, lat})
}

func coordinates_from_ewkb(data []byte) (float64, float64, error) {
        geometry, err := ewkb.Unmarshal(data)
        if err != nil {
                return 0, 0, err
        }
        
        point, ok := geometry.(*geom.Point)
        if !ok {
                return 0, 0, fmt.Errorf("geometry is not a point")
        }
        
        coords := point.Coords()
        return coords[1], coords[0], nil // Return lat, lng
}