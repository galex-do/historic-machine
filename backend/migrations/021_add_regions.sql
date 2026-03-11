-- +goose Up
CREATE TABLE regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    name_ru VARCHAR(255),
    description TEXT,
    description_en TEXT,
    description_ru TEXT,
    geojson JSONB NOT NULL,
    color VARCHAR(7) DEFAULT '#4f46e5',
    fill_opacity REAL DEFAULT 0.2,
    border_color VARCHAR(7) DEFAULT '#4f46e5',
    border_width REAL DEFAULT 2,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE template_regions (
    template_id INTEGER NOT NULL REFERENCES date_templates(id) ON DELETE CASCADE,
    region_id INTEGER NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, region_id)
);

CREATE INDEX idx_template_regions_template ON template_regions(template_id);
CREATE INDEX idx_template_regions_region ON template_regions(region_id);

-- +goose Down
DROP TABLE IF EXISTS template_regions;
DROP TABLE IF EXISTS regions;
