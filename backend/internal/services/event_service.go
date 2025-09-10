package services

import (
        "fmt"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/models"
)

// EventService provides business logic for events
type EventService struct {
        eventRepo *repositories.EventRepository
}

// NewEventService creates a new event service
func NewEventService(eventRepo *repositories.EventRepository) *EventService {
        return &EventService{
                eventRepo: eventRepo,
        }
}

// GetAllEvents retrieves all events
func (s *EventService) GetAllEvents() ([]models.HistoricalEvent, error) {
        return s.eventRepo.GetAll()
}

// GetEventByID retrieves a single event by ID
func (s *EventService) GetEventByID(id int) (*models.HistoricalEvent, error) {
        return s.eventRepo.GetByID(id)
}

// CreateEvent creates a new event with validation
func (s *EventService) CreateEvent(req *models.CreateEventRequest, userID int) (*models.HistoricalEvent, error) {
        // Validate coordinates
        if err := s.eventRepo.ValidateCoordinates(req.Latitude, req.Longitude); err != nil {
                return nil, err
        }
        
        // Convert to domain model
        event, err := req.ToHistoricalEvent(userID)
        if err != nil {
                return nil, fmt.Errorf("invalid event data: %w", err)
        }
        
        // Create event
        return s.eventRepo.Create(event)
}

// GetEventsInBoundingBox retrieves events within a geographical area
func (s *EventService) GetEventsInBoundingBox(minLat, minLng, maxLat, maxLng float64) ([]models.HistoricalEvent, error) {
        // Validate all coordinates
        coords := []float64{minLat, maxLat, minLng, maxLng}
        for i := 0; i < len(coords); i += 2 {
                if err := s.eventRepo.ValidateCoordinates(coords[i], coords[i+1]); err != nil {
                        return nil, err
                }
        }
        
        return s.eventRepo.GetInBoundingBox(minLat, minLng, maxLat, maxLng)
}