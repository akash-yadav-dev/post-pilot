package model

import (
	"time"

	"github.com/google/uuid"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusScheduled PostStatus = "scheduled"
	PostStatusPublished PostStatus = "published"
	PostStatusFailed    PostStatus = "failed"
)

type Post struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Content     string     `json:"content"`
	Status      PostStatus `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PostTarget struct {
	ID       uuid.UUID  `json:"id"`
	PostID   uuid.UUID  `json:"post_id"`
	Platform string     `json:"platform"`
	Status   PostStatus `json:"status"`
	ErrorMsg string     `json:"error_message,omitempty"`
}
