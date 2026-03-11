package models

import (
	"encoding/json"
	"time"
)

type Region struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	NameEn        string           `json:"name_en"`
	NameRu        string           `json:"name_ru"`
	Description   string           `json:"description"`
	DescriptionEn string           `json:"description_en"`
	DescriptionRu string           `json:"description_ru"`
	GeoJSON       json.RawMessage  `json:"geojson"`
	Color         string           `json:"color"`
	FillOpacity   float32          `json:"fill_opacity"`
	BorderColor   string           `json:"border_color"`
	BorderWidth   float32          `json:"border_width"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (r *Region) GetNameForLocale(locale string) string {
	if locale == "ru" {
		return r.NameRu
	}
	return r.NameEn
}

func (r *Region) GetDescriptionForLocale(locale string) string {
	if locale == "ru" {
		return r.DescriptionRu
	}
	return r.DescriptionEn
}

func (r *Region) PopulateLegacyFields(locale string) {
	r.Name = r.GetNameForLocale(locale)
	r.Description = r.GetDescriptionForLocale(locale)
}

type TemplateRegion struct {
	TemplateID int `json:"template_id"`
	RegionID   int `json:"region_id"`
}

type RegionCreate struct {
	Name          string          `json:"name"`
	NameEn        string          `json:"name_en"`
	NameRu        string          `json:"name_ru"`
	Description   string          `json:"description"`
	DescriptionEn string          `json:"description_en"`
	DescriptionRu string          `json:"description_ru"`
	GeoJSON       json.RawMessage `json:"geojson" validate:"required"`
	Color         string          `json:"color"`
	FillOpacity   *float32        `json:"fill_opacity"`
	BorderColor   string          `json:"border_color"`
	BorderWidth   *float32        `json:"border_width"`
	TemplateIDs   []int           `json:"template_ids,omitempty"`
}

type RegionUpdate struct {
	Name          *string          `json:"name,omitempty"`
	NameEn        *string          `json:"name_en,omitempty"`
	NameRu        *string          `json:"name_ru,omitempty"`
	Description   *string          `json:"description,omitempty"`
	DescriptionEn *string          `json:"description_en,omitempty"`
	DescriptionRu *string          `json:"description_ru,omitempty"`
	GeoJSON       *json.RawMessage `json:"geojson,omitempty"`
	Color         *string          `json:"color,omitempty"`
	FillOpacity   *float32         `json:"fill_opacity,omitempty"`
	BorderColor   *string          `json:"border_color,omitempty"`
	BorderWidth   *float32         `json:"border_width,omitempty"`
	TemplateIDs   []int            `json:"template_ids,omitempty"`
}
