package repositories

import (
        "database/sql"
        "fmt"
        "time"

        "historical-events-backend/internal/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
        db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
        return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
        query := `
                INSERT INTO users (username, email, password_hash, access_level, is_active, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id`

        err := r.db.QueryRow(
                query,
                user.Username,
                user.Email,
                user.PasswordHash,
                user.AccessLevel,
                user.IsActive,
                user.CreatedAt,
                user.UpdatedAt,
        ).Scan(&user.ID)

        if err != nil {
                return fmt.Errorf("failed to create user: %w", err)
        }

        return nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
        query := `
                SELECT id, username, email, password_hash, access_level, is_active, 
                       created_at, updated_at, last_login
                FROM users 
                WHERE username = $1 AND is_active = true`

        user := &models.User{}
        var lastLogin sql.NullTime

        err := r.db.QueryRow(query, username).Scan(
                &user.ID,
                &user.Username,
                &user.Email,
                &user.PasswordHash,
                &user.AccessLevel,
                &user.IsActive,
                &user.CreatedAt,
                &user.UpdatedAt,
                &lastLogin,
        )

        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("user not found")
                }
                return nil, fmt.Errorf("failed to get user: %w", err)
        }

        if lastLogin.Valid {
                user.LastLogin = &lastLogin.Time
        }

        return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
        query := `
                SELECT id, username, email, password_hash, access_level, is_active, 
                       created_at, updated_at, last_login
                FROM users 
                WHERE id = $1 AND is_active = true`

        user := &models.User{}
        var lastLogin sql.NullTime

        err := r.db.QueryRow(query, id).Scan(
                &user.ID,
                &user.Username,
                &user.Email,
                &user.PasswordHash,
                &user.AccessLevel,
                &user.IsActive,
                &user.CreatedAt,
                &user.UpdatedAt,
                &lastLogin,
        )

        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("user not found")
                }
                return nil, fmt.Errorf("failed to get user: %w", err)
        }

        if lastLogin.Valid {
                user.LastLogin = &lastLogin.Time
        }

        return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *models.User) error {
        query := `
                UPDATE users 
                SET email = $2, access_level = $3, is_active = $4, updated_at = $5
                WHERE id = $1`

        user.UpdatedAt = time.Now()

        _, err := r.db.Exec(
                query,
                user.ID,
                user.Email,
                user.AccessLevel,
                user.IsActive,
                user.UpdatedAt,
        )

        if err != nil {
                return fmt.Errorf("failed to update user: %w", err)
        }

        return nil
}

// UpdateUserPassword updates user password hash
func (r *UserRepository) UpdateUserPassword(userID int, passwordHash string) error {
        query := `
                UPDATE users 
                SET password_hash = $2, updated_at = $3
                WHERE id = $1`

        _, err := r.db.Exec(query, userID, passwordHash, time.Now())
        if err != nil {
                return fmt.Errorf("failed to update password: %w", err)
        }

        return nil
}

// UpdateLastLogin updates the user's last login time
func (r *UserRepository) UpdateLastLogin(userID int) error {
        query := `
                UPDATE users 
                SET last_login = $2, updated_at = $3
                WHERE id = $1`

        now := time.Now()
        _, err := r.db.Exec(query, userID, now, now)
        if err != nil {
                return fmt.Errorf("failed to update last login: %w", err)
        }

        return nil
}

// CreateSession creates a new user session
func (r *UserRepository) CreateSession(session *models.UserSession) error {
        query := `
                INSERT INTO user_sessions (user_id, token_hash, expires_at, created_at, is_active)
                VALUES ($1, $2, $3, $4, $5)
                RETURNING id`

        err := r.db.QueryRow(
                query,
                session.UserID,
                session.TokenHash,
                session.ExpiresAt,
                session.CreatedAt,
                session.IsActive,
        ).Scan(&session.ID)

        if err != nil {
                return fmt.Errorf("failed to create session: %w", err)
        }

        return nil
}

// GetActiveSession retrieves an active session by token hash
func (r *UserRepository) GetActiveSession(tokenHash string) (*models.UserSession, error) {
        query := `
                SELECT id, user_id, token_hash, expires_at, created_at, is_active
                FROM user_sessions 
                WHERE token_hash = $1 AND is_active = true AND expires_at > NOW()`

        session := &models.UserSession{}

        err := r.db.QueryRow(query, tokenHash).Scan(
                &session.ID,
                &session.UserID,
                &session.TokenHash,
                &session.ExpiresAt,
                &session.CreatedAt,
                &session.IsActive,
        )

        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("session not found or expired")
                }
                return nil, fmt.Errorf("failed to get session: %w", err)
        }

        return session, nil
}

// DeactivateSession deactivates a session by token hash
func (r *UserRepository) DeactivateSession(tokenHash string) error {
        query := `
                UPDATE user_sessions 
                SET is_active = false
                WHERE token_hash = $1`

        _, err := r.db.Exec(query, tokenHash)
        if err != nil {
                return fmt.Errorf("failed to deactivate session: %w", err)
        }

        return nil
}

// DeactivateUserSessions deactivates all sessions for a user
func (r *UserRepository) DeactivateUserSessions(userID int) error {
        query := `
                UPDATE user_sessions 
                SET is_active = false
                WHERE user_id = $1`

        _, err := r.db.Exec(query, userID)
        if err != nil {
                return fmt.Errorf("failed to deactivate user sessions: %w", err)
        }

        return nil
}

// CleanExpiredSessions removes expired sessions from database
func (r *UserRepository) CleanExpiredSessions() error {
        query := `
                DELETE FROM user_sessions 
                WHERE expires_at < NOW() OR is_active = false`

        _, err := r.db.Exec(query)
        if err != nil {
                return fmt.Errorf("failed to clean expired sessions: %w", err)
        }

        return nil
}

// ListUsers retrieves all users (admin function)
func (r *UserRepository) ListUsers() ([]*models.User, error) {
        query := `
                SELECT id, username, email, access_level, is_active, 
                       created_at, updated_at, last_login
                FROM users 
                ORDER BY created_at DESC`

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to list users: %w", err)
        }
        defer rows.Close()

        var users []*models.User
        for rows.Next() {
                user := &models.User{}
                var lastLogin sql.NullTime

                err := rows.Scan(
                        &user.ID,
                        &user.Username,
                        &user.Email,
                        &user.AccessLevel,
                        &user.IsActive,
                        &user.CreatedAt,
                        &user.UpdatedAt,
                        &lastLogin,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan user: %w", err)
                }

                if lastLogin.Valid {
                        user.LastLogin = &lastLogin.Time
                }

                users = append(users, user)
        }

        return users, nil
}