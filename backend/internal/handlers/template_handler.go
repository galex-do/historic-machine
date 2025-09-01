package handlers

import (
	"historical-events-backend/internal/database/repositories"
	"historical-events-backend/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// TemplateHandler handles HTTP requests for date templates
type TemplateHandler struct {
	templateRepo *repositories.TemplateRepository
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(templateRepo *repositories.TemplateRepository) *TemplateHandler {
	return &TemplateHandler{
		templateRepo: templateRepo,
	}
}

// GetAllGroups handles GET /api/date-template-groups
func (h *TemplateHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.templateRepo.GetAllGroups()
	if err != nil {
		log.Printf("Error fetching date template groups: %v", err)
		response.InternalError(w, "Failed to fetch date template groups")
		return
	}
	
	response.Success(w, groups)
}

// GetTemplatesByGroup handles GET /api/date-templates/{group_id}
func (h *TemplateHandler) GetTemplatesByGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupIDStr := vars["group_id"]
	
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid group ID")
		return
	}
	
	templates, err := h.templateRepo.GetTemplatesByGroup(groupID)
	if err != nil {
		log.Printf("Error fetching date templates for group %d: %v", groupID, err)
		response.InternalError(w, "Failed to fetch date templates")
		return
	}
	
	response.Success(w, templates)
}

// GetAllTemplates handles GET /api/date-templates
func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.templateRepo.GetAllTemplates()
	if err != nil {
		log.Printf("Error fetching all date templates: %v", err)
		response.InternalError(w, "Failed to fetch all date templates")
		return
	}
	
	response.Success(w, templates)
}