package handlers

import (
        "encoding/json"
        "fmt"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/models"
        "historical-events-backend/pkg/response"
        "log"
        "net/http"
        "strconv"
        "strings"
        "time"

        "github.com/gorilla/mux"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
        eventRepo *repositories.EventRepository
        tagRepo   *repositories.TagRepository
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventRepo *repositories.EventRepository, tagRepo *repositories.TagRepository) *EventHandler {
        return &EventHandler{
                eventRepo: eventRepo,
                tagRepo:   tagRepo,
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

// ImportEvents handles bulk importing of events from dataset
func (h *EventHandler) ImportEvents(w http.ResponseWriter, r *http.Request) {
        type ImportRequest struct {
                Events []struct {
                        Name        string   `json:"name"`
                        Description string   `json:"description"`
                        Date        string   `json:"date"`
                        Era         string   `json:"era"`
                        Latitude    float64  `json:"latitude"`
                        Longitude   float64  `json:"longitude"`
                        Type        string   `json:"type"`
                        Tags        []string `json:"tags"`
                } `json:"events"`
        }

        var req ImportRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid JSON format")
                return
        }

        if len(req.Events) == 0 {
                response.BadRequest(w, "No events provided for import")
                return
        }

        importedCount := 0
        for _, eventData := range req.Events {
                // Parse date with DD.MM.YYYY format
                parsedDate, err := time.Parse("02.01.2006", eventData.Date)
                if err != nil {
                        log.Printf("Failed to parse date %s: %v", eventData.Date, err)
                        continue
                }

                // Handle BC dates - convert to negative year for proper storage
                var eventDate time.Time
                if eventData.Era == "BC" {
                        // Convert BC year to negative astronomical year
                        year := parsedDate.Year()
                        bcYear := -(year - 1) // 753 BC becomes -752
                        eventDate = time.Date(bcYear, parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
                } else {
                        eventDate = parsedDate
                }

                // Create event
                event := &models.HistoricalEvent{
                        Name:        eventData.Name,
                        Description: eventData.Description,
                        Latitude:    eventData.Latitude,
                        Longitude:   eventData.Longitude,
                        EventDate:   eventDate,
                        Era:         eventData.Era,
                        LensType:    eventData.Type,
                        DisplayDate: formatDisplayDate(eventDate, eventData.Era),
                }

                // Save event
                createdEvent, err := h.eventRepo.Create(event)
                if err != nil {
                        log.Printf("Failed to create event %s: %v", eventData.Name, err)
                        continue
                }

                // Handle tags - find or create and associate with event
                if len(eventData.Tags) > 0 {
                        var tagIDs []int
                        
                        // Get all existing tags to check for duplicates
                        existingTags, err := h.tagRepo.GetAllTags()
                        if err != nil {
                                log.Printf("Failed to get existing tags: %v", err)
                                existingTags = []models.Tag{} // Continue with empty list
                        }
                        
                        for _, tagName := range eventData.Tags {
                                // Check if tag already exists
                                var foundTag *models.Tag
                                for _, existing := range existingTags {
                                        if strings.EqualFold(existing.Name, tagName) {
                                                foundTag = &existing
                                                break
                                        }
                                }
                                
                                if foundTag != nil {
                                        // Use existing tag
                                        tagIDs = append(tagIDs, foundTag.ID)
                                } else {
                                        // Create new tag
                                        tag := &models.Tag{
                                                Name:        tagName,
                                                Description: fmt.Sprintf("Auto-generated tag for %s", tagName),
                                                Color:       "#3B82F6", // Default blue color
                                        }
                                        
                                        createdTag, err := h.tagRepo.CreateTag(tag)
                                        if err != nil {
                                                log.Printf("Failed to create tag %s: %v", tagName, err)
                                                continue
                                        }
                                        
                                        tagIDs = append(tagIDs, createdTag.ID)
                                        // Add to existing tags list for next iteration
                                        existingTags = append(existingTags, *createdTag)
                                }
                        }
                        
                        // Associate all tags with the event
                        if len(tagIDs) > 0 {
                                err = h.tagRepo.SetEventTags(createdEvent.ID, tagIDs)
                                if err != nil {
                                        log.Printf("Failed to associate tags with event %s: %v", eventData.Name, err)
                                }
                        }
                }

                importedCount++
        }

        response.Success(w, map[string]interface{}{
                "success":        true,
                "imported_count": importedCount,
                "total_count":    len(req.Events),
                "message":        fmt.Sprintf("Successfully imported %d out of %d events", importedCount, len(req.Events)),
        })
}

// formatDisplayDate formats a date for display with era
func formatDisplayDate(date time.Time, era string) string {
        if era == "BC" {
                // For BC dates, use the original positive year
                year := -(date.Year() + 1)
                return fmt.Sprintf("%02d.%02d.%d BC", date.Day(), date.Month(), year)
        }
        return fmt.Sprintf("%02d.%02d.%d AD", date.Day(), date.Month(), date.Year())
}

// parseCoordinate parses a coordinate string to float64
func parseCoordinate(coord string) (float64, error) {
        return strconv.ParseFloat(coord, 64)
}