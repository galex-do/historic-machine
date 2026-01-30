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

// GetAllUsers retrieves all active users
func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
        query := `
                SELECT id, username, email, password_hash, access_level, is_active, 
                       created_at, updated_at, last_login
                FROM users 
                WHERE is_active = true
                ORDER BY created_at DESC`

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to query users: %w", err)
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
                        &user.PasswordHash,
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

        if err := rows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating users: %w", err)
        }

        return users, nil
}

// DeleteUser deactivates a user (soft delete)
func (r *UserRepository) DeleteUser(userID int) error {
        query := `
                UPDATE users 
                SET is_active = false, updated_at = $2
                WHERE id = $1`

        result, err := r.db.Exec(query, userID, time.Now())
        if err != nil {
                return fmt.Errorf("failed to delete user: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get affected rows: %w", err)
        }

        if rowsAffected == 0 {
                return fmt.Errorf("user not found")
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
                INSERT INTO user_sessions (user_id, token_hash, expires_at, created_at, last_seen_at, is_active)
                VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id`

        now := time.Now()
        session.LastSeenAt = &now

        err := r.db.QueryRow(
                query,
                session.UserID,
                session.TokenHash,
                session.ExpiresAt,
                session.CreatedAt,
                session.LastSeenAt,
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
                SELECT id, user_id, token_hash, expires_at, created_at, last_seen_at, ended_at, is_active
                FROM user_sessions 
                WHERE token_hash = $1 AND is_active = true AND expires_at > NOW()`

        session := &models.UserSession{}
        var lastSeenAt, endedAt sql.NullTime

        err := r.db.QueryRow(query, tokenHash).Scan(
                &session.ID,
                &session.UserID,
                &session.TokenHash,
                &session.ExpiresAt,
                &session.CreatedAt,
                &lastSeenAt,
                &endedAt,
                &session.IsActive,
        )

        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("session not found or expired")
                }
                return nil, fmt.Errorf("failed to get session: %w", err)
        }

        if lastSeenAt.Valid {
                session.LastSeenAt = &lastSeenAt.Time
        }
        if endedAt.Valid {
                session.EndedAt = &endedAt.Time
        }

        return session, nil
}

// DeactivateSession deactivates a session by token hash
func (r *UserRepository) DeactivateSession(tokenHash string) error {
        query := `
                UPDATE user_sessions 
                SET is_active = false, ended_at = $2
                WHERE token_hash = $1`

        _, err := r.db.Exec(query, tokenHash, time.Now())
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

// UpdateSessionHeartbeat updates the last_seen_at timestamp for an active session
func (r *UserRepository) UpdateSessionHeartbeat(tokenHash string) error {
        query := `
                UPDATE user_sessions 
                SET last_seen_at = $2
                WHERE token_hash = $1 AND is_active = true AND expires_at > NOW()`

        result, err := r.db.Exec(query, tokenHash, time.Now())
        if err != nil {
                return fmt.Errorf("failed to update session heartbeat: %w", err)
        }

        rows, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to check affected rows: %w", err)
        }

        if rows == 0 {
                return fmt.Errorf("session not found or expired")
        }

        return nil
}

// CreateOrUpdateAnonymousSession creates a new anonymous session or updates heartbeat if it exists
func (r *UserRepository) CreateOrUpdateAnonymousSession(sessionID string) error {
        query := `
                INSERT INTO anonymous_sessions (session_id, created_at, last_seen_at)
                VALUES ($1, $2, $3)
                ON CONFLICT (session_id)
                DO UPDATE SET last_seen_at = $3`

        now := time.Now()
        _, err := r.db.Exec(query, sessionID, now, now)
        if err != nil {
                return fmt.Errorf("failed to create or update anonymous session: %w", err)
        }

        return nil
}

// CountActiveAnonymousSessions counts sessions active in the last 5 minutes
func (r *UserRepository) CountActiveAnonymousSessions() (int, error) {
        query := `
                SELECT COUNT(*)
                FROM anonymous_sessions
                WHERE last_seen_at > NOW() - INTERVAL '5 minutes'`

        var count int
        err := r.db.QueryRow(query).Scan(&count)
        if err != nil {
                return 0, fmt.Errorf("failed to count active anonymous sessions: %w", err)
        }

        return count, nil
}

// CountAllActiveSessions counts both authenticated and anonymous sessions active in the last 5 minutes
func (r *UserRepository) CountAllActiveSessions() (int, error) {
        query := `
                SELECT 
                        (SELECT COUNT(*) 
                         FROM user_sessions 
                         WHERE is_active = true 
                           AND last_seen_at > NOW() - INTERVAL '5 minutes'
                           AND expires_at > NOW())
                        +
                        (SELECT COUNT(*) 
                         FROM anonymous_sessions 
                         WHERE last_seen_at > NOW() - INTERVAL '5 minutes')
                AS total_active`

        var count int
        err := r.db.QueryRow(query).Scan(&count)
        if err != nil {
                return 0, fmt.Errorf("failed to count all active sessions: %w", err)
        }

        return count, nil
}

// GetPeakConcurrentSessions retrieves the current peak value
func (r *UserRepository) GetPeakConcurrentSessions() (int, error) {
        query := `SELECT peak_concurrent_sessions FROM peak_stats WHERE id = 1`

        var peak int
        err := r.db.QueryRow(query).Scan(&peak)
        if err != nil {
                if err == sql.ErrNoRows {
                        return 0, nil
                }
                return 0, fmt.Errorf("failed to get peak concurrent sessions: %w", err)
        }

        return peak, nil
}

// UpdatePeakIfHigher updates the peak if the current count is higher
func (r *UserRepository) UpdatePeakIfHigher(currentCount int) error {
        query := `
                UPDATE peak_stats
                SET peak_concurrent_sessions = $1, updated_at = NOW()
                WHERE id = 1 AND peak_concurrent_sessions < $1`

        _, err := r.db.Exec(query, currentCount)
        if err != nil {
                return fmt.Errorf("failed to update peak concurrent sessions: %w", err)
        }

        return nil
}

