package main

import (
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

var events []HistoricalEvent

func main() {
        // Initialize with sample data
        events = []HistoricalEvent{
                {
                        ID:          1,
                        Name:        "Fall of Constantinople",
                        Description: "Ottoman Empire conquered Byzantine Empire",
                        Latitude:    41.0082,
                        Longitude:   28.9784,
                        EventDate:   time.Date(1453, 5, 29, 0, 0, 0, 0, time.UTC),
                        LensType:    "historic",
                },
        }

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
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events)
}

func create_event(w http.ResponseWriter, r *http.Request) {
        var new_event HistoricalEvent
        json.NewDecoder(r.Body).Decode(&new_event)
        new_event.ID = len(events) + 1
        events = append(events, new_event)
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(new_event)
}

func get_event_by_id(w http.ResponseWriter, r *http.Request) {
        // Simple implementation - in real app would parse ID and find by ID
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events[0])
}