package repositories

import (
        "database/sql"

        "historical-events-backend/internal/models"
)

type SupportRepository struct {
        db *sql.DB
}

func NewSupportRepository(db *sql.DB) *SupportRepository {
        return &SupportRepository{db: db}
}

func (r *SupportRepository) GetActiveCredentials() ([]models.SupportCredential, error) {
        query := `
                SELECT id, name, value, is_url, display_order, is_active, created_at, updated_at
                FROM support_credentials
                WHERE is_active = true
                ORDER BY display_order ASC, name ASC
        `

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var credentials []models.SupportCredential
        for rows.Next() {
                var c models.SupportCredential
                err := rows.Scan(
                        &c.ID,
                        &c.Name,
                        &c.Value,
                        &c.IsURL,
                        &c.DisplayOrder,
                        &c.IsActive,
                        &c.CreatedAt,
                        &c.UpdatedAt,
                )
                if err != nil {
                        return nil, err
                }
                credentials = append(credentials, c)
        }

        return credentials, nil
}

func (r *SupportRepository) Create(credential *models.SupportCredential) error {
        query := `
                INSERT INTO support_credentials (name, value, is_url, display_order, is_active)
                VALUES ($1, $2, $3, $4, $5)
                RETURNING id, created_at, updated_at
        `

        return r.db.QueryRow(
                query,
                credential.Name,
                credential.Value,
                credential.IsURL,
                credential.DisplayOrder,
                credential.IsActive,
        ).Scan(&credential.ID, &credential.CreatedAt, &credential.UpdatedAt)
}

func (r *SupportRepository) Update(credential *models.SupportCredential) error {
        query := `
                UPDATE support_credentials
                SET name = $1, value = $2, is_url = $3, display_order = $4, is_active = $5, updated_at = CURRENT_TIMESTAMP
                WHERE id = $6
                RETURNING updated_at
        `

        return r.db.QueryRow(
                query,
                credential.Name,
                credential.Value,
                credential.IsURL,
                credential.DisplayOrder,
                credential.IsActive,
                credential.ID,
        ).Scan(&credential.UpdatedAt)
}

func (r *SupportRepository) Delete(id int) error {
        query := `DELETE FROM support_credentials WHERE id = $1`
        _, err := r.db.Exec(query, id)
        return err
}
