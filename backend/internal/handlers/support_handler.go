package handlers

import (
	"encoding/json"
	"net/http"

	"historical-events-backend/internal/database/repositories"
	"historical-events-backend/internal/models"
	"historical-events-backend/pkg/response"
)

type SupportHandler struct {
	repo *repositories.SupportRepository
}

func NewSupportHandler(repo *repositories.SupportRepository) *SupportHandler {
	return &SupportHandler{repo: repo}
}

func (h *SupportHandler) GetSupportCredentials(w http.ResponseWriter, r *http.Request) {
	credentials, err := h.repo.GetActiveCredentials()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch support credentials")
		return
	}

	var result []models.SupportCredentialResponse
	for _, c := range credentials {
		result = append(result, *c.ToResponse())
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *SupportHandler) CreateSupportCredential(w http.ResponseWriter, r *http.Request) {
	var credential models.SupportCredential
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	credential.IsActive = true
	if err := h.repo.Create(&credential); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create support credential")
		return
	}

	response.JSON(w, http.StatusCreated, credential)
}

func (h *SupportHandler) UpdateSupportCredential(w http.ResponseWriter, r *http.Request) {
	var credential models.SupportCredential
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.Update(&credential); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update support credential")
		return
	}

	response.JSON(w, http.StatusOK, credential)
}

func (h *SupportHandler) DeleteSupportCredential(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.Delete(req.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete support credential")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Deleted successfully"})
}
