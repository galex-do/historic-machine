package handlers

import (
        "encoding/json"
        "fmt"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/models"
        "historical-events-backend/pkg/response"
        "log"
        "math/rand"
        "net/http"
        "strconv"
        "strings"
        "time"

        "github.com/gorilla/mux"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
        eventRepo   *repositories.EventRepository
        tagRepo     *repositories.TagRepository
        datasetRepo *repositories.DatasetRepository
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventRepo *repositories.EventRepository, tagRepo *repositories.TagRepository, datasetRepo *repositories.DatasetRepository) *EventHandler {
        return &EventHandler{
                eventRepo:   eventRepo,
                tagRepo:     tagRepo,
                datasetRepo: datasetRepo,
        }
}

// GetAllEvents handles GET /api/events with optional pagination
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        
        // Check if pagination parameters are provided
        pageStr := query.Get("page")
        limitStr := query.Get("limit")
        sortField := query.Get("sort")
        sortDir := query.Get("order")
        
        // If pagination parameters exist, use paginated query
        if pageStr != "" || limitStr != "" {
                page, err := strconv.Atoi(pageStr)
                if err != nil || page < 1 {
                        page = 1
                }
                
                limit, err := strconv.Atoi(limitStr)
                if err != nil || limit < 1 {
                        limit = 10 // Default page size
                }
                
                // Validate limit (max 100 to prevent abuse)
                if limit > 100 {
                        limit = 100
                }
                
                // Set default sorting if not provided
                if sortField == "" {
                        sortField = "date"
                }
                if sortDir == "" {
                        sortDir = "asc"
                }
                
                events, total, err := h.eventRepo.GetPaginatedWithSort(page, limit, sortField, sortDir)
                if err != nil {
                        log.Printf("Error fetching paginated events: %v", err)
                        response.InternalError(w, "Failed to fetch events")
                        return
                }
                
                totalPages := (total + limit - 1) / limit // Ceiling division
                
                paginatedResponse := map[string]interface{}{
                        "events": events,
                        "pagination": map[string]interface{}{
                                "current_page": page,
                                "page_size":    limit,
                                "total_items":  total,
                                "total_pages":  totalPages,
                        },
                }
                
                response.Success(w, paginatedResponse)
                return
        }
        
        // Fall back to non-paginated query for backward compatibility
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
        
        // Get current user from context
        user := getUserFromContext(r.Context())
        if user == nil {
                response.BadRequest(w, "User not found in context")
                return
        }
        
        // Convert request to event model
        event := req.ToHistoricalEvent(user.ID)
        
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
        
        // Get current user from context
        user := getUserFromContext(r.Context())
        if user == nil {
                response.BadRequest(w, "User not found in context")
                return
        }
        
        // Convert request to event model and set ID
        event := req.ToHistoricalEvent(user.ID)
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
                Filename string `json:"filename,omitempty"`
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

        // Get user ID from request context (set by auth middleware)
        userID := 1 // Default to admin user for now - in real app get from context
        if userCtx := r.Context().Value("user_id"); userCtx != nil {
                if uid, ok := userCtx.(int); ok {
                        userID = uid
                }
        }

        // Create dataset record
        filename := req.Filename
        if filename == "" {
                filename = "imported_dataset.json"
        }
        
        dataset := &models.EventDataset{
                Filename:    filename,
                Description: fmt.Sprintf("Dataset imported with %d events", len(req.Events)),
                EventCount:  0, // Will be updated as events are created
                UploadedBy:  userID,
        }
        
        createdDataset, err := h.datasetRepo.Create(dataset)
        if err != nil {
                log.Printf("Failed to create dataset record: %v", err)
                response.InternalError(w, "Failed to create dataset record")
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

                // Store BC dates as positive dates (same as predefined events)
                // The era field indicates BC/AD, astronomical conversion happens in views
                eventDate := parsedDate

                // Create event with dataset reference
                event := &models.HistoricalEvent{
                        Name:        eventData.Name,
                        Description: eventData.Description,
                        Latitude:    eventData.Latitude,
                        Longitude:   eventData.Longitude,
                        EventDate:   eventDate,
                        Era:         eventData.Era,
                        LensType:    eventData.Type,
                        DisplayDate: formatDisplayDate(eventDate, eventData.Era),
                        DatasetID:   &createdDataset.ID,
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
                                        // Create new tag with random color
                                        tag := &models.Tag{
                                                Name:        tagName,
                                                Description: fmt.Sprintf("Auto-generated tag for %s", tagName),
                                                Color:       h.generateRandomColor(),
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

        // Update dataset with final event count
        err = h.datasetRepo.UpdateEventCount(createdDataset.ID, importedCount)
        if err != nil {
                log.Printf("Failed to update dataset event count: %v", err)
        }

        response.Success(w, map[string]interface{}{
                "success":        true,
                "imported_count": importedCount,
                "total_count":    len(req.Events),
                "dataset_id":     createdDataset.ID,
                "dataset_name":   createdDataset.Filename,
                "message":        fmt.Sprintf("Successfully imported %d out of %d events", importedCount, len(req.Events)),
        })
}

// formatDisplayDate formats a date for display with era
func formatDisplayDate(date time.Time, era string) string {
        if era == "BC" {
                // BC dates are stored as positive dates, use them directly
                return fmt.Sprintf("%02d.%02d.%d BC", date.Day(), date.Month(), date.Year())
        }
        return fmt.Sprintf("%02d.%02d.%d AD", date.Day(), date.Month(), date.Year())
}

// parseCoordinate parses a coordinate string to float64
func parseCoordinate(coord string) (float64, error) {
        return strconv.ParseFloat(coord, 64)
}

// generateRandomColor generates a random hex color for tags
func (h *EventHandler) generateRandomColor() string {
        // Predefined vibrant colors that work well for tags
        colors := []string{
                "#3B82F6", // Blue
                "#EF4444", // Red
                "#10B981", // Green
                "#F59E0B", // Orange
                "#8B5CF6", // Purple
                "#06B6D4", // Cyan
                "#EC4899", // Pink
                "#84CC16", // Lime
                "#F97316", // Orange
                "#6366F1", // Indigo
                "#14B8A6", // Teal
                "#F43F5E", // Rose
        }
        return colors[rand.Intn(len(colors))]
}