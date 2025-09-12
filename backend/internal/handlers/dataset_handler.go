package handlers

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "strconv"
        "time"

        "historical-events-backend/internal/models"
        "historical-events-backend/internal/database/repositories"
        "historical-events-backend/pkg/response"

        "github.com/gorilla/mux"
)

type DatasetHandler struct {
        datasetRepo *repositories.DatasetRepository
        eventRepo   *repositories.EventRepository
}

// NewDatasetHandler creates a new dataset handler
func NewDatasetHandler(datasetRepo *repositories.DatasetRepository, eventRepo *repositories.EventRepository) *DatasetHandler {
        return &DatasetHandler{
                datasetRepo: datasetRepo,
                eventRepo:   eventRepo,
        }
}

// GetAllDatasets handles GET /api/datasets
func (h *DatasetHandler) GetAllDatasets(w http.ResponseWriter, r *http.Request) {
        datasets, err := h.datasetRepo.GetAll()
        if err != nil {
                log.Printf("Error retrieving datasets: %v", err)
                response.InternalError(w, "Failed to retrieve datasets")
                return
        }

        response.Success(w, datasets, "Datasets retrieved successfully")
}

// GetDatasetByID handles GET /api/datasets/{id}
func (h *DatasetHandler) GetDatasetByID(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr, exists := vars["id"]
        if !exists {
                response.BadRequest(w, "Dataset ID is required")
                return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid dataset ID")
                return
        }

        dataset, err := h.datasetRepo.GetByID(id)
        if err != nil {
                if err.Error() == "dataset not found" {
                        response.NotFound(w, "Dataset not found")
                        return
                }
                log.Printf("Error retrieving dataset %d: %v", id, err)
                response.InternalError(w, "Failed to retrieve dataset")
                return
        }

        response.Success(w, dataset, "Dataset retrieved successfully")
}

// CreateDataset handles POST /api/datasets
func (h *DatasetHandler) CreateDataset(w http.ResponseWriter, r *http.Request) {
        var req models.CreateDatasetRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                response.BadRequest(w, "Invalid JSON format")
                return
        }

        // Convert request to dataset model
        dataset := req.ToEventDataset()

        // Create dataset
        createdDataset, err := h.datasetRepo.Create(dataset)
        if err != nil {
                log.Printf("Error creating dataset: %v", err)
                response.InternalError(w, "Failed to create dataset")
                return
        }

        response.Created(w, createdDataset, "Dataset created successfully")
}

// DeleteDataset handles DELETE /api/datasets/{id}
// This will cascade delete all events associated with the dataset
func (h *DatasetHandler) DeleteDataset(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr, exists := vars["id"]
        if !exists {
                response.BadRequest(w, "Dataset ID is required")
                return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid dataset ID")
                return
        }

        // Check if dataset exists
        _, err = h.datasetRepo.GetByID(id)
        if err != nil {
                if err.Error() == "dataset not found" {
                        response.NotFound(w, "Dataset not found")
                        return
                }
                log.Printf("Error checking dataset %d: %v", id, err)
                response.InternalError(w, "Failed to check dataset")
                return
        }

        // Delete dataset and all associated events
        err = h.datasetRepo.Delete(id)
        if err != nil {
                log.Printf("Error deleting dataset %d: %v", id, err)
                response.InternalError(w, "Failed to delete dataset")
                return
        }

        response.Success(w, map[string]interface{}{
                "id": id,
                "message": "Dataset and all associated events deleted successfully",
        }, "Dataset deleted successfully")
}

// ExportDataset handles GET /api/datasets/{id}/export
func (h *DatasetHandler) ExportDataset(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr, exists := vars["id"]
        if !exists {
                response.BadRequest(w, "Dataset ID is required")
                return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
                response.BadRequest(w, "Invalid dataset ID")
                return
        }

        // Get dataset metadata
        dataset, err := h.datasetRepo.GetByID(id)
        if err != nil {
                if err.Error() == "dataset not found" {
                        response.NotFound(w, "Dataset not found")
                        return
                }
                log.Printf("Error retrieving dataset %d: %v", id, err)
                response.InternalError(w, "Failed to retrieve dataset")
                return
        }

        // Get all events for this dataset
        events, err := h.eventRepo.GetByDatasetID(id)
        if err != nil {
                log.Printf("Error retrieving events for dataset %d: %v", id, err)
                response.InternalError(w, "Failed to retrieve events for dataset")
                return
        }

        // Convert events to export format (matching import format)
        exportEvents := make([]map[string]interface{}, len(events))
        
        for i, event := range events {
                // Format date as DD.MM.YYYY
                eventDate, err := time.Parse(time.RFC3339, event.EventDate.Format(time.RFC3339))
                if err != nil {
                        log.Printf("Error parsing event date for event %d: %v", event.ID, err)
                        continue
                }
                
                dateStr := fmt.Sprintf("%02d.%02d.%04d", 
                        eventDate.Day(), eventDate.Month(), eventDate.Year())

                // Extract tag names only (not IDs)
                tagNames := make([]string, len(event.Tags))
                for j, tag := range event.Tags {
                        tagNames[j] = tag.Name
                }

                exportEvent := map[string]interface{}{
                        "date":        dateStr,
                        "era":         event.Era,
                        "latitude":    event.Latitude,
                        "longitude":   event.Longitude,
                        "type":        event.LensType, // Export as "type" for compatibility
                        "tags":        tagNames,
                }

                // Include locale-specific fields
                if event.NameEn != "" {
                        exportEvent["name_en"] = event.NameEn
                }
                if event.NameRu != "" {
                        exportEvent["name_ru"] = event.NameRu
                }
                if event.DescriptionEn != nil && *event.DescriptionEn != "" {
                        exportEvent["description_en"] = *event.DescriptionEn
                }
                if event.DescriptionRu != nil && *event.DescriptionRu != "" {
                        exportEvent["description_ru"] = *event.DescriptionRu
                }

                // Include legacy fields for backward compatibility
                exportEvent["name"] = event.Name
                exportEvent["description"] = event.Description

                // Include source if available
                if event.Source != nil {
                        exportEvent["source"] = *event.Source
                }

                exportEvents[i] = exportEvent
        }

        // Create export format matching the import structure
        exportData := map[string]interface{}{
                "filename":    dataset.Filename,
                "description": dataset.Description,
                "events":      exportEvents,
        }

        // Set appropriate headers for file download
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", dataset.Filename))
        
        response.Success(w, exportData, "Dataset exported successfully")
}