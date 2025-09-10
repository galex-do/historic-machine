package repositories

import (
        "database/sql"
        "fmt"
        "time"
        "historical-events-backend/internal/models"
)

type DatasetRepository struct {
        db *sql.DB
}

func NewDatasetRepository(db *sql.DB) *DatasetRepository {
        return &DatasetRepository{db: db}
}

// GetAll retrieves all datasets with user information
func (r *DatasetRepository) GetAll() ([]models.EventDataset, error) {
        query := `
                SELECT d.id, d.filename, d.description, d.event_count, d.uploaded_by, 
                       d.created_at, d.updated_at, u.username
                FROM event_datasets d
                LEFT JOIN users u ON d.uploaded_by = u.id
                ORDER BY d.created_at DESC`
        
        rows, err := r.db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to query datasets: %w", err)
        }
        defer rows.Close()
        
        var datasets []models.EventDataset
        
        for rows.Next() {
                var dataset models.EventDataset
                var username sql.NullString
                var uploadedBy sql.NullInt32
                
                err := rows.Scan(&dataset.ID, &dataset.Filename, &dataset.Description, 
                        &dataset.EventCount, &uploadedBy, &dataset.CreatedAt, &dataset.UpdatedAt, &username)
                if err != nil {
                        return nil, fmt.Errorf("error scanning dataset: %w", err)
                }
                
                // Handle nullable uploaded_by field
                if uploadedBy.Valid {
                        dataset.UploadedBy = int(uploadedBy.Int32)
                } else {
                        dataset.UploadedBy = 0 // Default value for unknown user
                }
                
                if username.Valid {
                        dataset.UploadedByUsername = username.String
                } else {
                        dataset.UploadedByUsername = "Unknown User"
                }
                
                datasets = append(datasets, dataset)
        }
        
        if err = rows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating over datasets: %w", err)
        }
        
        return datasets, nil
}

// Create creates a new dataset
func (r *DatasetRepository) Create(dataset *models.EventDataset) (*models.EventDataset, error) {
        query := `
                INSERT INTO event_datasets (filename, description, event_count, uploaded_by, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id`
        
        now := time.Now()
        dataset.CreatedAt = now
        dataset.UpdatedAt = now
        
        err := r.db.QueryRow(query, dataset.Filename, dataset.Description, 
                dataset.EventCount, dataset.UploadedBy, dataset.CreatedAt, dataset.UpdatedAt).Scan(&dataset.ID)
        if err != nil {
                return nil, fmt.Errorf("failed to create dataset: %w", err)
        }
        
        return dataset, nil
}

// GetByID retrieves a dataset by ID
func (r *DatasetRepository) GetByID(id int) (*models.EventDataset, error) {
        query := `
                SELECT id, filename, description, event_count, uploaded_by, created_at, updated_at
                FROM event_datasets 
                WHERE id = $1`
        
        var dataset models.EventDataset
        var uploadedBy sql.NullInt32
        err := r.db.QueryRow(query, id).Scan(&dataset.ID, &dataset.Filename, 
                &dataset.Description, &dataset.EventCount, &uploadedBy, 
                &dataset.CreatedAt, &dataset.UpdatedAt)
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("dataset not found")
                }
                return nil, fmt.Errorf("failed to get dataset: %w", err)
        }
        
        // Handle nullable uploaded_by field
        if uploadedBy.Valid {
                dataset.UploadedBy = int(uploadedBy.Int32)
        } else {
                dataset.UploadedBy = 0 // Default value for unknown user
        }
        
        return &dataset, nil
}

// Delete deletes a dataset and all its associated events
func (r *DatasetRepository) Delete(id int) error {
        tx, err := r.db.Begin()
        if err != nil {
                return fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()
        
        // First delete all events associated with this dataset
        _, err = tx.Exec("DELETE FROM event_tags WHERE event_id IN (SELECT id FROM events WHERE dataset_id = $1)", id)
        if err != nil {
                return fmt.Errorf("failed to delete event tags: %w", err)
        }
        
        // Delete events
        _, err = tx.Exec("DELETE FROM events WHERE dataset_id = $1", id)
        if err != nil {
                return fmt.Errorf("failed to delete events: %w", err)
        }
        
        // Delete the dataset itself
        _, err = tx.Exec("DELETE FROM event_datasets WHERE id = $1", id)
        if err != nil {
                return fmt.Errorf("failed to delete dataset: %w", err)
        }
        
        return tx.Commit()
}

// UpdateEventCount updates the event count for a dataset
func (r *DatasetRepository) UpdateEventCount(id int, count int) error {
        query := `UPDATE event_datasets SET event_count = $1, updated_at = $2 WHERE id = $3`
        
        _, err := r.db.Exec(query, count, time.Now(), id)
        if err != nil {
                return fmt.Errorf("failed to update dataset event count: %w", err)
        }
        
        return nil
}