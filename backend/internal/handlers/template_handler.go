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

// GetAllGroups handles GET /api/date-template-groups with locale support
func (h *TemplateHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
        // Get locale parameter (default to "en")
        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        
        groups, err := h.templateRepo.GetAllGroups()
        if err != nil {
                log.Printf("Error fetching date template groups: %v", err)
                response.InternalError(w, "Failed to fetch date template groups")
                return
        }
        
        // Populate legacy fields based on locale
        for i := range groups {
                groups[i].PopulateLegacyFields(locale)
        }
        
        response.Success(w, groups)
}

// GetTemplatesByGroup handles GET /api/date-templates/{group_id} with locale support
func (h *TemplateHandler) GetTemplatesByGroup(w http.ResponseWriter, r *http.Request) {
        // Get locale parameter (default to "en")
        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        
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
        
        // Populate legacy fields based on locale
        for i := range templates {
                templates[i].PopulateLegacyFields(locale)
        }
        
        response.Success(w, templates)
}

// GetAllTemplates handles GET /api/date-templates with locale support
func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
        // Get locale parameter (default to "en")
        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        
        templates, err := h.templateRepo.GetAllTemplates()
        if err != nil {
                log.Printf("Error fetching all date templates: %v", err)
                response.InternalError(w, "Failed to fetch all date templates")
                return
        }
        
        // Populate legacy fields based on locale
        for i := range templates {
                templates[i].PopulateLegacyFields(locale)
        }
        
        response.Success(w, templates)
}