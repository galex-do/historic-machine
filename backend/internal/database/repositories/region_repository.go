package repositories

import (
	"database/sql"
	"fmt"
	"historical-events-backend/internal/models"
)

type RegionRepository struct {
	db *sql.DB
}

func NewRegionRepository(db *sql.DB) *RegionRepository {
	return &RegionRepository{db: db}
}

func (r *RegionRepository) GetAll() ([]models.Region, error) {
	query := `
		SELECT id, name, name_en, name_ru, description, description_en, description_ru,
		       geojson, color, fill_opacity, border_color, border_width,
		       created_at, updated_at
		FROM regions
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query regions: %w", err)
	}
	defer rows.Close()

	var regions []models.Region
	for rows.Next() {
		var region models.Region
		err := rows.Scan(
			&region.ID, &region.Name, &region.NameEn, &region.NameRu,
			&region.Description, &region.DescriptionEn, &region.DescriptionRu,
			&region.GeoJSON, &region.Color, &region.FillOpacity,
			&region.BorderColor, &region.BorderWidth,
			&region.CreatedAt, &region.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan region: %w", err)
		}
		regions = append(regions, region)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over regions: %w", err)
	}

	return regions, nil
}

func (r *RegionRepository) GetByID(id int) (*models.Region, error) {
	query := `
		SELECT id, name, name_en, name_ru, description, description_en, description_ru,
		       geojson, color, fill_opacity, border_color, border_width,
		       created_at, updated_at
		FROM regions
		WHERE id = $1`

	var region models.Region
	err := r.db.QueryRow(query, id).Scan(
		&region.ID, &region.Name, &region.NameEn, &region.NameRu,
		&region.Description, &region.DescriptionEn, &region.DescriptionRu,
		&region.GeoJSON, &region.Color, &region.FillOpacity,
		&region.BorderColor, &region.BorderWidth,
		&region.CreatedAt, &region.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("region not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get region: %w", err)
	}

	return &region, nil
}

func (r *RegionRepository) GetByTemplateID(templateID int) ([]models.Region, error) {
	query := `
		SELECT r.id, r.name, r.name_en, r.name_ru, r.description, r.description_en, r.description_ru,
		       r.geojson, r.color, r.fill_opacity, r.border_color, r.border_width,
		       r.created_at, r.updated_at
		FROM regions r
		JOIN template_regions tr ON r.id = tr.region_id
		WHERE tr.template_id = $1
		ORDER BY r.name ASC`

	rows, err := r.db.Query(query, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to query regions for template %d: %w", templateID, err)
	}
	defer rows.Close()

	var regions []models.Region
	for rows.Next() {
		var region models.Region
		err := rows.Scan(
			&region.ID, &region.Name, &region.NameEn, &region.NameRu,
			&region.Description, &region.DescriptionEn, &region.DescriptionRu,
			&region.GeoJSON, &region.Color, &region.FillOpacity,
			&region.BorderColor, &region.BorderWidth,
			&region.CreatedAt, &region.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan region: %w", err)
		}
		regions = append(regions, region)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over regions: %w", err)
	}

	return regions, nil
}

func (r *RegionRepository) Create(region *models.RegionCreate) (*models.Region, error) {
	color := "#4f46e5"
	if region.Color != "" {
		color = region.Color
	}
	var fillOpacity float32 = 0.2
	if region.FillOpacity != nil {
		fillOpacity = *region.FillOpacity
	}
	borderColor := "#4f46e5"
	if region.BorderColor != "" {
		borderColor = region.BorderColor
	}
	var borderWidth float32 = 2
	if region.BorderWidth != nil {
		borderWidth = *region.BorderWidth
	}

	nameEn := region.NameEn
	if nameEn == "" {
		nameEn = region.Name
	}

	query := `
		INSERT INTO regions (name, name_en, name_ru, description, description_en, description_ru,
		                     geojson, color, fill_opacity, border_color, border_width)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, name_en, name_ru, description, description_en, description_ru,
		          geojson, color, fill_opacity, border_color, border_width,
		          created_at, updated_at`

	var created models.Region
	err := r.db.QueryRow(query,
		nameEn, nameEn, region.NameRu,
		region.DescriptionEn, region.DescriptionEn, region.DescriptionRu,
		region.GeoJSON, color, fillOpacity, borderColor, borderWidth,
	).Scan(
		&created.ID, &created.Name, &created.NameEn, &created.NameRu,
		&created.Description, &created.DescriptionEn, &created.DescriptionRu,
		&created.GeoJSON, &created.Color, &created.FillOpacity,
		&created.BorderColor, &created.BorderWidth,
		&created.CreatedAt, &created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create region: %w", err)
	}

	if len(region.TemplateIDs) > 0 {
		for _, templateID := range region.TemplateIDs {
			err := r.LinkToTemplate(created.ID, templateID)
			if err != nil {
				return nil, fmt.Errorf("failed to link region to template %d: %w", templateID, err)
			}
		}
	}

	return &created, nil
}

func (r *RegionRepository) Update(id int, region *models.RegionUpdate) (*models.Region, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if region.NameEn != nil {
		existing.NameEn = *region.NameEn
		existing.Name = *region.NameEn
	}
	if region.NameRu != nil {
		existing.NameRu = *region.NameRu
	}
	if region.DescriptionEn != nil {
		existing.DescriptionEn = *region.DescriptionEn
		existing.Description = *region.DescriptionEn
	}
	if region.DescriptionRu != nil {
		existing.DescriptionRu = *region.DescriptionRu
	}
	if region.GeoJSON != nil {
		existing.GeoJSON = *region.GeoJSON
	}
	if region.Color != nil {
		existing.Color = *region.Color
	}
	if region.FillOpacity != nil {
		existing.FillOpacity = *region.FillOpacity
	}
	if region.BorderColor != nil {
		existing.BorderColor = *region.BorderColor
	}
	if region.BorderWidth != nil {
		existing.BorderWidth = *region.BorderWidth
	}

	query := `
		UPDATE regions
		SET name = $2, name_en = $3, name_ru = $4,
		    description = $5, description_en = $6, description_ru = $7,
		    geojson = $8, color = $9, fill_opacity = $10,
		    border_color = $11, border_width = $12,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, name, name_en, name_ru, description, description_en, description_ru,
		          geojson, color, fill_opacity, border_color, border_width,
		          created_at, updated_at`

	var updated models.Region
	err = r.db.QueryRow(query, id,
		existing.Name, existing.NameEn, existing.NameRu,
		existing.Description, existing.DescriptionEn, existing.DescriptionRu,
		existing.GeoJSON, existing.Color, existing.FillOpacity,
		existing.BorderColor, existing.BorderWidth,
	).Scan(
		&updated.ID, &updated.Name, &updated.NameEn, &updated.NameRu,
		&updated.Description, &updated.DescriptionEn, &updated.DescriptionRu,
		&updated.GeoJSON, &updated.Color, &updated.FillOpacity,
		&updated.BorderColor, &updated.BorderWidth,
		&updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update region: %w", err)
	}

	if region.TemplateIDs != nil {
		_, err = r.db.Exec("DELETE FROM template_regions WHERE region_id = $1", id)
		if err != nil {
			return nil, fmt.Errorf("failed to clear template links: %w", err)
		}
		for _, templateID := range region.TemplateIDs {
			err := r.LinkToTemplate(id, templateID)
			if err != nil {
				return nil, fmt.Errorf("failed to link region to template %d: %w", templateID, err)
			}
		}
	}

	return &updated, nil
}

func (r *RegionRepository) Delete(id int) error {
	query := `DELETE FROM regions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete region: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("region not found")
	}

	return nil
}

func (r *RegionRepository) LinkToTemplate(regionID, templateID int) error {
	query := `
		INSERT INTO template_regions (region_id, template_id)
		VALUES ($1, $2)
		ON CONFLICT (template_id, region_id) DO NOTHING`

	_, err := r.db.Exec(query, regionID, templateID)
	if err != nil {
		return fmt.Errorf("failed to link region %d to template %d: %w", regionID, templateID, err)
	}

	return nil
}

func (r *RegionRepository) UnlinkFromTemplate(regionID, templateID int) error {
	query := `DELETE FROM template_regions WHERE region_id = $1 AND template_id = $2`

	result, err := r.db.Exec(query, regionID, templateID)
	if err != nil {
		return fmt.Errorf("failed to unlink region from template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("region-template link not found")
	}

	return nil
}

func (r *RegionRepository) GetTemplateIDsByRegion(regionID int) ([]int, error) {
	query := `SELECT template_id FROM template_regions WHERE region_id = $1 ORDER BY template_id`

	rows, err := r.db.Query(query, regionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query template IDs for region %d: %w", regionID, err)
	}
	defer rows.Close()

	var templateIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan template ID: %w", err)
		}
		templateIDs = append(templateIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over template IDs: %w", err)
	}

	return templateIDs, nil
}
