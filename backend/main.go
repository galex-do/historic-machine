package main

import (
        "database/sql"
        "encoding/json"
        "log"
        "net/http"
        "time"

        "github.com/gorilla/mux"
)

type HistoricalEvent struct {
        ID          int       `json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        Latitude    float64   `json:"latitude"`
        Longitude   float64   `json:"longitude"`
        EventDate   time.Time `json:"event_date"`
        LensType    string    `json:"lens_type"`
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

        log.Println("Server starting on :8080")
        log.Fatal(http.ListenAndServe(":8080", router))
}

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
        
        created_event, err := create_event_in_db(db, new_event)
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