package repositories

import (
        "database/sql"
        "encoding/json"
        "fmt"
        "historical-events-backend/internal/models"
        "log"
)

// EventRepository handles event data operations
type EventRepository struct {
        db *sql.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *sql.DB) *EventRepository {
        return &EventRepository{db: db}
}

// GetAll retrieves all events from the database
func (r *EventRepository) GetAll() ([]models.HistoricalEvent, error) {
        query := `
                SELECT id, name, description, latitude, longitude, event_date, era, lens_type, display_date, tags
                FROM events_with_display_dates 
                ORDER BY astronomical_year ASC`
        
        rows, err := r.db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to query events: %w", err)
        }
        defer rows.Close()
        
        var events []models.HistoricalEvent
        
        for rows.Next() {
                var event models.HistoricalEvent
                var tagsJSON []byte
                
                err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Latitude, 
                        &event.Longitude, &event.EventDate, &event.Era, &event.LensType, &event.DisplayDate, &tagsJSON)
                if err != nil {
                        log.Printf("Error scanning event: %v", err)
                        continue
                }
                
                // Parse tags JSON
                if len(tagsJSON) > 0 {
                        var tags []models.Tag
                        if err := json.Unmarshal(tagsJSON, &tags); err != nil {
                                log.Printf("Error unmarshaling tags for event %d: %v", event.ID, err)
                                tags = []models.Tag{}
                        }
                        event.Tags = tags
                } else {
                        event.Tags = []models.Tag{}
                }
                
                events = append(events, event)
        }
        
        if err = rows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating over events: %w", err)
        }
        
        return events, nil
}

// GetByID retrieves a single event by ID
func (r *EventRepository) GetByID(id int) (*models.HistoricalEvent, error) {
        query := `
                SELECT id, name, description, latitude, longitude, event_date, era, lens_type, display_date, tags
                FROM events_with_display_dates 
                WHERE id = $1`
        
        var event models.HistoricalEvent
        var tagsJSON []byte
        err := r.db.QueryRow(query, id).Scan(
                &event.ID, &event.Name, &event.Description, &event.Latitude,
                &event.Longitude, &event.EventDate, &event.Era, &event.LensType, &event.DisplayDate, &tagsJSON)
        
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("event with id %d not found", id)
                }
                return nil, fmt.Errorf("failed to get event by id: %w", err)
        }
        
        // Parse tags JSON
        if len(tagsJSON) > 0 {
                var tags []models.Tag
                if err := json.Unmarshal(tagsJSON, &tags); err != nil {
                        log.Printf("Error unmarshaling tags for event %d: %v", event.ID, err)
                        tags = []models.Tag{}
                }
                event.Tags = tags
        } else {
                event.Tags = []models.Tag{}
        }
        
        return &event, nil
}

// Create creates a new event in the database
func (r *EventRepository) Create(event *models.HistoricalEvent) (*models.HistoricalEvent, error) {
        query := `
                INSERT INTO events (name, description, latitude, longitude, event_date, era, lens_type, dataset_id) 
                VALUES ($1, $2, $3::double precision, $4::double precision, $5, $6, $7, $8) 
                RETURNING id`
        
        var createdEvent = *event
        
        err := r.db.QueryRow(query, event.Name, event.Description, event.Latitude, 
                event.Longitude, event.EventDate, event.Era, event.LensType, event.DatasetID).
                Scan(&createdEvent.ID)
        
        if err != nil {
                return nil, fmt.Errorf("failed to create event: %w", err)
        }
        
        return &createdEvent, nil
}

// GetInBoundingBox retrieves events within a geographical bounding box
func (r *EventRepository) GetInBoundingBox(minLat, minLng, maxLat, maxLng float64) ([]models.HistoricalEvent, error) {
        query := `
                SELECT id, name, description, latitude, longitude, event_date, era, lens_type,
                       ST_X(location::geometry) as lng, ST_Y(location::geometry) as lat
                FROM events 
                WHERE location && ST_MakeEnvelope($1, $2, $3, $4, 4326)
                ORDER BY event_date DESC`
        
        rows, err := r.db.Query(query, minLng, minLat, maxLng, maxLat)
        if err != nil {
                return nil, fmt.Errorf("bounding box query failed: %w", err)
        }
        defer rows.Close()
        
        var events []models.HistoricalEvent
        for rows.Next() {
                var event models.HistoricalEvent
                var lng, lat float64
                
                err := rows.Scan(&event.ID, &event.Name, &event.Description, 
                        &event.Latitude, &event.Longitude, &event.EventDate, &event.Era, &event.LensType,
                        &lng, &lat)
                if err != nil {
                        log.Printf("Error scanning bounding box event: %v", err)
                        continue
                }
                
                // Update with PostGIS coordinates for accuracy
                event.Longitude = lng
                event.Latitude = lat
                events = append(events, event)
        }
        
        if err = rows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating over bounding box events: %w", err)
        }
        
        return events, nil
}

// Update updates an existing event in the database
func (r *EventRepository) Update(event *models.HistoricalEvent) (*models.HistoricalEvent, error) {
        query := `
                UPDATE events 
                SET name = $2, description = $3, latitude = $4::double precision, longitude = $5::double precision, 
                    event_date = $6, era = $7, lens_type = $8
                WHERE id = $1
                RETURNING id, name, description, latitude, longitude, event_date, era, lens_type`
        
        var updatedEvent models.HistoricalEvent
        err := r.db.QueryRow(query, event.ID, event.Name, event.Description, 
                event.Latitude, event.Longitude, event.EventDate, event.Era, event.LensType).
                Scan(&updatedEvent.ID, &updatedEvent.Name, &updatedEvent.Description, 
                &updatedEvent.Latitude, &updatedEvent.Longitude, &updatedEvent.EventDate, 
                &updatedEvent.Era, &updatedEvent.LensType)
        
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("event with id %d not found", event.ID)
                }
                return nil, fmt.Errorf("failed to update event: %w", err)
        }
        
        return &updatedEvent, nil
}

// Delete removes an event from the database
func (r *EventRepository) Delete(id int) error {
        query := `DELETE FROM events WHERE id = $1`
        
        result, err := r.db.Exec(query, id)
        if err != nil {
                return fmt.Errorf("failed to delete event: %w", err)
        }
        
        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get affected rows: %w", err)
        }
        
        if rowsAffected == 0 {
                return fmt.Errorf("event with id %d not found", id)
        }
        
        return nil
}

// ValidateCoordinates validates latitude and longitude values
func (r *EventRepository) ValidateCoordinates(lat, lng float64) error {
        if lat < -90 || lat > 90 {
                return fmt.Errorf("latitude must be between -90 and 90 degrees, got: %f", lat)
        }
        if lng < -180 || lng > 180 {
                return fmt.Errorf("longitude must be between -180 and 180 degrees, got: %f", lng)
        }
        return nil
}