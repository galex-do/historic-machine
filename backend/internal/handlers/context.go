package handlers

import (
	"context"

	"historical-events-backend/internal/models"
)

// Context key type to avoid collisions
type contextKey string

const userContextKey contextKey = "user"

// setUserInContext adds a user to the request context
func setUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// getUserFromContext retrieves a user from the request context
func getUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}