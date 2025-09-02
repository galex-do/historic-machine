package services

import (
	"crypto/sha256"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"

	"historical-events-backend/internal/database/repositories"
	"historical-events-backend/internal/models"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret []byte
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// RegisterUser creates a new user account
func (s *AuthService) RegisterUser(req *models.CreateUserRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := req.ToUser()
	user.PasswordHash = string(hashedPassword)

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Remove password hash from response
	user.PasswordHash = ""
	return user, nil
}

// AuthenticateUser validates credentials and returns user with token
func (s *AuthService) AuthenticateUser(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session record
	tokenHash := s.hashToken(token)
	session := &models.UserSession{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	err = s.userRepo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: failed to update last login for user %d: %v\n", user.ID, err)
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Parse and validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Check if session is still active
	tokenHash := s.hashToken(tokenString)
	session, err := s.userRepo.GetActiveSession(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	// Get user details
	user, err := s.userRepo.GetUserByID(int(userID))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Verify session belongs to user
	if session.UserID != user.ID {
		return nil, fmt.Errorf("session mismatch")
	}

	// Remove password hash from response
	user.PasswordHash = ""
	return user, nil
}

// LogoutUser deactivates the user's session
func (s *AuthService) LogoutUser(tokenString string) error {
	tokenHash := s.hashToken(tokenString)
	return s.userRepo.DeactivateSession(tokenHash)
}

// LogoutAllSessions deactivates all sessions for a user
func (s *AuthService) LogoutAllSessions(userID int) error {
	return s.userRepo.DeactivateUserSessions(userID)
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID int, req *models.ChangePasswordRequest) error {
	// Get current user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	err = s.userRepo.UpdateUserPassword(userID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Logout all sessions to force re-login with new password
	err = s.userRepo.DeactivateUserSessions(userID)
	if err != nil {
		// Log error but don't fail the password change
		fmt.Printf("Warning: failed to deactivate sessions for user %d: %v\n", userID, err)
	}

	return nil
}

// generateJWT creates a new JWT token for the user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"username":     user.Username,
		"access_level": user.AccessLevel,
		"exp":          time.Now().Add(24 * time.Hour).Unix(), // 24 hour expiry
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// hashToken creates a hash of the token for database storage
func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// CleanExpiredSessions removes expired sessions (should be called periodically)
func (s *AuthService) CleanExpiredSessions() error {
	return s.userRepo.CleanExpiredSessions()
}