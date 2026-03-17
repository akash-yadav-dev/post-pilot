package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type AuthAccount struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	PasswordHash   string     `json:"-"`
	FailedCount    int        `json:"-"`
	LockedUntil    *time.Time `json:"-"`
}

type PasswordAuthIdentity struct {
	User
	PasswordHash string     `json:"-"`
	FailedCount  int        `json:"-"`
	LockedUntil  *time.Time `json:"-"`
}
