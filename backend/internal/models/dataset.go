package models

import "time"

// EventDataset represents an uploaded dataset of events
type EventDataset struct {
	ID          int       `json:"id"`
	Filename    string    `json:"filename"`
	Description string    `json:"description,omitempty"`
	EventCount  int       `json:"event_count"`
	UploadedBy  int       `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateDatasetRequest represents the request payload for creating a dataset
type CreateDatasetRequest struct {
	Filename    string `json:"filename" validate:"required"`
	Description string `json:"description,omitempty"`
	EventCount  int    `json:"event_count"`
	UploadedBy  int    `json:"uploaded_by"`
}

// ToEventDataset converts CreateDatasetRequest to EventDataset
func (req *CreateDatasetRequest) ToEventDataset() *EventDataset {
	return &EventDataset{
		Filename:    req.Filename,
		Description: req.Description,
		EventCount:  req.EventCount,
		UploadedBy:  req.UploadedBy,
	}
}