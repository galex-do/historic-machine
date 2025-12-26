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

// GetAllGroups retrieves all date template groups with localized fields
func (r *TemplateRepository) GetAllGroups() ([]models.DateTemplateGroup, error) {
        query := `
                SELECT id, name, description, name_en, name_ru, description_en, description_ru, display_order
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
                err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.NameEn, &group.NameRu, &group.DescriptionEn, &group.DescriptionRu, &group.DisplayOrder)
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

// GetTemplatesByGroup retrieves templates for a specific group with localized fields
func (r *TemplateRepository) GetTemplatesByGroup(groupID int) ([]models.DateTemplate, error) {
        query := `
                SELECT id, group_id, group_name, name, description,
                       name_en, name_ru, description_en, description_ru,
                       group_name_en, group_name_ru,
                       start_date, start_era, end_date, end_era, display_order,
                       start_display_date, end_display_date
                FROM date_templates_with_display 
                WHERE group_id = $1
                ORDER BY start_astronomical_year ASC`
        
        rows, err := r.db.Query(query, groupID)
        if err != nil {
                return nil, fmt.Errorf("failed to query templates for group %d: %w", groupID, err)
        }
        defer rows.Close()
        
        var templates []models.DateTemplate
        for rows.Next() {
                var template models.DateTemplate
                err := rows.Scan(&template.ID, &template.GroupID, &template.GroupName, 
                        &template.Name, &template.Description,
                        &template.NameEn, &template.NameRu, &template.DescriptionEn, &template.DescriptionRu,
                        &template.GroupNameEn, &template.GroupNameRu,
                        &template.StartDate, &template.StartEra,
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

// GetAllTemplates retrieves all date templates across all groups with localized fields
func (r *TemplateRepository) GetAllTemplates() ([]models.DateTemplate, error) {
        query := `
                SELECT id, group_id, group_name, name, description,
                       name_en, name_ru, description_en, description_ru,
                       group_name_en, group_name_ru,
                       start_date, start_era, end_date, end_era, display_order,
                       start_display_date, end_display_date
                FROM date_templates_with_display 
                ORDER BY group_id, start_astronomical_year ASC`
        
        rows, err := r.db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to query all templates: %w", err)
        }
        defer rows.Close()
        
        var templates []models.DateTemplate
        for rows.Next() {
                var template models.DateTemplate
                err := rows.Scan(&template.ID, &template.GroupID, &template.GroupName, 
                        &template.Name, &template.Description,
                        &template.NameEn, &template.NameRu, &template.DescriptionEn, &template.DescriptionRu,
                        &template.GroupNameEn, &template.GroupNameRu,
                        &template.StartDate, &template.StartEra,
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

// CreateGroup creates a new date template group
func (r *TemplateRepository) CreateGroup(group *models.DateTemplateGroup) (*models.DateTemplateGroup, error) {
        query := `
                INSERT INTO date_template_groups (name, description, name_en, name_ru, description_en, description_ru, display_order)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id`
        
        err := r.db.QueryRow(query, 
                group.NameEn, group.DescriptionEn, 
                group.NameEn, group.NameRu, 
                group.DescriptionEn, group.DescriptionRu, 
                group.DisplayOrder,
        ).Scan(&group.ID)
        
        if err != nil {
                return nil, fmt.Errorf("failed to create template group: %w", err)
        }
        
        return group, nil
}

// UpdateGroup updates an existing date template group
func (r *TemplateRepository) UpdateGroup(group *models.DateTemplateGroup) error {
        query := `
                UPDATE date_template_groups 
                SET name = $2, description = $3, name_en = $4, name_ru = $5, 
                    description_en = $6, description_ru = $7, display_order = $8
                WHERE id = $1`
        
        result, err := r.db.Exec(query, 
                group.ID,
                group.NameEn, group.DescriptionEn,
                group.NameEn, group.NameRu, 
                group.DescriptionEn, group.DescriptionRu, 
                group.DisplayOrder,
        )
        
        if err != nil {
                return fmt.Errorf("failed to update template group: %w", err)
        }
        
        rows, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to check affected rows: %w", err)
        }
        
        if rows == 0 {
                return fmt.Errorf("template group not found")
        }
        
        return nil
}

// DeleteGroup deletes a date template group
func (r *TemplateRepository) DeleteGroup(id int) error {
        query := `DELETE FROM date_template_groups WHERE id = $1`
        
        result, err := r.db.Exec(query, id)
        if err != nil {
                return fmt.Errorf("failed to delete template group: %w", err)
        }
        
        rows, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to check affected rows: %w", err)
        }
        
        if rows == 0 {
                return fmt.Errorf("template group not found")
        }
        
        return nil
}

// GetGroupByID retrieves a date template group by ID
func (r *TemplateRepository) GetGroupByID(id int) (*models.DateTemplateGroup, error) {
        query := `
                SELECT id, name, description, name_en, name_ru, description_en, description_ru, display_order
                FROM date_template_groups 
                WHERE id = $1`
        
        var group models.DateTemplateGroup
        err := r.db.QueryRow(query, id).Scan(
                &group.ID, &group.Name, &group.Description, 
                &group.NameEn, &group.NameRu, 
                &group.DescriptionEn, &group.DescriptionRu, 
                &group.DisplayOrder,
        )
        
        if err == sql.ErrNoRows {
                return nil, fmt.Errorf("template group not found")
        }
        if err != nil {
                return nil, fmt.Errorf("failed to get template group: %w", err)
        }
        
        return &group, nil
}

// GetTemplateCount returns the number of templates in a group
func (r *TemplateRepository) GetTemplateCount(groupID int) (int, error) {
        query := `SELECT COUNT(*) FROM date_templates WHERE group_id = $1`
        
        var count int
        err := r.db.QueryRow(query, groupID).Scan(&count)
        if err != nil {
                return 0, fmt.Errorf("failed to count templates: %w", err)
        }
        
        return count, nil
}

// CreateTemplate creates a new date template
func (r *TemplateRepository) CreateTemplate(template *models.DateTemplate) (*models.DateTemplate, error) {
        query := `
                INSERT INTO date_templates (group_id, name, description, name_en, name_ru, description_en, description_ru, 
                                            start_date, start_era, end_date, end_era, display_order)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
                RETURNING id`
        
        err := r.db.QueryRow(query, 
                template.GroupID,
                template.NameEn, template.DescriptionEn,
                template.NameEn, template.NameRu, 
                template.DescriptionEn, template.DescriptionRu,
                template.StartDate, template.StartEra,
                template.EndDate, template.EndEra,
                template.DisplayOrder,
        ).Scan(&template.ID)
        
        if err != nil {
                return nil, fmt.Errorf("failed to create template: %w", err)
        }
        
        return template, nil
}

// UpdateTemplate updates an existing date template
func (r *TemplateRepository) UpdateTemplate(template *models.DateTemplate) error {
        query := `
                UPDATE date_templates 
                SET group_id = $2, name = $3, description = $4, name_en = $5, name_ru = $6, 
                    description_en = $7, description_ru = $8, start_date = $9, start_era = $10, 
                    end_date = $11, end_era = $12, display_order = $13
                WHERE id = $1`
        
        result, err := r.db.Exec(query, 
                template.ID,
                template.GroupID,
                template.NameEn, template.DescriptionEn,
                template.NameEn, template.NameRu, 
                template.DescriptionEn, template.DescriptionRu,
                template.StartDate, template.StartEra,
                template.EndDate, template.EndEra,
                template.DisplayOrder,
        )
        
        if err != nil {
                return fmt.Errorf("failed to update template: %w", err)
        }
        
        rows, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to check affected rows: %w", err)
        }
        
        if rows == 0 {
                return fmt.Errorf("template not found")
        }
        
        return nil
}

// DeleteTemplate deletes a date template
func (r *TemplateRepository) DeleteTemplate(id int) error {
        query := `DELETE FROM date_templates WHERE id = $1`
        
        result, err := r.db.Exec(query, id)
        if err != nil {
                return fmt.Errorf("failed to delete template: %w", err)
        }
        
        rows, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to check affected rows: %w", err)
        }
        
        if rows == 0 {
                return fmt.Errorf("template not found")
        }
        
        return nil
}

// GetTemplateByID retrieves a date template by ID
func (r *TemplateRepository) GetTemplateByID(id int) (*models.DateTemplate, error) {
        query := `
                SELECT id, group_id, group_name, name, description,
                       name_en, name_ru, description_en, description_ru,
                       group_name_en, group_name_ru,
                       start_date, start_era, end_date, end_era, display_order,
                       start_display_date, end_display_date
                FROM date_templates_with_display 
                WHERE id = $1`
        
        var template models.DateTemplate
        err := r.db.QueryRow(query, id).Scan(
                &template.ID, &template.GroupID, &template.GroupName,
                &template.Name, &template.Description,
                &template.NameEn, &template.NameRu, &template.DescriptionEn, &template.DescriptionRu,
                &template.GroupNameEn, &template.GroupNameRu,
                &template.StartDate, &template.StartEra,
                &template.EndDate, &template.EndEra, &template.DisplayOrder,
                &template.StartDisplayDate, &template.EndDisplayDate,
        )
        
        if err == sql.ErrNoRows {
                return nil, fmt.Errorf("template not found")
        }
        if err != nil {
                return nil, fmt.Errorf("failed to get template: %w", err)
        }
        
        return &template, nil
}