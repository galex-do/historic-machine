package models

import "time"

// Tag represents a tag that can be associated with events
type Tag struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTagRequest represents the request payload for creating a tag
type CreateTagRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
	Color       string `json:"color,omitempty"`
}

// UpdateTagRequest represents the request payload for updating a tag
type UpdateTagRequest struct {
	Name        string `json:"name,omitempty" validate:"max=100"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

// EventTag represents the many-to-many relationship between events and tags
type EventTag struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	TagID     int       `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ToTag converts CreateTagRequest to Tag
func (req *CreateTagRequest) ToTag() *Tag {
	color := req.Color
	if color == "" {
		color = "#3b82f6" // Default blue color
	}
	
	return &Tag{
		Name:        req.Name,
		Description: req.Description,
		Color:       color,
	}
}