package handlers

import (
	"net/http"
	"os"

	"historical-events-backend/pkg/response"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

type PublicConfigResponse struct {
	ContactEmail string `json:"contact_email,omitempty"`
}

func (h *ConfigHandler) GetPublicConfig(w http.ResponseWriter, r *http.Request) {
	config := PublicConfigResponse{
		ContactEmail: os.Getenv("CONTACT_EMAIL"),
	}

	response.JSON(w, http.StatusOK, config)
}
