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

// GetNameForLocale returns the name in the specified locale
func (g *DateTemplateGroup) GetNameForLocale(locale string) string {
        if locale == "ru" {
                return g.NameRu
        }
        return g.NameEn
}

// GetDescriptionForLocale returns the description in the specified locale
func (g *DateTemplateGroup) GetDescriptionForLocale(locale string) string {
        if locale == "ru" {
                return g.DescriptionRu
        }
        return g.DescriptionEn
}

// PopulateLegacyFields sets name and description based on locale
func (g *DateTemplateGroup) PopulateLegacyFields(locale string) {
        g.Name = g.GetNameForLocale(locale)
        g.Description = g.GetDescriptionForLocale(locale)
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

// GetNameForLocale returns the name in the specified locale
func (t *DateTemplate) GetNameForLocale(locale string) string {
        if locale == "ru" {
                return t.NameRu
        }
        return t.NameEn
}

// GetDescriptionForLocale returns the description in the specified locale
func (t *DateTemplate) GetDescriptionForLocale(locale string) string {
        if locale == "ru" {
                return t.DescriptionRu
        }
        return t.DescriptionEn
}

// GetGroupNameForLocale returns the group name in the specified locale
func (t *DateTemplate) GetGroupNameForLocale(locale string) string {
        if locale == "ru" {
                return t.GroupNameRu
        }
        return t.GroupNameEn
}

// PopulateLegacyFields sets name, description, and group_name based on locale
func (t *DateTemplate) PopulateLegacyFields(locale string) {
        t.Name = t.GetNameForLocale(locale)
        t.Description = t.GetDescriptionForLocale(locale)
        t.GroupName = t.GetGroupNameForLocale(locale)
}