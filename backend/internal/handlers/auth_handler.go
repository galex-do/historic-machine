package handlers

import (
        "encoding/json"
        "net/http"
        "strings"

        "historical-events-backend/internal/models"
        "historical-events-backend/internal/services"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
        authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
        return &AuthHandler{
                authService: authService,
        }
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        var loginReq models.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Basic validation
        if loginReq.Username == "" || loginReq.Password == "" {
                http.Error(w, "Username and password are required", http.StatusBadRequest)
                return
        }

        response, err := h.authService.AuthenticateUser(&loginReq)
        if err != nil {
                http.Error(w, "Invalid credentials", http.StatusUnauthorized)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

// Register handles user registration requests
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        var createReq models.CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Basic validation
        if createReq.Username == "" || createReq.Password == "" {
                http.Error(w, "Username and password are required", http.StatusBadRequest)
                return
        }

        if len(createReq.Username) < 3 || len(createReq.Username) > 50 {
                http.Error(w, "Username must be between 3 and 50 characters", http.StatusBadRequest)
                return
        }

        if len(createReq.Password) < 8 {
                http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
                return
        }

        // Only allow admin users to set access levels other than 'user'
        if createReq.AccessLevel != "" && createReq.AccessLevel != models.AccessLevelUser {
                // For now, disallow setting custom access levels during registration
                // This should be done by admin users later
                createReq.AccessLevel = models.AccessLevelUser
        }

        user, err := h.authService.RegisterUser(&createReq)
        if err != nil {
                if strings.Contains(err.Error(), "username already exists") {
                        http.Error(w, "Username already exists", http.StatusConflict)
                        return
                }
                http.Error(w, "Failed to create user", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
}

// Logout handles user logout requests
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        token := h.extractTokenFromHeader(r)
        if token == "" {
                http.Error(w, "No token provided", http.StatusUnauthorized)
                return
        }

        err := h.authService.LogoutUser(token)
        if err != nil {
                http.Error(w, "Failed to logout", http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}

// Me returns the current authenticated user's information
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        user := h.getCurrentUser(r)
        if user == nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user.ToProfile())
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        user := h.getCurrentUser(r)
        if user == nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
        }

        var changeReq models.ChangePasswordRequest
        if err := json.NewDecoder(r.Body).Decode(&changeReq); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        if changeReq.CurrentPassword == "" || changeReq.NewPassword == "" {
                http.Error(w, "Current password and new password are required", http.StatusBadRequest)
                return
        }

        if len(changeReq.NewPassword) < 8 {
                http.Error(w, "New password must be at least 8 characters", http.StatusBadRequest)
                return
        }

        err := h.authService.ChangePassword(user.ID, &changeReq)
        if err != nil {
                if strings.Contains(err.Error(), "current password is incorrect") {
                        http.Error(w, "Current password is incorrect", http.StatusBadRequest)
                        return
                }
                http.Error(w, "Failed to change password", http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}

// AuthMiddleware validates JWT tokens and adds user to request context
func (h *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                token := h.extractTokenFromHeader(r)
                if token == "" {
                        http.Error(w, "No token provided", http.StatusUnauthorized)
                        return
                }

                user, err := h.authService.ValidateToken(token)
                if err != nil {
                        http.Error(w, "Invalid token", http.StatusUnauthorized)
                        return
                }

                // Add user to request context
                ctx := r.Context()
                ctx = setUserInContext(ctx, user)
                r = r.WithContext(ctx)

                next(w, r)
        }
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func (h *AuthHandler) OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                token := h.extractTokenFromHeader(r)
                if token != "" {
                        user, err := h.authService.ValidateToken(token)
                        if err == nil {
                                // Add user to request context
                                ctx := r.Context()
                                ctx = setUserInContext(ctx, user)
                                r = r.WithContext(ctx)
                        }
                }

                next(w, r)
        }
}

// RequireAccessLevel middleware that requires specific access level
func (h *AuthHandler) RequireAccessLevel(accessLevel models.AccessLevel) func(http.HandlerFunc) http.HandlerFunc {
        return func(next http.HandlerFunc) http.HandlerFunc {
                return h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
                        user := h.getCurrentUser(r)
                        if user == nil {
                                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                                return
                        }

                        // Check if user has required access level
                        if !h.hasAccessLevel(user.AccessLevel, accessLevel) {
                                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                                return
                        }

                        next(w, r)
                })
        }
}

// extractTokenFromHeader extracts JWT token from Authorization header
func (h *AuthHandler) extractTokenFromHeader(r *http.Request) string {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
                return ""
        }

        // Expect format: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
                return ""
        }

        return parts[1]
}

// getCurrentUser gets the authenticated user from request context
func (h *AuthHandler) getCurrentUser(r *http.Request) *models.User {
        return getUserFromContext(r.Context())
}

// hasAccessLevel checks if user's access level is sufficient
func (h *AuthHandler) hasAccessLevel(userLevel, requiredLevel models.AccessLevel) bool {
        accessLevels := map[models.AccessLevel]int{
                models.AccessLevelGuest: 0,
                models.AccessLevelUser:  1,
                models.AccessLevelAdmin: 2,
                models.AccessLevelSuper: 3,
        }

        userLevelNum, ok1 := accessLevels[userLevel]
        requiredLevelNum, ok2 := accessLevels[requiredLevel]

        if !ok1 || !ok2 {
                return false
        }

        return userLevelNum >= requiredLevelNum
}

// User Management Methods (for admin interfaces)

// GetAllUsers returns all users (super users only)
func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        users, err := h.authService.GetAllUsers()
        if err != nil {
                http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
                "data": users,
                "message": "Users retrieved successfully",
        })
}

// CreateUser creates a new user (super users only)
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        var createReq models.CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Basic validation
        if createReq.Username == "" || createReq.Password == "" {
                http.Error(w, "Username and password are required", http.StatusBadRequest)
                return
        }

        if len(createReq.Username) < 3 || len(createReq.Username) > 50 {
                http.Error(w, "Username must be between 3 and 50 characters", http.StatusBadRequest)
                return
        }

        if len(createReq.Password) < 8 {
                http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
                return
        }

        // Default to user access level if not specified
        if createReq.AccessLevel == "" {
                createReq.AccessLevel = models.AccessLevelUser
        }

        user, err := h.authService.RegisterUser(&createReq)
        if err != nil {
                if strings.Contains(err.Error(), "username already exists") {
                        http.Error(w, "Username already exists", http.StatusConflict)
                        return
                }
                http.Error(w, "Failed to create user", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]interface{}{
                "data": user.ToProfile(),
                "message": "User created successfully",
        })
}

// UpdateUser updates an existing user (super users only)
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPut {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Extract user ID from URL path
        userID := strings.TrimPrefix(r.URL.Path, "/api/users/")
        if userID == "" {
                http.Error(w, "User ID is required", http.StatusBadRequest)
                return
        }

        var updateReq models.UpdateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        user, err := h.authService.UpdateUser(userID, &updateReq)
        if err != nil {
                if strings.Contains(err.Error(), "user not found") {
                        http.Error(w, "User not found", http.StatusNotFound)
                        return
                }
                http.Error(w, "Failed to update user", http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
                "data": user.ToProfile(),
                "message": "User updated successfully",
        })
}

// DeleteUser deletes a user (super users only)
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodDelete {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Extract user ID from URL path
        userID := strings.TrimPrefix(r.URL.Path, "/api/users/")
        if userID == "" {
                http.Error(w, "User ID is required", http.StatusBadRequest)
                return
        }

        err := h.authService.DeleteUser(userID)
        if err != nil {
                if strings.Contains(err.Error(), "user not found") {
                        http.Error(w, "User not found", http.StatusNotFound)
                        return
                }
                if strings.Contains(err.Error(), "cannot delete") {
                        http.Error(w, "Cannot delete this user", http.StatusForbidden)
                        return
                }
                http.Error(w, "Failed to delete user", http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}