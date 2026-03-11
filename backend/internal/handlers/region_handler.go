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

type RegionHandler struct {
        regionRepo *repositories.RegionRepository
}

func NewRegionHandler(regionRepo *repositories.RegionRepository) *RegionHandler {
        return &RegionHandler{
                regionRepo: regionRepo,
        }
}

func (h *RegionHandler) GetAllRegions(w http.ResponseWriter, r *http.Request) {
        regions, err := h.regionRepo.GetAll()
        if err != nil {
                log.Printf("Error fetching regions: %v", err)
                response.InternalError(w, "Failed to fetch regions")
                return
        }

        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        for i := range regions {
                regions[i].PopulateLegacyFields(locale)
        }

        type regionWithTemplates struct {
                models.Region
                TemplateIDs []int `json:"template_ids"`
        }

        result := make([]regionWithTemplates, 0, len(regions))
        for _, region := range regions {
                templateIDs, err := h.regionRepo.GetTemplateIDsByRegion(region.ID)
                if err != nil {
                        log.Printf("Error fetching template IDs for region %d: %v", region.ID, err)
                        templateIDs = []int{}
                }
                if templateIDs == nil {
                        templateIDs = []int{}
                }
                result = append(result, regionWithTemplates{
                        Region:      region,
                        TemplateIDs: templateIDs,
                })
        }

        response.Success(w, result)
}

func (h *RegionHandler) GetRegionByID(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
                response.BadRequest(w, "Invalid region ID")
                return
        }

        region, err := h.regionRepo.GetByID(id)
        if err != nil {
                log.Printf("Error fetching region %d: %v", id, err)
                response.NotFound(w, "Region not found")
                return
        }

        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        region.PopulateLegacyFields(locale)

        templateIDs, err := h.regionRepo.GetTemplateIDsByRegion(id)
        if err != nil {
                log.Printf("Error fetching template IDs for region %d: %v", id, err)
                templateIDs = []int{}
        }
        if templateIDs == nil {
                templateIDs = []int{}
        }

        type regionWithTemplates struct {
                models.Region
                TemplateIDs []int `json:"template_ids"`
        }

        response.Success(w, regionWithTemplates{
                Region:      *region,
                TemplateIDs: templateIDs,
        })
}

func (h *RegionHandler) GetRegionsByTemplate(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        templateID, err := strconv.Atoi(vars["id"])
        if err != nil {
                response.BadRequest(w, "Invalid template ID")
                return
        }

        regions, err := h.regionRepo.GetByTemplateID(templateID)
        if err != nil {
                log.Printf("Error fetching regions for template %d: %v", templateID, err)
                response.InternalError(w, "Failed to fetch regions for template")
                return
        }

        locale := r.URL.Query().Get("locale")
        if locale == "" {
                locale = "en"
        }
        for i := range regions {
                regions[i].PopulateLegacyFields(locale)
        }

        if regions == nil {
                regions = []models.Region{}
        }

        response.Success(w, regions)
}

func (h *RegionHandler) CreateRegion(w http.ResponseWriter, r *http.Request) {
        var req models.RegionCreate
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }

        if len(req.GeoJSON) == 0 {
                response.BadRequest(w, "GeoJSON is required")
                return
        }

        created, err := h.regionRepo.Create(&req)
        if err != nil {
                log.Printf("Error creating region: %v", err)
                response.InternalError(w, "Failed to create region")
                return
        }

        response.Created(w, created)
}

func (h *RegionHandler) UpdateRegion(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
                response.BadRequest(w, "Invalid region ID")
                return
        }

        var req models.RegionUpdate
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }

        updated, err := h.regionRepo.Update(id, &req)
        if err != nil {
                log.Printf("Error updating region %d: %v", id, err)
                response.InternalError(w, "Failed to update region")
                return
        }

        response.Success(w, updated)
}

func (h *RegionHandler) DeleteRegion(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
                response.BadRequest(w, "Invalid region ID")
                return
        }

        if err := h.regionRepo.Delete(id); err != nil {
                log.Printf("Error deleting region %d: %v", id, err)
                response.InternalError(w, "Failed to delete region")
                return
        }

        response.Success(w, map[string]string{"message": "Region deleted successfully"})
}

func (h *RegionHandler) LinkRegionToTemplates(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        regionID, err := strconv.Atoi(vars["id"])
        if err != nil {
                response.BadRequest(w, "Invalid region ID")
                return
        }

        var req struct {
                TemplateIDs []int `json:"template_ids"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid request body")
                return
        }

        for _, templateID := range req.TemplateIDs {
                if err := h.regionRepo.LinkToTemplate(regionID, templateID); err != nil {
                        log.Printf("Error linking region %d to template %d: %v", regionID, templateID, err)
                        response.InternalError(w, "Failed to link region to template")
                        return
                }
        }

        response.Success(w, map[string]string{"message": "Region linked to templates successfully"})
}

func (h *RegionHandler) UnlinkRegionFromTemplate(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        regionID, err := strconv.Atoi(vars["id"])
        if err != nil {
                response.BadRequest(w, "Invalid region ID")
                return
        }

        templateID, err := strconv.Atoi(vars["templateId"])
        if err != nil {
                response.BadRequest(w, "Invalid template ID")
                return
        }

        if err := h.regionRepo.UnlinkFromTemplate(regionID, templateID); err != nil {
                log.Printf("Error unlinking region %d from template %d: %v", regionID, templateID, err)
                response.InternalError(w, "Failed to unlink region from template")
                return
        }

        response.Success(w, map[string]string{"message": "Region unlinked from template successfully"})
}
