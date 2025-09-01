package models

import "time"

// HistoricalEvent represents a historical event with geographical and temporal data
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

// CreateEventRequest represents the request payload for creating an event
type CreateEventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Latitude    float64   `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   float64   `json:"longitude" validate:"required,min=-180,max=180"`
	EventDate   time.Time `json:"event_date" validate:"required"`
	Era         string    `json:"era"`
	LensType    string    `json:"lens_type" validate:"required"`
}

// ToHistoricalEvent converts CreateEventRequest to HistoricalEvent
func (req *CreateEventRequest) ToHistoricalEvent() *HistoricalEvent {
	era := req.Era
	if era == "" {
		era = "AD"
	}
	
	return &HistoricalEvent{
		Name:        req.Name,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		EventDate:   req.EventDate,
		Era:         era,
		LensType:    req.LensType,
	}
}