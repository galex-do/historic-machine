package repositories

import (
	"database/sql"
	"fmt"
	"historical-events-backend/internal/models"
	"log"
)

// TemplateRepository handles template data operations
type TemplateRepository struct {
	db *sql.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

// GetAllGroups retrieves all date template groups
func (r *TemplateRepository) GetAllGroups() ([]models.DateTemplateGroup, error) {
	query := `
		SELECT id, name, description, display_order
		FROM date_template_groups 
		ORDER BY display_order`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query template groups: %w", err)
	}
	defer rows.Close()
	
	var groups []models.DateTemplateGroup
	for rows.Next() {
		var group models.DateTemplateGroup
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.DisplayOrder)
		if err != nil {
			log.Printf("Error scanning date template group: %v", err)
			continue
		}
		groups = append(groups, group)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over template groups: %w", err)
	}
	
	return groups, nil
}

// GetTemplatesByGroup retrieves templates for a specific group
func (r *TemplateRepository) GetTemplatesByGroup(groupID int) ([]models.DateTemplate, error) {
	query := `
		SELECT id, group_id, group_name, name, description,
		       start_date, start_era, end_date, end_era, display_order,
		       start_display_date, end_display_date
		FROM date_templates_with_display 
		WHERE group_id = $1
		ORDER BY display_order`
	
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates for group %d: %w", groupID, err)
	}
	defer rows.Close()
	
	var templates []models.DateTemplate
	for rows.Next() {
		var template models.DateTemplate
		err := rows.Scan(&template.ID, &template.GroupID, &template.GroupName, 
			&template.Name, &template.Description, &template.StartDate, &template.StartEra,
			&template.EndDate, &template.EndEra, &template.DisplayOrder,
			&template.StartDisplayDate, &template.EndDisplayDate)
		if err != nil {
			log.Printf("Error scanning date template: %v", err)
			continue
		}
		templates = append(templates, template)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over templates: %w", err)
	}
	
	return templates, nil
}

// GetAllTemplates retrieves all date templates across all groups
func (r *TemplateRepository) GetAllTemplates() ([]models.DateTemplate, error) {
	query := `
		SELECT id, group_id, group_name, name, description,
		       start_date, start_era, end_date, end_era, display_order,
		       start_display_date, end_display_date
		FROM date_templates_with_display 
		ORDER BY group_id, display_order`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all templates: %w", err)
	}
	defer rows.Close()
	
	var templates []models.DateTemplate
	for rows.Next() {
		var template models.DateTemplate
		err := rows.Scan(&template.ID, &template.GroupID, &template.GroupName, 
			&template.Name, &template.Description, &template.StartDate, &template.StartEra,
			&template.EndDate, &template.EndEra, &template.DisplayOrder,
			&template.StartDisplayDate, &template.EndDisplayDate)
		if err != nil {
			log.Printf("Error scanning date template: %v", err)
			continue
		}
		templates = append(templates, template)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over all templates: %w", err)
	}
	
	return templates, nil
}