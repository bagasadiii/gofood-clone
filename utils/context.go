package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ctxKey string

const (
	UserIDKey   ctxKey = "user_id_key"
	UsernameKey ctxKey = "username_key"
	RoleKey     ctxKey = "role_key"
)

type ContextValues struct {
	UserID   uuid.UUID
	Username string
	Role     string
}

func CheckContextValue(ctx context.Context) (*ContextValues, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if userID == uuid.Nil || !ok {
		return nil, fmt.Errorf("missing or invalid user ID: %w", ErrUnauthorized)
	}

	username, ok := ctx.Value(UsernameKey).(string)
	if username == "" || !ok {
		return nil, fmt.Errorf("missing or invalid username: %w", ErrUnauthorized)
	}

	role, ok := ctx.Value(RoleKey).(string)
	if role == "" || !ok {
		return nil, fmt.Errorf("missing or invalid role: %w", ErrUnauthorized)
	}

	return &ContextValues{
		UserID:   userID,
		Username: username,
		Role:     role,
	}, nil
}
