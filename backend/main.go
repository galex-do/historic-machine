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

type DateTemplateGroup struct {
        ID           int    `json:"id"`
        Name         string `json:"name"`
        Description  string `json:"description"`
        DisplayOrder int    `json:"display_order"`
}

type GroupWithTemplates struct {
        ID           int            `json:"id"`
        Name         string         `json:"name"`
        Description  string         `json:"description"`
        DisplayOrder int            `json:"display_order"`
        Templates    []DateTemplate `json:"templates"`
}

type DateTemplate struct {
        ID               int    `json:"id"`
        GroupID          int    `json:"group_id"`
        GroupName        string `json:"group_name"`
        Name             string `json:"name"`
        Description      string `json:"description"`
        StartDate        string `json:"start_date"`
        StartEra         string `json:"start_era"`
        EndDate          string `json:"end_date"`
        EndEra           string `json:"end_era"`
        DisplayOrder     int    `json:"display_order"`
        StartDisplayDate string `json:"start_display_date"`
        EndDisplayDate   string `json:"end_display_date"`
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
        router.HandleFunc("/api/events/{id}", update_event).Methods("PUT")
        router.HandleFunc("/api/events/{id}", delete_event).Methods("DELETE")
        
        // Spatial query endpoints
        router.HandleFunc("/api/events/bbox", get_events_in_bbox).Methods("GET")
        router.HandleFunc("/api/events/radius", get_events_in_radius).Methods("GET")
        
        // Date template group endpoints
        router.HandleFunc("/api/date-template-groups", get_date_template_groups).Methods("GET")
        router.HandleFunc("/api/date-template-groups", create_date_template_group).Methods("POST")
        router.HandleFunc("/api/date-template-groups/{id}", update_date_template_group).Methods("PUT")
        router.HandleFunc("/api/date-template-groups/{id}", delete_date_template_group).Methods("DELETE")
        
        // Date template endpoints
        router.HandleFunc("/api/date-templates", get_all_date_templates).Methods("GET")
        router.HandleFunc("/api/date-templates", create_date_template).Methods("POST")
        router.HandleFunc("/api/date-templates/{group_id}", get_date_templates_by_group).Methods("GET")
        router.HandleFunc("/api/date-templates/{id}", update_date_template).Methods("PUT")
        router.HandleFunc("/api/date-templates/{id}", delete_date_template).Methods("DELETE")

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

// Date template handlers
func get_date_template_groups(w http.ResponseWriter, r *http.Request) {
        // First get all groups
        group_query := `
                SELECT id, name, description, display_order
                FROM date_template_groups 
                ORDER BY display_order`
        
        group_rows, err := db.Query(group_query)
        if err != nil {
                http.Error(w, "Failed to fetch date template groups", http.StatusInternalServerError)
                log.Printf("Error fetching date template groups: %v", err)
                return
        }
        defer group_rows.Close()
        
        var groups []GroupWithTemplates
        for group_rows.Next() {
                var group GroupWithTemplates
                err := group_rows.Scan(&group.ID, &group.Name, &group.Description, &group.DisplayOrder)
                if err != nil {
                        log.Printf("Error scanning date template group: %v", err)
                        continue
                }
                
                // Get templates for this group
                template_query := `
                        SELECT id, group_id, name, start_date, start_era, end_date, end_era, display_order
                        FROM date_templates 
                        WHERE group_id = $1
                        ORDER BY display_order`
                
                log.Printf("Fetching templates for group %d (%s)", group.ID, group.Name)
                template_rows, err := db.Query(template_query, group.ID)
                if err != nil {
                        log.Printf("Error fetching templates for group %d: %v", group.ID, err)
                        groups = append(groups, group) // Add group even without templates
                        continue
                }
                
                var templates []DateTemplate
                for template_rows.Next() {
                        var template DateTemplate
                        err := template_rows.Scan(&template.ID, &template.GroupID, &template.Name,
                                &template.StartDate, &template.StartEra, &template.EndDate, 
                                &template.EndEra, &template.DisplayOrder)
                        if err != nil {
                                log.Printf("Error scanning template: %v", err)
                                continue
                        }
                        log.Printf("Loaded template: %s (ID: %d)", template.Name, template.ID)
                        templates = append(templates, template)
                }
                template_rows.Close()
                
                log.Printf("Group %s has %d templates", group.Name, len(templates))
                group.Templates = templates
                groups = append(groups, group)
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(groups)
}

func get_date_templates_by_group(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        group_id := vars["group_id"]
        
        query := `
                SELECT id, group_id, group_name, name, description,
                       start_date, start_era, end_date, end_era, display_order,
                       start_display_date, end_display_date
                FROM date_templates_with_display 
                WHERE group_id = $1
                ORDER BY display_order`
        
        rows, err := db.Query(query, group_id)
        if err != nil {
                http.Error(w, "Failed to fetch date templates", http.StatusInternalServerError)
                log.Printf("Error fetching date templates for group %s: %v", group_id, err)
                return
        }
        defer rows.Close()
        
        var templates []DateTemplate
        for rows.Next() {
                var template DateTemplate
                err := rows.Scan(&template.ID, &template.GroupID, &template.GroupName, 
                        &template.Name, &template.Description, &template.StartDate, &template.StartEra,
                        &template.EndDate, &template.EndEra, &template.DisplayOrder,
                        &template.StartDisplayDate, &template.EndDisplayDate)
                if err != nil {
                        log.Printf("Error scanning date template: %v", err)
                        continue
                }
                templates = append(templates, template)
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(templates)
}

func get_all_date_templates(w http.ResponseWriter, r *http.Request) {
        query := `
                SELECT id, group_id, group_name, name, description,
                       start_date, start_era, end_date, end_era, display_order,
                       start_display_date, end_display_date
                FROM date_templates_with_display 
                ORDER BY group_id, display_order`
        
        rows, err := db.Query(query)
        if err != nil {
                http.Error(w, "Failed to fetch all date templates", http.StatusInternalServerError)
                log.Printf("Error fetching all date templates: %v", err)
                return
        }
        defer rows.Close()
        
        var templates []DateTemplate
        for rows.Next() {
                var template DateTemplate
                err := rows.Scan(&template.ID, &template.GroupID, &template.GroupName, 
                        &template.Name, &template.Description, &template.StartDate, &template.StartEra,
                        &template.EndDate, &template.EndEra, &template.DisplayOrder,
                        &template.StartDisplayDate, &template.EndDisplayDate)
                if err != nil {
                        log.Printf("Error scanning date template: %v", err)
                        continue
                }
                templates = append(templates, template)
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(templates)
}

// Admin CRUD handlers for events
func update_event(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        event_id := vars["id"]
        
        var updated_event HistoricalEvent
        err := json.NewDecoder(r.Body).Decode(&updated_event)
        if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        // Validate coordinates
        err = validate_coordinates(updated_event.Latitude, updated_event.Longitude)
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }
        
        // Set default era if not provided
        if updated_event.Era == "" {
                updated_event.Era = "AD"
        }
        
        query := `
                UPDATE events 
                SET name = $1, description = $2, latitude = $3, longitude = $4, 
                    event_date = $5, era = $6, lens_type = $7
                WHERE id = $8`
        
        _, err = db.Exec(query, updated_event.Name, updated_event.Description, 
                         updated_event.Latitude, updated_event.Longitude, 
                         updated_event.EventDate, updated_event.Era, updated_event.LensType, event_id)
        if err != nil {
                http.Error(w, "Failed to update event", http.StatusInternalServerError)
                log.Printf("Error updating event: %v", err)
                return
        }
        
        // Return the updated event
        updated_event.ID, _ = strconv.Atoi(event_id)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updated_event)
}

func delete_event(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        event_id := vars["id"]
        
        query := "DELETE FROM events WHERE id = $1"
        result, err := db.Exec(query, event_id)
        if err != nil {
                http.Error(w, "Failed to delete event", http.StatusInternalServerError)
                log.Printf("Error deleting event: %v", err)
                return
        }
        
        rows_affected, err := result.RowsAffected()
        if err != nil {
                http.Error(w, "Error checking deletion", http.StatusInternalServerError)
                return
        }
        
        if rows_affected == 0 {
                http.Error(w, "Event not found", http.StatusNotFound)
                return
        }
        
        w.WriteHeader(http.StatusNoContent)
}

// Admin CRUD handlers for date template groups
func create_date_template_group(w http.ResponseWriter, r *http.Request) {
        var new_group DateTemplateGroup
        err := json.NewDecoder(r.Body).Decode(&new_group)
        if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        // Set default display order if not provided
        if new_group.DisplayOrder == 0 {
                new_group.DisplayOrder = 1
        }
        
        query := `
                INSERT INTO date_template_groups (name, description, display_order) 
                VALUES ($1, $2, $3) 
                RETURNING id`
        
        err = db.QueryRow(query, new_group.Name, new_group.Description, new_group.DisplayOrder).
                Scan(&new_group.ID)
        if err != nil {
                http.Error(w, "Failed to create date template group", http.StatusInternalServerError)
                log.Printf("Error creating date template group: %v", err)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(new_group)
}

func update_date_template_group(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        group_id := vars["id"]
        
        var updated_group DateTemplateGroup
        err := json.NewDecoder(r.Body).Decode(&updated_group)
        if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        query := `
                UPDATE date_template_groups 
                SET name = $1, description = $2, display_order = $3
                WHERE id = $4`
        
        _, err = db.Exec(query, updated_group.Name, updated_group.Description, 
                         updated_group.DisplayOrder, group_id)
        if err != nil {
                http.Error(w, "Failed to update date template group", http.StatusInternalServerError)
                log.Printf("Error updating date template group: %v", err)
                return
        }
        
        updated_group.ID, _ = strconv.Atoi(group_id)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updated_group)
}

func delete_date_template_group(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        group_id := vars["id"]
        
        // First check if any templates are using this group
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM date_templates WHERE group_id = $1", group_id).Scan(&count)
        if err != nil {
                http.Error(w, "Error checking group dependencies", http.StatusInternalServerError)
                return
        }
        
        if count > 0 {
                http.Error(w, "Cannot delete group with existing templates", http.StatusBadRequest)
                return
        }
        
        query := "DELETE FROM date_template_groups WHERE id = $1"
        result, err := db.Exec(query, group_id)
        if err != nil {
                http.Error(w, "Failed to delete date template group", http.StatusInternalServerError)
                log.Printf("Error deleting date template group: %v", err)
                return
        }
        
        rows_affected, err := result.RowsAffected()
        if err != nil {
                http.Error(w, "Error checking deletion", http.StatusInternalServerError)
                return
        }
        
        if rows_affected == 0 {
                http.Error(w, "Date template group not found", http.StatusNotFound)
                return
        }
        
        w.WriteHeader(http.StatusNoContent)
}

// Admin CRUD handlers for date templates
func create_date_template(w http.ResponseWriter, r *http.Request) {
        var new_template DateTemplate
        err := json.NewDecoder(r.Body).Decode(&new_template)
        if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        // Set default values
        if new_template.DisplayOrder == 0 {
                new_template.DisplayOrder = 1
        }
        if new_template.StartEra == "" {
                new_template.StartEra = "AD"
        }
        if new_template.EndEra == "" {
                new_template.EndEra = "AD"
        }
        
        query := `
                INSERT INTO date_templates (group_id, name, description, start_date, start_era, 
                                          end_date, end_era, display_order, start_display_date, end_display_date) 
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
                RETURNING id`
        
        err = db.QueryRow(query, new_template.GroupID, new_template.Name, new_template.Description,
                         new_template.StartDate, new_template.StartEra, new_template.EndDate, new_template.EndEra,
                         new_template.DisplayOrder, new_template.StartDisplayDate, new_template.EndDisplayDate).
                Scan(&new_template.ID)
        if err != nil {
                http.Error(w, "Failed to create date template", http.StatusInternalServerError)
                log.Printf("Error creating date template: %v", err)
                return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(new_template)
}

func update_date_template(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        template_id := vars["id"]
        
        var updated_template DateTemplate
        err := json.NewDecoder(r.Body).Decode(&updated_template)
        if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }
        
        query := `
                UPDATE date_templates 
                SET group_id = $1, name = $2, description = $3, start_date = $4, start_era = $5,
                    end_date = $6, end_era = $7, display_order = $8, start_display_date = $9, end_display_date = $10
                WHERE id = $11`
        
        _, err = db.Exec(query, updated_template.GroupID, updated_template.Name, updated_template.Description,
                         updated_template.StartDate, updated_template.StartEra, updated_template.EndDate, 
                         updated_template.EndEra, updated_template.DisplayOrder, 
                         updated_template.StartDisplayDate, updated_template.EndDisplayDate, template_id)
        if err != nil {
                http.Error(w, "Failed to update date template", http.StatusInternalServerError)
                log.Printf("Error updating date template: %v", err)
                return
        }
        
        updated_template.ID, _ = strconv.Atoi(template_id)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updated_template)
}

func delete_date_template(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        template_id := vars["id"]
        
        query := "DELETE FROM date_templates WHERE id = $1"
        result, err := db.Exec(query, template_id)
        if err != nil {
                http.Error(w, "Failed to delete date template", http.StatusInternalServerError)
                log.Printf("Error deleting date template: %v", err)
                return
        }
        
        rows_affected, err := result.RowsAffected()
        if err != nil {
                http.Error(w, "Error checking deletion", http.StatusInternalServerError)
                return
        }
        
        if rows_affected == 0 {
                http.Error(w, "Date template not found", http.StatusNotFound)
                return
        }
        
        w.WriteHeader(http.StatusNoContent)
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