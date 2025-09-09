package models

import (
        "encoding/json"
        "fmt"
        "time"
)

// HistoricalEvent represents a historical event with geographical and temporal data
type HistoricalEvent struct {
        ID          int       `json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        Latitude    float64   `json:"latitude"`
        Longitude   float64   `json:"longitude"`
        EventDate   time.Time `json:"-"` // Don't auto-marshal this field
        Era         string    `json:"era"`
        LensType    string    `json:"lens_type"`
        DisplayDate string    `json:"display_date,omitempty"`
        DatasetID   *int      `json:"dataset_id,omitempty"`
        CreatedBy   int       `json:"created_by"` // User ID who created this event
        Tags        []Tag     `json:"tags,omitempty"`
}

// MarshalJSON custom JSON marshaler to handle BC dates properly
func (e HistoricalEvent) MarshalJSON() ([]byte, error) {
        // Create alias to avoid infinite recursion
        type Alias HistoricalEvent
        
        // Convert EventDate to string format that can handle BC dates
        var eventDateStr string
        if e.Era == "BC" {
                // For BC dates, use positive year with format YYYY-MM-DD
                year := -(e.EventDate.Year() + 1)
                eventDateStr = formatBCDate(year, e.EventDate.Month(), e.EventDate.Day())
        } else {
                // For AD dates, use standard formatting
                eventDateStr = e.EventDate.Format("2006-01-02T15:04:05Z07:00")
        }
        
        return json.Marshal(&struct {
                EventDate string `json:"event_date"`
                *Alias
        }{
                EventDate: eventDateStr,
                Alias:     (*Alias)(&e),
        })
}

// formatBCDate formats BC dates in a safe string format
func formatBCDate(year int, month time.Month, day int) string {
        return fmt.Sprintf("%04d-%02d-%02dT00:00:00Z", year, int(month), day)
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
        DatasetID   *int      `json:"dataset_id,omitempty"`
        TagIDs      []int     `json:"tag_ids,omitempty"`
}

// ToHistoricalEvent converts CreateEventRequest to HistoricalEvent
func (req *CreateEventRequest) ToHistoricalEvent(createdBy int) *HistoricalEvent {
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
                DatasetID:   req.DatasetID,
                CreatedBy:   createdBy,
        }
}