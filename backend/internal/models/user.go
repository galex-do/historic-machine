package models

import (
        "time"
)

// AccessLevel represents the user access level
type AccessLevel string

const (
        AccessLevelGuest AccessLevel = "guest"
        AccessLevelUser  AccessLevel = "user"
        AccessLevelEditor AccessLevel = "editor"
        AccessLevelAdmin AccessLevel = "admin"
        AccessLevelSuper AccessLevel = "super"
)

// User represents a user in the system
type User struct {
        ID           int         `json:"id"`
        Username     string      `json:"username"`
        Email        string      `json:"email,omitempty"`
        PasswordHash string      `json:"-"` // Exclude password from JSON
        AccessLevel  AccessLevel `json:"access_level"`
        IsActive     bool        `json:"is_active"`
        CreatedAt    time.Time   `json:"created_at"`
        UpdatedAt    time.Time   `json:"updated_at"`
        LastLogin    *time.Time  `json:"last_login,omitempty"`
}

// UserSession represents a user session with JWT token
type UserSession struct {
        ID         int        `json:"id"`
        UserID     int        `json:"user_id"`
        TokenHash  string     `json:"-"` // Exclude token hash from JSON
        ExpiresAt  time.Time  `json:"expires_at"`
        CreatedAt  time.Time  `json:"created_at"`
        LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
        EndedAt    *time.Time `json:"ended_at,omitempty"`
        IsActive   bool       `json:"is_active"`
}

// AnonymousSession represents an anonymous user session
type AnonymousSession struct {
        SessionID  string     `json:"session_id"`
        CreatedAt  time.Time  `json:"created_at"`
        LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
        EndedAt    *time.Time `json:"ended_at,omitempty"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
        Username    string      `json:"username" validate:"required,min=3,max=50"`
        Email       string      `json:"email" validate:"email"`
        Password    string      `json:"password" validate:"required,min=8,max=72"`
        AccessLevel AccessLevel `json:"access_level,omitempty"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
        Username string `json:"username" validate:"required"`
        Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
        User  User   `json:"user"`
        Token string `json:"token"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
        Email       string      `json:"email,omitempty" validate:"email"`
        AccessLevel AccessLevel `json:"access_level,omitempty"`
        IsActive    *bool       `json:"is_active,omitempty"`
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
        CurrentPassword string `json:"current_password" validate:"required"`
        NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
}

// UserProfile represents public user profile information
type UserProfile struct {
        ID          int         `json:"id"`
        Username    string      `json:"username"`
        AccessLevel AccessLevel `json:"access_level"`
        CreatedAt   time.Time   `json:"created_at"`
        LastLogin   *time.Time  `json:"last_login,omitempty"`
}

// ToUser converts CreateUserRequest to User (without password processing)
func (req *CreateUserRequest) ToUser() *User {
        accessLevel := req.AccessLevel
        if accessLevel == "" {
                accessLevel = AccessLevelUser
        }

        return &User{
                Username:    req.Username,
                Email:       req.Email,
                AccessLevel: accessLevel,
                IsActive:    true,
                CreatedAt:   time.Now(),
                UpdatedAt:   time.Now(),
        }
}

// ToProfile converts User to UserProfile
func (u *User) ToProfile() *UserProfile {
        return &UserProfile{
                ID:          u.ID,
                Username:    u.Username,
                AccessLevel: u.AccessLevel,
                CreatedAt:   u.CreatedAt,
                LastLogin:   u.LastLogin,
        }
}

// CanCreateEvents checks if user has permission to create events
func (u *User) CanCreateEvents() bool {
        return u.AccessLevel != AccessLevelGuest && u.IsActive
}

// CanManageUsers checks if user has permission to manage other users
func (u *User) CanManageUsers() bool {
        return (u.AccessLevel == AccessLevelAdmin || u.AccessLevel == AccessLevelSuper) && u.IsActive
}

// CanManageSystem checks if user has system-level permissions
func (u *User) CanManageSystem() bool {
        return u.AccessLevel == AccessLevelSuper && u.IsActive
}