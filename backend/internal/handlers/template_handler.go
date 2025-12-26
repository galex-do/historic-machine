package handlers

import (
        "encoding/json"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/internal/models"
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

// CreateGroup handles POST /api/date-template-groups
func (h *TemplateHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
        var group models.DateTemplateGroup
        if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }
        
        created, err := h.templateRepo.CreateGroup(&group)
        if err != nil {
                log.Printf("Error creating template group: %v", err)
                response.InternalError(w, "Failed to create template group")
                return
        }
        
        response.Created(w, created)
}

// UpdateGroup handles PUT /api/date-template-groups/{id}
func (h *TemplateHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid group ID")
                return
        }
        
        var group models.DateTemplateGroup
        if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }
        group.ID = id
        
        if err := h.templateRepo.UpdateGroup(&group); err != nil {
                log.Printf("Error updating template group: %v", err)
                response.InternalError(w, "Failed to update template group")
                return
        }
        
        response.Success(w, group)
}

// DeleteGroup handles DELETE /api/date-template-groups/{id}
func (h *TemplateHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid group ID")
                return
        }
        
        // Check if group has templates
        count, err := h.templateRepo.GetTemplateCount(id)
        if err != nil {
                log.Printf("Error checking template count: %v", err)
                response.InternalError(w, "Failed to check template count")
                return
        }
        
        if count > 0 {
                response.BadRequest(w, "Cannot delete group with templates. Delete templates first.")
                return
        }
        
        if err := h.templateRepo.DeleteGroup(id); err != nil {
                log.Printf("Error deleting template group: %v", err)
                response.InternalError(w, "Failed to delete template group")
                return
        }
        
        response.Success(w, map[string]string{"message": "Template group deleted successfully"})
}

// GetGroupByID handles GET /api/date-template-groups/{id}
func (h *TemplateHandler) GetGroupByID(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid group ID")
                return
        }
        
        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        
        group, err := h.templateRepo.GetGroupByID(id)
        if err != nil {
                log.Printf("Error fetching template group: %v", err)
                response.NotFound(w, "Template group not found")
                return
        }
        
        group.PopulateLegacyFields(locale)
        response.Success(w, group)
}

// CreateTemplate handles POST /api/date-templates
func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
        var template models.DateTemplate
        if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }
        
        created, err := h.templateRepo.CreateTemplate(&template)
        if err != nil {
                log.Printf("Error creating template: %v", err)
                response.InternalError(w, "Failed to create template")
                return
        }
        
        response.Created(w, created)
}

// UpdateTemplate handles PUT /api/date-templates/single/{id}
func (h *TemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid template ID")
                return
        }
        
        var template models.DateTemplate
        if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }
        template.ID = id
        
        if err := h.templateRepo.UpdateTemplate(&template); err != nil {
                log.Printf("Error updating template: %v", err)
                response.InternalError(w, "Failed to update template")
                return
        }
        
        response.Success(w, template)
}

// DeleteTemplate handles DELETE /api/date-templates/single/{id}
func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid template ID")
                return
        }
        
        if err := h.templateRepo.DeleteTemplate(id); err != nil {
                log.Printf("Error deleting template: %v", err)
                response.InternalError(w, "Failed to delete template")
                return
        }
        
        response.Success(w, map[string]string{"message": "Template deleted successfully"})
}

// GetTemplateByID handles GET /api/date-templates/single/{id}
func (h *TemplateHandler) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        
        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid template ID")
                return
        }
        
        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        
        template, err := h.templateRepo.GetTemplateByID(id)
        if err != nil {
                log.Printf("Error fetching template: %v", err)
                response.NotFound(w, "Template not found")
                return
        }
        
        template.PopulateLegacyFields(locale)
        response.Success(w, template)
}