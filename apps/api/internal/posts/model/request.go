package model

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Content     string     `json:"content"      binding:"required,min=1,max=5000"`
	Platforms   []string   `json:"platforms"    binding:"required,min=1"`
	ScheduledAt *time.Time `json:"scheduled_at" binding:"omitempty"`
}

type UpdatePostRequest struct {
	Content     string     `json:"content"      binding:"omitempty,min=1,max=5000"`
	ScheduledAt *time.Time `json:"scheduled_at" binding:"omitempty"`
}

type PostResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Content     string     `json:"content"`
	Status      PostStatus `json:"status"`
	Platforms   []string   `json:"platforms"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
