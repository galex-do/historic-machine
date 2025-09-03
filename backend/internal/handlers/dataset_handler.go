package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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