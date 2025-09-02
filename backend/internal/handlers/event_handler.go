package handlers

import (
        "encoding/json"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/models"
        "historical-events-backend/pkg/response"
        "log"
        "net/http"
        "strconv"
        "strings"

        "github.com/gorilla/mux"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
        eventRepo *repositories.EventRepository
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventRepo *repositories.EventRepository) *EventHandler {
        return &EventHandler{
                eventRepo: eventRepo,
        }
}

// GetAllEvents handles GET /api/events
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
        events, err := h.eventRepo.GetAll()
        if err != nil {
                log.Printf("Error fetching events: %v", err)
                response.InternalError(w, "Failed to fetch events")
                return
        }
        
        response.Success(w, events)
}

// GetEventByID handles GET /api/events/{id}
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid event ID")
                return
        }
        
        event, err := h.eventRepo.GetByID(id)
        if err != nil {
                log.Printf("Error fetching event by ID %d: %v", id, err)
                response.NotFound(w, "Event not found")
                return
        }
        
        response.Success(w, event)
}

// CreateEvent handles POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
        var req models.CreateEventRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid JSON format")
                return
        }
        
        // Validate coordinates
        if err := h.eventRepo.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
                response.BadRequest(w, err.Error())
                return
        }
        
        // Convert request to event model
        event := req.ToHistoricalEvent()
        
        // Create event
        createdEvent, err := h.eventRepo.Create(event)
        if err != nil {
                log.Printf("Error creating event: %v", err)
                response.InternalError(w, "Failed to create event")
                return
        }
        
        response.Created(w, createdEvent, "Event created successfully")
}

// UpdateEvent handles PUT /api/events/{id}
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid event ID")
                return
        }
        
        var req models.CreateEventRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid JSON format")
                return
        }
        
        // Validate coordinates
        if err := h.eventRepo.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
                response.BadRequest(w, err.Error())
                return
        }
        
        // Convert request to event model and set ID
        event := req.ToHistoricalEvent()
        event.ID = id
        
        // Update event
        updatedEvent, err := h.eventRepo.Update(event)
        if err != nil {
                log.Printf("Error updating event: %v", err)
                if strings.Contains(err.Error(), "not found") {
                        response.NotFound(w, "Event not found")
                        return
                }
                response.InternalError(w, "Failed to update event")
                return
        }
        
        response.Success(w, updatedEvent)
}

// DeleteEvent handles DELETE /api/events/{id}
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid event ID")
                return
        }
        
        err = h.eventRepo.Delete(id)
        if err != nil {
                log.Printf("Error deleting event: %v", err)
                if strings.Contains(err.Error(), "not found") {
                        response.NotFound(w, "Event not found")
                        return
                }
                response.InternalError(w, "Failed to delete event")
                return
        }
        
        response.Success(w, map[string]string{"message": "Event deleted successfully"})
}

// GetEventsInBBox handles GET /api/events/bbox
func (h *EventHandler) GetEventsInBBox(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        
        minLat, err := parseCoordinate(query.Get("min_lat"))
        if err != nil {
                response.BadRequest(w, "Invalid min_lat parameter")
                return
        }
        
        minLng, err := parseCoordinate(query.Get("min_lng"))
        if err != nil {
                response.BadRequest(w, "Invalid min_lng parameter")
                return
        }
        
        maxLat, err := parseCoordinate(query.Get("max_lat"))
        if err != nil {
                response.BadRequest(w, "Invalid max_lat parameter")
                return
        }
        
        maxLng, err := parseCoordinate(query.Get("max_lng"))
        if err != nil {
                response.BadRequest(w, "Invalid max_lng parameter")
                return
        }
        
        events, err := h.eventRepo.GetInBoundingBox(minLat, minLng, maxLat, maxLng)
        if err != nil {
                log.Printf("Bounding box query error: %v", err)
                response.InternalError(w, "Failed to fetch events in bounding box")
                return
        }
        
        response.Success(w, events)
}

// GetEventsInRadius handles GET /api/events/radius
func (h *EventHandler) GetEventsInRadius(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        
        centerLat, err := parseCoordinate(query.Get("lat"))
        if err != nil {
                response.BadRequest(w, "Invalid lat parameter")
                return
        }
        
        centerLng, err := parseCoordinate(query.Get("lng"))
        if err != nil {
                response.BadRequest(w, "Invalid lng parameter")
                return
        }
        
        radius, err := parseCoordinate(query.Get("radius"))
        if err != nil || radius <= 0 {
                response.BadRequest(w, "Invalid radius parameter (must be > 0)")
                return
        }
        
        // Convert radius to rough bounding box for demo
        // TODO: Implement proper radius query with PostGIS ST_DWithin
        degreeRadius := radius / 111000 // Rough conversion from meters to degrees
        events, err := h.eventRepo.GetInBoundingBox(
                centerLat-degreeRadius, centerLng-degreeRadius,
                centerLat+degreeRadius, centerLng+degreeRadius)
        
        if err != nil {
                log.Printf("Radius query error: %v", err)
                response.InternalError(w, "Failed to fetch events within radius")
                return
        }
        
        response.Success(w, events)
}

// parseCoordinate parses a coordinate string to float64
func parseCoordinate(coord string) (float64, error) {
        return strconv.ParseFloat(coord, 64)
}