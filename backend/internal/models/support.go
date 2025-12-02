package models

import (
	"time"
)

type SupportCredential struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	IsURL        bool      `json:"is_url"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SupportCredentialResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	IsURL bool   `json:"is_url"`
}

func (s *SupportCredential) ToResponse() *SupportCredentialResponse {
	return &SupportCredentialResponse{
		Name:  s.Name,
		Value: s.Value,
		IsURL: s.IsURL,
	}
}
