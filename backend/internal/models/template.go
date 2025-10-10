package models

// DateTemplateGroup represents a group of date templates
type DateTemplateGroup struct {
        ID             int    `json:"id"`
        Name           string `json:"name"`
        Description    string `json:"description"`
        NameEn         string `json:"name_en"`
        NameRu         string `json:"name_ru"`
        DescriptionEn  string `json:"description_en"`
        DescriptionRu  string `json:"description_ru"`
        DisplayOrder   int    `json:"display_order"`
}

// DateTemplate represents a date range template with display formatting
type DateTemplate struct {
        ID               int    `json:"id"`
        GroupID          int    `json:"group_id"`
        GroupName        string `json:"group_name"`
        GroupNameEn      string `json:"group_name_en"`
        GroupNameRu      string `json:"group_name_ru"`
        Name             string `json:"name"`
        Description      string `json:"description"`
        NameEn           string `json:"name_en"`
        NameRu           string `json:"name_ru"`
        DescriptionEn    string `json:"description_en"`
        DescriptionRu    string `json:"description_ru"`
        StartDate        string `json:"start_date"`
        StartEra         string `json:"start_era"`
        EndDate          string `json:"end_date"`
        EndEra           string `json:"end_era"`
        DisplayOrder     int    `json:"display_order"`
        StartDisplayDate string `json:"start_display_date"`
        EndDisplayDate   string `json:"end_display_date"`
}